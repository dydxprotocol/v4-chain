package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func newMarketIdWithValue(id types.MarketId) *types.MarketId {
	ptr := new(types.MarketId)
	*ptr = id
	return ptr
}

func TestEqual_Mixed(t *testing.T) {
	tests := map[string]struct {
		A, B          types.MarketConfig
		expectedEqual bool
	}{
		"Equal: adjustBy market defined": {
			A: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			B: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			expectedEqual: true,
		},
		"Equal: adjustBy market nil": {
			A: types.MarketConfig{
				Ticker: "ABC-USD",
				Invert: true,
			},
			B: types.MarketConfig{
				Ticker: "ABC-USD",
				Invert: true,
			},
			expectedEqual: true,
		},
		"Not equal: adjustBy markets differ": {
			A: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			B: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(3),
				Invert:         true,
			},
			expectedEqual: false,
		},
		"Not equal: tickers differ": {
			A: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			B: types.MarketConfig{
				Ticker:         "DEF-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			expectedEqual: false,
		},
		"Not equal: invert differs": {
			A: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			B: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         false,
			},
			expectedEqual: false,
		},
		"Not equal: adjustBy markets nil/non-nil": {
			A: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
			B: types.MarketConfig{
				Ticker: "ABC-USD",
				Invert: true,
			},
			expectedEqual: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			equal := tc.A.Equal(tc.B)
			require.Equal(t, tc.expectedEqual, equal)
		})
	}
}

func TestMarketConfig_Copy(t *testing.T) {
	tests := map[string]struct {
		config types.MarketConfig
	}{
		"Copy: adjustBy market defined": {
			config: types.MarketConfig{
				Ticker:         "ABC-USD",
				AdjustByMarket: newMarketIdWithValue(2),
				Invert:         true,
			},
		},
		"Copy: adjustBy market nil": {
			config: types.MarketConfig{
				Ticker: "ABC-USD",
				Invert: true,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			copied := tc.config.Copy()
			require.True(t, tc.config.Equal(copied))
		})
	}
}
