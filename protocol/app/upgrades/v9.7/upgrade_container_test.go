//go:build all || container_test

package v_9_7_test

import (
	"testing"
	"time"

	v_9_7 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.7"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
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

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_9_7.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

// preExistingConditionalClientId identifies the conditional order placed BEFORE the upgrade.
const preExistingConditionalClientId = 0

// postUpgradeConditionalClientId identifies the conditional order placed AFTER the upgrade.
const postUpgradeConditionalClientId = 1

// buildRestingConditionalOrder returns a conditional order that provably rests UNTRIGGERED
// regardless of the exact genesis oracle scale: it is a TAKE_PROFIT BUY (LTE trigger direction —
// triggers only when the oracle price falls to or below the trigger) with a trigger of 1 subtick,
// so any positive oracle price stays above it and the order never crosses. GoodTilBlockTime is set
// an hour into the future so the order survives the (minutes-long) upgrade without expiring.
// Quantums/Subticks mirror the known-good container place-order test so collateral checks pass.
func buildRestingConditionalOrder(clientId uint32) clobtypes.Order {
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: satypes.SubaccountId{
				Owner:  constants.AliceAccAddress.String(),
				Number: 0,
			},
			ClientId:   clientId,
			OrderFlags: clobtypes.OrderIdFlags_Conditional,
			ClobPairId: 0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 10_000_000,
		Subticks: 1_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(time.Now().Unix() + 3600),
		},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 1,
	}
}

// placeAndCommitOrder broadcasts a MsgPlaceOrder and waits for it to be committed.
func placeAndCommitOrder(node *containertest.Node, t *testing.T, order clobtypes.Order) {
	require.NoError(t, containertest.BroadcastTx(
		node,
		&clobtypes.MsgPlaceOrder{Order: order},
		constants.AliceAccAddress.String(),
	))
	require.NoError(t, node.Wait(3))
}

// requireUntriggeredConditionalExists asserts the conditional order with the given client id is
// present in stateful order state and is NOT triggered.
func requireUntriggeredConditionalExists(node *containertest.Node, t *testing.T, clientId uint32) {
	orderId := clobtypes.OrderId{
		SubaccountId: satypes.SubaccountId{
			Owner:  constants.AliceAccAddress.String(),
			Number: 0,
		},
		ClientId:   clientId,
		OrderFlags: clobtypes.OrderIdFlags_Conditional,
		ClobPairId: 0,
	}
	resp, err := containertest.Query(
		node,
		clobtypes.NewQueryClient,
		clobtypes.QueryClient.StatefulOrder,
		&clobtypes.QueryStatefulOrderRequest{OrderId: orderId},
	)
	require.NoError(t, err, "conditional order (client id %d) should exist in stateful state", clientId)
	statefulOrderResp, ok := resp.(*clobtypes.QueryStatefulOrderResponse)
	require.True(t, ok)
	require.Equal(t, orderId, statefulOrderResp.OrderPlacement.Order.OrderId)
	require.False(t, statefulOrderResp.Triggered, "conditional order (client id %d) should be untriggered", clientId)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {
	// Rest an untriggered conditional order on the pre-upgrade (old) binary. The old binary has no
	// trigger-price index, so this order is the migration's input: the v9.7 upgrade handler must
	// build the index from it and preserve it across the state-breaking upgrade.
	placeAndCommitOrder(node, t, buildRestingConditionalOrder(preExistingConditionalClientId))
}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Baseline: the seeded conditional order rests untriggered on the old binary.
	requireUntriggeredConditionalExists(node, t, preExistingConditionalClientId)
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// 1. Migration integrity: the pre-existing untriggered conditional order survived the
	//    state-breaking upgrade (its untriggered-store entry — which the migration also indexed —
	//    was neither lost nor corrupted, and it was not spuriously triggered). Without the v9.7
	//    migration this order would be absent from the new index and could never trigger.
	requireUntriggeredConditionalExists(node, t, preExistingConditionalClientId)

	// 2. Always-on path on the upgraded binary: a NEW conditional order is admitted (the admission
	//    caps do not reject a normal placement), maintained in the index by the placement hooks, and
	//    rests untriggered — and the chain keeps producing blocks throughout (Wait succeeds).
	placeAndCommitOrder(node, t, buildRestingConditionalOrder(postUpgradeConditionalClientId))
	requireUntriggeredConditionalExists(node, t, postUpgradeConditionalClientId)
}
