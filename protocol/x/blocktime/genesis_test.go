package blocktime_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/x/blocktime"
	"github.com/dydxprotocol/v4/x/blocktime/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	got := blocktime.ExportGenesis(ctx, tApp.App.BlockTimeKeeper)
	require.NotNil(t, got)
	require.Equal(t, types.DefaultGenesis(), got)
}
