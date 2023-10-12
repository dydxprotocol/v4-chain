package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "rewards", types.ModuleName)
	require.Equal(t, "rewards", types.StoreKey)
	require.Equal(t, "tmp_rewards", types.TransientStoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "Shares:", types.RewardShareKeyPrefix)
	require.Equal(t, "Params", types.ParamsKey)
}

func TestModuleAccountKeys(t *testing.T) {
	require.Equal(t, "rewards_treasury", types.TreasuryAccountName)
	require.Equal(t, "rewards_vester", types.VesterAccountName)
}
