package types

import (
	"fmt"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"sync"
	"time"
)

const (
	MaximumAcceptableUpdateDelay = 5 * time.Minute
)

// HealthCheckable is a common interface for services that can be health checked.
type HealthCheckable interface {
	// HealthCheck returns an error if a service is unhealthy. If the service is healthy, this method returns nil.
	HealthCheck(provider libtime.TimeProvider) (err error)
}

// HealthCheckableImpl implements the HealthCheckable interface by tracking the timestamps of the last successful and
// failed updates. If:
// - no update has occurred
// - the most recent update failed, or
// - the daemon has not seen a successful update in at least 5 minutes,
// the service is considered unhealthy.
type HealthCheckableImpl struct {
	sync.Mutex

	// lastSuccessfulUpdate is the timestamp of the last successful update.
	lastSuccessfulUpdate time.Time
	// lastFailedUpdate is the timestamp of the last failed update. After the HealthCheckableImpl is initialized,
	// this should never be a zero value.
	lastFailedUpdate time.Time
	// lastUpdateError is the error describing the failure reason for the last failed update.
	lastUpdateError error

	// initialized is true if the HealthCheckableImpl has been initialized. If false, the service will panic
	// when the HealthCheck method is called. This error should occur during or soon after app initialization.
	initialized bool
}

// NewHealthCheckableImpl creates a new HealthCheckableImpl instance.
func NewHealthCheckableImpl(daemon string, timeProvider libtime.TimeProvider) *HealthCheckableImpl {
	hc := &HealthCheckableImpl{}
	hc.InitializeHealthStatus(daemon, timeProvider)
	return hc
}

// InitializeHealthStatus initializes the health status of a service as unhealthy. The service will become healthy
// as soon as it reports it's first successful update. This method must be called, or the service will panic when
// the HealthCheck method is called. This method is synchronized.
func (h *HealthCheckableImpl) InitializeHealthStatus(serviceName string, timeProvider libtime.TimeProvider) {
	h.RecordUpdateFailure(timeProvider, fmt.Errorf("%v is initializing", serviceName))

	h.Lock()
	defer h.Unlock()
	h.initialized = true
}

// RecordUpdateSuccess records a successful update. This method is synchronized.
func (h *HealthCheckableImpl) RecordUpdateSuccess(timeProvider libtime.TimeProvider) {
	h.Lock()
	defer h.Unlock()

	h.lastSuccessfulUpdate = timeProvider.Now()
}

// RecordUpdateFailure records a failed update. This method is synchronized.
func (h *HealthCheckableImpl) RecordUpdateFailure(timeProvider libtime.TimeProvider, err error) {
	h.Lock()
	defer h.Unlock()

	h.lastFailedUpdate = timeProvider.Now()
	h.lastUpdateError = err
}

// HealthCheck returns an error if a service is unhealthy.
// The service is unhealthy if any of the following are true:
// - no successful update has occurred
// - the most recent update failed, or
// - the daemon has not seen a successful update in at least 5 minutes,
// This method is synchronized.
func (h *HealthCheckableImpl) HealthCheck(timeProvider libtime.TimeProvider) error {
	h.Lock()
	defer h.Unlock()

	if !h.initialized {
		panic("HealthCheckableImpl has not been initialized")
	}

	if h.lastSuccessfulUpdate.IsZero() {
		return fmt.Errorf("no successful update has occurred")
	}

	if h.lastFailedUpdate.After(h.lastSuccessfulUpdate) {
		return fmt.Errorf(
			"last update failed at %v with error: %w",
			h.lastFailedUpdate,
			h.lastUpdateError,
		)
	}

	// If the last successful update was more than 5 minutes ago, report the specific error.
	if timeProvider.Now().Sub(h.lastSuccessfulUpdate) > MaximumAcceptableUpdateDelay {
		return fmt.Errorf(
			"last successful update occurred at %v, which is more than %v ago",
			h.lastSuccessfulUpdate,
			MaximumAcceptableUpdateDelay,
		)
	}

	return nil
}
