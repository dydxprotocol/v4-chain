package gov_test

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestWindDownMarketProposal(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount

		expectedSubaccounts []satypes.Subaccount
	}{
		`Succeeds with final settlement deleveraging, non-negative TNC accounts deleveraged
			at oracle price`: {
			subaccounts: []satypes.Subaccount{
				// well-collateralized long and short positions
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_50_000,
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
		},
		`Succeeds with final settlement deleveraging, negative TNC accounts deleveraged at
			bankruptcy price`: {
			subaccounts: []satypes.Subaccount{
				// negative TNC position
				constants.Carl_Num0_1BTC_Short_49999USD,
				// offsetting position
				constants.Dave_Num0_1BTC_Long_50001USD,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize test app
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{
							constants.BtcUsd_20PercentInitial_10PercentMaintenance,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						// Set oracle prices in the genesis.
						pricesGenesis := constants.TestPricesGenesisState
						*genesisState = pricesGenesis
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{constants.ClobPair_Btc}
						genesisState.LiquidationsConfig = constants.LiquidationsConfig_No_Limit
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				genesis.GenesisTime = GenesisTime
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Build MsgUpdateClobPair
			clobPairId := 0
			clobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
			require.True(t, exists)
			clobPair.Status = clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT
			msgUpdateClobPairToFinalSettlement := &clobtypes.MsgUpdateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair:  clobPair,
			}

			// Submit and Tally Proposal
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{
					msgUpdateClobPairToFinalSettlement,
				},
				false,
				false,
				govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
			)

			// Advance to next block to trigger proposal execution, executed in EndBlocker
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{})

			// Verify clob pair is transitioned to final settlement
			updatedClobPair, exists := tApp.App.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
			require.True(t, exists)
			require.Equal(t, clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT, updatedClobPair.Status)

			// Verify that open stateful orders are cancelled

			// Set liquidation daemon info, to simulate liquidations daemon updating SubaccountOpenPositionInfo
			_, err := tApp.App.Server.LiquidateSubaccounts(
				ctx,
				&api.LiquidateSubaccountsRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
						},
					},
				},
			)
			require.NoError(t, err)

			// Advance block again to trigger final settlement deleveraging in PrepareCheckState
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+5, testapp.AdvanceToBlockOptions{})

			// Verify that final settlement deleveraging occurs
			for _, expectedSubaccount := range tc.expectedSubaccounts {
				fmt.Printf("expectedSubaccount: %+v\n", expectedSubaccount.AssetPositions)
				subaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id)
				fmt.Printf("subaccount: %+v\n", subaccount.AssetPositions)
				require.Equal(t, expectedSubaccount, subaccount)
			}
		})
	}
}
