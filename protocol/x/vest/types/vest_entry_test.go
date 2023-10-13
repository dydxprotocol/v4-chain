package types_test

import (
	"testing"
	time "time"

	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	now := time.Now().In(time.UTC)
	tests := []struct {
		desc        string
		entry       types.VestEntry
		expectedErr error
	}{
		{
			desc: "valid",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "test_treasury",
				Denom:           "testdenom",
				StartTime:       now.Add(-1 * time.Hour).In(time.UTC),
				EndTime:         now,
			},
			expectedErr: nil,
		},
		{
			desc: "empty vester account",
			entry: types.VestEntry{
				VesterAccount:   "",
				TreasuryAccount: "test_treasury",
				Denom:           "testdenom",
				StartTime:       now.Add(-1 * time.Hour).In(time.UTC),
				EndTime:         now,
			},
			expectedErr: types.ErrInvalidVesterAccount,
		},
		{
			desc: "empty treasury account",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "",
				Denom:           "testdenom",
				StartTime:       now.Add(-1 * time.Hour).In(time.UTC),
				EndTime:         now.In(time.UTC),
			},
			expectedErr: types.ErrInvalidTreasuryAccount,
		},
		{
			desc: "invalid denom",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "test_treasury",
				Denom:           "invalid denom!",
				StartTime:       now.Add(-1 * time.Hour).In(time.UTC),
				EndTime:         now,
			},
			expectedErr: types.ErrInvalidDenom,
		},
		{
			desc: "end_time < start_time",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "test_treasury",
				Denom:           "testdenom",
				StartTime:       now,
				EndTime:         now.Add(-1 * time.Hour).In(time.UTC),
			},
			expectedErr: types.ErrInvalidStartAndEndTimes,
		},
		{
			desc: "end_time = start_time",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "test_treasury",
				Denom:           "testdenom",
				StartTime:       now,
				EndTime:         now,
			},
			expectedErr: types.ErrInvalidStartAndEndTimes,
		},
		{
			desc: "start_time not utc",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "test_treasury",
				Denom:           "testdenom",
				StartTime:       now.Add(-1 * time.Hour).In(time.FixedZone("EST", -5*60*60)),
				EndTime:         now,
			},
			expectedErr: types.ErrInvalidTimeZone,
		},
		{
			desc: "end_time not utc",
			entry: types.VestEntry{
				VesterAccount:   "test_vester",
				TreasuryAccount: "test_treasury",
				Denom:           "testdenom",
				StartTime:       now.Add(-1 * time.Hour),
				EndTime:         now.In(time.FixedZone("EST", -5*60*60)),
			},
			expectedErr: types.ErrInvalidTimeZone,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.entry.Validate()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
