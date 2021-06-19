package cache

import (
	"sync"

	"github.com/siacentral/sia-host-dashboard/daemon/types"
)

var (
	hostStatus types.HostStatus
	alerts     = make(map[types.HostAlertID][]types.HostAlert)
	mu         sync.RWMutex
)

// GetHostStatus returns the last connectivity report
func GetHostStatus() (s types.HostStatus) {
	mu.RLock()
	s = hostStatus
	mu.RUnlock()

	return
}

// SetHostStatus updates the node's last connectivity report
func SetHostStatus(s types.HostStatus) {
	mu.Lock()
	hostStatus = s
	mu.Unlock()
}

//GetAlerts returns all active alerts
func GetAlerts() (active []types.HostAlert) {
	mu.Lock()
	defer mu.Unlock()

	for _, alert := range alerts {
		active = append(active, alert...)
	}

	return
}

// AddAlert adds a new alert to the dashboard
func AddAlert(id types.HostAlertID, alert types.HostAlert) {
	mu.Lock()
	alerts[id] = append(alerts[id], alert)
	mu.Unlock()
}

// ClearAlerts removes the specified ids from the dashboard, if no ids are
// specified removes all alerts
func ClearAlerts(ids ...types.HostAlertID) {
	mu.Lock()
	defer mu.Unlock()

	if len(ids) == 0 {
		for id := range alerts {
			delete(alerts, id)
		}
		return
	}

	for _, id := range ids {
		delete(alerts, id)
	}
}
