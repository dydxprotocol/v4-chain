package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestPerMarketFeeDiscountParams_Validate(t *testing.T) {
	// Set a fixed current time for testing
	currentTime := time.Unix(1000, 0).UTC()

	tests := []struct {
		name    string
		params  types.PerMarketFeeDiscountParams
		wantErr error
	}{
		{
			name: "valid params",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1200, 0).UTC(),
				ChargePpm:  500_000, // 50% discount
			},
			wantErr: nil,
		},
		{
			name: "valid params - start time in past but end time in future",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(900, 0).UTC(),
				EndTime:    time.Unix(1100, 0).UTC(),
				ChargePpm:  500_000, // 50% discount
			},
			wantErr: nil,
		},
		{
			name: "valid params - zero charge (100% discount)",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1200, 0).UTC(),
				ChargePpm:  0, // 100% discount (free)
			},
			wantErr: nil,
		},
		{
			name: "valid params - max charge (no discount)",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1200, 0).UTC(),
				ChargePpm:  types.MaxChargePpm, // 100% charge (no discount)
			},
			wantErr: nil,
		},
		{
			name: "valid params - maximum duration",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1100, 0).Add(time.Duration(types.MaxFeeDiscountDuration) * time.Second).UTC(),
				ChargePpm:  500_000, // 50% discount
			},
			wantErr: nil,
		},
		{
			name: "invalid params - start time equals end time",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1100, 0).UTC(), // Same as start time
				ChargePpm:  500_000,
			},
			wantErr: types.ErrInvalidTimeRange,
		},
		{
			name: "invalid params - start time after end time",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1200, 0).UTC(),
				EndTime:    time.Unix(1100, 0).UTC(), // Before start time
				ChargePpm:  500_000,
			},
			wantErr: types.ErrInvalidTimeRange,
		},
		{
			name: "invalid params - end time in past",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(900, 0).UTC(),
				EndTime:    time.Unix(950, 0).UTC(), // Before current time (1000)
				ChargePpm:  500_000,
			},
			wantErr: types.ErrInvalidTimeRange,
		},
		{
			name: "invalid params - end time equals current time",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(900, 0).UTC(),
				EndTime:    time.Unix(1000, 0).UTC(), // Equal to current time
				ChargePpm:  500_000,
			},
			wantErr: types.ErrInvalidTimeRange,
		},
		{
			name: "invalid params - duration exceeds maximum",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1100, 0).Add(time.Duration(types.MaxFeeDiscountDuration+1) * time.Second).UTC(),
				ChargePpm:  500_000,
			},
			wantErr: types.ErrInvalidTimeRange,
		},
		{
			name: "invalid params - charge PPM exceeds maximum",
			params: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0).UTC(),
				EndTime:    time.Unix(1200, 0).UTC(),
				ChargePpm:  types.MaxChargePpm + 1, // Exceeds maximum charge PPM
			},
			wantErr: types.ErrInvalidChargePpm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate(currentTime)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test validations across different current times
func TestPerMarketFeeDiscountParams_ValidateWithDifferentTimes(t *testing.T) {
	// Define a fixed discount params
	discountParams := types.PerMarketFeeDiscountParams{
		ClobPairId: 1,
		StartTime:  time.Unix(1100, 0),
		EndTime:    time.Unix(1200, 0),
		ChargePpm:  500_000, // 50% discount
	}

	tests := []struct {
		name        string
		currentTime time.Time
		wantErr     error
	}{
		{
			name:        "current time before start time",
			currentTime: time.Unix(1050, 0),
			wantErr:     nil,
		},
		{
			name:        "current time at start time",
			currentTime: time.Unix(1100, 0),
			wantErr:     nil,
		},
		{
			name:        "current time between start and end",
			currentTime: time.Unix(1150, 0),
			wantErr:     nil,
		},
		{
			name:        "current time at end time",
			currentTime: time.Unix(1200, 0),
			wantErr:     types.ErrInvalidTimeRange,
		},
		{
			name:        "current time after end time",
			currentTime: time.Unix(1250, 0),
			wantErr:     types.ErrInvalidTimeRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := discountParams.Validate(tt.currentTime)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test for edge cases around the MaxChargePpm constant
func TestPerMarketFeeDiscountParams_ChargePpmEdgeCases(t *testing.T) {
	currentTime := time.Unix(1000, 0).UTC()

	tests := []struct {
		name      string
		chargePpm uint32
		wantErr   error
	}{
		{
			name:      "minimum charge (0)",
			chargePpm: 0,
			wantErr:   nil,
		},
		{
			name:      "mid-range charge (500,000)",
			chargePpm: 500_000,
			wantErr:   nil,
		},
		{
			name:      "maximum charge (1,000,000)",
			chargePpm: types.MaxChargePpm,
			wantErr:   nil,
		},
		{
			name:      "charge just over maximum (1,000,001)",
			chargePpm: types.MaxChargePpm + 1,
			wantErr:   types.ErrInvalidChargePpm,
		},
		{
			name:      "large charge (2,000,000)",
			chargePpm: 2_000_000,
			wantErr:   types.ErrInvalidChargePpm,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1100, 0),
				EndTime:    time.Unix(1200, 0),
				ChargePpm:  tt.chargePpm,
			}
			err := params.Validate(currentTime)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test for edge cases around the MaxFeeDiscountDuration constant
func TestPerMarketFeeDiscountParams_DurationEdgeCases(t *testing.T) {
	currentTime := time.Unix(1000, 0).UTC()

	tests := []struct {
		name     string
		duration time.Duration
		wantErr  error
	}{
		{
			name:     "minimum duration (1 second)",
			duration: 1 * time.Second,
			wantErr:  nil,
		},
		{
			name:     "1 day duration",
			duration: 24 * time.Hour,
			wantErr:  nil,
		},
		{
			name:     "30 days duration",
			duration: 30 * 24 * time.Hour,
			wantErr:  nil,
		},
		{
			name:     "maximum duration (90 days)",
			duration: time.Duration(types.MaxFeeDiscountDuration) * time.Second,
			wantErr:  nil,
		},
		{
			name:     "duration just over maximum (90 days + 1 second)",
			duration: time.Duration(types.MaxFeeDiscountDuration+1) * time.Second,
			wantErr:  types.ErrInvalidTimeRange,
		},
		{
			name:     "large duration (180 days)",
			duration: 180 * 24 * time.Hour,
			wantErr:  types.ErrInvalidTimeRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Unix(1100, 0).UTC().UTC()
			params := types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  startTime,
				EndTime:    startTime.Add(tt.duration),
				ChargePpm:  500_000,
			}
			err := params.Validate(currentTime)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
