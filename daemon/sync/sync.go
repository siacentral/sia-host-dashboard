package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	siaapi "gitlab.com/NebulousLabs/Sia/node/api/client"
)

var (
	apiClient *siaapi.Client
)

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

//Start begins syncing data from Sia
func Start(siaAddr string) error {
	apiClient = &siaapi.Client{
		Address:   siaAddr,
		UserAgent: "Sia-Agent",
	}

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
