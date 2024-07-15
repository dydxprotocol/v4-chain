//go:build all || container_test

package v_6_0_0_test

import (
	"testing"
	"time"

	"github.com/cosmos/gogoproto/proto"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	marketmapmoduletypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	v_6_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v6.0.0"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	AliceBobBTCQuantums = 1_000_000
)

func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	preUpgradeSetups(node, t)
	preUpgradeChecks(node, t)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_6_0_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {
	placeOrders(node, t)
}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
	preUpgradeStatefulOrderCheck(node, t)
	preUpgradeMarketMapState(node, t)
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
	postUpgradeMarketMapState(node, t)
	postUpgradeStatefulOrderCheck(node, t)
	postUpgradeMarketMapperRevShareChecks(node, t)
}

func placeOrders(node *containertest.Node, t *testing.T) {
	// FOK order setups.
	require.NoError(t, containertest.BroadcastTxWithoutValidateBasic(node, &clobtypes.MsgPlaceOrder{
		Order: clobtypes.Order{
			OrderId: clobtypes.OrderId{
				ClientId: 100,
				SubaccountId: satypes.SubaccountId{
					Owner:  constants.AliceAccAddress.String(),
					Number: 0,
				},
				ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_Conditional,
			},
			Side:                            clobtypes.Order_SIDE_BUY,
			Quantums:                        AliceBobBTCQuantums,
			Subticks:                        6_000_000,
			TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_FILL_OR_KILL,
			ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
			ConditionalOrderTriggerSubticks: 6_000_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(time.Now().Unix() + 300),
			},
		},
	}, constants.AliceAccAddress.String()))

	err := node.Wait(1)
	require.NoError(t, err)

	require.NoError(t, containertest.BroadcastTxWithoutValidateBasic(node, &clobtypes.MsgPlaceOrder{
		Order: clobtypes.Order{
			OrderId: clobtypes.OrderId{
				ClientId: 101,
				SubaccountId: satypes.SubaccountId{
					Owner:  constants.AliceAccAddress.String(),
					Number: 0,
				},
				ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_Conditional,
			},
			Side:                            clobtypes.Order_SIDE_SELL,
			Quantums:                        AliceBobBTCQuantums,
			Subticks:                        6_000_000,
			TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_FILL_OR_KILL,
			ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
			ConditionalOrderTriggerSubticks: 6_000_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
				GoodTilBlockTime: uint32(time.Now().Unix() + 300),
			},
		},
	}, constants.AliceAccAddress.String()))

	err = node.Wait(2)
	require.NoError(t, err)
}

func preUpgradeStatefulOrderCheck(node *containertest.Node, t *testing.T) {
	// Check that all stateful orders are present.
	_, err := containertest.Query(node, clobtypes.NewQueryClient, clobtypes.QueryClient.StatefulOrder,
		&clobtypes.QueryStatefulOrderRequest{
			OrderId: clobtypes.OrderId{
				ClientId: 100,
				SubaccountId: satypes.SubaccountId{
					Owner:  constants.AliceAccAddress.String(),
					Number: 0,
				},
				ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_Conditional,
			},
		})
	require.NoError(t, err)

	_, err = containertest.Query(node, clobtypes.NewQueryClient, clobtypes.QueryClient.StatefulOrder,
		&clobtypes.QueryStatefulOrderRequest{
			OrderId: clobtypes.OrderId{
				ClientId: 101,
				SubaccountId: satypes.SubaccountId{
					Owner:  constants.AliceAccAddress.String(),
					Number: 0,
				},
				ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_Conditional,
			},
		})
	require.NoError(t, err)
}

func preUpgradeMarketMapState(node *containertest.Node, t *testing.T) {
	// check that the market map state does not exist
	_, err := containertest.Query(node, marketmapmoduletypes.NewQueryClient, marketmapmoduletypes.QueryClient.MarketMap,
		&marketmapmoduletypes.MarketMapRequest{})
	require.Error(t, err)
}

func postUpgradeStatefulOrderCheck(node *containertest.Node, t *testing.T) {
	// Check that all stateful orders are removed.
	_, err := containertest.Query(node, clobtypes.NewQueryClient, clobtypes.QueryClient.StatefulOrder,
		&clobtypes.QueryStatefulOrderRequest{
			OrderId: clobtypes.OrderId{
				ClientId: 100,
				SubaccountId: satypes.SubaccountId{
					Owner:  constants.AliceAccAddress.String(),
					Number: 0,
				},
				ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_Conditional,
			},
		})
	require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))

	_, err = containertest.Query(node, clobtypes.NewQueryClient, clobtypes.QueryClient.StatefulOrder,
		&clobtypes.QueryStatefulOrderRequest{
			OrderId: clobtypes.OrderId{
				ClientId: 101,
				SubaccountId: satypes.SubaccountId{
					Owner:  constants.AliceAccAddress.String(),
					Number: 0,
				},
				ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_Conditional,
			},
		})
	require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))
}

func postUpgradeMarketMapperRevShareChecks(node *containertest.Node, t *testing.T) {
	// Check that all rev share params are set to the default value
	resp, err := containertest.Query(node, revsharetypes.NewQueryClient,
		revsharetypes.QueryClient.MarketMapperRevenueShareParams,
		&revsharetypes.QueryMarketMapperRevenueShareParams{})
	require.NoError(t, err)

	params := revsharetypes.QueryMarketMapperRevenueShareParamsResponse{}
	err = proto.UnmarshalText(resp.String(), &params)
	require.NoError(t, err)
	require.Equal(t, params.Params.Address, authtypes.NewModuleAddress(authtypes.FeeCollectorName).String())
	require.Equal(t, params.Params.RevenueSharePpm, uint32(0))
	require.Equal(t, params.Params.ValidDays, uint32(0))

	// Get all markets list
	marketParams := pricestypes.QueryAllMarketParamsResponse{}
	resp, err = containertest.Query(node, pricestypes.NewQueryClient, pricestypes.QueryClient.AllMarketParams,
		&pricestypes.QueryAllMarketParamsRequest{})
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), &marketParams)
	require.NoError(t, err)

	// Check that all rev share details are set to the default value
	for _, market := range marketParams.MarketParams {
		revShareDetails := revsharetypes.QueryMarketMapperRevShareDetailsResponse{}
		resp, err := containertest.Query(node, revsharetypes.NewQueryClient,
			revsharetypes.QueryClient.MarketMapperRevShareDetails,
			&revsharetypes.QueryMarketMapperRevShareDetails{
				MarketId: market.Id,
			})
		require.NoError(t, err)
		err = proto.UnmarshalText(resp.String(), &revShareDetails)
		require.NoError(t, err)
		require.Equal(t, revShareDetails.Details.ExpirationTs, uint64(0))
	}
}

func postUpgradeMarketMapState(node *containertest.Node, t *testing.T) {
	// check that the market map state has been initialized
	resp, err := containertest.Query(node, marketmapmoduletypes.NewQueryClient,
		marketmapmoduletypes.QueryClient.MarketMap, &marketmapmoduletypes.MarketMapRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	marketMapResp := marketmapmoduletypes.MarketMapResponse{}
	err = proto.UnmarshalText(resp.String(), &marketMapResp)
	require.NoError(t, err)

	require.Equal(t, "localdydxprotocol", marketMapResp.ChainId)
	require.Equal(t, expectedMarketMap, marketMapResp.MarketMap)

	// check that the market map params have been initialized
	resp, err = containertest.Query(node, marketmapmoduletypes.NewQueryClient,
		marketmapmoduletypes.QueryClient.Params, &marketmapmoduletypes.ParamsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	paramsResp := marketmapmoduletypes.ParamsResponse{}
	err = proto.UnmarshalText(resp.String(), &paramsResp)
	require.NoError(t, err)

	require.Equal(t, expectedParams, paramsResp.Params)
}

var (
	expectedParams = marketmapmoduletypes.Params{
		MarketAuthorities: []string{"dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"},
		Admin:             "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
	}

	expectedMarketMap = marketmapmoduletypes.MarketMap{
		Markets: map[string]marketmapmoduletypes.Market{
			"ADA/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "ADA", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "ADAUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "ADAUSDT", NormalizeByPair: &slinkytypes.
						CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "ADA-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "ADA_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "adausdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "ADAUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "ADA-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "ADAUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "ADA-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"APE/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "APE", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "APEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "APE-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "APE_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "APEUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "APE-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "APEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "APE-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"APT/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "APT", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "APTUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "APTUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "APT-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "APT_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "aptusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "APT-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "APTUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "APT-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			}, "ARB/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "ARB", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "ARBUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "ARBUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "ARB-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "ARB_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "arbusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "ARB-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "ARBUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "ARB-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			}, "ATOM/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "ATOM", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "ATOMUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "ATOMUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "ATOM-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "ATOM_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "ATOMUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "ATOM-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "ATOMUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "ATOM-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"AVAX/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "AVAX", Quote: "USD"},
					Decimals:         0x8,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "AVAXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "AVAXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "AVAX-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "AVAX_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "avaxusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "AVAXUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "AVAX-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "AVAX-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			}, "BCH/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "BCH", Quote: "USD"},
					Decimals:         0x7,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "BCHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "BCHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "BCH-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "BCH_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "bchusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "BCHUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "BCH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "BCHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "BCH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			}, "BLUR/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "BLUR", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "coinbase_ws", OffChainTicker: "BLUR-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "BLUR_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "BLURUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "BLUR-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "BLURUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "BLUR-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"BTC/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "BTC", Quote: "USD"},
					Decimals:         0x5,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "BTC-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "btcusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XXBTZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			}, "COMP/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "COMP", Quote: "USD"},
					Decimals:         0x8,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "COMPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "COMP-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "COMP_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "COMPUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "COMPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "COMP-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"CRV/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "CRV", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "CRVUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "CRV-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "CRV_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "CRVUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "CRV-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "CRVUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "CRV-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"DOGE/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "DOGE", Quote: "USD"},
					Decimals:         0xb,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "DOGEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "DOGEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "DOGE-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "DOGE_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "dogeusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XDGUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "DOGE-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "DOGEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "DOGE-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"DOT/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "DOT", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "DOTUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "DOTUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "DOT-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "DOT_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "DOTUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "DOT-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "DOTUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "DOT-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"DYDX/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "DYDX", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "DYDXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "DYDXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "DYDX_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "DYDX-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "DYDXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "DYDX-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"ETC/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "ETC", Quote: "USD"},
					Decimals:         0x8,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "ETCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "ETC-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "ETC_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "etcusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "ETC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "ETCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "ETC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"ETH/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "ETH", Quote: "USD"},
					Decimals:         0x6,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "ETH-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "ethusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XETHZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "ETH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "ETH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"FIL/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "FIL", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "FILUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "FIL-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "FIL_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "filusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "FILUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "FILUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "FIL-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"LDO/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "LDO", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "LDOUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "LDO-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "LDOUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "LDO-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "LDOUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "LDO-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"LINK/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "LINK", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "LINKUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "LINKUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "LINK-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "LINKUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "LINK-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "LINKUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "LINK-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"LTC/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "LTC", Quote: "USD"},
					Decimals:         0x8,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "LTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "LTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "LTC-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "ltcusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XLTCZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "LTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "LTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "LTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"MATIC/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "MATIC", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "MATICUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "MATICUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "MATIC-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "MATIC_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "maticusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "MATICUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "MATIC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "MATICUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "MATIC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"MKR/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "MKR", Quote: "USD"},
					Decimals:         0x6,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "MKRUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "MKR-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "MKRUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "MKR-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "MKRUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "MKR-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"NEAR/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "NEAR", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "NEARUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "NEAR-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "NEAR_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "nearusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "NEAR-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "NEARUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "NEAR-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"OP/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "OP", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3, Enabled: true,
					Metadata_JSON: "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "OPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "OP-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "OP_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "OP-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "OPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "OP-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"PEPE/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "PEPE", Quote: "USD"},
					Decimals:         0x10,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "PEPEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "PEPEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "PEPE_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "PEPEUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "PEPE-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "PEPEUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "PEPE-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"SEI/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "SEI", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "SEIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "SEIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "SEI-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "SEI_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "seiusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "SEI-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "SEIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"SHIB/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "SHIB", Quote: "USD"},
					Decimals:         0xf,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "SHIBUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "SHIBUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "SHIB-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "SHIB_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "SHIBUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "SHIB-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "SHIBUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "SHIB-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"SOL/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "SOL", Quote: "USD"},
					Decimals:         0x8,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "SOLUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "SOLUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "SOL-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "solusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "SOLUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "SOL-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "SOLUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "SOL-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"SUI/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "SUI", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "SUIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "SUIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "SUI-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "SUI_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "suiusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "SUI-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "SUIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "SUI-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"TEST/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "TEST", Quote: "USD"},
					Decimals:         0x5,
					MinProviderCount: 0x1,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "volatile-exchange-provider", OffChainTicker: "TEST-USD", NormalizeByPair: nil,
						Invert: false, Metadata_JSON: ""},
				},
			},
			"TRX/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "TRX", Quote: "USD"},
					Decimals:         0xb,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "TRXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "TRXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "TRX_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "trxusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "TRXUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "TRX-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "TRXUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "TRX-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"UNI/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "UNI", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "UNIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "UNIUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "UNI-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "UNI_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "UNIUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "UNI-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "UNI-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"USDT/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "USDCUSDT",
						NormalizeByPair: nil, Invert: true, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "USDCUSDT",
						NormalizeByPair: nil, Invert: true, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "USDT-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "ethusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "ETH", Quote: "USD"}, Invert: true, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "USDTZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "BTC", Quote: "USD"}, Invert: true, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "USDC-USDT",
						NormalizeByPair: nil, Invert: true, Metadata_JSON: ""},
				},
			},
			"WLD/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "WLD", Quote: "USD"},
					Decimals:         0x9,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "WLDUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "WLDUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "WLD_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "wldusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "WLD-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "WLDUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "WLD-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"XLM/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "XLM", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "XLMUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "XLMUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "XLM-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XXLMZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "XLM-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "XLMUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "XLM-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			"XRP/USD": {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "XRP", Quote: "USD"},
					Decimals:         0xa,
					MinProviderCount: 0x3,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "XRPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "XRPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "XRP-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "gate_ws", OffChainTicker: "XRP_USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "xrpusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XXRPZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "XRP-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "XRPUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "XRP-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
		},
	}
)
