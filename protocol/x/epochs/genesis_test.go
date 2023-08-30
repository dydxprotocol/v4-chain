package epochs_test

import (
	"testing"
	"time"

	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		EpochInfoList: []types.EpochInfo{
			{
				Name:                "0",
				Duration:            60,
				FastForwardNextTick: true,
			},
			{
				Name:                "1",
				Duration:            60,
				FastForwardNextTick: true,
			},
		},
	}

	expectedExportState := types.GenesisState{
		EpochInfoList: []types.EpochInfo{
			{
				Name:                "0",
				Duration:            60,
				FastForwardNextTick: true,
			},
			{
				Name:                "1",
				Duration:            60,
				FastForwardNextTick: true,
			},
		},
	}

	ctx, k, _ := keepertest.EpochsKeeper(t)
	initCtx := ctx.WithBlockTime(time.Unix(1800000000, 0))
	epochs.InitGenesis(initCtx, *k, genesisState)
	got := epochs.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, expectedExportState.EpochInfoList, got.EpochInfoList)
	// this line is used by starport scaffolding # genesis/test/assert
}
