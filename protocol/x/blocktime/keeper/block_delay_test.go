package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/stretchr/testify/require"
)

func Test_Set_GetSynchronyParams_GetBlockDelay(t *testing.T) {
	tests := map[string]struct {
		setUp                   func(keeper.Keeper, sdk.Context)
		expectedSynchronyParams types.SynchronyParams
		expectedBlockDelay      time.Duration
	}{
		"No set-up, empty synchrony params, block_delay = 0": {
			setUp: func(k keeper.Keeper, ctx sdk.Context) {
			},
			expectedSynchronyParams: types.DefaultSynchronyParams(),
			expectedBlockDelay:      0,
		},
		"Non-nil synchrony param": {
			setUp: func(k keeper.Keeper, ctx sdk.Context) {
				k.SetSynchronyParams(ctx, types.SynchronyParams{
					NextBlockDelay: 300 * time.Millisecond,
				})
			},
			expectedSynchronyParams: types.SynchronyParams{
				NextBlockDelay: 300 * time.Millisecond,
			},
			expectedBlockDelay: 300 * time.Millisecond,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.BlockTimeKeeper
			tc.setUp(k, ctx)
			require.Equal(t, tc.expectedSynchronyParams, k.GetSynchronyParams(ctx))
			require.Equal(t, tc.expectedBlockDelay, k.GetBlockDelay(ctx))
		})
	}
}
