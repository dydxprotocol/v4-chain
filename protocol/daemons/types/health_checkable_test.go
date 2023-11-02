package types_test

import (
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

	TestError = fmt.Errorf("test error")
)

func TestHealthCheckableImpl_PanicsWithoutInitilization(t *testing.T) {
	hc := types.HealthCheckableImpl{}
	require.Panics(
		t,
		func() {
			hc.HealthCheck(&libtime.TimeProviderImpl{}) // nolint:errcheck
		},
		"HealthCheckableImpl.HealthCheck should panic if not initialized",
	)
}

// singleUseTimeProvider returns a TimeProvider that returns the given time on the first call to Now.
func singleUseTimeProvider(time time.Time) libtime.TimeProvider {
	m := mocks.TimeProvider{}
	m.On("Now").Return(time).Once()
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
		"unhealthy: no updates": {
			healthCheckTime:      Time1,
			expectedHealthStatus: fmt.Errorf("no successful update has occurred"),
		},
		"unhealthy: no successful updates": {
			updates: []struct {
				timestamp time.Time
				err       error
			}{
				{Time1, TestError}, // failed update
			},
			healthCheckTime:      Time2,
			expectedHealthStatus: fmt.Errorf("no successful update has occurred"),
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
				"last update failed at %v with error: %w",
				Time2,
				TestError,
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
		"unhealthy: last successful update was more than 5 minutes ago": {
			updates: []struct {
				timestamp time.Time
				err       error
			}{
				{Time1, nil}, // successful update
			},
			healthCheckTime: Time_5Minutes_And_2Seconds,
			expectedHealthStatus: fmt.Errorf(
				"last successful update occurred at %v, which is more than %v ago",
				Time1,
				types.MaximumAcceptableUpdateDelay,
			),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hci := types.NewHealthCheckableImpl("test", singleUseTimeProvider(Time0))
			for _, update := range tc.updates {
				timeProvider := singleUseTimeProvider(update.timestamp)
				if update.err == nil {
					hci.RecordUpdateSuccess(timeProvider)
				} else {
					hci.RecordUpdateFailure(timeProvider, update.err)
				}
			}

			err := hci.HealthCheck(singleUseTimeProvider(tc.healthCheckTime))
			if tc.expectedHealthStatus == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedHealthStatus.Error())
			}
		})
	}
}
