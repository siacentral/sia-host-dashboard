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
		UsedStorage         uint64            `json:"used_storage"`
		TotalStorage        uint64            `json:"total_storage"`
		UploadBandwidth     uint64            `json:"upload_bandwidth"`
		DownloadBandwidth   uint64            `json:"download_bandwidth"`
		Payout              siatypes.Currency `json:"payout"`
		EarnedRevenue       BigNumber         `json:"earned_revenue"`
		PotentialRevenue    siatypes.Currency `json:"potential_revenue"`
		BurntCollateral     siatypes.Currency `json:"burnt_collateral"`
		Settings            HostSettings      `json:"host_settings"`
		FirstSeen           time.Time         `json:"first_seen"`
		Timestamp           time.Time         `json:"timestamp"`
	}
)
