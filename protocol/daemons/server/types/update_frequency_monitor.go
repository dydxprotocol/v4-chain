package types

import (
	"fmt"
	"time"
)

type updateMetadata struct {
	timer           *time.Timer
	updateFrequency time.Duration
}

type UpdateFrequencyMonitor struct {
	serviceToUpdateMetadata map[string]updateMetadata
}

// NewUpdateFrequencyMonitor creates a new update frequency monitor.
func NewUpdateFrequencyMonitor() *UpdateFrequencyMonitor {
	return &UpdateFrequencyMonitor{
		serviceToUpdateMetadata: make(map[string]updateMetadata),
	}
}

// RegisterDaemonService registers a new daemon service with the update frequency monitor.
func (ufm *UpdateFrequencyMonitor) RegisterDaemonService(service string, maximumAcceptableUpdateDelay time.Duration) {
	ufm.serviceToUpdateMetadata[service] = updateMetadata{
		timer: time.AfterFunc(maximumAcceptableUpdateDelay, func() {
			panic(fmt.Sprintf("%v daemon not responding", service))
		}),
		updateFrequency: maximumAcceptableUpdateDelay,
	}
}

// Stop stops the update frequency monitor.
func (ufm *UpdateFrequencyMonitor) Stop() {
	for _, metadata := range ufm.serviceToUpdateMetadata {
		metadata.timer.Stop()
	}
}

// RegisterValidResponse registers a valid response from the daemon service. This will reset the timer for the
// daemon service.
func (ufm *UpdateFrequencyMonitor) RegisterValidResponse(service string) error {
	metadata, ok := ufm.serviceToUpdateMetadata[service]
	if !ok {
		return fmt.Errorf("service not registered")
	}

	metadata.timer.Reset(metadata.updateFrequency)
	return nil
}
