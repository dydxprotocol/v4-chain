package types_test

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	Time0                      = time.Unix(0, 0)
	Time1                      = Time0.Add(time.Second)
	Time2                      = Time0.Add(2 * time.Second)
	Time3                      = Time0.Add(3 * time.Second)
	Time4                      = Time0.Add(4 * time.Second)
	Time_5Minutes_And_2Seconds = Time0.Add(5*time.Minute + 2*time.Second)

	TestError          = fmt.Errorf("test error")
	InitializingStatus = fmt.Errorf("test is initializing")
)

// mockTimeProviderWithTimestamps returns a TimeProvider that returns the given timestamps in order.
func mockTimeProviderWithTimestamps(times []time.Time) libtime.TimeProvider {
	m := mocks.TimeProvider{}
	for _, timestamp := range times {
		m.On("Now").Return(timestamp).Once()
	}
	return &m
}

func TestHealthCheckableImpl_Mixed(t *testing.T) {
	tests := map[string]struct {
		updates []struct {
			timestamp time.Time
			// leave error nil for a successful update
			err error
		}
		healthCheckTime      time.Time
		expectedHealthStatus error
	}{
		"unhealthy: no updates, returns initializing error": {
			healthCheckTime: Time1,
			expectedHealthStatus: fmt.Errorf(
				"no successful update has occurred; last failed update occurred at %v with error '%w'",
				Time0,
				InitializingStatus,
			),
		},
		"unhealthy: no successful updates": {
			updates: []struct {
				timestamp time.Time
				err       error
			}{
				{Time1, TestError}, // failed update
			},
			healthCheckTime: Time2,
			expectedHealthStatus: fmt.Errorf(
				"no successful update has occurred; last failed update occurred at %v with error '%w'",
				Time1,
				TestError,
			),
		},
		"healthy: one recent successful update": {
			updates: []struct {
				timestamp time.Time
				err       error
			}{
				{Time1, nil}, // successful update
			},
			healthCheckTime:      Time2,
			expectedHealthStatus: nil, // expect healthy
		},
		"unhealthy: one recent successful update, followed by a failed update": {
			updates: []struct {
				timestamp time.Time
				err       error
			}{
				{Time1, nil},       // successful update
				{Time2, TestError}, // failed update
			},
			healthCheckTime: Time3,
			expectedHealthStatus: fmt.Errorf(
				"last update failed at %v with error: '%w', most recent successful update occurred at %v",
				Time2,
				TestError,
				Time1,
			),
		},
		"healthy: one recent failed update followed by a successful update": {
			updates: []struct {
				timestamp time.Time
				err       error
			}{
				{Time1, TestError}, // failed update
				{Time2, nil},       // successful update
			},
			healthCheckTime:      Time3,
			expectedHealthStatus: nil, // expect healthy
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.

			// Construct list of timestamps to provide for the timeProvider stored within the health checkable instance.
			timestamps := make([]time.Time, 0, len(tc.updates)+2)

			// The first timestamp is used during HealthCheckable initialization.
			timestamps = append(timestamps, Time0)
			// One timestamp used for each update.
			for _, update := range tc.updates {
				timestamps = append(timestamps, update.timestamp)
			}
			// A final timestamp is consumed by the HealthCheck call.
			timestamps = append(timestamps, tc.healthCheckTime)

			// Create a new time-bounded health checkable instance.
			hci := types.NewTimeBoundedHealthCheckable(
				"test",
				mockTimeProviderWithTimestamps(timestamps),
				log.NewNopLogger(),
			)

			// Act.
			// Report the test sequence of successful / failed updates.
			for _, update := range tc.updates {
				if update.err == nil {
					hci.ReportSuccess()
				} else {
					hci.ReportFailure(update.err)
				}
			}

			// Assert.
			// Check the health status after all updates have been reported.
			err := hci.HealthCheck()
			if tc.expectedHealthStatus == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedHealthStatus.Error())
			}
		})
	}
}
