//go:build all || container_test

package v_5_2_0_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	v_5_2_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v5.2.0"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

const (
	AliceBobBTCQuantums = 1_000_000
	CarlDaveETHQuantums = 2_000_000
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

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_5_2_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

/*** Preupgrade Setup ***/
func preUpgradeSetups(node *containertest.Node, t *testing.T) {
	preUpgradePlaceLongTermOrders(node, t)
	preUpgradeSetupVaults(node, t)
}

func preUpgradePlaceLongTermOrders(node *containertest.Node, t *testing.T) {
	latestBlockTime, err := node.LatestBlockTime()
	require.NoError(t, err)
	goodTilBlockTime := uint32(latestBlockTime.Unix()) + 600

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
					OrderFlags: 64,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: AliceBobBTCQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
					GoodTilBlockTime: goodTilBlockTime,
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
					OrderFlags: 64,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: AliceBobBTCQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
					GoodTilBlockTime: goodTilBlockTime,
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
					ClobPairId: 1,
					OrderFlags: 64,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: CarlDaveETHQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
					GoodTilBlockTime: goodTilBlockTime,
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
					OrderFlags: 64,
				},
				Side:     clobtypes.Order_SIDE_BUY,
				Quantums: CarlDaveETHQuantums,
				Subticks: 5_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
					GoodTilBlockTime: goodTilBlockTime,
				},
			},
		},
		constants.DaveAccAddress.String(),
	))

	err = node.Wait(2)
	require.NoError(t, err)
}

func preUpgradeSetupVaults(node *containertest.Node, t *testing.T) {
	// Query vault activation threshold quote quantums.
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

	// Wait until subaccount transfers are enabled.
	err = node.WaitUntilBlockHeight(51)
	require.NoError(t, err)

	// Deposit to vaults.
	require.NoError(t, containertest.BroadcastTx(
		node,
		&vaulttypes.MsgDepositToVault{
			VaultId:       &constants.Vault_Clob_0,
			SubaccountId:  &constants.Alice_Num0,
			QuoteQuantums: params.Params.ActivationThresholdQuoteQuantums,
		},
		constants.AliceAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&vaulttypes.MsgDepositToVault{
			VaultId:       &constants.Vault_Clob_1,
			SubaccountId:  &constants.Bob_Num0,
			QuoteQuantums: params.Params.ActivationThresholdQuoteQuantums,
		},
		constants.BobAccAddress.String(),
	))
	require.NoError(t, containertest.BroadcastTx(
		node,
		&vaulttypes.MsgDepositToVault{
			VaultId: &vaulttypes.VaultId{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 17,
			},
			SubaccountId:  &constants.Carl_Num0,
			QuoteQuantums: params.Params.ActivationThresholdQuoteQuantums,
		},
		constants.CarlAccAddress.String(),
	))

	err = node.Wait(2)
	require.NoError(t, err)
}

/*** Preupgrade Check ***/
func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
}

/*** Postupgrade Check ***/
func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	postUpgradeCheckVaultModuleParams(node, t)
	postUpgradeCheckVaultBestFeeTier(node, t)
}

func postUpgradeCheckVaultModuleParams(node *containertest.Node, t *testing.T) {
	// Get vault module params.
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

	// Ensure that `OrderExpirationSeconds` is updated to 60.
	require.Equal(t, uint32(60), params.Params.OrderExpirationSeconds)
}

func postUpgradeCheckVaultBestFeeTier(node *containertest.Node, t *testing.T) {
	// Get all vaults.
	allVaults := &vaulttypes.QueryAllVaultsResponse{}
	resp, err := containertest.Query(
		node,
		vaulttypes.NewQueryClient,
		vaulttypes.QueryClient.AllVaults,
		&vaulttypes.QueryAllVaultsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), allVaults)
	require.NoError(t, err)
	// Ensure that there are vaults.
	require.NotEmpty(t, allVaults.Vaults)

	// Get fee tiers.
	feeParams := &feetierstypes.QueryPerpetualFeeParamsResponse{}
	resp, err = containertest.Query(
		node,
		feetierstypes.NewQueryClient,
		feetierstypes.QueryClient.PerpetualFeeParams,
		&feetierstypes.QueryPerpetualFeeParamsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), feeParams)
	require.NoError(t, err)

	// Verify that every vault is of the best fee tier.
	bestFeeTierIndex := len(feeParams.Params.Tiers) - 1
	for _, vault := range allVaults.Vaults {
		userFeeTier := &feetierstypes.QueryUserFeeTierResponse{}
		resp, err = containertest.Query(
			node,
			feetierstypes.NewQueryClient,
			feetierstypes.QueryClient.UserFeeTier,
			&feetierstypes.QueryUserFeeTierRequest{
				User: vault.SubaccountId.Owner,
			},
		)
		require.NoError(t, err)
		err = proto.UnmarshalText(resp.String(), userFeeTier)
		require.NoError(t, err)

		require.Equal(t, bestFeeTierIndex, int(userFeeTier.Index))
	}
}
