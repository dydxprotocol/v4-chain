package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	require.Equal(t, types.DefaultGenesis().Params, k.GetParams(ctx))
}

func TestSetParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	params := types.Params{
		TreasuryAccount: "dydx12345",
		Denom:           "newdenom",
	}
	require.NoError(t, params.Validate())

	require.NoError(t, k.SetParams(ctx, params))
	require.Equal(t, params, k.GetParams(ctx))
}
