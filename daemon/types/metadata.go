package types

import (
	"time"

	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

type (
	// HostSettings the host's current pricing and settings
	HostSettings struct {
		BaseRPCPrice           siatypes.Currency `json:"base_rpc_price"`
		SectorAccessPrice      siatypes.Currency `json:"sector_access_price"`
		Collateral             siatypes.Currency `json:"collateral"`
		MaxCollateral          siatypes.Currency `json:"max_collateral"`
		ContractPrice          siatypes.Currency `json:"contract_price"`
		DownloadBandwidthPrice siatypes.Currency `json:"download_price"`
		StoragePrice           siatypes.Currency `json:"storage_price"`
		UploadBandwidthPrice   siatypes.Currency `json:"upload_price"`
	}

	//HostMeta a snapshot of a host at a point in time
	HostMeta struct {
		ActiveContracts     uint64            `json:"active_contracts"`
		SuccessfulContracts uint64            `json:"successful_contracts"`
		FailedContracts     uint64            `json:"failed_contracts"`
		Payout              siatypes.Currency `json:"payout"`
		EarnedRevenue       BigNumber         `json:"earned_revenue"`
		PotentialRevenue    siatypes.Currency `json:"potential_revenue"`
		BurntCollateral     siatypes.Currency `json:"burnt_collateral"`
		UsedStorage         siatypes.Currency `json:"used_storage"`
		TotalStorage        siatypes.Currency `json:"total_storage"`
		UploadBandwidth     siatypes.Currency `json:"upload_bandwidth"`
		DownloadBandwidth   siatypes.Currency `json:"download_bandwidth"`
		Settings            HostSettings      `json:"host_settings"`
		FirstSeen           time.Time         `json:"first_seen"`
		Timestamp           time.Time         `json:"timestamp"`
	}
)
