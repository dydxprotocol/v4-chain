package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		// Params to validate.
		params types.Params
		// Expected error
		expectedErr error
	}{
		"Success": {
			params:      types.DefaultParams(),
			expectedErr: nil,
		},
		"Failure - SkewFactorPpm is 0": {
			params: types.Params{
				Layers:                 2,
				SpreadMinPpm:           3_000,
				SpreadBufferPpm:        1_500,
				SkewFactorPpm:          0,
				OrderSizePpm:           100_000,
				OrderExpirationSeconds: 5,
			},
			expectedErr: types.ErrInvalidSkewFactorPpm,
		},
		"Failure - SkewFactorPpm is 1": {
			params: types.Params{
				Layers:                 2,
				SpreadMinPpm:           3_000,
				SpreadBufferPpm:        1_500,
				SkewFactorPpm:          1_000_000,
				OrderSizePpm:           100_000,
				OrderExpirationSeconds: 5,
			},
			expectedErr: types.ErrInvalidSkewFactorPpm,
		},
		"Failure - OrderSizePpm is 0": {
			params: types.Params{
				Layers:                 2,
				SpreadMinPpm:           3_000,
				SpreadBufferPpm:        1_500,
				SkewFactorPpm:          500_000,
				OrderSizePpm:           0,
				OrderExpirationSeconds: 5,
			},
			expectedErr: types.ErrInvalidOrderSizePpm,
		},
		"Failure - OrderExpirationSeconds is 0": {
			params: types.Params{
				Layers:                 2,
				SpreadMinPpm:           3_000,
				SpreadBufferPpm:        1_500,
				SkewFactorPpm:          500_000,
				OrderSizePpm:           100_000,
				OrderExpirationSeconds: 0,
			},
			expectedErr: types.ErrInvalidOrderExpirationSeconds,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
