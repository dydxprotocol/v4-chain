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
// monitor has already been stopped. If the monitor has been stopped, this method will proactively stop the
// health-checkable service before returning.
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
		// If the service is stoppable, stop it. This helps us to clean up daemon services in test cases
		// where the monitor is stopped before all daemon services have been registered.
		if stoppable, ok := hc.(Stoppable); ok {
			stoppable.Stop()
		}

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
		hm.logger,
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
