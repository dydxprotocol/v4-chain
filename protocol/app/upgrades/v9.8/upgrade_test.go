package v_9_8_test

import (
	"testing"

	v_9_8 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.8"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestUpgradeAddsCanonicalUsdcStore(t *testing.T) {
	require.Equal(t, []string{canonicalusdctypes.StoreKey}, v_9_8.Upgrade.StoreUpgrades.Added)
}
