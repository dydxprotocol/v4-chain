//go:build all || container_test

package v_7_0_0_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"

	v_7_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v7.0.0"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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
