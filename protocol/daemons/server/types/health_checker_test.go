package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
		"unhealthy with multiple errors: most recent error returned": {
			healthCheckResponses: []error{
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError1,
				TestError2,
			},
			expectedUnhealthy: TestError2,
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
				&mocks.Logger{},
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
