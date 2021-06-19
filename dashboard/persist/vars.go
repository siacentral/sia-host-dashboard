package persist

import "gitlab.com/NebulousLabs/bolt"

var (
	db *bolt.DB

	bucketHostMeta      = []byte("hostmeta")
	bucketHostSnapshots = []byte("hostsnapshots")
)
