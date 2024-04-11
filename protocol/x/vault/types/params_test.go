package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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
		"Failure - Layer is greater than MaxUint8": {
			params: types.Params{
				Layers:                           256,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    500_000,
				OrderSizePctPpm:                  100_000,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidLayers,
		},
		"Failure - SpreadMinPpm is 0": {
			params: types.Params{
				Layers:                           2,
				SpreadMinPpm:                     0,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    500_000,
				OrderSizePctPpm:                  100_000,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidSpreadMinPpm,
		},
		"Failure - OrderSizePctPpm is 0": {
			params: types.Params{
				Layers:                           2,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    500_000,
				OrderSizePctPpm:                  0,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidOrderSizePctPpm,
		},
		"Failure - OrderExpirationSeconds is 0": {
			params: types.Params{
				Layers:                           2,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    500_000,
				OrderSizePctPpm:                  100_000,
				OrderExpirationSeconds:           0,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidOrderExpirationSeconds,
		},
		"Failure - ActivationThresholdQuoteQuantums is negative": {
			params: types.Params{
				Layers:                           2,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    500_000,
				OrderSizePctPpm:                  100_000,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(-1),
			},
			expectedErr: types.ErrInvalidActivationThresholdQuoteQuantums,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
