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
	// HealthCheckPollFrequency is the frequency at which the health checkable service is polled.
	HealthCheckPollFrequency = 5 * time.Second
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

// NewHealthMonitorMutableState creates a new health monitor mutable state.
func NewHealthMonitorMutableState() *healthMonitorMutableState {
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

// RegisterHealthChecker registers a new health checker for a health checkable with the health monitor. The health
// checker is lazily created using the provided function if needed. This method is synchronized. It returns an error if
// the service was already registered.
func (ms *healthMonitorMutableState) RegisterHealthChecker(
	checkable types.HealthCheckable,
	lazyHealthCheckerCreator func() *healthChecker,
) error {
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
		// If the service is stoppable, stop it. This helps us to clean up daemon services in test cases
		// where the monitor is stopped before all daemon services have been registered.
		if stoppable, ok := checkable.(Stoppable); ok {
			stoppable.Stop()
		}
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
		mutableState:       NewHealthMonitorMutableState(),
		logger:             logger.With(cosmoslog.ModuleKey, "health-monitor"),
		startupGracePeriod: startupGracePeriod,
		pollingFrequency:   pollingFrequency,
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
	maximumAcceptableUnhealthyDuration time.Duration,
	callback func(error),
) error {
	if maximumAcceptableUnhealthyDuration <= 0 {
		return fmt.Errorf(
			"health check registration failure for service %v: "+
				"maximum acceptable unhealthy duration %v must be positive",
			hc.ServiceName(),
			maximumAcceptableUnhealthyDuration,
		)
	}

	return hm.mutableState.RegisterHealthChecker(hc, func() *healthChecker {
		return StartNewHealthChecker(
			hc,
			hm.pollingFrequency,
			callback,
			&libtime.TimeProviderImpl{},
			maximumAcceptableUnhealthyDuration,
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
// already been stopped. If the monitor has been stopped, this method will proactively stop the health-checkable
// service before returning.
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
	hm.mutableState.Stop()
}
