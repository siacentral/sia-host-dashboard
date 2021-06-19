package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/siacentral/apisdkgo"
	"github.com/siacentral/host-dashboard/daemon/cache"
	"github.com/siacentral/host-dashboard/daemon/persist"
	"github.com/siacentral/host-dashboard/daemon/types"
)

var siacentralapi = apisdkgo.NewSiaClient()

func calcHostContracts(contracts []mergedContract, meta *types.HostMeta) {
	for _, contract := range contracts {
		meta.Payout = meta.Payout.Add(contract.Payout)
		meta.EarnedRevenue = meta.EarnedRevenue.Add(contract.EarnedRevenue)
		meta.BurntCollateral = meta.BurntCollateral.Add(contract.BurntCollateral)

		switch contract.Status {
		case contractStatusSucceeded:
			meta.SuccessfulContracts++
		case contractStatusFailed:
			meta.FailedContracts++
		case contractStatusUnresolved:
			meta.PotentialRevenue = meta.PotentialRevenue.Add(contract.PotentialRevenue)
			meta.ActiveContracts++
		}
	}
}

func getFirstSeen(pubkey string, meta *types.HostMeta) error {
	public, err := siacentralapi.GetHost(pubkey)

	if err != nil {
		return fmt.Errorf("get host history: %w", err)
	}

	meta.FirstSeen = public.FirstSeenTimestamp

	return nil
}

func syncHostMeta(contracts []mergedContract) {
	var meta types.HostMeta

	cache.ClearAlerts(AlertSyncError)

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
	}
}
