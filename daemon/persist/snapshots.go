package persist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/siacentral/host-dashboard/daemon/types"
	"gitlab.com/NebulousLabs/bolt"
)

//SaveHostSnapshot SaveHostSnapshot
func SaveHostSnapshot(snapshot types.HostSnapshot) error {
	snapshot.Timestamp = snapshot.Timestamp.Truncate(time.Hour).UTC()

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketHostSnapshots)

		buf, err := json.Marshal(snapshot)

		if err != nil {
			return fmt.Errorf("json encode: %s", err)
		}

		bucket.Put(timeID(snapshot.Timestamp), buf)

		return nil
	})
}

//GetHostSnapshots returns all snapshots between two timestamps (inclusive)
func GetHostSnapshots(startTime, endTime time.Time) (snapshots []types.HostSnapshot, err error) {
	startID := timeID(startTime)
	endID := timeID(endTime)

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketHostSnapshots).Cursor()

		for key, buf := c.First(); key != nil; key, buf = c.Next() {
			if bytes.Compare(key, startID) < 0 {
				continue
			}

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
