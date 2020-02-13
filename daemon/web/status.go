package web

import (
	"log"
	"net/http"

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

	status := cache.GetHostStatus()
	settings := status.Settings
	status.HostMeta = meta
	status.Settings = settings

	router.SendJSONResponse(hostStatusResponse{
		APIResponse: router.APIResponse{
			Type: "success",
		},
		Status: status,
		Alerts: cache.GetAlerts(),
	}, 200, w, r)
}
