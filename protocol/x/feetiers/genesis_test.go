package feetiers_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	feetiers "github.com/dydxprotocol/v4/x/feetiers"
	"github.com/dydxprotocol/v4/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	got := feetiers.ExportGenesis(ctx, tApp.App.FeeTiersKeeper)
	require.NotNil(t, got)
	require.Equal(t, types.DefaultGenesis(), got)
}
