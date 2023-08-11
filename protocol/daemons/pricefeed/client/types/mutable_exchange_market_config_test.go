package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	ExchangeId1 = "Exchange1"
	ExchangeId2 = "Exchange2"
)

func TestCopy(t *testing.T) {
	mutableMarketExchangeConfig := &types.MutableExchangeMarketConfig{
		Id: ExchangeId1,
		MarketToTicker: map[types.MarketId]string{
			exchange_common.MARKET_ETH_USD: "ETHUSD",
			exchange_common.MARKET_BTC_USD: "BTCUSD",
		},
	}
	mmecCopy := mutableMarketExchangeConfig.Copy()
	require.NotSame(t, mutableMarketExchangeConfig, mmecCopy)
	require.Equal(t, mutableMarketExchangeConfig, mmecCopy)
}

func TestGetMarketIds_Success(t *testing.T) {
	tests := map[string]struct {
		marketToTicker  map[types.MarketId]string
		expectedMarkets []types.MarketId
	}{
		"Empty map": {
			marketToTicker:  map[types.MarketId]string{},
			expectedMarkets: []types.MarketId{},
		},
		"One market": {
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_ETH_USD: "ETHUSD",
			},
			expectedMarkets: []types.MarketId{
				exchange_common.MARKET_ETH_USD,
			},
		},
		"Multiple markets": {
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_ETH_USD: "ETHUSD",
				exchange_common.MARKET_BTC_USD: "BTCUSD",
			},
			expectedMarkets: []types.MarketId{
				exchange_common.MARKET_BTC_USD,
				exchange_common.MARKET_ETH_USD,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memc := types.MutableExchangeMarketConfig{
				Id:             ExchangeId1,
				MarketToTicker: tc.marketToTicker,
			}
			actualMarkets := memc.GetMarketIds()
			require.ElementsMatch(t, tc.expectedMarkets, actualMarkets)
		})
	}
}

func TestEqual_Mixed(t *testing.T) {
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
				MarketToTicker: map[types.MarketId]string{
					1: "ETHUSD",
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "BTCUSD",
				},
			},
			expectedEqual: false,
		},
		"False: non-matching markets": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					2: "ETHUSD",
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "BTCUSD",
				},
			},
			expectedEqual: false,
		},
		"True": {
			A: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "ETHUSD",
				},
			},
			B: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "ETHUSD",
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

func TestValidate_Mixed(t *testing.T) {
	tests := map[string]struct {
		mutableExchangeConfig *types.MutableExchangeMarketConfig
		marketConfigs         []*types.MutableMarketConfig
		expectedError         error
	}{
		"Success: 0 markets": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id:             ExchangeId1,
				MarketToTicker: map[types.MarketId]string{},
			},
			marketConfigs: []*types.MutableMarketConfig{},
		},
		"Success: 1 market": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "ETHUSD",
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:       1,
					Pair:     "ETHUSD",
					Exponent: -5,
				},
			},
		},
		"Failure: Missing market config": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "ETHUSD",
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:       3,
					Pair:     "BTCUSD",
					Exponent: -6,
				},
			},
			expectedError: fmt.Errorf("no market config for market 1 on exchange 'Exchange1'"),
		},
		"Failure: Invalid market config": {
			mutableExchangeConfig: &types.MutableExchangeMarketConfig{
				Id: ExchangeId1,
				MarketToTicker: map[types.MarketId]string{
					1: "ETHUSD",
				},
			},
			marketConfigs: []*types.MutableMarketConfig{
				{
					Id:       1,
					Exponent: -5,
				},
			},
			expectedError: fmt.Errorf("invalid market config for market 1 on exchange 'Exchange1': pair cannot be empty"),
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
