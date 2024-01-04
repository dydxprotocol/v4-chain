package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	TestError1 = fmt.Errorf("test error 1")
	TestError2 = fmt.Errorf("test error 2")

	Time0 = time.Unix(0, 0)
	Time1 = Time0.Add(time.Second)
	Time2 = Time0.Add(2 * time.Second)
	Time3 = Time0.Add(3 * time.Second)
	Time4 = Time0.Add(4 * time.Second)
	Time5 = Time0.Add(5 * time.Second)

	// Use a maximum unhealthy duration of 5 seconds for testing, simulating a poll frequency of 1s with timestamps.
	TestMaximumUnhealthyDuration = 5 * time.Second

	// TestLargeDuration is used to ensure that the health checker does not trigger a callback through the timer.
	TestLargeDuration = 5 * time.Minute

	SmallDuration = 10 * time.Millisecond
)

// createTestMonitor creates a health monitor with a poll frequency of 10ms and a zero duration grace period.
func createTestMonitor() (*types.HealthMonitor, *mocks.Logger) {
	logger := &mocks.Logger{}
	logger.On("With", "module", "daemon-health-monitor").Return(logger).Once()
	return types.NewHealthMonitor(
		SmallDuration,
		10*time.Millisecond,
		logger,
		true, // enable panics here for stricter testing - a panic will definitely cause a test failure.
	), logger
}

// mockFailingHealthCheckerWithError creates a mock health checkable service that returns the given error.
func mockFailingHealthCheckerWithError(name string, err error) *mocks.HealthCheckable {
	hc := &mocks.HealthCheckable{}
	hc.On("ServiceName").Return(name)
	hc.On("HealthCheck").Return(err)
	return hc
}

// The following tests may still intermittently fail or report false positives on an overloaded system as they rely
// on callbacks to execute before the termination of the `time.Sleep` call, which is not guaranteed.
func TestRegisterService_Healthy(t *testing.T) {
	// Setup.
	ufm, logger := createTestMonitor()
	hc := mockFailingHealthCheckerWithError("test-service", nil)

	// Act.
	err := ufm.RegisterService(hc, 50*time.Millisecond)
	require.NoError(t, err)

	// Cleanup.
	defer func() {
		ufm.Stop()
	}()

	// Give the monitor time to poll the health checkable service. Polls occur once every 10ms.
	time.Sleep(100 * time.Millisecond)

	// Assert: no calls to the logger were made.
	mock.AssertExpectationsForObjects(t, hc, logger)
}

func TestRegisterServiceWithCallback_Mixed(t *testing.T) {
	tests := map[string]struct {
		healthCheckResponse error
	}{
		"healthy: no callback triggered": {
			healthCheckResponse: nil,
		},
		"unhealthy: callback triggered": {
			healthCheckResponse: TestError1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ufm, logger := createTestMonitor()
			hc := mockFailingHealthCheckerWithError("test-service", test.healthCheckResponse)
			errs := make(chan (error), 100)

			// Act.
			err := ufm.RegisterServiceWithCallback(
				hc,
				50*time.Millisecond,
				func(err error) {
					errs <- err
				},
			)
			require.NoError(t, err)

			// Give the monitor time to poll the health checkable service. Polls occur once every 10ms.
			select {
			case err = <-errs:
			case <-time.After(1 * time.Second):
			}

			// Cleanup.
			ufm.Stop()

			// Assert: no calls to the logger were made.
			mock.AssertExpectationsForObjects(t, hc, logger)

			// Assert: the callback was called or not called as expected.
			require.Equal(t, test.healthCheckResponse, err)
		})
	}
}

func TestHealthMonitor_DisablePanics_DoesNotPanic(t *testing.T) {
	logger := &mocks.Logger{}
	logger.On("With", "module", "daemon-health-monitor").Return(logger).Once()
	logger.On(
		"Error",
		"health-checked service is unhealthy",
		"service",
		"test-service",
		"error",
		mock.Anything,
	).Return()

	hm := types.NewHealthMonitor(
		SmallDuration,
		10*time.Millisecond,
		logger,
		false,
	)

	hc := mockFailingHealthCheckerWithError("test-service", TestError1)

	err := hm.RegisterService(hc, 10*time.Millisecond)
	require.NoError(t, err)

	defer func() {
		hm.Stop()
	}()

	// A 100ms sleep should be sufficient for the health monitor to detect the unhealthy service and trigger a callback.
	time.Sleep(100 * time.Millisecond)

	// Assert.
	// This test is confirmed to panic when panics are not disabled - but because the panic occurs in a separate
	// go-routine, it cannot be easily captured with an assert. Instead, we do not try to capture the panic, but
	// assert that the logger was called with the expected arguments.
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterServiceWithCallback_DoubleRegistrationFails(t *testing.T) {
	// Setup.
	ufm, logger := createTestMonitor()

	hc := mockFailingHealthCheckerWithError("test-service", TestError1)
	hc2 := mockFailingHealthCheckerWithError("test-service", TestError2)

	errs := make(chan (error), 100)

	err := ufm.RegisterServiceWithCallback(hc, 50*time.Millisecond, func(err error) {
		errs <- err
	})
	require.NoError(t, err)

	// Register a service with the same name. This registration should fail.
	err = ufm.RegisterServiceWithCallback(hc2, 50*time.Millisecond, func(err error) {
		errs <- err
	})
	require.ErrorContains(t, err, "service already registered")

	// Expect that the first service is still operating and will produce a callback after a sustained unhealthy period.
	select {
	case err = <-errs:
	case <-time.After(1 * time.Second):
		t.Fatalf("Failed to receive callback before timeout")
	}

	ufm.Stop()

	// Assert no calls to the logger were made.
	mock.AssertExpectationsForObjects(t, logger, hc)
	hc2.AssertNotCalled(t, "HealthCheck")

	// Assert the callback was called with the expected error.
	require.Equal(t, TestError1, err)
}

// Create a struct that implements HealthCheckable and Stoppable to check that the monitor stops the service.
type stoppableFakeHealthChecker struct {
	stopped bool
}

// Implement stub methods to conform to interfaces.
func (f *stoppableFakeHealthChecker) ServiceName() string   { return "test-service" }
func (f *stoppableFakeHealthChecker) HealthCheck() error    { return fmt.Errorf("unhealthy") }
func (f *stoppableFakeHealthChecker) ReportSuccess()        {}
func (f *stoppableFakeHealthChecker) ReportFailure(_ error) {}

// Stop stub tracks whether the service was stopped.
func (f *stoppableFakeHealthChecker) Stop() {
	f.stopped = true
}

var _ types.Stoppable = (*stoppableFakeHealthChecker)(nil)
var _ daemontypes.HealthCheckable = (*stoppableFakeHealthChecker)(nil)

func TestRegisterService_RegistrationFailsAfterStop(t *testing.T) {
	ufm, _ := createTestMonitor()
	ufm.Stop()

	stoppableHc := &stoppableFakeHealthChecker{}
	hc2 := mockFailingHealthCheckerWithError("test-service-2", TestError1)

	// Register unhealthy services. These services are confirmed to trigger a panic if registered when the monitor is
	// not stopped.
	// Register a stoppable unhealthy service.
	err := ufm.RegisterService(stoppableHc, 10*time.Millisecond)
	require.Nil(t, err)

	// Register a non-stoppable unhealthy service.
	err = ufm.RegisterService(hc2, 10*time.Millisecond)
	require.Nil(t, err)

	// Since the max allowable unhealthy duration is 10ms, and the polling period is 10ms, 100ms is long enough to wait
	// in order to trigger a panic if a service is polled.
	time.Sleep(100 * time.Millisecond)

	// Assert that the monitor proactively stops any stoppable service that was registered after the monitor was
	// stopped.
	require.True(t, stoppableHc.stopped)
}

func TestRegisterValidResponseWithCallback_NegativeUnhealthyDuration(t *testing.T) {
	ufm, _ := createTestMonitor()
	hc := mockFailingHealthCheckerWithError("test-service", TestError1)
	err := ufm.RegisterServiceWithCallback(hc, -50*time.Millisecond, func(error) {})
	require.ErrorContains(t, err, "maximum unhealthy duration -50ms must be positive")
}

func TestPanicServiceNotResponding(t *testing.T) {
	panicFunc := types.PanicServiceNotResponding(&mocks.HealthCheckable{})
	require.Panics(t, func() {
		panicFunc(TestError1)
	})
}

func TestLogErrorServiceNotResponding(t *testing.T) {
	logger := &mocks.Logger{}
	hc := &mocks.HealthCheckable{}

	hc.On("ServiceName").Return("test-service")
	logger.On(
		"Error",
		"health-checked service is unhealthy",
		"service",
		"test-service",
		"error",
		TestError1,
	).Return()
	logFunc := types.LogErrorServiceNotResponding(hc, logger)
	logFunc(TestError1)

	// Assert: the logger was called with the expected arguments.
	mock.AssertExpectationsForObjects(t, logger)
}
