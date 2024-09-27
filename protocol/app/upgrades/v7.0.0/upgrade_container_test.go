//go:build all || container_test

package v_7_0_0_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/gogoproto/proto"

	v_7_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v7.0.0"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliatestypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
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

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_7_0_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeSetups(node *containertest.Node, t *testing.T) {}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	// Add test for your upgrade handler logic below
	postUpgradeVaultParamsCheck(node, t)
	postUpgradeMegavaultSharesCheck(node, t)

	// Check that the affiliates module has been initialized with the default tiers.
	postUpgradeAffiliatesModuleTiersCheck(node, t)
}

func postUpgradeVaultParamsCheck(node *containertest.Node, t *testing.T) {
	// Check that a vault with quoting params is successfully migrated and the quoting params are
	// successfully migrated to the vault params.
	expectedQuotingParams := &vaulttypes.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     1500,
		SpreadBufferPpm:                  500,
		SkewFactorPpm:                    1000000,
		OrderSizePctPpm:                  50000,
		OrderExpirationSeconds:           30,
		ActivationThresholdQuoteQuantums: dtypes.NewIntFromUint64(500_000_000),
	}

	checkVaultParams(node, t, 0, vaulttypes.VaultStatus_VAULT_STATUS_QUOTING, expectedQuotingParams)

	// Check that a vault without quoting params is successfully migrated and the quoting params are
	// not set in the migrated vault params.
	checkVaultParams(node, t, 1, vaulttypes.VaultStatus_VAULT_STATUS_QUOTING, nil)
}

func checkVaultParams(
	node *containertest.Node,
	t *testing.T,
	vaultNumber uint32,
	expectedStatus vaulttypes.VaultStatus,
	expectedQuotingParams *vaulttypes.QuotingParams,
) {
	resp, err := containertest.Query(
		node,
		vaulttypes.NewQueryClient,
		vaulttypes.QueryClient.VaultParams,
		&vaulttypes.QueryVaultParamsRequest{
			Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
			Number: vaultNumber,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	vaultParamsResp := vaulttypes.QueryVaultParamsResponse{}
	err = proto.UnmarshalText(resp.String(), &vaultParamsResp)
	require.NoError(t, err)

	require.Equal(t, expectedStatus, vaultParamsResp.VaultParams.Status)
	require.Equal(t, expectedQuotingParams, vaultParamsResp.VaultParams.QuotingParams)
}

func postUpgradeMegavaultSharesCheck(node *containertest.Node, t *testing.T) {
	// Alice equity = vault_0_equity * 1 + vault_1_equity * 1/3 + vault_2_equity * 123_456/556_677
	// = 1_000 + 2_000 * 1/3 + 3_000 * 123_456/556_677
	// ~= 2331.99
	// Bob equity = vault_1_equity * 1/3 + vault_2_equity * 433_221/556_677
	// = 2_000 * 1/3 + 3_000 * 433_221/556_677
	// ~= 3001.35
	// Carl equity = vault_1_equity * 1/3
	// = 2_000 * 1/3
	// ~= 666.67
	// 1 USDC in equity should be granted 1 megavault share and round down to nearest integer.
	expectedOwnerShares := map[string]*big.Int{
		constants.AliceAccAddress.String(): big.NewInt(2_331),
		constants.BobAccAddress.String():   big.NewInt(3_001),
		constants.CarlAccAddress.String():  big.NewInt(666),
	}
	// 2331 + 3001 + 666 = 5998
	expectedTotalShares := big.NewInt(5_998)

	// Check MegaVault total shares.
	resp, err := containertest.Query(
		node,
		vaulttypes.NewQueryClient,
		vaulttypes.QueryClient.MegavaultTotalShares,
		&vaulttypes.QueryMegavaultTotalSharesRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	totalSharesResp := vaulttypes.QueryMegavaultTotalSharesResponse{}
	err = proto.UnmarshalText(resp.String(), &totalSharesResp)
	require.NoError(t, err)

	require.Equal(
		t,
		expectedTotalShares,
		totalSharesResp.TotalShares.NumShares.BigInt(),
	)

	// Check MegaVault owner shares.
	resp, err = containertest.Query(
		node,
		vaulttypes.NewQueryClient,
		vaulttypes.QueryClient.MegavaultAllOwnerShares,
		&vaulttypes.QueryMegavaultAllOwnerSharesRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	allOwnerSharesResp := vaulttypes.QueryMegavaultAllOwnerSharesResponse{}
	err = proto.UnmarshalText(resp.String(), &allOwnerSharesResp)
	require.NoError(t, err)

	require.Len(t, allOwnerSharesResp.OwnerShares, 3)
	gotOwnerShares := make(map[string]*big.Int)
	for _, ownerShare := range allOwnerSharesResp.OwnerShares {
		gotOwnerShares[ownerShare.Owner] = ownerShare.Shares.NumShares.BigInt()
	}
	for owner, expectedShares := range expectedOwnerShares {
		require.Contains(t, gotOwnerShares, owner)
		require.Equal(t, expectedShares, gotOwnerShares[owner])
	}
}

func postUpgradeAffiliatesModuleTiersCheck(node *containertest.Node, t *testing.T) {
	resp, err := containertest.Query(
		node,
		affiliatestypes.NewQueryClient,
		affiliatestypes.QueryClient.AllAffiliateTiers,
		&affiliatestypes.AllAffiliateTiersRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	affiliateTiersResp := affiliatestypes.AllAffiliateTiersResponse{}
	err = proto.UnmarshalText(resp.String(), &affiliateTiersResp)
	require.NoError(t, err)
	require.Equal(t, affiliatestypes.DefaultAffiliateTiers, affiliateTiersResp.Tiers)
}
