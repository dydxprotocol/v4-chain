package types

import (
	cosmoslog "cosmossdk.io/log"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"sync"
	"time"
)

const (
	HealthCheckPollFrequency = 5 * time.Second
)

// timestampWithError couples a timestamp and error to make it easier to update them in tandem.
type timestampWithError struct {
	timestamp time.Time
	err       error
}

func (u *timestampWithError) Update(timestamp time.Time, err error) {
	u.timestamp = timestamp
	u.err = err
}

func (u *timestampWithError) Reset() {
	u.Update(time.Time{}, nil)
}

func (u *timestampWithError) IsZero() bool {
	return u.timestamp.IsZero()
}

func (u *timestampWithError) Timestamp() time.Time {
	return u.timestamp
}

func (u *timestampWithError) Error() error {
	return u.err
}

// healthChecker encapsulates the logic for monitoring the health of a health checkable service.
type healthChecker struct {
	// healthCheckable is the health checkable service to be monitored.
	healthCheckable types.HealthCheckable

	// timer triggers a health check poll for a health checkable service.
	timer *time.Timer

	// pollFrequency is the frequency at which the health checkable service is polled.
	pollFrequency time.Duration

	// mostRecentSuccess is the timestamp of the most recent successful health check.
	mostRecentSuccess time.Time

	// firstFailureInStreak is the timestamp of the first error in the most recent streak of errors. It is set
	// whenever the service toggles from healthy to an unhealthy state, and used to determine how long the daemon has
	// been unhealthy. If this timestamp is nil, then the error streak ended before it could trigger a callback.
	firstFailureInStreak timestampWithError

	// unhealthyCallback is the callback function to be executed if the health checkable service remains
	// unhealthy for a period of time greater than or equal to the maximum acceptable unhealthy duration.
	// This callback function is executed with the error that caused the service to become unhealthy.
	unhealthyCallback func(error)

	// timeProvider is used to get the current time. It is added as a field so that it can be mocked in tests.
	timeProvider libtime.TimeProvider

	// maximumAcceptableUnhealthyDuration is the maximum acceptable duration for a health checkable service to
	// remain unhealthy. If the service remains unhealthy for this duration, the monitor will execute the
	// specified callback function.
	maximumAcceptableUnhealthyDuration time.Duration
}

// Poll executes a health check for the health checkable service. If the service has been unhealthy for longer than the
// maximum acceptable unhealthy duration, the callback function is executed.
// This method is publicly exposed for testing.
func (hc *healthChecker) Poll() {
	// Don't return an error if the monitor has been disabled.
	err := hc.healthCheckable.HealthCheck()
	now := hc.timeProvider.Now()
	if err == nil {
		hc.mostRecentSuccess = now
		// Whenever the service is healthy, reset the first failure in streak timestamp.
		hc.firstFailureInStreak.Reset()
	} else if hc.firstFailureInStreak.IsZero() {
		// Capture the timestamp of the first failure in a new streak.
		hc.firstFailureInStreak.Update(now, err)
	}

	// If the service has been unhealthy for longer than the maximum acceptable unhealthy duration, execute the
	// callback function.
	if !hc.firstFailureInStreak.IsZero() &&
		now.Sub(hc.firstFailureInStreak.Timestamp()) >= hc.maximumAcceptableUnhealthyDuration {
		hc.unhealthyCallback(hc.firstFailureInStreak.Error())
	} else {
		// If we do not execute the callback, schedule the next poll.
		hc.timer.Reset(hc.pollFrequency)
	}
}

func (hc *healthChecker) Stop() {
	hc.timer.Stop()
}

// StartNewHealthChecker creates and starts a new health checker for a health checkable service.
func StartNewHealthChecker(
	healthCheckable types.HealthCheckable,
	pollFrequency time.Duration,
	unhealthyCallback func(error),
	timeProvider libtime.TimeProvider,
	maximumAcceptableUnhealthyDuration time.Duration,
	startupGracePeriod time.Duration,
) *healthChecker {
	checker := &healthChecker{
		healthCheckable:                    healthCheckable,
		pollFrequency:                      pollFrequency,
		unhealthyCallback:                  unhealthyCallback,
		timeProvider:                       timeProvider,
		maximumAcceptableUnhealthyDuration: maximumAcceptableUnhealthyDuration,
	}
	// The first poll is scheduled after the startup grace period to allow the service to initialize.
	checker.timer = time.AfterFunc(startupGracePeriod, checker.Poll)

	return checker
}

// HealthMonitor monitors the health of daemon services, which implement the HealthCheckable interface. If a
// registered health-checkable service sustains an unhealthy state for the maximum acceptable unhealthy duration,
// the monitor will execute a callback function.
type HealthMonitor struct {
	// serviceToHealthChecker maps daemon service names to their update metadata.
	serviceToHealthChecker map[string]*healthChecker
	// stopped indicates whether the monitor has been stopped. Additional daemon services cannot be registered
	// after the monitor has been stopped.
	stopped bool
	// disabled indicates whether the monitor has been disabled. This is used to disable the monitor in testApp
	// tests, where app.New is not executed.
	disabled bool
	// lock is used to synchronize access to the monitor.
	lock sync.Mutex

	// These fields are initialized in NewHealthMonitor and are not modified after initialization.
	logger             log.Logger
	startupGracePeriod time.Duration
	pollingFrequency   time.Duration
}

// NewHealthMonitor creates a new health monitor.
func NewHealthMonitor(
	startupGracePeriod time.Duration,
	pollingFrequency time.Duration,
	logger log.Logger,
) *HealthMonitor {
	return &HealthMonitor{
		serviceToHealthChecker: make(map[string]*healthChecker),
		logger:                 logger.With(cosmoslog.ModuleKey, "health-monitor"),
		startupGracePeriod:     startupGracePeriod,
		pollingFrequency:       pollingFrequency,
	}
}

func (hm *HealthMonitor) DisableForTesting() {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	hm.disabled = true
}

// RegisterServiceWithCallback registers a HealthCheckable with the health monitor. If the service
// stays unhealthy every time it is polled during the maximum acceptable unhealthy duration, the monitor will
// execute the callback function.
// This method is synchronized. The method returns an error if the service was already registered or the
// monitor has already been stopped.
func (hm *HealthMonitor) RegisterServiceWithCallback(
	hc types.HealthCheckable,
	maximumAcceptableUnhealthyDuration time.Duration,
	callback func(error),
) error {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	if maximumAcceptableUnhealthyDuration <= 0 {
		return fmt.Errorf(
			"health check registration failure for service %v: "+
				"maximum acceptable unhealthy duration %v must be positive",
			hc.ServiceName(),
			maximumAcceptableUnhealthyDuration,
		)
	}

	// Don't register daemon services if the monitor has been disabled.
	if hm.disabled {
		return nil
	}

	// Don't register additional daemon services if the monitor has already been stopped.
	// This could be a concern for short-running integration test cases, where the network
	// stops before all daemon services have been registered.
	if hm.stopped {
		return fmt.Errorf(
			"health check registration failure for service %v: monitor has been stopped",
			hc.ServiceName(),
		)
	}

	if _, ok := hm.serviceToHealthChecker[hc.ServiceName()]; ok {
		return fmt.Errorf("service %v already registered", hc.ServiceName())
	}

	hm.serviceToHealthChecker[hc.ServiceName()] = StartNewHealthChecker(
		hc,
		hm.pollingFrequency,
		callback,
		&libtime.TimeProviderImpl{},
		maximumAcceptableUnhealthyDuration,
		hm.startupGracePeriod,
	)
	return nil
}

// PanicServiceNotResponding returns a function that panics with a message indicating that the specified daemon
// service is not responding. This is ideal for creating a callback function when registering a daemon service.
func PanicServiceNotResponding(hc types.HealthCheckable) func(error) {
	return func(err error) {
		panic(fmt.Sprintf("%v unhealthy: %v", hc.ServiceName(), err))
	}
}

// LogErrorServiceNotResponding returns a function that logs an error indicating that the specified service
// is not responding. This is ideal for creating a callback function when registering a health checkable service.
func LogErrorServiceNotResponding(hc types.HealthCheckable, logger log.Logger) func(error) {
	return func(err error) {
		logger.Error(
			"health-checked service is unhealthy",
			"service",
			hc.ServiceName(),
			"error",
			err,
		)
	}
}

// RegisterService registers a new health checkable service with the health check monitor. If the service
// is unhealthy every time it is polled for a duration greater than or equal to the maximum acceptable unhealthy
// duration, the monitor will panic.
// This method is synchronized. It returns an error if the service was already registered or the monitor has
// already been stopped.
func (hm *HealthMonitor) RegisterService(
	hc types.HealthCheckable,
	maximumAcceptableUnhealthyDuration time.Duration,
) error {
	return hm.RegisterServiceWithCallback(
		hc,
		maximumAcceptableUnhealthyDuration,
		PanicServiceNotResponding(hc),
	)
}

// Stop stops the update frequency monitor. This method is synchronized.
func (hm *HealthMonitor) Stop() {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	// Don't stop the monitor if it has already been stopped.
	if hm.stopped {
		return
	}

	for _, checker := range hm.serviceToHealthChecker {
		checker.Stop()
	}

	hm.stopped = true
}
