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
)

const (
	AliceBobBTCQuantums = 1_000_000
	CarlDaveBTCQuantums = 2_000_000
	CarlDaveETHQuantums = 4_000_000
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
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
}

func placeOrders(node *containertest.Node, t *testing.T) {
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{
			Order: clobtypes.Order{
				OrderId: clobtypes.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.AliceAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 0,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: AliceBobBTCQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.AliceAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{
			Order: clobtypes.Order{
				OrderId: clobtypes.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.BobAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 0,
				},
				Side:     clobtypes.Order_SIDE_SELL,
				Quantums: AliceBobBTCQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.BobAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{
			Order: clobtypes.Order{
				OrderId: clobtypes.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.CarlAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 0,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: CarlDaveBTCQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.CarlAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{
			Order: clobtypes.Order{
				OrderId: clobtypes.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.DaveAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 0,
				},
				Side:     clobtypes.Order_SIDE_SELL,
				Quantums: CarlDaveBTCQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.DaveAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{
			Order: clobtypes.Order{
				OrderId: clobtypes.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.CarlAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 1,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: CarlDaveETHQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.CarlAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{
			Order: clobtypes.Order{
				OrderId: clobtypes.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.DaveAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 1,
				},
				Side:     clobtypes.Order_SIDE_SELL,
				Quantums: CarlDaveETHQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.DaveAccAddress.String(),
	))

	// FOK order setups.
	require.NoError(t, containertest.BroadcastTx(
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
				Subticks:                        5_000_000,
				TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_FILL_OR_KILL,
				ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
				ConditionalOrderTriggerSubticks: 4_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
					GoodTilBlockTime: uint32(time.Now().Unix() + 300),
				},
			},
		},
		constants.AliceAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
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
				Subticks:                        5_050_000,
				ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
				ConditionalOrderTriggerSubticks: 6_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
					GoodTilBlockTime: uint32(time.Now().Unix() + 300),
				},
			},
		},
		constants.AliceAccAddress.String(),
	))

	err := node.Wait(2)
	require.NoError(t, err)
}
