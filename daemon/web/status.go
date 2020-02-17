package web

import (
	"log"
	"net/http"
	"time"

	"github.com/siacentral/host-dashboard/daemon/cache"
	"github.com/siacentral/host-dashboard/daemon/persist"
	"github.com/siacentral/host-dashboard/daemon/types"
	"github.com/siacentral/host-dashboard/daemon/web/router"
)

type (
	hostStatusResponse struct {
		router.APIResponse
		Status types.HostStatus  `json:"status"`
		Alerts []types.HostAlert `json:"alerts"`
	}
)

func handleGetHostStatus(w http.ResponseWriter, r *router.APIRequest) {
	meta, err := persist.GetLastMetadata()
	if err != nil {
		log.Println(err)
		router.HandleError("unable to retrieve metadata", 500, w, r)
	}

	usage, err := persist.GetClosestMeta(time.Now().AddDate(0, 0, -30))
	if err != nil {
		log.Println(err)
		router.HandleError("unable to retrieve past usage", 500, w, r)
	}

	status := cache.GetHostStatus()

	status.UploadBandwidth -= usage.UploadBandwidth
	status.DownloadBandwidth -= usage.DownloadBandwidth

	status.ActiveContracts = meta.ActiveContracts
	status.SuccessfulContracts = meta.SuccessfulContracts
	status.FailedContracts = meta.FailedContracts
	status.Payout = meta.Payout
	status.EarnedRevenue = meta.EarnedRevenue
	status.PotentialRevenue = meta.PotentialRevenue
	status.BurntCollateral = meta.BurntCollateral
	status.FirstSeen = meta.FirstSeen

	router.SendJSONResponse(hostStatusResponse{
		APIResponse: router.APIResponse{
			Type: "success",
		},
		Status: status,
		Alerts: cache.GetAlerts(),
	}, 200, w, r)
}
