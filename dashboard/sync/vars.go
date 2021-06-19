package sync

import "github.com/siacentral/sia-host-dashboard/dashboard/types"

var (
	// AlertSyncError alerts the user there was an issue syncing the status
	AlertSyncError = types.HostAlertID("hostAlertSyncError")

	// AlertFolderReadWriteError alerts the user to read and/or write failures from the host
	AlertFolderReadWriteError = types.HostAlertID("hostAlertFolderError")

	// AlertStorageUtilization alerts the user to storage utilization warnings from the host
	AlertStorageUtilization = types.HostAlertID("hostAlertStorageUtilization")

	// AlertWalletLocked alerts the user that their wallet is currently locked
	AlertWalletLocked = types.HostAlertID("hostAlertWalletLocked")

	// AlertWalletBalance alerts the user that their wallet is currently locked
	AlertWalletBalance = types.HostAlertID("hostAlertWalletBalance")

	// AlertCollateralBudget alerts the user that the collateral budget is almost fully utilized
	AlertCollateralBudget = types.HostAlertID("hostAlertCollateralBudgetUtilization")

	// AlertConnectionStatus alerts the user that a connection issue occurred
	AlertConnectionStatus = types.HostAlertID("hostAlertConnectionStatus")
)
