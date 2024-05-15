package metrics_test

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	grpc_util "github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"unsafe"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/stretchr/testify/require"
)

const (
	INVALID_ID = 10000000
)

// Used to clear the marketToPair map for testing purposes.
//
//go:linkname marketToPair github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics.marketToPair
var marketToPair map[types.MarketId]string

// clearMarketToPair resets the backing marketToPair map to an empty map.
func clearMarketToPair() {
	marketToPairValue := reflect.ValueOf(&marketToPair).Elem()
	rp := reflect.NewAt(marketToPairValue.Type(), unsafe.Pointer(marketToPairValue.UnsafeAddr())).Elem()
	newMap := map[types.MarketId]string{}
	rp.Set(reflect.ValueOf(newMap))
}

func TestGetLabelForMarketIdSuccess(t *testing.T) {
	pricefeedmetrics.SetMarketPairForTelemetry(exchange_config.MARKET_BTC_USD, "BTCUSD")
	require.Equal(
		t,
		metrics.GetLabelForStringValue(metrics.MarketId, "BTCUSD"),
		pricefeedmetrics.GetLabelForMarketId(exchange_config.MARKET_BTC_USD),
	)
}

func TestGetLabelForMarketIdFailure(t *testing.T) {
	require.Equal(
		t,
		metrics.GetLabelForStringValue(metrics.MarketId, fmt.Sprintf("invalid_id:%d", INVALID_ID)),
		pricefeedmetrics.GetLabelForMarketId(INVALID_ID),
	)
}

func TestGetLabelForExchangeId(t *testing.T) {
	require.Equal(
		t,
		metrics.GetLabelForStringValue(metrics.ExchangeId, exchange_common.EXCHANGE_ID_BINANCE_US),
		pricefeedmetrics.GetLabelForExchangeId(exchange_common.EXCHANGE_ID_BINANCE_US),
	)
}

// The market id -> pair mapping is automatically populated in the pricefeed daemon.
func TestGetLabelForMarketIdAutoPopulated(t *testing.T) {
	tests := map[string]struct {
		MarketParams  []pricetypes.MarketParam
		ExpectedLabel string
	}{
		"Without params, label is 'invalid_id'": {
			MarketParams:  []pricetypes.MarketParam{},
			ExpectedLabel: "invalid_id:1",
		},
		"With params, label is 'BTC-USD'": {
			MarketParams: []pricetypes.MarketParam{
				{
					Id:           1,
					MinExchanges: 1,
					Exponent:     -9,
					Pair:         "BTC-USD",
				},
			},
			ExpectedLabel: "BTC-USD",
		},
	}

	for name, tc := range tests {
		// Clear the mapping from market to pair to start the test with a blank state
		clearMarketToPair()

		t.Run(name, func(t *testing.T) {
			// Mock the `QueryAllMarketParams` call to return the given `MarketParams`
			r := &pricetypes.QueryAllMarketParamsResponse{MarketParams: tc.MarketParams}

			pricesQueryClient := &mocks.QueryClient{}
			pricesQueryClient.On(
				"AllMarketParams",
				grpc_util.Ctx,
				mock.Anything,
			).Return(r, nil)

			configs := &mocks.PricefeedMutableMarketConfigs{}
			configs.On(
				"UpdateMarkets",
				mock.Anything,
			).Return(map[types.MarketId]error{}, nil)

			// Run the same market param updater task loop that the daemon calls
			client.RunMarketParamUpdaterTaskLoop(
				grpc_util.Ctx,
				configs,
				pricesQueryClient,
				log.NewNopLogger(),
				true,
			)

			// Now check that the label is as expected
			haveLabel := pricefeedmetrics.GetLabelForMarketId(1).Value
			if haveLabel != tc.ExpectedLabel {
				t.Errorf("Expected label: %v, got: %v", tc.ExpectedLabel, haveLabel)
			}
		})
	}
}
