package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/siacentral/host-dashboard/daemon/persist"
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

// defaultSiaDir returns the default data directory of siad. The values for
// supported operating systems are:
//
// Linux:   $HOME/.sia
// MacOS:   $HOME/Library/Application Support/Sia
// Windows: %LOCALAPPDATA%\Sia
func defaultSiaDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "Sia")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Sia")
	default:
		return filepath.Join(os.Getenv("HOME"), ".sia")
	}
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
	bw, err := apiClient.HostBandwidthGet()

	if err != nil {
		log.Printf("warn: unable to retrieve bandwidth: %s", err)
		return
	}

	bandwidthMu.Lock()
	upload = counters.totalUpload + (bw.Upload - counters.lastUpload)
	download = counters.totalDownload + (bw.Download - counters.lastDownload)

	counters.totalUpload = upload
	counters.totalDownload = download
	counters.lastUpload = bw.Upload
	counters.lastDownload = bw.Download
	bandwidthMu.Unlock()

	return
}

// initializes the bandwidth counters from the database to count the total bandwidth usage of the
// host. This will cause us to lose a few bytes on initialization, but should prevent counting
// bandwidth twice
func initBandwidthCounters() {
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

	bandwidthMu.Lock()
	counters = new(bandwidthCounters)
	counters.totalUpload = meta.UploadBandwidth
	counters.totalDownload = meta.DownloadBandwidth
	counters.lastUpload = bw.Upload
	counters.lastDownload = bw.Download
	bandwidthMu.Unlock()
}

//Start begins syncing data from Sia
func Start(siaAddr string) error {
	apiClient = &siaapi.Client{
		Address:   siaAddr,
		UserAgent: "Sia-Agent",
	}

	initBandwidthCounters()

	if err := syncContracts(); err != nil {
		return fmt.Errorf("refreshing contracts: %s", err)
	}

	if err := syncHostConnectivity(); err != nil {
		return fmt.Errorf("refreshing connectivity: %s", err)
	}

	if err := syncHostStatus(); err != nil {
		return fmt.Errorf("refreshing status: %s", err)
	}

	go refreshContracts()
	go refreshStatus()
	go refreshConnectivity()

	return nil
}
