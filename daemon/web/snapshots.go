package web

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/siacentral/host-dashboard/daemon/persist"
	"github.com/siacentral/host-dashboard/daemon/types"
	"github.com/siacentral/host-dashboard/daemon/web/router"
)

type (
	hostResponse struct {
		router.APIResponse
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	hostSnapshotResponse struct {
		hostResponse
		Snapshots []types.HostSnapshot `json:"snapshots"`
	}

	hostTotalResponse struct {
		hostResponse
		Day   types.HostSnapshot `json:"day"`
		Month types.HostSnapshot `json:"month"`
		Year  types.HostSnapshot `json:"year"`
		Total types.HostSnapshot `json:"total"`
	}
)

func parseTimeParams(r *router.APIRequest, params ...string) (timestamps []time.Time) {
	for _, param := range params {
		if len(param) == 0 {
			timestamps = append(timestamps, time.Time{})
			continue
		}

		seconds, err := strconv.ParseInt(r.Request.URL.Query().Get(param), 10, 64)

		if err != nil {
			timestamps = append(timestamps, time.Time{})
			continue
		}

		timestamps = append(timestamps, time.Unix(seconds, 0).UTC())
	}

	return
}

func handleGetHostSnapshots(w http.ResponseWriter, r *router.APIRequest) {
	date := parseTimeParams(r, "end")[0]
	if date.IsZero() {
		date = time.Now().UTC()
	}

	date = date.Truncate(time.Hour)
	start := time.Date(date.Year(), date.Month(), 1, date.Hour(), 0, 0, 0, time.UTC).
		AddDate(-1, 0, 0)
	end := start.AddDate(1, 4, -1)

	if end.Before(start) {
		router.HandleError("end must be after start", 400, w, r)
		return
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, time.UTC)

	snapshots, err := persist.GetDailySnapshots(start, end)

	if err != nil {
		router.HandleError("unable to get snapshots", 400, w, r)
		return
	}

	router.SendJSONResponse(hostSnapshotResponse{
		hostResponse: hostResponse{
			APIResponse: router.APIResponse{
				Message: "successfully retrieved snapshots",
				Type:    "success",
			},
			Start: start,
			End:   end,
		},
		Snapshots: snapshots,
	}, 200, w, r)
}

func handleGetHostTotals(w http.ResponseWriter, r *router.APIRequest) {
	var start, end time.Time

	current := time.Now()
	date := parseTimeParams(r, "date")[0]
	if date.IsZero() {
		date = current
	}

	date = date.Truncate(time.Hour).UTC()
	start = time.Date(date.Year(), 1, 1, date.Hour(), 0, 0, 0, time.UTC)
	end = start.AddDate(1, 0, 0)

	dailySnapshots, err := persist.GetDailySnapshots(start, end)
	if err != nil {
		router.HandleError("unable to retrieve snapshots", 400, w, r)
		log.Println(err)
		return
	}

	lastMetadata, err := persist.GetLastMetadata()
	if err != nil {
		router.HandleError("unable to retrieve metadata", 400, w, r)
		log.Println(err)
		return
	}

	resp := hostTotalResponse{
		hostResponse: hostResponse{
			APIResponse: router.APIResponse{
				Message: "successfully retrieved totals",
				Type:    "success",
			},
			Start: start,
			End:   end,
		},
		Day: types.HostSnapshot{
			Timestamp: date,
		},
		Month: types.HostSnapshot{
			Timestamp: date,
		},
		Year: types.HostSnapshot{
			Timestamp: end,
		},
		Total: types.HostSnapshot{
			ActiveContracts:     lastMetadata.ActiveContracts,
			SuccessfulContracts: lastMetadata.SuccessfulContracts,
			FailedContracts:     lastMetadata.FailedContracts,
			Payout:              lastMetadata.Payout,
			EarnedRevenue:       lastMetadata.EarnedRevenue,
			PotentialRevenue:    lastMetadata.PotentialRevenue,
			BurntCollateral:     lastMetadata.BurntCollateral,
			Timestamp:           lastMetadata.Timestamp,
		},
	}

	dy, dm, _ := date.Date()

	for _, snapshot := range dailySnapshots {
		sy, sm, _ := snapshot.Timestamp.Date()

		if snapshot.Timestamp.Equal(date) {
			resp.Day = snapshot
		}

		if snapshot.Timestamp.Before(current) {
			if sy == dy && sm == dm {
				resp.Month.ActiveContracts = snapshot.ActiveContracts
				resp.Month.ExpiredContracts += snapshot.ExpiredContracts
			}

			if sy == dy {
				resp.Year.ActiveContracts = snapshot.ActiveContracts
				resp.Year.ExpiredContracts += snapshot.ExpiredContracts
			}
		}

		if sy == dy && sm == dm {
			resp.Month.NewContracts += snapshot.NewContracts
			resp.Month.SuccessfulContracts += snapshot.SuccessfulContracts
			resp.Month.FailedContracts += snapshot.FailedContracts

			resp.Month.Payout = resp.Month.Payout.Add(snapshot.Payout)
			resp.Month.EarnedRevenue = resp.Month.EarnedRevenue.Add(snapshot.EarnedRevenue)
			resp.Month.PotentialRevenue = resp.Month.PotentialRevenue.Add(snapshot.PotentialRevenue)
			resp.Month.BurntCollateral = resp.Month.BurntCollateral.Add(snapshot.BurntCollateral)
		}

		if sy == dy {
			resp.Year.NewContracts += snapshot.NewContracts
			resp.Year.SuccessfulContracts += snapshot.SuccessfulContracts
			resp.Year.FailedContracts += snapshot.FailedContracts

			resp.Year.Payout = resp.Year.Payout.Add(snapshot.Payout)
			resp.Year.EarnedRevenue = resp.Year.EarnedRevenue.Add(snapshot.EarnedRevenue)
			resp.Year.PotentialRevenue = resp.Year.PotentialRevenue.Add(snapshot.PotentialRevenue)
			resp.Year.BurntCollateral = resp.Year.BurntCollateral.Add(snapshot.BurntCollateral)
		}
	}

	router.SendJSONResponse(resp, 200, w, r)
}
