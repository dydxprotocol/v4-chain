package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGetEventParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	require.Equal(t, types.DefaultGenesis().EventParams, k.GetEventParams(ctx))
}

func TestSetEventParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	params := types.EventParams{
		EthAddress: "0xtest",
		EthChainId: uint64(123),
		Denom:      "test-token",
	}
	require.NoError(t, params.Validate())

	require.NoError(t, k.SetEventParams(ctx, params))
	require.Equal(t, params, k.GetEventParams(ctx))
}

func TestGetProposeParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	require.Equal(t, types.DefaultGenesis().ProposeParams, k.GetProposeParams(ctx))
}

func TestSetProposeParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	params := types.ProposeParams{
		MaxBridgesPerBlock:           12,
		ProposeDelayDuration:         34,
		SkipRatePpm:                  56,
		SkipIfBlockDelayedByDuration: 78,
	}
	require.NoError(t, params.Validate())

	require.NoError(t, k.SetProposeParams(ctx, params))
	require.Equal(t, params, k.GetProposeParams(ctx))
}

func TestSetProposeParams_ValidationError(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	params := types.ProposeParams{
		MaxBridgesPerBlock:           12,
		ProposeDelayDuration:         -1,
		SkipRatePpm:                  56,
		SkipIfBlockDelayedByDuration: 78,
	}
	require.ErrorIs(t, params.Validate(), k.SetProposeParams(ctx, params))
	require.NotEqual(t, params, k.GetProposeParams(ctx))
}

func TestGetSafetyParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	require.Equal(t, types.DefaultGenesis().SafetyParams, k.GetSafetyParams(ctx))
}

func TestSetSafetyParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	params := types.SafetyParams{
		IsDisabled:  true,
		DelayBlocks: 1234,
	}
	require.NoError(t, params.Validate())

	require.NoError(t, k.SetSafetyParams(ctx, params))
	require.Equal(t, params, k.GetSafetyParams(ctx))
}
