package sync

import (
	"fmt"
	"log"
	"time"

	siacentralapi "github.com/siacentral/apisdkgo"
	apitypes "github.com/siacentral/apisdkgo/types"
	"github.com/siacentral/host-dashboard/daemon/cache"
	"github.com/siacentral/host-dashboard/daemon/persist"
	"github.com/siacentral/host-dashboard/daemon/types"
	"gitlab.com/NebulousLabs/Sia/modules"
	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

//merges fields from the local contract db and the blockchain
type mergedContract struct {
	apitypes.StorageContract
	LockedCollateral     siatypes.Currency
	TransactionFeesAdded siatypes.Currency
	PotentialRevenue     siatypes.Currency
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

func getContracts() ([]mergedContract, error) {
	siaContracts, err := apiClient.HostContractInfoGet()

	if err != nil {
		return nil, fmt.Errorf("get sia contracts: %s", err)
	}

	if len(siaContracts.Contracts) == 0 {
		return []mergedContract{}, nil
	}

	contractIDs := []string{}
	contractMap := make(map[string]modules.StorageObligation)

	for _, contract := range siaContracts.Contracts {
		contractIDs = append(contractIDs, contract.ObligationId.String())
		contractMap[contract.ObligationId.String()] = contract
	}

	chainContracts, err := siacentralapi.FindContractsByID(contractIDs...)

	if err != nil {
		return nil, fmt.Errorf("get chain contracts: %s", err)
	}

	contracts := make([]mergedContract, 0, len(chainContracts))

	for _, contract := range chainContracts {
		siaContract := contractMap[contract.ID]

		contracts = append(contracts, mergedContract{
			StorageContract:      contract,
			LockedCollateral:     siaContract.LockedCollateral,
			TransactionFeesAdded: siaContract.TransactionFeesAdded,
			PotentialRevenue: siaContract.ContractCost.Add(siaContract.PotentialStorageRevenue).
				Add(siaContract.PotentialUploadRevenue).Add(siaContract.PotentialDownloadRevenue),
		})
	}

	return contracts, nil
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
		case "obligationSucceeded":
			var successfulID uint64
			var payout siatypes.Currency

			if contract.ProofConfirmed {
				successfulID = snapshotID(contract.ProofTimestamp)
				payout = contract.ValidProofOutputs[1].Value
			} else {
				successfulID = snapshotID(contract.ProofDeadlineTimestamp)
				payout = contract.MissedProofOutputs[1].Value
			}

			snapshot := snapshotMap[successfulID]
			snapshot.Timestamp = time.Unix(int64(successfulID), 0)
			snapshot.SuccessfulContracts++

			snapshot.EarnedRevenue = snapshot.EarnedRevenue.AddCurrency(payout).
				SubCurrency(contract.LockedCollateral).
				SubCurrency(contract.TransactionFeesAdded)

			snapshotMap[successfulID] = snapshot
			break
		case "obligationFailed":
			failedID := snapshotID(contract.ProofDeadlineTimestamp)
			snapshot := snapshotMap[failedID]

			snapshot.FailedContracts++
			snapshot.EarnedRevenue = snapshot.EarnedRevenue.
				SubCurrency(contract.LockedCollateral)
			snapshot.BurntCollateral = snapshot.BurntCollateral.
				Add(contract.LockedCollateral)
			snapshot.Timestamp = time.Unix(int64(failedID), 0)

			snapshotMap[failedID] = snapshot
			break
		case "obligationUnresolved":
			expirationID := snapshotID(contract.ExpirationTimestamp)
			snapshot := snapshotMap[expirationID]
			snapshot.PotentialRevenue = snapshot.PotentialRevenue.Add(contract.PotentialRevenue)
			snapshot.ExpiredContracts++
			snapshot.ActiveContracts--
			snapshot.Timestamp = time.Unix(int64(expirationID), 0)
			snapshotMap[expirationID] = snapshot

			break
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
