package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
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

	ZeroDuration = 0 * time.Second
)

// createTestMonitor creates a health monitor with a poll frequency of 10ms and a zero duration grace period.
func createTestMonitor() (*types.HealthMonitor, *mocks.Logger) {
	logger := &mocks.Logger{}
	logger.On("With", "module", "health-monitor").Return(logger).Once()
	return types.NewHealthMonitor(
		ZeroDuration,
		10*time.Millisecond,
		logger,
	), logger
}

// mockFailingHealthCheckerWithError creates a mock health checkable service that returns the given error.
func mockFailingHealthCheckerWithError(name string, err error) *mocks.HealthCheckable {
	hc := &mocks.HealthCheckable{}
	hc.On("ServiceName").Return(name)
	hc.On("HealthCheck").Return(err)
	return hc
}

// callbackWithErrorPointer returns a callback function and an error pointer that tracks the error passed to the
// callback. This can be used to validate that a service was considered unhealthy for the maximum allowable duration.
func callbackWithErrorPointer() (func(error), *error) {
	var callbackError error
	callback := func(err error) {
		callbackError = err
	}
	return callback, &callbackError
}

func TestHealthChecker(t *testing.T) {
	tests := map[string]struct {
		healthCheckResponses []error
		expectedUnhealthy    error
	}{
		"initialized: no callback triggered": {
			healthCheckResponses: []error{},
			expectedUnhealthy:    nil,
		},
		"healthy, then unhealthy for < maximum unhealthy duration: no callback triggered": {
			healthCheckResponses: []error{
				nil,
				TestError1,
			},
			expectedUnhealthy: nil,
		},
		"unhealthy for < maximum unhealthy duration: no callback triggered": {
			healthCheckResponses: []error{
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
			},
			expectedUnhealthy: nil,
		},
		"unhealthy, healthy, unhealthy: no callback triggered": {
			healthCheckResponses: []error{
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				nil,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
			},
			expectedUnhealthy: nil,
		},
		"unhealthy for maximum unhealthy duration: callback triggered": {
			healthCheckResponses: []error{
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
			},
			expectedUnhealthy: TestError1,
		},
		"unhealthy with multiple errors: first error returned": {
			healthCheckResponses: []error{
				TestError1,
				TestError2,
				TestError2,
				TestError2,
				TestError2,
				TestError2,
			},
			expectedUnhealthy: TestError1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			// Set up callback to track error passed to it.
			callback, callbackError := callbackWithErrorPointer()

			// Set up health checkable service.
			checkable := &mocks.HealthCheckable{}
			for _, response := range test.healthCheckResponses {
				checkable.On("HealthCheck").Return(response).Once()
			}

			// Set up time provider to return a sequence of timestamps one second apart starting at Time0.
			timeProvider := &mocks.TimeProvider{}
			for i := range test.healthCheckResponses {
				timeProvider.On("Now").Return(Time0.Add(time.Duration(i) * time.Second)).Once()
			}

			healthChecker := types.StartNewHealthChecker(
				checkable,
				TestLargeDuration, // set to a >> value so that poll is never auto-triggered by the timer
				callback,
				timeProvider,
				TestMaximumUnhealthyDuration,
				types.DaemonStartupGracePeriod,
			)

			// Cleanup.
			defer func() {
				healthChecker.Stop()
			}()

			// Act - simulate the health checker polling for updates.
			for i := 0; i < len(test.healthCheckResponses); i++ {
				healthChecker.Poll()
			}

			// Assert.
			// Validate the expected polls occurred according to the mocks.
			checkable.AssertExpectations(t)
			timeProvider.AssertExpectations(t)

			// Validate the callback was called with the expected error.
			if test.expectedUnhealthy == nil {
				require.NoError(t, *callbackError)
			} else {
				require.ErrorContains(t, *callbackError, test.expectedUnhealthy.Error())
			}
		})
	}
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
			callback, callbackError := callbackWithErrorPointer()

			// Act.
			err := ufm.RegisterServiceWithCallback(
				hc,
				50*time.Millisecond,
				callback,
			)
			require.NoError(t, err)

			// Cleanup.
			defer func() {
				ufm.Stop()
			}()

			// Give the monitor time to poll the health checkable service. Polls occur once every 10ms.
			time.Sleep(100 * time.Millisecond)

			// Assert: no calls to the logger were made.
			mock.AssertExpectationsForObjects(t, hc, logger)

			// Assert: the callback was called or not called as expected.
			require.Equal(t, test.healthCheckResponse, *callbackError)
		})
	}
}

func TestRegisterServiceWithCallback_DoubleRegistrationFails(t *testing.T) {
	// Setup.
	ufm, logger := createTestMonitor()

	hc := mockFailingHealthCheckerWithError("test-service", TestError1)
	hc2 := mockFailingHealthCheckerWithError("test-service", TestError2)

	callback, callbackError := callbackWithErrorPointer()

	err := ufm.RegisterServiceWithCallback(hc, 50*time.Millisecond, callback)
	require.NoError(t, err)

	// Register a service with the same name. This registration should fail.
	err = ufm.RegisterServiceWithCallback(hc2, 50*time.Millisecond, callback)
	require.ErrorContains(t, err, "service already registered")

	// Expect that the first service is still operating and will produce a callback after a sustained unhealthy period.
	time.Sleep(100 * time.Millisecond)
	ufm.Stop()

	// Assert no calls to the logger were made.
	mock.AssertExpectationsForObjects(t, logger, hc)
	hc2.AssertNotCalled(t, "HealthCheck")

	// Assert the callback was called with the expected error.
	require.Equal(t, TestError1, *callbackError)
}

func TestRegisterService_RegistrationFailsAfterStop(t *testing.T) {
	ufm, logger := createTestMonitor()
	ufm.Stop()

	hc := mockFailingHealthCheckerWithError("test-service", TestError1)
	err := ufm.RegisterService(hc, 50*time.Millisecond)
	require.ErrorContains(t, err, "monitor has been stopped")

	// Any scheduled functions with error logs that were not cleaned up should trigger before this sleep finishes.
	time.Sleep(100 * time.Millisecond)
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterValidResponseWithCallback_NegativeUpdateDuration(t *testing.T) {
	ufm, _ := createTestMonitor()
	hc := mockFailingHealthCheckerWithError("test-service", TestError1)
	err := ufm.RegisterServiceWithCallback(hc, -50*time.Millisecond, func(error) {})
	require.ErrorContains(t, err, "maximum acceptable unhealthy duration -50ms must be positive")
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
