package price_encoder

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func newMarketIdWithValue(id types.MarketId) *types.MarketId {
	ptr := new(types.MarketId)
	*ptr = id
	return ptr
}

func TestGetPriceConversionDetailsForMarket(t *testing.T) {
	tests := map[string]struct {
		mutableExchangeConfig *types.MutableExchangeMarketConfig
		marketToMutableConfig map[types.MarketId]*types.MutableMarketConfig
		expected              priceConversionDetails
		expectedErr           error
	}{
		"Error: Market config not found": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id:                   "exchange1",
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{},
			},
			marketToMutableConfig: map[types.MarketId]*types.MutableMarketConfig{},
			expectedErr:           fmt.Errorf("market config for market 1 not found on exchange 'exchange1'"),
		},
		"Error: Mutable market config for adjustment market not found": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: "exchange1",
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "PAIR-USD",
						AdjustByMarket: newMarketIdWithValue(2),
					},
				},
			},
			marketToMutableConfig: map[types.MarketId]*types.MutableMarketConfig{},
			expectedErr: fmt.Errorf(
				"mutable market config for adjust-by market 2 not found on exchange 'exchange1'",
			),
		},
		"Error: Mutable market config for market not found": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: "exchange1",
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "PAIR-USD",
					},
				},
			},
			marketToMutableConfig: map[types.MarketId]*types.MutableMarketConfig{},
			expectedErr:           fmt.Errorf("mutable market config for market 1 not found on exchange 'exchange1'"),
		},
		"Success: no adjustment market": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: "exchange1",
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "PAIR-USD",
					},
				},
			},
			marketToMutableConfig: map[types.MarketId]*types.MutableMarketConfig{
				1: {
					Id:           1,
					Pair:         "PAIR-USD",
					Exponent:     -5,
					MinExchanges: 2,
				},
			},
			expected: priceConversionDetails{
				Exponent: -5,
			},
		},
		"Success: with adjustment market": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: "exchange1",
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "PAIR-USD",
						AdjustByMarket: newMarketIdWithValue(2),
						Invert:         true,
					},
				},
			},
			marketToMutableConfig: map[types.MarketId]*types.MutableMarketConfig{
				1: {
					Id:           1,
					Pair:         "PAIR-USD",
					Exponent:     -5,
					MinExchanges: 2,
				},
				2: {
					Id:           2,
					Pair:         "ADJ-USD",
					Exponent:     -6,
					MinExchanges: 3,
				},
			},
			expected: priceConversionDetails{
				Invert:   true,
				Exponent: -5,
				AdjustByMarketDetails: &adjustByMarketDetails{
					MarketId:     2,
					Exponent:     -6,
					MinExchanges: 3,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mutableState := &mutableState{
				mutableExchangeConfig: tc.mutableExchangeConfig,
				marketToMutableConfig: tc.marketToMutableConfig,
			}
			actual, actualErr := mutableState.GetPriceConversionDetailsForMarket(1)
			if tc.expectedErr != nil {
				require.Error(t, tc.expectedErr, actualErr)
				require.Zero(t, actual)
			} else {
				require.NoError(t, actualErr)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}
