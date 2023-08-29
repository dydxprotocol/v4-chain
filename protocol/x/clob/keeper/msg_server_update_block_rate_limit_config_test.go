package keeper_test

import (
	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateBlockRateLimitConfig(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).WithGenesisDocFn(func() types.GenesisDoc {
		genesis := testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
			state.Subaccounts = []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_100_000USD,
			}
		})
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
			state.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
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
				MaxShortTermOrderCancellationsPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 5,
						Limit:     6,
					},
				},
			}
		})
		return genesis
	}).Build()

	expectedConfig := clobtypes.BlockRateLimitConfiguration{
		MaxShortTermOrdersPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
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
		MaxShortTermOrderCancellationsPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
			{
				NumBlocks: 11,
				Limit:     12,
			},
		},
	}

	ctx := tApp.InitChain()
	require.NotEqual(t, expectedConfig, tApp.App.ClobKeeper.GetBlockRateLimitConfiguration(ctx))

	request := clobtypes.MsgUpdateBlockRateLimitConfiguration{
		BlockRateLimitConfig: expectedConfig,
	}
	handler := tApp.App.MsgServiceRouter().Handler(&request)
	_, err := handler(ctx, &request)
	require.NoError(t, err)

	require.Equal(t, expectedConfig, tApp.App.ClobKeeper.GetBlockRateLimitConfiguration(ctx))
}
