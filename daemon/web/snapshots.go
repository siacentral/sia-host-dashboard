package web

import (
	"log"
	"net/http"
	"sort"
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

	hostMetaResponse struct {
		hostResponse
		Metadata []types.HostMeta `json:"metadata"`
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

func parseTimeParams(params ...string) ([]time.Time, error) {
	timestamps := []time.Time{}

	for _, param := range params {
		if len(param) == 0 {
			timestamps = append(timestamps, time.Time{})
			continue
		}

		seconds, err := strconv.ParseInt(param, 10, 64)

		if err != nil {
			return []time.Time{}, err
		}

		timestamps = append(timestamps, time.Unix(seconds, 0))
	}

	return timestamps, nil
}

func handleGetHostMetadata(w http.ResponseWriter, r *router.APIRequest) {
	timestamps, err := parseTimeParams(r.Request.URL.Query().Get("start"), r.Request.URL.Query().Get("end"))

	if err != nil {
		router.HandleError("unable to parse start or end time", 400, w, r)
		return
	}

	start := timestamps[0]
	end := timestamps[1]

	if end.IsZero() && start.IsZero() {
		end = time.Now().UTC().Truncate(time.Hour)
		start = end.AddDate(0, 0, -30)
	} else if end.IsZero() {
		end = start.AddDate(0, 0, 30)
	} else if start.IsZero() {
		start = end.AddDate(0, 0, -30)
	}

	if end.Before(start) {
		router.HandleError("end must be after start", 400, w, r)
		return
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), start.Hour(), 0, 0, 0, time.UTC)

	metaHours := make(map[int64]types.HostMeta)

	// fill the time series with all possible values even if no data was stored
	for d := start; d.Before(end); d = d.Add(time.Hour) {
		metaHours[d.Unix()] = types.HostMeta{
			Timestamp: d,
		}
	}

	hostMeta, err := persist.GetHostMetadata(start, end)

	if err != nil {
		router.HandleError("unable to retrieve snapshots", 400, w, r)
		log.Println(err)
		return
	}

	for _, meta := range hostMeta {
		metaHours[meta.Timestamp.Unix()] = meta
	}

	metadata := make([]types.HostMeta, 0, len(metaHours))

	for _, meta := range metaHours {
		metadata = append(metadata, meta)
	}

	sort.Slice(metadata, func(i, j int) bool {
		return metadata[i].Timestamp.Before(metadata[j].Timestamp)
	})

	router.SendJSONResponse(hostMetaResponse{
		hostResponse: hostResponse{
			APIResponse: router.APIResponse{
				Message: "",
				Type:    "success",
			},
			Start: start,
			End:   end,
		},
		Metadata: metadata,
	}, 200, w, r)
}

func handleGetHostSnapshots(w http.ResponseWriter, r *router.APIRequest) {
	timestamps, err := parseTimeParams(r.Request.URL.Query().Get("start"), r.Request.URL.Query().Get("end"))

	if err != nil {
		router.HandleError("unable to parse start or end time", 400, w, r)
		return
	}

	start := timestamps[0]
	end := timestamps[1]

	if end.IsZero() && start.IsZero() {
		current := time.Now().UTC()
		end = time.Date(current.Year(), current.Month(), current.Day(), current.Hour(), 0, 0, 0, time.UTC)
		start = end.AddDate(0, -12, 0)
	} else if end.IsZero() {
		end = start.AddDate(0, 12, 0)
	} else if start.IsZero() {
		start = end.AddDate(0, -12, 0)
	}

	end = time.Date(end.Year(), end.Month()+4, 1, end.Hour(), 0, 0, 0, time.UTC)
	end = end.AddDate(0, 0, -1)
	start = time.Date(start.Year(), start.Month(), 1, start.Hour(), 0, 0, 0, time.UTC)

	if end.Before(start) {
		router.HandleError("end must be after start", 400, w, r)
		return
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, time.UTC)

	snapshotDays := make(map[int64]types.HostSnapshot)

	// fill the time series with all possible values even if no data was stored
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		snapshotDays[d.Unix()] = types.HostSnapshot{
			Timestamp: d,
		}
	}

	hostSnapshots, err := persist.GetHostSnapshots(start, end)

	if err != nil {
		router.HandleError("unable to retrieve snapshots", 400, w, r)
		log.Println(err)
		return
	}

	for _, hourSnap := range hostSnapshots {
		timestamp := hourSnap.Timestamp.UTC()
		timestamp = time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), start.Hour(), 0, 0, 0, time.UTC)

		daySnap := snapshotDays[timestamp.Unix()]

		daySnap.Timestamp = timestamp
		daySnap.ActiveContracts += hourSnap.ActiveContracts
		daySnap.NewContracts += hourSnap.NewContracts
		daySnap.ExpiredContracts += hourSnap.ExpiredContracts
		daySnap.SuccessfulContracts += hourSnap.SuccessfulContracts
		daySnap.FailedContracts += hourSnap.FailedContracts

		daySnap.Payout = daySnap.Payout.Add(hourSnap.Payout)
		daySnap.EarnedRevenue = daySnap.EarnedRevenue.Add(hourSnap.EarnedRevenue)
		daySnap.PotentialRevenue = daySnap.PotentialRevenue.Add(hourSnap.PotentialRevenue)
		daySnap.BurntCollateral = daySnap.BurntCollateral.Add(hourSnap.BurntCollateral)

		snapshotDays[timestamp.Unix()] = daySnap
	}

	snapshots := make([]types.HostSnapshot, 0, len(snapshotDays))

	for _, snapshot := range snapshotDays {
		snapshots = append(snapshots, snapshot)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.Before(snapshots[j].Timestamp)
	})

	router.SendJSONResponse(hostSnapshotResponse{
		hostResponse: hostResponse{
			APIResponse: router.APIResponse{
				Message: "",
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
	timestamps, err := parseTimeParams(r.Request.URL.Query().Get("date"))

	if err != nil {
		router.HandleError("unable to parse unix timestamp", 400, w, r)
		return
	}

	date := timestamps[0]

	if date.IsZero() {
		date = time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	start = time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	end = start.AddDate(1, 0, -1)

	dailySnapshots, err := persist.GetHostSnapshots(start, end)
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
				Message: "",
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
			Timestamp: start,
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

	dy, dm, dd := date.Date()

	for _, snapshot := range dailySnapshots {
		sy, sm, sd := snapshot.Timestamp.Date()

		if sy == dy && sm == dm && sd == dd {
			resp.Day.ActiveContracts = snapshot.ActiveContracts
			resp.Day.ExpiredContracts = snapshot.ExpiredContracts

			resp.Day.NewContracts += snapshot.NewContracts
			resp.Day.SuccessfulContracts += snapshot.SuccessfulContracts
			resp.Day.FailedContracts += snapshot.FailedContracts

			resp.Day.Payout = resp.Day.Payout.Add(snapshot.Payout)
			resp.Day.EarnedRevenue = resp.Day.EarnedRevenue.Add(snapshot.EarnedRevenue)
			resp.Day.PotentialRevenue = resp.Day.PotentialRevenue.Add(snapshot.PotentialRevenue)
			resp.Day.BurntCollateral = resp.Day.BurntCollateral.Add(snapshot.BurntCollateral)
		}

		if sy == dy && sm == dm {
			resp.Month.ActiveContracts = snapshot.ActiveContracts
			resp.Month.ExpiredContracts = snapshot.ExpiredContracts

			resp.Month.NewContracts += snapshot.NewContracts
			resp.Month.SuccessfulContracts += snapshot.SuccessfulContracts
			resp.Month.FailedContracts += snapshot.FailedContracts

			resp.Month.Payout = resp.Month.Payout.Add(snapshot.Payout)
			resp.Month.EarnedRevenue = resp.Month.EarnedRevenue.Add(snapshot.EarnedRevenue)
			resp.Month.PotentialRevenue = resp.Month.PotentialRevenue.Add(snapshot.PotentialRevenue)
			resp.Month.BurntCollateral = resp.Month.BurntCollateral.Add(snapshot.BurntCollateral)
		}

		if sy == dy {
			resp.Year.ActiveContracts += snapshot.ActiveContracts
			resp.Year.ExpiredContracts += snapshot.ExpiredContracts

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
