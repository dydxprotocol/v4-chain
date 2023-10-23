package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	require.Equal(t, types.DefaultGenesis().Params, k.GetParams(ctx))
}

func TestSetParams_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
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

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		input     types.Params
		expErrMsg string
	}{
		{
			name: "empty treasure account name",
			input: types.Params{
				TreasuryAccount: "",
			},
			expErrMsg: "treasury account cannot have empty name",
		},
		{
			name: "invalid denom",
			input: types.Params{
				TreasuryAccount: "treasury_account",
				Denom:           "invalid dnom !!!",
			},
			expErrMsg: "invalid denom",
		},
		{
			name: "invalid FeeMultiplierPpm",
			input: types.Params{
				TreasuryAccount:  "treasury_account",
				Denom:            "foo",
				FeeMultiplierPpm: 1_000_001,
			},
			expErrMsg: "FeeMultiplierPpm cannot be greater than 1_000_000 (100%)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.expErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expErrMsg)
			}
		})
	}
}
