package sync

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/siacentral/sia-host-dashboard/daemon/persist"
	siaapi "gitlab.com/NebulousLabs/Sia/node/api/client"
)

var (
	apiClient   *siaapi.Client
	counters    *bandwidthCounters
	bandwidthMu sync.Mutex
)

type bandwidthCounters struct {
	lastUpload    uint64
	lastDownload  uint64
	totalUpload   uint64
	totalDownload uint64
}

func waitInterval(d time.Duration) {
	current := time.Now()
	sleepTime := current.Add(d).Truncate(d).Sub(current)

	time.Sleep(sleepTime)
}

func refreshContracts() {
	for {
		if err := syncContracts(); err != nil {
			log.Println("refreshing contracts:", err)
			time.Sleep(time.Second * 30)
			continue
		}

		waitInterval(time.Minute * 10)
	}
}

func refreshConnectivity() {
	for {
		if err := syncHostConnectivity(); err != nil {
			log.Println("refreshing connectivity:", err)
			time.Sleep(time.Second * 30)
			continue
		}

		waitInterval(time.Minute * 10)
	}
}

func refreshStatus() {
	for {
		if err := syncHostStatus(); err != nil {
			log.Println("refreshing status:", err)
			time.Sleep(time.Second * 30)
			continue
		}

		waitInterval(time.Second * 10)
	}
}

func getBandwidthUsage() (upload, download uint64) {
	bandwidthMu.Lock()
	defer bandwidthMu.Unlock()

	bw, err := apiClient.HostBandwidthGet()

	if err != nil {
		log.Printf("warn: unable to retrieve bandwidth: %s", err)
		return
	}

	dUp := bw.Upload
	dDown := bw.Download

	if dUp >= counters.lastUpload {
		dUp -= counters.lastUpload
	}

	if dDown >= counters.lastDownload {
		dDown -= counters.lastDownload
	}

	upload = counters.totalUpload + dUp
	download = counters.totalDownload + dDown

	counters.totalUpload = upload
	counters.totalDownload = download
	counters.lastUpload = bw.Upload
	counters.lastDownload = bw.Download

	return
}

// initializes the bandwidth counters from the database to count the total bandwidth usage of the
// host. This will cause us to lose any existing bytes on initialization, but should prevent counting
// bandwidth twice
func initBandwidthCounters() {
	bandwidthMu.Lock()
	defer bandwidthMu.Unlock()

	meta, err := persist.GetLastMetadata()
	if err != nil {
		log.Printf("warn: unable to load bandwidth: %s", err)
		return
	}

	bw, err := apiClient.HostBandwidthGet()

	if err != nil {
		log.Printf("warn: unable to retrieve bandwidth: %s", err)
		return
	}

	counters = new(bandwidthCounters)
	counters.totalUpload = meta.UploadBandwidth
	counters.totalDownload = meta.DownloadBandwidth
	counters.lastUpload = bw.Upload
	counters.lastDownload = bw.Download
}

//Start begins syncing data from Sia
func Start(siaAddr string) error {
	apiClient = &siaapi.Client{
		Options: siaapi.Options{
			Address:   siaAddr,
			UserAgent: "Sia-Agent",
		},
	}

	initBandwidthCounters()

	if err := syncContracts(); err != nil {
		return fmt.Errorf("refreshing contracts: %w", err)
	}

	if err := syncHostConnectivity(); err != nil {
		log.Println(fmt.Errorf("refreshing connectivity: %w", err))
	}

	if err := syncHostStatus(); err != nil {
		return fmt.Errorf("refreshing status: %w", err)
	}

	go refreshContracts()
	go refreshStatus()
	go refreshConnectivity()

	return nil
}
