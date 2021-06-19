package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/siacentral/sia-host-dashboard/dashboard/cache"
	"github.com/siacentral/sia-host-dashboard/dashboard/types"
	"gitlab.com/NebulousLabs/Sia/node/api"
	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

func calcPercentage(a, b siatypes.Currency) (v uint8) {
	if a.Cmp(b) != -1 {
		return 100
	}

	if b.Cmp64(0) != 1 {
		b = siatypes.NewCurrency64(1)
	}

	p := a.Mul64(100).Div(b)
	n, _ := p.Uint64()

	return uint8(n)
}

func calcPercentage64(a, b uint64) uint8 {
	if a == b || a > b {
		return 100
	}

	if b == 0 {
		b = 1
	}

	return uint8((a * 100) / b)
}

func syncStorageStatus(status *types.HostStatus) error {
	var totalStorage, usedStorage uint64

	folders, err := apiClient.HostStorageGet()

	if err != nil {
		return fmt.Errorf("get storage folders: %w", err)
	}

	cache.ClearAlerts(AlertFolderReadWriteError, AlertStorageUtilization)

	for _, folder := range folders.Folders {
		totalStorage += folder.Capacity
		usedStorage += folder.Capacity - folder.CapacityRemaining

		if folder.FailedReads != 0 && folder.FailedWrites != 0 {
			cache.AddAlert(AlertFolderReadWriteError, types.HostAlert{
				Severity: "severe",
				Text:     fmt.Sprintf("Folder %s has read and write errors", folder.Path),
				Type:     "storage",
			})
		} else if folder.FailedWrites != 0 {
			cache.AddAlert(AlertFolderReadWriteError, types.HostAlert{
				Severity: "severe",
				Text:     fmt.Sprintf("Folder %s has write errors", folder.Path),
				Type:     "storage",
			})
		} else if folder.FailedReads != 0 {
			cache.AddAlert(AlertFolderReadWriteError, types.HostAlert{
				Severity: "severe",
				Text:     fmt.Sprintf("Folder %s has read errors", folder.Path),
				Type:     "storage",
			})
		}
	}

	status.UsedStorage = usedStorage
	status.TotalStorage = totalStorage

	usagePerc := calcPercentage64(status.UsedStorage, status.TotalStorage)

	if status.TotalStorage <= 0 {
		cache.AddAlert(AlertStorageUtilization, types.HostAlert{
			Severity: "severe",
			Text:     "No storage added, add storage folders to accept contracts.",
			Type:     "storage",
		})
	} else if usagePerc >= 98 {
		cache.AddAlert(AlertStorageUtilization, types.HostAlert{
			Severity: "severe",
			Text:     "Storage almost full, add more storage to avoid low storage penalty.",
			Type:     "storage",
		})
	} else if usagePerc >= 85 {
		cache.AddAlert(AlertStorageUtilization, types.HostAlert{
			Severity: "warning",
			Text:     fmt.Sprintf("Storage %d%% utilized, add more storage to avoid low storage penalty.", usagePerc),
			Type:     "storage",
		})
	}

	return nil
}

func syncHostConnectivity() error {
	cache.ClearAlerts(AlertConnectionStatus, AlertSyncError)

	host, err := apiClient.HostGet()

	if err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to check host connectivity",
			Type:     "sync",
		})
		return fmt.Errorf("sia api get failed: %w", err)
	}

	netaddress := string(host.ExternalSettings.NetAddress)

	log.Printf("host netaddress %s", netaddress)

	if len(netaddress) == 0 {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to check host connectivity",
			Type:     "sync",
		})
		return fmt.Errorf("unable to netaddress")
	}

	report, err := siacentralapi.GetHostConnectivity(netaddress)

	if err != nil {
		cache.AddAlert(AlertConnectionStatus, types.HostAlert{
			Severity: "severe",
			Text:     fmt.Sprintf("Failed to check connectivity: %s", err.Error()),
			Type:     "connection",
		})
		return fmt.Errorf("failed to check connection: %w", err)
	}

	for _, err := range report.Errors {
		cache.AddAlert(AlertConnectionStatus, types.HostAlert{
			Severity: err.Severity,
			Text:     err.Message,
			Type:     err.Type,
		})
	}

	return nil
}

func buildStatus(host api.HostGET, wallet api.WalletGET, start time.Time, up, down uint64) types.HostStatus {
	return types.HostStatus{
		HostMeta: types.HostMeta{
			UploadBandwidth:   up,
			DownloadBandwidth: down,
			Settings: types.HostSettings{
				BaseRPCPrice:           host.ExternalSettings.BaseRPCPrice,
				SectorAccessPrice:      host.ExternalSettings.SectorAccessPrice,
				Collateral:             host.ExternalSettings.Collateral,
				MaxCollateral:          host.ExternalSettings.MaxCollateral,
				ContractPrice:          host.ExternalSettings.ContractPrice,
				DownloadBandwidthPrice: host.ExternalSettings.DownloadBandwidthPrice,
				StoragePrice:           host.ExternalSettings.StoragePrice,
				UploadBandwidthPrice:   host.ExternalSettings.UploadBandwidthPrice,
			},
		},
		AcceptingContracts: host.InternalSettings.AcceptingContracts,
		Version:            host.ExternalSettings.Version,
		WalletUnlocked:     wallet.Unlocked,
		StartTime:          start,
	}
}

func syncHostStatus() error {
	host, err := apiClient.HostGet()

	cache.ClearAlerts(AlertSyncError)

	if err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host. Check your Sia connection.",
			Type:     "sync",
		})
		return fmt.Errorf("get host: %w", err)
	}

	wallet, err := apiClient.WalletGet()

	if err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host. Check your Sia connection.",
			Type:     "sync",
		})
		return fmt.Errorf("get wallet: %w", err)
	}

	gbw, err := apiClient.GatewayBandwidthGet()

	if err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host. Check your Sia connection.",
			Type:     "sync",
		})
		return fmt.Errorf("get wallet: %w", err)
	}

	up, down := getBandwidthUsage()
	status := buildStatus(host, wallet, gbw.StartTime, up, down)

	if err := syncStorageStatus(&status); err != nil {
		cache.AddAlert(AlertSyncError, types.HostAlert{
			Severity: "severe",
			Text:     "Unable to sync host. Check your Sia connection.",
			Type:     "sync",
		})
		return fmt.Errorf("get storage folders: %w", err)
	}

	cache.ClearAlerts(AlertWalletLocked, AlertWalletBalance, AlertCollateralBudget)

	if !status.WalletUnlocked {
		cache.AddAlert(AlertWalletLocked, types.HostAlert{
			Severity: "severe",
			Text:     "Wallet is locked. Wallet must be unlocked to form new contracts.",
			Type:     "wallet",
		})
	}

	if wallet.ConfirmedSiacoinBalance.Cmp64(0) != 1 {
		cache.AddAlert(AlertWalletBalance, types.HostAlert{
			Severity: "severe",
			Text:     "Wallet has no more Siacoins to use for collateral. Send more Siacoins.",
			Type:     "wallet",
		})
	}

	usedCollateral := calcPercentage(host.FinancialMetrics.LockedStorageCollateral,
		host.InternalSettings.CollateralBudget)

	if usedCollateral >= 98 {
		cache.AddAlert(AlertCollateralBudget, types.HostAlert{
			Severity: "severe",
			Text:     "Collateral budget fully utilized. Restart host or increase collateral budget.",
			Type:     "collateral",
		})
	} else if usedCollateral >= 85 {
		cache.AddAlert(AlertCollateralBudget, types.HostAlert{
			Severity: "warning",
			Text:     fmt.Sprintf("Collateral budget %d%% utilized. Restart host or increase collateral budget.", usedCollateral),
			Type:     "collateral",
		})
	}

	cache.SetHostStatus(status)

	return nil
}
