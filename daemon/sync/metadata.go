package sync

import (
	"fmt"
	"log"
	"time"

	siacentralapi "github.com/siacentral/apisdkgo"
	"github.com/siacentral/host-dashboard/daemon/cache"
	"github.com/siacentral/host-dashboard/daemon/persist"
	"github.com/siacentral/host-dashboard/daemon/types"
	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

func calcHostContracts(contracts []mergedContract, meta *types.HostMeta) {
	for _, contract := range contracts {
		switch contract.Status {
		case "obligationSucceeded":
			var payout siatypes.Currency

			if contract.ProofConfirmed {
				payout = contract.ValidProofOutputs[1].Value
			} else {
				payout = contract.MissedProofOutputs[1].Value
			}

			meta.SuccessfulContracts++

			meta.Payout = meta.Payout.Add(payout)
			meta.EarnedRevenue = meta.EarnedRevenue.AddCurrency(payout).
				SubCurrency(contract.LockedCollateral).
				SubCurrency(contract.TransactionFeesAdded)

			break
		case "obligationFailed":
			meta.FailedContracts++

			meta.Payout = meta.Payout.Add(contract.MissedProofOutputs[1].Value)
			meta.BurntCollateral = meta.BurntCollateral.Add(contract.LockedCollateral)
			meta.EarnedRevenue = meta.EarnedRevenue.
				SubCurrency(contract.LockedCollateral)
			break
		case "obligationUnresolved":
			meta.PotentialRevenue = meta.PotentialRevenue.Add(contract.PotentialRevenue)
			meta.ActiveContracts++
			break
		}
	}
}

func getFirstSeen(pubkey string, meta *types.HostMeta) error {
	public, err := siacentralapi.GetHost(pubkey)

	if err != nil {
		return fmt.Errorf("get host history: %s", err)
	}

	meta.FirstSeen = public.FirstSeenTimestamp

	return nil
}

func syncHostMeta(contracts []mergedContract) {
	var meta types.HostMeta

	calcHostContracts(contracts, &meta)

	host, err := apiClient.HostGet()

	if err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host. Check your Sia connection.",
			Type:     "sync",
		})
		return
	}

	if err := getFirstSeen(host.PublicKey.String(), &meta); err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host. Check your Sia connection.",
			Type:     "sync",
		})
		return
	}

	up, down := getBandwidthUsage()

	meta.UploadBandwidth = up
	meta.DownloadBandwidth = down
	meta.UsedStorage = host.ExternalSettings.TotalStorage - host.ExternalSettings.RemainingStorage
	meta.TotalStorage = host.ExternalSettings.TotalStorage
	meta.Settings.BaseRPCPrice = host.ExternalSettings.BaseRPCPrice
	meta.Settings.Collateral = host.ExternalSettings.Collateral
	meta.Settings.MaxCollateral = host.ExternalSettings.MaxCollateral
	meta.Settings.ContractPrice = host.ExternalSettings.ContractPrice
	meta.Settings.DownloadBandwidthPrice = host.ExternalSettings.DownloadBandwidthPrice
	meta.Settings.SectorAccessPrice = host.ExternalSettings.SectorAccessPrice
	meta.Settings.StoragePrice = host.ExternalSettings.StoragePrice
	meta.Settings.UploadBandwidthPrice = host.ExternalSettings.UploadBandwidthPrice
	meta.Timestamp = time.Now().UTC().Truncate(time.Hour)

	if err := persist.SaveHostMeta(meta); err != nil {
		log.Printf("sync error: save meta: %s", err)
		return
	}

	return
}
