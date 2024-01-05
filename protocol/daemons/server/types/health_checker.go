package types

import (
	"cosmossdk.io/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"sync"
	"time"
)

// errorStreak tracks two relevant statistics for an error streak returned by a HealthCheckable - the timestamp of the
// beginning of the error streak, and the most recent error. This is useful for determining how long a service has been
// unhealthy, as well as the current state of the service.
type errorStreak struct {
	startOfStreak   time.Time
	mostRecentError error
}

// UpdateLastError updates the errorStreak to reflect the current error. If the startOfStreak timestamp is zero, this
// error the first error in a new error streak, so the startOfStreak timestamp is set to the current timestamp.
func (u *errorStreak) UpdateLastError(timestamp time.Time, err error) {
	// If the startOfStreak is zero, this is the first update, so set the startOfStreak.
	if u.startOfStreak.IsZero() {
		u.startOfStreak = timestamp
	}

	u.mostRecentError = err
}

// Reset resets the errorStreak to its zero value, indicating that the service has no active error streak.
func (u *errorStreak) Reset() {
	u.startOfStreak = time.Time{}
	u.mostRecentError = nil
}

// IsUnset returns true if the errorStreak is unset, indicating that the service has no active error streak.
func (u *errorStreak) IsUnset() bool {
	return u.startOfStreak.IsZero() && u.mostRecentError == nil
}

// StartOfStreak returns the timestamp of th start of the most recent error streak.
func (u *errorStreak) StartOfStreak() time.Time {
	return u.startOfStreak
}

// MostRecentError returns the most recent error associated with the current error streak.
func (u *errorStreak) MostRecentError() error {
	return u.mostRecentError
}

// healthCheckerMutableState tracks the current health state of the HealthCheckable, encapsulating all mutable state
// into a single struct for ease of synchronization.
type healthCheckerMutableState struct {
	// lock is used to synchronize access to mutable state fields.
	lock sync.Mutex

	// lastSuccessTimestamp is the startOfStreak of the most recent successful health check.
	// Access to lastSuccessTimestamp is synchronized.
	lastSuccessTimestamp time.Time

	// mostRecentErrorStreak tracks the beginning of the most recent streak, as well as the current error in the streak.
	// It is updated on every error and reset every time the service sees a healthy response.
	// This field is used to determine how long the daemon has been unhealthy. If the mostRecentErrorStreak is unset,
	// then either the service has never been unhealthy, or the most recent error streak ended before it could trigger
	// a callback.
	// Access to mostRecentErrorStreak is synchronized.
	mostRecentErrorStreak errorStreak

	// timer triggers a health check poll for a health-checkable service.
	timer *time.Timer

	// stopped indicates whether the health checker has been stopped. Additional health checks cannot be scheduled
	// after the health checker has been stopped.
	stopped bool
}

// ReportSuccess updates the health checker's mutable state to reflect a successful health check and schedules the next
// poll as an atomic operation.
func (u *healthCheckerMutableState) ReportSuccess(now time.Time) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.lastSuccessTimestamp = now

	// Whenever the service is healthy, reset the first failure in streak startOfStreak.
	u.mostRecentErrorStreak.Reset()
}

// ReportFailure updates the health checker's mutable state to reflect a failed health check and schedules the next
// poll as an atomic operation. The method returns the duration of the current failure streak.
func (u *healthCheckerMutableState) ReportFailure(now time.Time, err error) time.Duration {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.mostRecentErrorStreak.UpdateLastError(now, err)

	return now.Sub(u.mostRecentErrorStreak.StartOfStreak())
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

// InitializePolling schedules the first poll for the health-checkable service. This method is meant to be called
// immediately after initializing the health checker mutable state. This method is synchronized.
func (u *healthCheckerMutableState) InitializePolling(firstPollDelay time.Duration, pollFunc func()) {
	u.lock.Lock()
	defer u.lock.Unlock()

	// If the timer is already initialized, don't initialize it again.
	if u.timer != nil {
		return
	}

	// The first poll is scheduled after a custom delay to allow the service to initialize.
	u.timer = time.AfterFunc(firstPollDelay, pollFunc)
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

	// maxUnhealthyDuration is the maximum acceptable duration for a health-checkable service to
	// remain unhealthy. If the service remains unhealthy for this duration, the monitor will execute the
	// specified callback function.
	maxUnhealthyDuration time.Duration

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

	if err == nil { // Capture healthy response.
		hc.mutableState.ReportSuccess(now)
	} else { // Capture unhealthy response.
		streakDuration := hc.mutableState.ReportFailure(now, err)
		// If the service has been unhealthy for longer than the maximum acceptable unhealthy duration, execute the
		// callback function.
		if streakDuration >= hc.maxUnhealthyDuration {
			hc.unhealthyCallback(err)
		}
	}

	// Schedule next poll. We schedule another poll whether the callback was invoked or not, as callbacks are not
	// guaranteed to panic or otherwise halt the daemon. In such cases, we may end up invoking the callback several
	// times once the service exceeds the maximum unhealthy duration. For example, a callback that emits error logs
	// will continue to emit error logs every 5s until the service becomes healthy again.
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
	maxUnhealthyDuration time.Duration,
	startupGracePeriod time.Duration,
	logger log.Logger,
) *healthChecker {
	checker := &healthChecker{
		healthCheckable:      healthCheckable,
		pollFrequency:        pollFrequency,
		unhealthyCallback:    unhealthyCallback,
		timeProvider:         timeProvider,
		maxUnhealthyDuration: maxUnhealthyDuration,
		logger:               logger,
		mutableState:         &healthCheckerMutableState{},
	}

	// The first poll is scheduled after the startup grace period to allow the service to initialize.
	// We initialize the timer and schedule a poll outside of object creation in order to avoid data races for
	// extremely short startup grace periods.
	checker.mutableState.InitializePolling(startupGracePeriod, checker.Poll)

	return checker
}
