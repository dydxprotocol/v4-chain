package types

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"sync"
	"time"
)

// timestampWithError couples a timestamp and error to make it easier to update them in tandem. The
// timestampWithError will track the timestamp of the first error in a streak of errors, but keeps a record of the
// most recent error. This is useful for determining how long a service has been unhealthy, as well as the current
// state of the service.
type timestampWithError struct {
	timestamp time.Time
	err       error
}

// Update updates the timeStampWithError to reflect the current error. If the timestamp is zero, this is the first
// update, so set the timestamp.
func (u *timestampWithError) Update(timestamp time.Time, err error) {
	// If the timestamp is zero, this is the first update, so set the timestamp.
	if u.timestamp.IsZero() {
		u.timestamp = timestamp
	}

	u.err = err
}

// Reset resets the timestampWithError to its zero value, indicating that the service is healthy.
func (u *timestampWithError) Reset() {
	u.timestamp = time.Time{}
	u.err = nil
}

// IsZero returns true if the timestampWithError is zero, indicating that the service is healthy.
func (u *timestampWithError) IsZero() bool {
	return u.timestamp.IsZero() && u.err == nil
}

// Timestamp returns the timestamp associated with the timestampWithError, which is the timestamp of the first error
// in the current error streak.
func (u *timestampWithError) Timestamp() time.Time {
	return u.timestamp
}

// Error returns the error associated with the timestampWithError, which is the most recent error in the current error
// streak.
func (u *timestampWithError) Error() error {
	return u.err
}

// healthCheckerMutableState tracks the current health state of the HealthCheckable, encapsulating all mutable state
// into a single struct for ease of synchronization.
type healthCheckerMutableState struct {
	// lock is used to synchronize access to mutable state fields.
	lock sync.Mutex

	// lastSuccessTimestamp is the timestamp of the most recent successful health check.
	// Access to lastSuccessTimestamp is synchronized.
	lastSuccessTimestamp time.Time

	// mostRecentFailureStreakError tracks the timestamp of the first error in the most recent streak of errors, as well
	// as the most recent error. It is updated on every error and reset every time the service sees a healthy response.
	// This field is used to determine how long the daemon has been unhealthy. If this timestamp is nil, then either
	// the service has never been unhealthy, or the most recent error streak ended before it could trigger a callback.
	// Access to mostRecentFailureStreakError is synchronized.
	mostRecentFailureStreakError timestampWithError

	// timer triggers a health check poll for a health-checkable service.
	timer *time.Timer

	// stopped indicates whether the health checker has been stopped. Additional health checks cannot be scheduled
	// after the health checker has been stopped.
	stopped bool
}

// newHealthCheckerMutableState creates a new health checker mutable state scheduled to trigger a poll after the
// initial poll delay.
func newHealthCheckerMutableState(initialPollDelay time.Duration, pollFunc func()) *healthCheckerMutableState {
	return &healthCheckerMutableState{
		timer: time.AfterFunc(initialPollDelay, pollFunc),
	}
}

// ReportSuccess updates the health checker's mutable state to reflect a successful health check and schedules the next
// poll as an atomic operation.
func (u *healthCheckerMutableState) ReportSuccess(now time.Time) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.lastSuccessTimestamp = now

	// Whenever the service is healthy, reset the first failure in streak timestamp.
	u.mostRecentFailureStreakError.Reset()
}

// ReportFailure updates the health checker's mutable state to reflect a failed health check and schedules the next
// poll as an atomic operation. The method returns the duration of the current failure streak.
func (u *healthCheckerMutableState) ReportFailure(now time.Time, err error) time.Duration {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.mostRecentFailureStreakError.Update(now, err)

	return now.Sub(u.mostRecentFailureStreakError.Timestamp())
}

// SchedulePoll schedules the next poll for the health-checkable service. If the service is stopped, the next poll
// will not be scheduled. This method is synchronized.
func (u *healthCheckerMutableState) SchedulePoll(nextPollDelay time.Duration) {
	u.lock.Lock()
	defer u.lock.Unlock()

	// Don't schedule a poll if the health checker has been stopped.
	if u.stopped {
		return
	}

	// Schedule the next poll.
	u.timer.Reset(nextPollDelay)
}

// Stop stops the health checker. This method is synchronized.
func (u *healthCheckerMutableState) Stop() {
	u.lock.Lock()
	defer u.lock.Unlock()

	// Don't stop the health checker if it has already been stopped.
	if u.stopped {
		return
	}

	u.timer.Stop()
	u.stopped = true
}

// healthChecker encapsulates the logic for monitoring the health of a health-checkable service.
type healthChecker struct {
	// mutableState is the mutable state of the health checker. Access to these fields is synchronized.
	mutableState *healthCheckerMutableState

	// healthCheckable is the health-checkable service to be monitored.
	healthCheckable types.HealthCheckable

	// pollFrequency is the frequency at which the health-checkable service is polled.
	pollFrequency time.Duration

	// maxAcceptableUnhealthyDuration is the maximum acceptable duration for a health-checkable service to
	// remain unhealthy. If the service remains unhealthy for this duration, the monitor will execute the
	// specified callback function.
	maxAcceptableUnhealthyDuration time.Duration

	// unhealthyCallback is the callback function to be executed if the health-checkable service remains
	// unhealthy for a period of time greater than or equal to the maximum acceptable unhealthy duration.
	// This callback function is executed with the error that caused the service to become unhealthy.
	unhealthyCallback func(error)

	// timeProvider is used to get the current time. It is added as a field so that it can be mocked in tests.
	timeProvider libtime.TimeProvider

	logger log.Logger
}

// Poll executes a health check for the health-checkable service. If the service has been unhealthy for longer than the
// maximum acceptable unhealthy duration, the callback function is executed.
// This method is publicly exposed for testing. This method is synchronized.
func (hc *healthChecker) Poll() {
	err := hc.healthCheckable.HealthCheck()
	now := hc.timeProvider.Now()

	// Capture healthy response.

	if err == nil {
		hc.mutableState.ReportSuccess(now)
	} else {
		streakDuration := hc.mutableState.ReportFailure(now, err)
		// If the service has been unhealthy for longer than the maximum acceptable unhealthy duration, execute the
		// callback function.
		if streakDuration >= hc.maxAcceptableUnhealthyDuration {
			hc.unhealthyCallback(err)
		}
	}

	// Schedule next poll.
	hc.mutableState.SchedulePoll(hc.pollFrequency)
}

// Stop stops the health checker. This method is not synchronized, as the timer does not need synchronization.
func (hc *healthChecker) Stop() {
	hc.mutableState.Stop()
}

// StartNewHealthChecker creates and starts a new health checker for a health-checkable service.
func StartNewHealthChecker(
	healthCheckable types.HealthCheckable,
	pollFrequency time.Duration,
	unhealthyCallback func(error),
	timeProvider libtime.TimeProvider,
	maximumAcceptableUnhealthyDuration time.Duration,
	startupGracePeriod time.Duration,
	logger log.Logger,
) *healthChecker {
	checker := &healthChecker{
		healthCheckable:                healthCheckable,
		pollFrequency:                  pollFrequency,
		unhealthyCallback:              unhealthyCallback,
		timeProvider:                   timeProvider,
		maxAcceptableUnhealthyDuration: maximumAcceptableUnhealthyDuration,
		logger:                         logger,
	}

	// The first poll is scheduled after the startup grace period to allow the service to initialize.
	checker.mutableState = newHealthCheckerMutableState(startupGracePeriod, checker.Poll)

	return checker
}
