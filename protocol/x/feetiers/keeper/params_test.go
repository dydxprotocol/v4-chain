package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	require.Equal(t, types.DefaultGenesis().Params, k.GetPerpetualFeeParams(ctx))
}

func TestSetPerpetualFeeParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	params := types.PerpetualFeeParams{
		Tiers: []*types.PerpetualFeeTier{
			{},
			{
				AbsoluteVolumeRequirement: 100,
				MakerFeePpm:               1,
				TakerFeePpm:               1,
			},
		},
	}
	require.NoError(t, params.Validate())

	require.NoError(t, k.SetPerpetualFeeParams(ctx, params))
	require.Equal(t, params, k.GetPerpetualFeeParams(ctx))
}
