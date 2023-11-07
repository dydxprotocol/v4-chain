package vest_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"

	"github.com/dydxprotocol/v4-chain/protocol/x/vest"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.DefaultGenesis()

	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	vest.InitGenesis(ctx, k, *genesisState)
	got := vest.ExportGenesis(ctx, k)
	require.NotNil(t, got)
	require.Equal(t, *genesisState, *got)
}

func TestInvalidGenesis_Panics(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	genesisState := types.GenesisState{
		VestEntries: []types.VestEntry{
			{}, // invalid - empty vester account
		},
	}

	require.Panics(t, func() {
		vest.InitGenesis(ctx, k, genesisState)
	})
}
