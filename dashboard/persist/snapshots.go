package persist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/siacentral/sia-host-dashboard/dashboard/types"
	"gitlab.com/NebulousLabs/bolt"
)

//SaveHostSnapshots SaveHostSnapshot
func SaveHostSnapshots(snapshots ...types.HostSnapshot) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketHostSnapshots)

		for _, snapshot := range snapshots {
			snapshot.Timestamp = snapshot.Timestamp.Truncate(time.Hour).UTC()
			buf, err := json.Marshal(snapshot)

			if err != nil {
				return fmt.Errorf("json encode: %w", err)
			}

			if err := bucket.Put(timeID(snapshot.Timestamp), buf); err != nil {
				return fmt.Errorf("unable to put snapshot: %w", err)
			}
		}

		return nil
	})
}

//GetHostSnapshots returns all snapshots between two timestamps (inclusive)
func GetHostSnapshots(start, end time.Time) (snapshots []types.HostSnapshot, err error) {
	if start.After(end) {
		err = errors.New("start must be before end")
		return
	}

	startID := timeID(start)
	endID := timeID(end)

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketHostSnapshots).Cursor()

		for key, buf := c.Seek(startID); key != nil; key, buf = c.Next() {
			if bytes.Compare(key, endID) > 0 {
				break
			}

			var snapshot types.HostSnapshot

			if err := json.Unmarshal(buf, &snapshot); err != nil {
				return err
			}

			snapshots = append(snapshots, snapshot)
		}

		return nil
	})

	return
}

// GetDailySnapshots returns snapshot totals for every day between two timestamps (inclusive)
func GetDailySnapshots(start, end time.Time) (snapshots []types.HostSnapshot, err error) {
	if start.After(end) {
		err = errors.New("start must be before end")
		return
	}

	snapshots = append(snapshots, types.HostSnapshot{
		Timestamp: start,
	})

	end = end.AddDate(0, 0, -1)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketHostSnapshots)

		for current := start.AddDate(0, 0, -1); current.Before(end); current = current.Add(time.Hour) {
			var snapshot types.HostSnapshot

			id := timeID(current)
			i := len(snapshots) - 1
			buf := b.Get(id)

			if current.Equal(snapshots[i].Timestamp) {
				snapshots = append(snapshots, types.HostSnapshot{
					Timestamp: snapshots[i].Timestamp.AddDate(0, 0, 1),
				})
			}

			if buf == nil {
				continue
			}

			if err := json.Unmarshal(buf, &snapshot); err != nil {
				return fmt.Errorf("unable to decode snaphot: %w", err)
			}

			snapshots[i].ActiveContracts = snapshot.ActiveContracts
			snapshots[i].NewContracts += snapshot.NewContracts
			snapshots[i].ExpiredContracts += snapshot.ExpiredContracts
			snapshots[i].SuccessfulContracts += snapshot.SuccessfulContracts
			snapshots[i].FailedContracts += snapshot.FailedContracts

			snapshots[i].Payout = snapshots[i].Payout.Add(snapshot.Payout)
			snapshots[i].EarnedRevenue = snapshots[i].EarnedRevenue.Add(snapshot.EarnedRevenue)
			snapshots[i].PotentialRevenue = snapshots[i].PotentialRevenue.Add(snapshot.PotentialRevenue)
			snapshots[i].BurntCollateral = snapshots[i].BurntCollateral.Add(snapshot.BurntCollateral)
		}

		return nil
	})

	return
}
