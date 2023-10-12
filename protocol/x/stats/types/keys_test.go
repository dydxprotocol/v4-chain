package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "stats", types.ModuleName)
	require.Equal(t, "stats", types.StoreKey)
	require.Equal(t, "tmp_stats", types.TransientStoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "Epoch:", types.EpochStatsKeyPrefix)
	require.Equal(t, "User:", types.UserStatsKeyPrefix)
	require.Equal(t, "Metadata", types.StatsMetadataKey)
	require.Equal(t, "Global", types.GlobalStatsKey)
	require.Equal(t, "Block", types.BlockStatsKey)
	require.Equal(t, "Params", types.ParamsKey)
}
