package types_test

import (
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestPerpetual_GetYieldIndexAsRat(t *testing.T) {
	tests := []struct {
		desc        string
		perpetual   *types.Perpetual
		expectedRat *big.Rat
		expectedErr error
	}{
		{
			desc: "Valid yield index",
			perpetual: &types.Perpetual{
				YieldIndex: "0.05",
			},
			expectedRat: big.NewRat(5, 100),
			expectedErr: nil,
		},
		{
			desc:        "Nil perpetual",
			perpetual:   nil,
			expectedRat: nil,
			expectedErr: types.ErrPerpIsNil,
		},
		{
			desc: "Empty yield index",
			perpetual: &types.Perpetual{
				YieldIndex: "",
			},
			expectedRat: nil,
			expectedErr: types.ErrYieldIndexDoesNotExist,
		},
		{
			desc: "Invalid yield index format",
			perpetual: &types.Perpetual{
				YieldIndex: "not_a_number",
			},
			expectedRat: nil,
			expectedErr: types.ErrRatToStringConversion,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			rat, err := tc.perpetual.GetYieldIndexAsRat()

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Nil(t, rat)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rat)
				require.Equal(t, 0, tc.expectedRat.Cmp(rat), "Expected %v, got %v", tc.expectedRat, rat)
			}
		})
	}
}

func TestPerpetualParams_Validate(t *testing.T) {
	tests := []struct {
		desc        string
		params      types.PerpetualParams
		expectedErr string
	}{
		{
			desc: "Valid param",
			params: types.PerpetualParams{
				Ticker:            "test",
				DefaultFundingPpm: 1_000_000,
				MarketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			},
			expectedErr: "",
		},
		{
			desc: "Empty ticker",
			params: types.PerpetualParams{
				Ticker:            "",
				DefaultFundingPpm: 1_000_000,
				MarketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			},
			expectedErr: "Ticker must be non-empty string",
		},
		{
			desc: "Invalid DefaultFundingPpm",
			params: types.PerpetualParams{
				Ticker:            "test",
				DefaultFundingPpm: 100_000_000,
				MarketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			},
			expectedErr: "DefaultFundingPpm magnitude exceeds maximum value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
