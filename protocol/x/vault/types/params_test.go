package types_test

import (
	"math"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestValidateQuotingParams(t *testing.T) {
	tests := map[string]struct {
		// Params to validate.
		params types.QuotingParams
		// Expected error
		expectedErr error
	}{
		"Success": {
			params:      types.DefaultQuotingParams(),
			expectedErr: nil,
		},
		"Failure - Layer is greater than MaxUint8": {
			params: types.QuotingParams{
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
			params: types.QuotingParams{
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
			params: types.QuotingParams{
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
			params: types.QuotingParams{
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
			params: types.QuotingParams{
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
		"Failure - SkewFactorPpm is high. Product of SkewFactorPpm and OrderSizePctPpm is above threshold.": {
			params: types.QuotingParams{
				Layers:                           2,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    20_000_000,
				OrderSizePctPpm:                  100_000,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidSkewFactor,
		},
		"Failure - OrderSizePctPpm is high. Product of SkewFactorPpm and OrderSizePctPpm is above threshold.": {
			params: types.QuotingParams{
				Layers:                           2,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    2_000_000,
				OrderSizePctPpm:                  1_000_000,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidSkewFactor,
		},
		"Failure - Both OrderSizePctPpm and SkewFactorPpm are MaxUint32. The product is above the threshold.": {
			params: types.QuotingParams{
				Layers:                           2,
				SpreadMinPpm:                     3_000,
				SpreadBufferPpm:                  1_500,
				SkewFactorPpm:                    math.MaxUint32,
				OrderSizePctPpm:                  math.MaxUint32,
				OrderExpirationSeconds:           5,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: types.ErrInvalidSkewFactor,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateVaultParams(t *testing.T) {
	tests := map[string]struct {
		// Params to validate.
		params types.VaultParams
		// Expected error
		expectedErr error
	}{
		"Success": {
			params: types.VaultParams{
				Status:        types.VaultStatus_VAULT_STATUS_QUOTING,
				QuotingParams: &constants.QuotingParams,
			},
			expectedErr: nil,
		},
		"Success - nil quoting params": {
			params: types.VaultParams{
				Status: types.VaultStatus_VAULT_STATUS_QUOTING,
			},
			expectedErr: nil,
		},
		"Failure - unspecified vault status": {
			params: types.VaultParams{
				QuotingParams: &constants.QuotingParams,
			},
			expectedErr: types.ErrUnspecifiedVaultStatus,
		},
		"Failure - invalid quoting params": {
			params: types.VaultParams{
				Status:        types.VaultStatus_VAULT_STATUS_QUOTING,
				QuotingParams: &constants.InvalidQuotingParams,
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

func TestValidateOperatorParams(t *testing.T) {
	tests := map[string]struct {
		// Params to validate.
		params types.OperatorParams
		// Expected error
		expectedErr error
	}{
		"Success: Alice as operator": {
			params: types.OperatorParams{
				Operator: constants.AliceAccAddress.String(),
			},
			expectedErr: nil,
		},
		"Success: Gov module account as operator": {
			params: types.OperatorParams{
				Operator: constants.GovAuthority,
			},
			expectedErr: nil,
		},
		"Failure - empty operator": {
			params: types.OperatorParams{
				Operator: "",
			},
			expectedErr: types.ErrEmptyOperator,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
