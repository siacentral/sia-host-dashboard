package sync

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/siacentral/sia-host-dashboard/dashboard/cache"
	"github.com/siacentral/sia-host-dashboard/dashboard/persist"
	"github.com/siacentral/sia-host-dashboard/dashboard/types"
	"gitlab.com/NebulousLabs/Sia/node/api"
	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

const (
	contractStatusSucceeded  = "obligationSucceeded"
	contractStatusFailed     = "obligationFailed"
	contractStatusUnresolved = "obligationUnresolved"
)

var (
	blockMetaCache = make(map[uint64]blockMeta)
	blockMetaMu    sync.Mutex
)

type (
	blockMeta struct {
		ID        string
		Timestamp time.Time
	}
)

//merges fields from the local contract db and the blockchain
type mergedContract struct {
	ID                     string            `json:"id"`
	BlockID                string            `json:"blockID"`
	TransactionID          string            `json:"transactionID"`
	MerkleRoot             string            `json:"merkleRoot"`
	UnlockHash             string            `json:"unlockHash"`
	Status                 string            `json:"status"`
	RevisionNumber         uint64            `json:"revisionNumber"`
	NegotiationHeight      uint64            `json:"negotiationHeight"`
	ExpirationHeight       uint64            `json:"expirationHeight"`
	ProofDeadline          uint64            `json:"proofDeadline"`
	ProofHeight            uint64            `json:"proofHeight"`
	DataSize               uint64            `json:"fileSize"`
	ProofConfirmed         bool              `json:"proofConfirmed"`
	NegotiationTimestamp   time.Time         `json:"negotiationTimestamp"`
	ExpirationTimestamp    time.Time         `json:"expirationTimestamp"`
	ProofDeadlineTimestamp time.Time         `json:"proofDeadlineTimestamp"`
	ProofTimestamp         time.Time         `json:"proofTimestamp"`
	Payout                 siatypes.Currency `json:"payout"`
	LockedCollateral       siatypes.Currency `json:"lockedCollateral"`
	PotentialRevenue       siatypes.Currency `json:"potentialRevenue"`
	EarnedRevenue          types.BigNumber   `json:"earnedRevenue"`
	LostRevenue            siatypes.Currency `json:"lostRevenue"`
	BurntCollateral        siatypes.Currency `json:"burntCollateral"`
}

func syncContracts() error {
	cache.ClearAlerts(AlertSyncError)
	contracts, err := getContracts()

	if err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host contracts. Check your Sia connection",
			Type:     "sync",
		})
		return err
	}

	if len(contracts) == 0 {
		return nil
	}

	syncHostMeta(contracts)
	syncHostSnapshots(contracts)

	return nil
}

func snapshotID(timestamp time.Time) uint64 {
	return uint64(snapshotTime(timestamp).Unix())
}

func snapshotTime(timestamp time.Time) time.Time {
	return timestamp.UTC().Truncate(time.Hour)
}

func getSyncedHeight() (height uint64, err error) {
	consensus, err := apiClient.ConsensusGet()
	if err != nil {
		return 0, fmt.Errorf("error getting consensus: %w", err)
	}

	return uint64(consensus.Height), nil
}

func txnConfirmed(id siatypes.TransactionID) (bool, error) {
	var resp api.TpoolConfirmedGET
	err := apiClient.Get("/tpool/confirmed/"+url.PathEscape(id.String()), &resp)
	return resp.Confirmed, err
}

func getBlockMeta(height uint64) (time.Time, error) {
	blockMetaMu.Lock()

	defer blockMetaMu.Unlock()

	if block, exists := blockMetaCache[height]; exists {
		return block.Timestamp, nil
	}

	block, err := apiClient.ConsensusBlocksHeightGet(siatypes.BlockHeight(height))
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting block: %w", err)
	}

	blockMetaCache[height] = blockMeta{
		ID:        block.ID.String(),
		Timestamp: time.Unix(int64(block.Timestamp), 0),
	}

	return blockMetaCache[height].Timestamp, nil
}

func getContracts() (contracts []mergedContract, err error) {
	siaContracts, err := apiClient.HostContractInfoGet()

	if err != nil {
		return nil, fmt.Errorf("get sia contracts: %w", err)
	}

	currentHeight, err := getSyncedHeight()
	if err != nil {
		return nil, fmt.Errorf("get current height: %w", err)
	}

	if len(siaContracts.Contracts) == 0 {
		return []mergedContract{}, nil
	}

	for _, siaContract := range siaContracts.Contracts {
		confirmed, err := txnConfirmed(siaContract.TransactionID)
		if err != nil {
			return nil, fmt.Errorf("unable to check confirmed txn %s for contract %s: %w", siaContract.TransactionID, siaContract.ObligationId, err)
		}
		if !confirmed {
			continue
		}

		contract := mergedContract{
			ID:                siaContract.ObligationId.String(),
			TransactionID:     siaContract.TransactionID.String(),
			RevisionNumber:    siaContract.RevisionNumber,
			NegotiationHeight: uint64(siaContract.NegotiationHeight),
			ExpirationHeight:  uint64(siaContract.ExpirationHeight),
			ProofDeadline:     uint64(siaContract.ProofDeadLine),
			ProofConfirmed:    siaContract.ProofConfirmed,
			DataSize:          siaContract.DataSize,
			LockedCollateral:  siaContract.LockedCollateral,
			PotentialRevenue:  siaContract.ValidProofOutputs[1].Value.Sub(siaContract.LockedCollateral),
		}

		proofRequired := siaContract.ValidProofOutputs[1].Value.Cmp(siaContract.MissedProofOutputs[1].Value) == 1

		if contract.ProofConfirmed || (!proofRequired && contract.ExpirationHeight < currentHeight) {
			contract.Status = contractStatusSucceeded

			if contract.ProofConfirmed {
				contract.Payout = siaContract.ValidProofOutputs[1].Value
			} else {
				contract.Payout = siaContract.MissedProofOutputs[1].Value
			}

			contract.EarnedRevenue = contract.EarnedRevenue.AddCurrency(contract.Payout).SubCurrency(contract.LockedCollateral)
		} else if !contract.ProofConfirmed && contract.ProofDeadline < currentHeight {
			contract.Status = contractStatusFailed
			contract.Payout = siaContract.MissedProofOutputs[1].Value
			contract.LostRevenue = siaContract.ValidProofOutputs[1].Value.Sub(contract.LockedCollateral)
			contract.EarnedRevenue = contract.EarnedRevenue.AddCurrency(siaContract.MissedProofOutputs[1].Value).SubCurrency(contract.LockedCollateral)

			if siaContract.MissedProofOutputs[1].Value.Cmp(contract.LockedCollateral) == -1 {
				contract.BurntCollateral = contract.LockedCollateral.Sub(siaContract.MissedProofOutputs[1].Value)
			}
		} else {
			contract.Status = contractStatusUnresolved
			contract.Payout = siaContract.ValidProofOutputs[1].Value
			contract.PotentialRevenue = contract.Payout.Sub(contract.LockedCollateral)
		}

		negotiationTimestamp, err := getBlockMeta(contract.NegotiationHeight)
		if err != nil {
			return nil, fmt.Errorf("get negotiation height: %w", err)
		}

		contract.NegotiationTimestamp = negotiationTimestamp

		if currentHeight < contract.ExpirationHeight {
			hoursRemaining := time.Duration((contract.ExpirationHeight-currentHeight)/uint64(siatypes.BlocksPerHour)) * time.Hour
			contract.ExpirationTimestamp = time.Now().Truncate(time.Hour).Add(hoursRemaining)
		} else {
			expirationTimestamp, err := getBlockMeta(contract.ExpirationHeight)
			if err != nil {
				return nil, fmt.Errorf("get expiration height: %w", err)
			}

			contract.ExpirationTimestamp = expirationTimestamp
		}

		if currentHeight < contract.ProofDeadline {
			hoursRemaining := time.Duration((contract.ProofDeadline-currentHeight)/uint64(siatypes.BlocksPerHour)) * time.Hour
			contract.ProofDeadlineTimestamp = time.Now().Truncate(time.Hour).Add(hoursRemaining)
		} else {
			proofDeadlineTimestamp, err := getBlockMeta(contract.ProofDeadline)
			if err != nil {
				return nil, fmt.Errorf("get negotiation height: %w", err)
			}

			contract.ProofDeadlineTimestamp = proofDeadlineTimestamp
		}

		contracts = append(contracts, contract)
	}

	return
}

func syncHostSnapshots(contracts []mergedContract) {
	snapshotMap := make(map[uint64]types.HostSnapshot)

	for _, contract := range contracts {
		endTimestamp := snapshotTime(contract.ProofDeadlineTimestamp)

		for current := snapshotTime(contract.NegotiationTimestamp); current.Before(endTimestamp); current = current.Add(time.Hour) {
			activeID := snapshotID(current)
			snapshot := snapshotMap[activeID]
			snapshot.ActiveContracts++
			snapshot.Timestamp = time.Unix(int64(activeID), 0)
			snapshotMap[activeID] = snapshot
		}

		switch contract.Status {
		case contractStatusSucceeded:
			var successfulID uint64

			if contract.ProofConfirmed {
				successfulID = snapshotID(contract.ExpirationTimestamp)
			} else {
				successfulID = snapshotID(contract.ProofDeadlineTimestamp)
			}

			snapshot := snapshotMap[successfulID]
			snapshot.Timestamp = time.Unix(int64(successfulID), 0)
			snapshot.SuccessfulContracts++

			snapshot.EarnedRevenue = snapshot.EarnedRevenue.Add(contract.EarnedRevenue)
			snapshotMap[successfulID] = snapshot
		case contractStatusFailed:
			failedID := snapshotID(contract.ProofDeadlineTimestamp)
			snapshot := snapshotMap[failedID]

			snapshot.FailedContracts++
			snapshot.EarnedRevenue = snapshot.EarnedRevenue.Add(contract.EarnedRevenue)
			snapshot.BurntCollateral = snapshot.BurntCollateral.Add(contract.BurntCollateral)
			snapshot.Timestamp = time.Unix(int64(failedID), 0)

			snapshotMap[failedID] = snapshot
		case contractStatusUnresolved:
			expirationID := snapshotID(contract.ExpirationTimestamp)
			snapshot := snapshotMap[expirationID]
			snapshot.PotentialRevenue = snapshot.PotentialRevenue.Add(contract.PotentialRevenue)
			snapshot.ExpiredContracts++
			snapshot.ActiveContracts--
			snapshot.Timestamp = time.Unix(int64(expirationID), 0)
			snapshotMap[expirationID] = snapshot
		}

		formationID := snapshotID(contract.NegotiationTimestamp)
		formationSnapshot := snapshotMap[formationID]

		formationSnapshot.NewContracts++
		formationSnapshot.Timestamp = time.Unix(int64(formationID), 0)
		snapshotMap[formationID] = formationSnapshot
	}

	snapshots := []types.HostSnapshot{}

	for _, snapshot := range snapshotMap {
		snapshots = append(snapshots, snapshot)
	}

	if err := persist.SaveHostSnapshots(snapshots...); err != nil {
		log.Printf("sync error: save snapshotMap: %s", err)
	}
}
