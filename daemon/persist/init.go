package persist

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/NebulousLabs/bolt"
)

//InitializeDB opens or creates the database at the specified path
func InitializeDB(dataPath string) error {
	var err error

	if _, err := os.Stat(dataPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat datapath: %w", err)
		}

		if err := os.MkdirAll(dataPath, 0770); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
	}

	db, err = bolt.Open(filepath.Join(dataPath, "hoststats.db"), 0600, &bolt.Options{Timeout: 5 * time.Second})

	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketHostMeta)

		if err != nil {
			return fmt.Errorf("create metadata bucket: %w", err)
		}

		_, err = tx.CreateBucketIfNotExists(bucketHostSnapshots)

		if err != nil {
			return fmt.Errorf("create snapshots bucket: %w", err)
		}

		return nil
	})
}

//CloseDB closes the database
func CloseDB() error {
	return db.Close()
}
