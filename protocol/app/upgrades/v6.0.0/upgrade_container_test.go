//go:build all || container_test

package v_6_0_0_test

import (
	"testing"
	"time"

	"github.com/cosmos/gogoproto/proto"
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
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
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
	postUpgradeVaultDefaultQuotingParams(node, t)
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
	require.Equal(t, v_6_0_0.DefaultMarketMap, marketMapResp.MarketMap)

	// check that the market map params have been initialized
	resp, err = containertest.Query(node, marketmapmoduletypes.NewQueryClient,
		marketmapmoduletypes.QueryClient.Params, &marketmapmoduletypes.ParamsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	paramsResp := marketmapmoduletypes.ParamsResponse{}
	err = proto.UnmarshalText(resp.String(), &paramsResp)
	require.NoError(t, err)

	require.Equal(t, v_6_0_0.DefaultMarketMapParams, paramsResp.Params)
}

func postUpgradeVaultDefaultQuotingParams(node *containertest.Node, t *testing.T) {
	params := &vaulttypes.QueryParamsResponse{}
	resp, err := containertest.Query(
		node,
		vaulttypes.NewQueryClient,
		vaulttypes.QueryClient.Params,
		&vaulttypes.QueryParamsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), params)
	require.NoError(t, err)

	require.Nil(t, params.DefaultQuotingParams.Validate())
}
