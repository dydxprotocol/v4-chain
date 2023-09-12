package types

import (
	"fmt"
	"sync"
	"time"
)

type updateMetadata struct {
	timer           *time.Timer
	updateFrequency time.Duration
}

// UpdateMonitor monitors the update frequency of daemon services. If a daemon service does not respond within
// the maximum acceptable update delay set when the daemon is registered, the monitor will panic and halt the
// protocol. This was judged to be the best solution for network performance because it prevents any validator from
// interacting with the network at all if a daemon service is not responding.
type UpdateMonitor struct {
	// serviceToUpdateMetadata maps daemon service names to their update metadata.
	serviceToUpdateMetadata map[string]updateMetadata
	stopped                 bool
	sync.Mutex
}

// NewUpdateFrequencyMonitor creates a new update frequency monitor.
func NewUpdateFrequencyMonitor() *UpdateMonitor {
	return &UpdateMonitor{
		serviceToUpdateMetadata: make(map[string]updateMetadata),
	}
}

// RegisterDaemonService registers a new daemon service with the update frequency monitor.
// This method is synchronized.
func (ufm *UpdateMonitor) RegisterDaemonService(service string, maximumAcceptableUpdateDelay time.Duration) {
	ufm.Lock()
	defer ufm.Unlock()

	// Don't register additional daemon services if the monitor has already been stopped.
	// This could be a concern for short-running integration test cases, where the network
	// stops before all daemon services have been registered.
	if ufm.stopped {
		return
	}

	ufm.serviceToUpdateMetadata[service] = updateMetadata{
		timer: time.AfterFunc(maximumAcceptableUpdateDelay, func() {
			panic(fmt.Sprintf("%v daemon not responding", service))
		}),
		updateFrequency: maximumAcceptableUpdateDelay,
	}
}

// Stop stops the update frequency monitor. This method is synchronized.
func (ufm *UpdateMonitor) Stop() {
	ufm.Lock()
	defer ufm.Unlock()

	// Don't stop the monitor if it has already been stopped.
	if ufm.stopped {
		return
	}

	for _, metadata := range ufm.serviceToUpdateMetadata {
		metadata.timer.Stop()
	}
	ufm.stopped = true
}

// RegisterValidResponse registers a valid response from the daemon service. This will reset the timer for the
// daemon service. This method is synchronized.
func (ufm *UpdateMonitor) RegisterValidResponse(service string) error {
	ufm.Lock()
	defer ufm.Unlock()

	// Don't bother to reset the timer if the monitor has already been stopped.
	if ufm.stopped {
		return nil
	}

	metadata, ok := ufm.serviceToUpdateMetadata[service]
	if !ok {
		return fmt.Errorf("service not registered")
	}

	metadata.timer.Reset(metadata.updateFrequency)
	return nil
}
