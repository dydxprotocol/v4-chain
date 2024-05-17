package types

import (
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
)

const (
	MaxAcceptableUpdateDelay = 5 * time.Minute
)

// HealthCheckable is a common interface for services that can be health checked.
//
// Instances of this type are thread-safe.
type HealthCheckable interface {
	// HealthCheck returns an error if a service is unhealthy. If the service is healthy, this method returns nil.
	HealthCheck() (err error)
	// ReportFailure records a failed update.
	ReportFailure(err error)
	// ReportSuccess records a successful update.
	ReportSuccess()
	// ServiceName returns the name of the service being monitored. This name is expected to be unique.
	ServiceName() string
}

// timestampWithError couples a timestamp and error to make it easier to update them in tandem.
type timestampWithError struct {
	timestamp time.Time
	err       error
}

func (u *timestampWithError) Update(timestamp time.Time, err error) {
	u.timestamp = timestamp
	u.err = err
}

func (u *timestampWithError) Timestamp() time.Time {
	return u.timestamp
}

func (u *timestampWithError) Error() error {
	return u.err
}

// timeBoundedHealthCheckable implements the HealthCheckable interface by tracking the timestamps of the last successful
// and failed updates.
// If any of the following occurs, then the service should be considered unhealthy:
// - no update has occurred
// - the most recent update failed, or
// - the daemon has not seen a successful update within `MaxAcceptableUpdateDelay`.
//
// This object is thread-safe.
type timeBoundedHealthCheckable struct {
	sync.Mutex

	// lastSuccessfulUpdate is the timestamp of the last successful update.
	lastSuccessfulUpdate time.Time
	// lastFailedUpdate is the timestamp, error pair of the last failed update.
	lastFailedUpdate timestampWithError

	// timeProvider is the time provider used to determine the current time. This is used for timestamping
	// creation and checking for update staleness during HealthCheck.
	timeProvider libtime.TimeProvider

	// logger is the logger used to log errors.
	logger log.Logger

	// serviceName is the name of the service being monitored. This field is read-only and not synchronized.
	serviceName string
}

// NewTimeBoundedHealthCheckable creates a new HealthCheckable instance.
func NewTimeBoundedHealthCheckable(
	serviceName string,
	timeProvider libtime.TimeProvider,
	logger log.Logger,
) HealthCheckable {
	hc := &timeBoundedHealthCheckable{
		timeProvider: timeProvider,
		logger:       logger,
		serviceName:  serviceName,
	}
	// Initialize the timeBoudnedHealthCheckable to an unhealthy state by reporting an error.
	hc.ReportFailure(fmt.Errorf("%v is initializing", serviceName))
	return hc
}

// ServiceName returns the name of the service being monitored.
func (hc *timeBoundedHealthCheckable) ServiceName() string {
	return hc.serviceName
}

// ReportSuccess records a successful update. This method is thread-safe.
func (h *timeBoundedHealthCheckable) ReportSuccess() {
	h.Lock()
	defer h.Unlock()

	h.lastSuccessfulUpdate = h.timeProvider.Now()
}

// ReportFailure records a failed update. This method is thread-safe.
func (h *timeBoundedHealthCheckable) ReportFailure(err error) {
	h.Lock()
	defer h.Unlock()
	h.lastFailedUpdate.Update(h.timeProvider.Now(), err)
}

// HealthCheck returns an error if a service is unhealthy.
// The service is unhealthy if any of the following are true:
// - no successful update has occurred
// - the most recent update failed, or
// - the daemon has not seen a successful update in at least 5 minutes,
// Note: since the timeBoundedHealthCheckable is not exposed and can only be created via
// NewTimeBoundedHealthCheckable, we expect that the lastFailedUpdate is never a zero value.
// This method is thread-safe.
func (h *timeBoundedHealthCheckable) HealthCheck() error {
	h.Lock()
	defer h.Unlock()

	if h.lastSuccessfulUpdate.IsZero() {
		return fmt.Errorf(
			"no successful update has occurred; last failed update occurred at %v with error '%w'",
			h.lastFailedUpdate.Timestamp(),
			h.lastFailedUpdate.Error(),
		)
	}

	if h.lastFailedUpdate.Timestamp().After(h.lastSuccessfulUpdate) {
		return fmt.Errorf(
			"last update failed at %v with error: '%w', most recent successful update occurred at %v",
			h.lastFailedUpdate.Timestamp(),
			h.lastFailedUpdate.Error(),
			h.lastSuccessfulUpdate,
		)
	}

	// If the last successful update was more than 5 minutes ago, log the specific error.
	if h.timeProvider.Now().Sub(h.lastSuccessfulUpdate) > MaxAcceptableUpdateDelay {
		h.logger.Error(
			fmt.Sprintf(
				"last successful update occurred at %v, which is more than %v ago. "+
					"Last failure occurred at %v with error '%v'",
				h.lastSuccessfulUpdate,
				MaxAcceptableUpdateDelay,
				h.lastFailedUpdate.Timestamp(),
				h.lastFailedUpdate.Error(),
			),
		)
	}

	return nil
}
