package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateBlockRateLimitConfig(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		genesis := testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
			state.Subaccounts = []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_100_000USD,
			}
		})
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
			state.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 1,
						Limit:     2,
					},
				},
				MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 3,
						Limit:     4,
					},
				},
			}
		})
		return genesis
	}).Build()

	expectedConfig := clobtypes.BlockRateLimitConfiguration{
		MaxShortTermOrdersAndCancelsPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 7,
				Limit:     8,
			},
		},
		MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 9,
				Limit:     10,
			},
		},
	}

	ctx := tApp.InitChain()
	originalConfig := tApp.App.ClobKeeper.GetBlockRateLimitConfiguration(ctx)
	require.NotEqual(t, expectedConfig, originalConfig)
	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateBlockRateLimitConfiguration{})

	requestWithoutAuthority := clobtypes.MsgUpdateBlockRateLimitConfiguration{
		Authority:            "fake authority",
		BlockRateLimitConfig: expectedConfig,
	}
	_, err := handler(ctx, &requestWithoutAuthority)
	require.Error(t, err, "invalid authority")
	require.Equal(t, originalConfig, tApp.App.ClobKeeper.GetBlockRateLimitConfiguration(ctx))

	requestWithAuthority := clobtypes.MsgUpdateBlockRateLimitConfiguration{
		Authority:            lib.GovModuleAddress.String(),
		BlockRateLimitConfig: expectedConfig,
	}
	_, err = handler(ctx, &requestWithAuthority)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, tApp.App.ClobKeeper.GetBlockRateLimitConfiguration(ctx))
}
