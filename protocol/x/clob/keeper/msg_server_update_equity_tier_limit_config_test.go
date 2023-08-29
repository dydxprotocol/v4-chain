package keeper_test

import (
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assets "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateEquityTierLimitConfig(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).WithGenesisDocFn(func() types.GenesisDoc {
		genesis := testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
			state.Subaccounts = []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_100_000USD,
			}
		})
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
			state.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{
				ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_001_000_000), // $5,001
						Limit:          1,
					},
				},
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_002_000_000), // $5,002
						Limit:          2,
					},
				},
			}
		})
		return genesis
	}).Build()

	expectedConfig := clobtypes.EquityTierLimitConfiguration{
		ShortTermOrderEquityTiers: []clobtypes.EquityTierLimit{
			{
				UsdTncRequired: dtypes.NewInt(0),
				Limit:          0,
			},
			{
				UsdTncRequired: dtypes.NewInt(5_003_000_000), // $5,003
				Limit:          3,
			},
		},
		StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
			{
				UsdTncRequired: dtypes.NewInt(0),
				Limit:          0,
			},
			{
				UsdTncRequired: dtypes.NewInt(5_004_000_000), // $5,004
				Limit:          4,
			},
		},
	}

	ctx := tApp.InitChain()
	require.NotEqual(t, expectedConfig, tApp.App.ClobKeeper.GetEquityTierLimitConfiguration(ctx))

	proposal, err := gov.NewMsgSubmitProposal(
		[]sdk.Msg{
			&clobtypes.MsgUpdateEquityTierLimitConfiguration{
				//Authority:             constants.AliceAccAddress.String(),
				EquityTierLimitConfig: expectedConfig,
			},
		},
		sdk.NewCoins(sdk.NewInt64Coin(assets.AssetUsdc.Denom, 1_000_000)),
		tApp.App.GovKeeper.GetGovernanceAccount(ctx).GetAddress().String(),
		"metadata",
		"title",
		"summary",
	)
	require.NoError(t, err)

	response := tApp.CheckTx(testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning:        tApp.App.GovKeeper.GetGovernanceAccount(ctx).GetAddress().String(),
			AccSequenceNumberForSigning: 0,
			Gas:                         1_000_000,
		},
		proposal,
	))
	require.True(t, response.IsOK())

	ctx = tApp.AdvanceToBlock(
		2,
		testapp.AdvanceToBlockOptions{
			//DeliverTxsOverride: [][]byte{
			//	testapp.MustMakeCheckTxsWithClobMsg(
			//		ctx,
			//		tApp.App,
			//		clobtypes.MsgUpdateEquityTierLimitConfiguration{
			//			Authority:             constants.AliceAccAddress.String(),
			//			EquityTierLimitConfig: expectedConfig,
			//		},
			//	)[0].Tx,
			//},
		},
	)
	require.Equal(t, expectedConfig, tApp.App.ClobKeeper.GetEquityTierLimitConfiguration(ctx))

}
