package persist

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/siacentral/sia-host-dashboard/daemon/types"
	"gitlab.com/NebulousLabs/bolt"
)

func timeID(timestamp time.Time) []byte {
	buf := make([]byte, 8)
	seconds := uint64(timestamp.Truncate(time.Hour).Unix())

	binary.BigEndian.PutUint64(buf, seconds)

	return buf
}

//SaveHostMeta SaveHostMeta
func SaveHostMeta(meta types.HostMeta) error {
	meta.Timestamp = meta.Timestamp.Truncate(time.Hour)

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketHostMeta)

		buf, err := json.Marshal(meta)

		if err != nil {
			return fmt.Errorf("json encode: %w", err)
		}

		if err := bucket.Put(timeID(meta.Timestamp), buf); err != nil {
			return fmt.Errorf("unable to put meta: %w", err)
		}

		return nil
	})
}

//GetHostMetadata returns all metadata snapshots between two timestamps (inclusive)
func GetHostMetadata(startTime, endTime time.Time) (metadata []types.HostMeta, err error) {
	startID := timeID(startTime)
	endID := timeID(endTime)

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketHostMeta).Cursor()

		for key, buf := c.First(); key != nil; key, buf = c.Next() {
			if bytes.Compare(key, startID) < 0 {
				continue
			}

			if bytes.Compare(key, endID) > 0 {
				break
			}

			var meta types.HostMeta

			if err := json.Unmarshal(buf, &meta); err != nil {
				return err
			}

			metadata = append(metadata, meta)
		}

		return nil
	})

	return
}

//GetLastMetadata returns the last metadata snapshot stored in the database
func GetLastMetadata() (metadata types.HostMeta, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		_, buf := tx.Bucket(bucketHostMeta).Cursor().Last()

		if buf == nil {
			return nil
		}

		if err := json.Unmarshal(buf, &metadata); err != nil {
			return err
		}

		return nil
	})

	return
}

//GetClosestMeta returns the metadata snapshot closest to the specified time
func GetClosestMeta(timestamp time.Time) (metadata types.HostMeta, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketHostMeta).Cursor()
		_, buf := c.Seek(timeID(timestamp))

		if buf == nil {
			return nil
		}

		if err := json.Unmarshal(buf, &metadata); err != nil {
			return err
		}

		return nil
	})

	return
}
