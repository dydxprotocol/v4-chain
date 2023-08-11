package stats_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/x/stats"
	"github.com/dydxprotocol/v4/x/stats/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	got := stats.ExportGenesis(ctx, tApp.App.StatsKeeper)
	require.NotNil(t, got)
	require.Equal(t, types.DefaultGenesis(), got)
}
