package types

import (
	"time"

	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

type (
	//HostSnapshot a snapshot formed from the host's contract set
	HostSnapshot struct {
		ActiveContracts     uint64            `json:"active_contracts"`
		NewContracts        uint64            `json:"new_contracts"`
		ExpiredContracts    uint64            `json:"expired_contracts"`
		SuccessfulContracts uint64            `json:"successful_contracts"`
		FailedContracts     uint64            `json:"failed_contracts"`
		Payout              siatypes.Currency `json:"payout"`
		EarnedRevenue       BigNumber         `json:"earned_revenue"`
		PotentialRevenue    siatypes.Currency `json:"potential_revenue"`
		BurntCollateral     siatypes.Currency `json:"burnt_collateral"`
		Timestamp           time.Time         `json:"timestamp"`
	}
)
