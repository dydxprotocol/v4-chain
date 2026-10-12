//go:build all || container_test

package v_9_8_test

import (
	"testing"

	v_9_8 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.8"
	"github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

// TestStateUpgrade is a smoke test for the v9.8 upgrade: the chain applies the upgrade and keeps
// producing blocks. The orphan-TWAP-state migration cannot be seeded through the preupgrade genesis
// (orphan state only arises from in-flight suborders), so the migration's effect is covered by the
// keeper and app-level tests; here we only verify the node advances across the upgrade.
func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_9_8.UpgradeName)
	require.NoError(t, err)

	// The chain continues to produce blocks after the migration runs.
	require.NoError(t, node.Wait(1))
}
