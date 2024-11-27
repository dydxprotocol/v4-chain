package types_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	"github.com/stretchr/testify/require"
)

const (
	ExchangeId1 = "Exchange1"
	ExchangeId2 = "Exchange2"
)

func newUint32WithValue(val uint32) *uint32 {
	ptr := new(uint32)
	*ptr = val
	return ptr
}

func TestMutableExchangeMarketConfig_Copy(t *testing.T) {
	mutableMarketExchangeConfig := &types.MutableExchangeMarketConfig{
		Id: ExchangeId1,
		MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
			exchange_config.MARKET_ETH_USD: {
				Ticker: "ETHUSD",
				Invert: false,
			},
			exchange_config.MARKET_BTC_USD: {
				Ticker:         "BTCUSD",
				AdjustByMarket: newUint32WithValue(exchange_config.MARKET_ETH_USD),
				Invert:         true,
			},
		},
	}
	mmecCopy := mutableMarketExchangeConfig.Copy()
	require.NotSame(t, mutableMarketExchangeConfig, mmecCopy)
	require.True(t, mutableMarketExchangeConfig.Equal(mmecCopy))
}

func TestGetMarketIds_Success(t *testing.T) {
	tests := map[string]struct {
		marketToConfig  map[types.MarketId]types.MarketConfig
		expectedMarkets []types.MarketId
	}{
		"Empty map": {
			marketToConfig:  map[types.MarketId]types.MarketConfig{},
			expectedMarkets: []types.MarketId{},
		},
		"One market": {
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_config.MARKET_ETH_USD: {
					Ticker: "ETHUSD",
				},
			},
			expectedMarkets: []types.MarketId{
				exchange_config.MARKET_ETH_USD,
			},
		},
		"Multiple markets": {
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_config.MARKET_ETH_USD: {
					Ticker: "ETHUSD",
				},
				exchange_config.MARKET_BTC_USD: {
					Ticker: "BTCUSD",
				},
			},
			expectedMarkets: []types.MarketId{
				exchange_config.MARKET_BTC_USD,
				exchange_config.MARKET_ETH_USD,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memc := types.MutableExchangeMarketConfig{
				Id:                   ExchangeId1,
				MarketToMarketConfig: tc.marketToConfig,
			}
			actualMarkets := memc.GetMarketIds()
			require.ElementsMatch(t, tc.expectedMarkets, actualMarkets)
		})
	}
}

func TestMutableExchangeMarketConfig_Equal(t *testing.T) {
	tests := map[string]struct {
		A, B          *types.MutableExchangeMarketConfig
		expectedEqual bool
	}{
		"False: non-matching IDs": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId2,
			},
			expectedEqual: false,
		},
		"False: non-matching tickers": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "BTCUSD", // Non-matching ticker.
					},
				},
			},
			expectedEqual: false,
		},
		"False: non-matching markets": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					2: { // Non-matching market id.
						Ticker: "ETHUSD",
					},
				},
			},
			expectedEqual: false,
		},
		"False: non-matching adjust-by markets": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(3),
					},
				},
			},
			expectedEqual: false,
		},
		"False: adjust-by market defined on one config": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(3),
					},
				},
			},
			expectedEqual: false,
		},
		"False: non-matching inversions": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
						Invert:         true,
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
						Invert:         false,
					},
				},
			},
			expectedEqual: false,
		},
		"True: populated adjust-by markets": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
						Invert:         true,
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
						Invert:         true,
					},
				},
			},
			expectedEqual: true,
		},
		"True: nil adjust-by markets": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			expectedEqual: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedEqual, tc.A.Equal(tc.B))
		})
	}
}

func TestMutableExchangeMarketConfig_Validate(t *testing.T) {
	tests := map[string]struct {
		mutableExchangeConfig *types.MutableExchangeMarketConfig
		marketConfigs         []*types.MutableMarketConfig
		expectedError         error
	}{
		"Success: 0 markets": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id:                   ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{},
			},
			marketConfigs: []*types.MutableMarketConfig{},
		},
		"Success: 1 market": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
						Invert: true,
					},
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:           1,
					Pair:         "ETHUSD",
					Exponent:     -5,
					MinExchanges: 1,
				},
			},
		},
		"Success: Multiple markets with conversion details": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
					},
					2: {
						Ticker: "USDBTC",
						Invert: true,
					},
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:           1,
					Pair:         "ETHUSD",
					Exponent:     -5,
					MinExchanges: 1,
				},
				{
					Id:           2,
					Pair:         "BTCUSD",
					Exponent:     -6,
					MinExchanges: 1,
				},
			},
		},
		"Failure: Missing market config": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:           3,
					Pair:         "BTCUSD",
					Exponent:     -6,
					MinExchanges: 1,
				},
			},
			expectedError: fmt.Errorf("no market config for market 1 on exchange 'Exchange1'"),
		},
		"Failure: Invalid market config": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker: "ETHUSD",
					},
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					// Missing pair
					Id:           1,
					Exponent:     -5,
					MinExchanges: 1,
				},
			},
			expectedError: fmt.Errorf("invalid market config for market 1 on exchange 'Exchange1': pair cannot be empty"),
		},
		"Failure: no market config exists for adjust-by market": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
					1: {
						Ticker:         "ETHUSD",
						AdjustByMarket: newUint32WithValue(2),
						Invert:         false,
					},
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:           1,
					Pair:         "ETHUSD",
					Exponent:     -5,
					MinExchanges: 1,
				},
			},
			expectedError: fmt.Errorf(
				"no market config for adjust-by market 2 used to convert market 1 price on exchange 'Exchange1'",
			),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.mutableExchangeConfig.Validate(tc.marketConfigs)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.expectedError.Error(), err.Error())
			}
		})
	}
}
