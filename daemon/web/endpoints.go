package web

import "github.com/siacentral/host-dashboard/daemon/web/router"

var endpoints = []router.APIEndpoint{
	router.APIEndpoint{
		Name:    "Get Snapshots",
		Method:  "GET",
		Pattern: "/snapshots",
		Secure:  false,
		Handler: handleGetHostSnapshots,
	},
	router.APIEndpoint{
		Name:    "Get Totals",
		Method:  "GET",
		Pattern: "/totals",
		Secure:  false,
		Handler: handleGetHostTotals,
	},
	router.APIEndpoint{
		Name:    "Get Status",
		Method:  "GET",
		Pattern: "/status",
		Secure:  false,
		Handler: handleGetHostStatus,
	},
}
