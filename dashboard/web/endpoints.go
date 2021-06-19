package web

import "github.com/siacentral/sia-host-dashboard/dashboard/web/router"

var endpoints = []router.APIEndpoint{
	{
		Name:    "Get Snapshots",
		Method:  "GET",
		Pattern: "/snapshots",
		Secure:  false,
		Handler: handleGetHostSnapshots,
	},
	{
		Name:    "Get Totals",
		Method:  "GET",
		Pattern: "/totals",
		Secure:  false,
		Handler: handleGetHostTotals,
	},
	{
		Name:    "Get Status",
		Method:  "GET",
		Pattern: "/status",
		Secure:  false,
		Handler: handleGetHostStatus,
	},
}
