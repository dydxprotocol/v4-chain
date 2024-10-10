package ratelimit_test

import (
	"testing"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	genState := types.DefaultGenesis()
	ratelimit.InitGenesis(ctx, tApp.App.RatelimitKeeper, *genState)

	// Verify the state after initialization
	exportedGenState := ratelimit.ExportGenesis(ctx, tApp.App.RatelimitKeeper)
	require.NotNil(t, exportedGenState)
	require.Equal(t, genState, exportedGenState)
}

func TestGenesis(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	got := ratelimit.ExportGenesis(ctx, tApp.App.RatelimitKeeper)
	require.NotNil(t, got)
	require.Equal(t, types.DefaultGenesis(), got)
}
