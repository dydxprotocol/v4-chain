package types

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"sync"
	"time"
)

const (
	// HealthCheckPollFrequency is the frequency at which the health-checkable service is polled.
	HealthCheckPollFrequency = 5 * time.Second

	// HealthMonitorLogModuleName is the module name used for logging within the health monitor.
	HealthMonitorLogModuleName = "daemon-health-monitor"
)

// healthMonitorMutableState tracks all mutable state associated with the health monitor. This state is gathered into
// a single struct for ease of synchronization.
type healthMonitorMutableState struct {
	sync.Mutex

	// serviceToHealthChecker maps daemon service names to their update metadata.
	serviceToHealthChecker map[string]*healthChecker
	// stopped indicates whether the monitor has been stopped. Additional daemon services cannot be registered
	// after the monitor has been stopped.
	stopped bool
	// disabled indicates whether the monitor has been disabled. This is used to disable the monitor in testApp
	// tests, where app.New is not executed.
	disabled bool
}

// newHealthMonitorMutableState creates a new health monitor mutable state.
func newHealthMonitorMutableState() *healthMonitorMutableState {
	return &healthMonitorMutableState{
		serviceToHealthChecker: make(map[string]*healthChecker),
	}
}

// DisableForTesting disables the health monitor mutable state from receiving updates. This prevents the monitor
// from registering services when called before app initialization and is used for testing.
func (ms *healthMonitorMutableState) DisableForTesting() {
	ms.Lock()
	defer ms.Unlock()

	ms.disabled = true
}

// Stop stops the update frequency monitor. This method is synchronized.
func (ms *healthMonitorMutableState) Stop() {
	ms.Lock()
	defer ms.Unlock()

	// Don't stop the monitor if it has already been stopped.
	if ms.stopped {
		return
	}

	// Stop all health checkers.
	for _, checker := range ms.serviceToHealthChecker {
		checker.Stop()
	}

	ms.stopped = true
}

// RegisterHealthChecker registers a new health checker for a health-checkable with the health monitor. The health
// checker is lazily created using the provided function if needed. This method is synchronized. It returns an error if
// the service was already registered.
func (ms *healthMonitorMutableState) RegisterHealthChecker(
	checkable types.HealthCheckable,
	lazyHealthCheckerCreator func() *healthChecker,
) error {
	stopService := false

	// If the monitor has already been stopped, we want to stop the checkable service before returning.
	// However, we'd prefer not to stop the service within the critical section in order to prevent deadlocks.
	// This defer will be called last, after the lock is released.
	defer func() {
		if stopService {
			// If the service is stoppable, stop it. This helps us to clean up daemon services in test cases
			// where the monitor is stopped before all daemon services have been registered.
			if stoppable, ok := checkable.(Stoppable); ok {
				stoppable.Stop()
			}
		}
	}()

	// Enter into the critical section.
	ms.Lock()
	defer ms.Unlock()

	// Don't register daemon services if the monitor has been disabled.
	if ms.disabled {
		return nil
	}

	// Don't register additional daemon services if the monitor has already been stopped.
	// This could be a concern for short-running integration test cases, where the network
	// stops before all daemon services have been registered.
	if ms.stopped {
		// Toggle the stopService flag to true so that the service is stopped after the lock is released.
		stopService = true
		return nil
	}

	if _, ok := ms.serviceToHealthChecker[checkable.ServiceName()]; ok {
		return fmt.Errorf("service %v already registered", checkable.ServiceName())
	}

	ms.serviceToHealthChecker[checkable.ServiceName()] = lazyHealthCheckerCreator()
	return nil
}

// HealthMonitor monitors the health of daemon services, which implement the HealthCheckable interface. If a
// registered health-checkable service sustains an unhealthy state for the maximum acceptable unhealthy duration,
// the monitor will execute a callback function.
type HealthMonitor struct {
	mutableState *healthMonitorMutableState

	// These fields are initialized in NewHealthMonitor and are not modified after initialization.
	logger log.Logger
	// startupGracePeriod is the grace period before the monitor starts polling the health-checkable services.
	startupGracePeriod time.Duration
	// pollingFrequency is the frequency at which the health-checkable services are polled.
	pollingFrequency time.Duration
	// enablePanics is used to toggle between panics or error logs when a daemon sustains an unhealthy state past the
	// maximum allowable duration.
	enablePanics bool
}

// NewHealthMonitor creates a new health monitor.
func NewHealthMonitor(
	startupGracePeriod time.Duration,
	pollingFrequency time.Duration,
	logger log.Logger,
	enablePanics bool,
) *HealthMonitor {
	return &HealthMonitor{
		mutableState:       newHealthMonitorMutableState(),
		logger:             logger.With(log.ModuleKey, HealthMonitorLogModuleName),
		startupGracePeriod: startupGracePeriod,
		pollingFrequency:   pollingFrequency,
		enablePanics:       enablePanics,
	}
}

func (hm *HealthMonitor) DisableForTesting() {
	hm.mutableState.DisableForTesting()
}

// RegisterServiceWithCallback registers a HealthCheckable with the health monitor. If the service
// stays unhealthy every time it is polled during the maximum acceptable unhealthy duration, the monitor will
// execute the callback function.
// This method is synchronized. The method returns an error if the service was already registered or the
// monitor has already been stopped. If the monitor has been stopped, this method will proactively stop the
// health-checkable service before returning.
func (hm *HealthMonitor) RegisterServiceWithCallback(
	hc types.HealthCheckable,
	maxUnhealthyDuration time.Duration,
	callback func(error),
) error {
	if maxUnhealthyDuration <= 0 {
		return fmt.Errorf(
			"health check registration failure for service %v: "+
				"maximum unhealthy duration %v must be positive",
			hc.ServiceName(),
			maxUnhealthyDuration,
		)
	}

	return hm.mutableState.RegisterHealthChecker(hc, func() *healthChecker {
		return StartNewHealthChecker(
			hc,
			hm.pollingFrequency,
			callback,
			&libtime.TimeProviderImpl{},
			maxUnhealthyDuration,
			hm.startupGracePeriod,
			hm.logger,
		)
	})
}

// PanicServiceNotResponding returns a function that panics with a message indicating that the specified daemon
// service is not responding. This is ideal for creating a callback function when registering a daemon service.
func PanicServiceNotResponding(hc types.HealthCheckable) func(error) {
	return func(err error) {
		panic(fmt.Sprintf("%v unhealthy: %v", hc.ServiceName(), err))
	}
}

// LogErrorServiceNotResponding returns a function that logs an error indicating that the specified service
// is not responding. This is ideal for creating a callback function when registering a health-checkable service.
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

// RegisterService registers a new health-checkable service with the health check monitor. If the service
// is unhealthy every time it is polled for a duration greater than or equal to the maximum acceptable unhealthy
// duration, the monitor will panic or log an error, depending on the app configuration via the
// `panic-on-daemon-failure-enabled` flag.
// This method is synchronized. It returns an error if the service was already registered or the monitor has
// already been stopped. If the monitor has been stopped, this method will proactively stop the health-checkable
// service before returning.
func (hm *HealthMonitor) RegisterService(
	hc types.HealthCheckable,
	maxDaemonUnhealthyDuration time.Duration,
) error {
	// If the monitor is configured to panic, use the panic callback. Otherwise, use the error log callback.
	// This behavior is configured via flag and defaults to panicking on daemon failure.
	callback := LogErrorServiceNotResponding(hc, hm.logger)
	if hm.enablePanics {
		callback = PanicServiceNotResponding(hc)
	}

	return hm.RegisterServiceWithCallback(
		hc,
		maxDaemonUnhealthyDuration,
		callback,
	)
}

// Stop stops the update frequency monitor. This method is synchronized.
func (hm *HealthMonitor) Stop() {
	hm.mutableState.Stop()
}
