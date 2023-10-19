package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/stretchr/testify/require"
)

func TestDowntimeParams_Validate(t *testing.T) {
	tests := map[string]struct {
		params        *types.DowntimeParams
		expectedError error
	}{
		"valid": {
			params: &types.DowntimeParams{
				Durations: []time.Duration{
					time.Minute,
					time.Minute * 2,
					time.Minute * 3,
				},
			},
			expectedError: nil,
		},
		"not ascending": {
			params: &types.DowntimeParams{
				Durations: []time.Duration{
					time.Minute,
					time.Minute * 3,
					time.Minute * 2,
				},
			},
			expectedError: types.ErrUnorderedDurations,
		},
		"invalid duration": {
			params: &types.DowntimeParams{
				Durations: []time.Duration{
					time.Minute,
					time.Minute * 2,
					time.Minute * 0,
				},
			},
			expectedError: types.ErrNonpositiveDuration,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, tc.expectedError, err)
			}
		})
	}
}
