package types_test

import (
	"cosmossdk.io/log"
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
			// Set up channel to track errors passed to it from the callback.
			errs := make(chan (error), 100)

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
				func(err error) {
					errs <- err
				},
				timeProvider,
				TestMaximumUnhealthyDuration,
				types.DaemonStartupGracePeriod,
				log.NewNopLogger(),
			)

			// Act - simulate the health checker polling for updates.
			for i := 0; i < len(test.healthCheckResponses); i++ {
				healthChecker.Poll()
			}

			// Assert.
			// Validate the expected polls occurred according to the mocks.
			checkable.AssertExpectations(t)
			timeProvider.AssertExpectations(t)

			var err error
			select {
			case err = <-errs:
			case <-time.After(1 * time.Second):
			}

			// Cleanup.
			healthChecker.Stop()

			// Validate the callback was called with the expected error.
			if test.expectedUnhealthy == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.expectedUnhealthy.Error())
			}
		})
	}
}
