//go:build all || container_test

package v_6_0_0_test

import (
	"testing"
	"time"

	v_6_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v6.0.0"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
	postUpgradeStatefulOrderCheck(node, t)
}

func placeOrders(node *containertest.Node, t *testing.T) {
	// FOK order setups.
	require.NoError(t, containertest.BroadcastTxWithoutValidateBasic(
		node,
		&clobtypes.MsgPlaceOrder{
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
		},
		constants.AliceAccAddress.String(),
	))

	err := node.Wait(1)
	require.NoError(t, err)

	require.NoError(t, containertest.BroadcastTxWithoutValidateBasic(
		node,
		&clobtypes.MsgPlaceOrder{
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
		},
		constants.AliceAccAddress.String(),
	))

	err = node.Wait(2)
	require.NoError(t, err)
}

func preUpgradeStatefulOrderCheck(node *containertest.Node, t *testing.T) {
	// Check that all stateful orders are present.
	_, err := containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.QueryClient.StatefulOrder,
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
		},
	)
	require.NoError(t, err)

	_, err = containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.QueryClient.StatefulOrder,
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
		},
	)
	require.NoError(t, err)
}

func postUpgradeStatefulOrderCheck(node *containertest.Node, t *testing.T) {
	// Check that all stateful orders are removed.
	_, err := containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.QueryClient.StatefulOrder,
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
		},
	)
	require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))

	_, err = containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.QueryClient.StatefulOrder,
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
		},
	)
	require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))
}
