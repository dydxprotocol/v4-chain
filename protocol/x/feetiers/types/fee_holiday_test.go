package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestFeeHolidayParams_Validate(t *testing.T) {
	// Set a fixed current time for testing
	currentTime := time.Unix(1000, 0)

	tests := map[string]struct {
		params *types.FeeHolidayParams
		err    error
	}{
		"valid fee holiday params": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1200,
			},
			err: nil,
		},
		"start time equal to end time is invalid": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1100,
			},
			err: types.ErrInvalidTimeRange,
		},
		"start time after end time is invalid": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1200,
				EndTimeUnix:   1100,
			},
			err: types.ErrInvalidTimeRange,
		},
		"end time before current time is invalid": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 900,
				EndTimeUnix:   950,
			},
			err: types.ErrInvalidTimeRange,
		},
		"end time equal to current time is invalid": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 900,
				EndTimeUnix:   1000,
			},
			err: types.ErrInvalidTimeRange,
		},
		"duration exceeding 30 days is invalid": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1100 + 31*24*60*60, // 31 days
			},
			err: types.ErrInvalidTimeRange,
		},
		"start time in the past is valid as long as end time is in the future": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 900,  // Before current time
				EndTimeUnix:   1100, // After current time
			},
			err: nil,
		},
		"max allowed duration (exactly 30 days) is valid": {
			params: &types.FeeHolidayParams{
				ClobPairId:    1,
				StartTimeUnix: 1100,
				EndTimeUnix:   1100 + 30*24*60*60, // 30 days
			},
			err: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate(currentTime)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}

// Test that validates the FeeHolidayParams across different current times
func TestFeeHolidayParams_ValidateWithDifferentCurrentTimes(t *testing.T) {
	// Define a valid fee holiday
	feeHoliday := &types.FeeHolidayParams{
		ClobPairId:    1,
		StartTimeUnix: 1100,
		EndTimeUnix:   1200,
	}

	tests := map[string]struct {
		currentTimeUnix int64
		expectError     bool
	}{
		"current time before start time": {
			currentTimeUnix: 1000,
			expectError:     false,
		},
		"current time equal to start time": {
			currentTimeUnix: 1100,
			expectError:     false,
		},
		"current time between start and end": {
			currentTimeUnix: 1150,
			expectError:     false,
		},
		"current time equal to end time": {
			currentTimeUnix: 1200,
			expectError:     true,
		},
		"current time after end time": {
			currentTimeUnix: 1300,
			expectError:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			currentTime := time.Unix(tc.currentTimeUnix, 0)
			err := feeHoliday.Validate(currentTime)

			if tc.expectError {
				require.Error(t, err)
				require.ErrorIs(t, err, types.ErrInvalidTimeRange)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
