package clob_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"

	"github.com/cometbft/cometbft/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalGating_NegativeTncSubaccount_BlocksThenUnblocks(t *testing.T) {
	Alice_Num0_AfterWithdrawal := satypes.Subaccount{
		Id: &constants.Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			testutil.CreateSingleAssetPosition(
				uint32(0),
				big.NewInt(9_999_999_999),
			),
		},
		PerpetualPositions: nil,
	}
	tests := map[string]struct {
		// State.
		subaccounts                   []satypes.Subaccount
		marketIdToOraclePriceOverride map[uint32]uint64
		collateralPoolBalances        map[string]int64

		// Parameters.
		placedMatchableOrders     []clobtypes.MatchableOrder
		liquidatableSubaccountIds []satypes.SubaccountId
		negativeTncSubaccountIds  []satypes.SubaccountId

		// Configuration.
		liquidationConfig            clobtypes.LiquidationsConfig
		liquidityTiers               []perptypes.LiquidityTier
		perpetuals                   []perptypes.Perpetual
		clobPairs                    []clobtypes.ClobPair
		transferOrWithdrawSubaccount satypes.SubaccountId
		isWithdrawal                 bool
		gatedPerpetualId             uint32

		// Expectations.
		expectedSubaccounts                      []satypes.Subaccount
		expectedSubaccountsAfterWithdrawal       []satypes.Subaccount
		expectedWithdrawalsGated                 map[uint32]bool
		expectedNegativeTncSubaccountSeenAtBlock map[uint32]uint32
		expectedErr                              string
	}{
		`Can place a liquidation order that is unfilled and cannot be deleveraged due to
		non-overlapping bankruptcy prices, withdrawals are gated`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_10_000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 1_000_000_000_000, // $1,000,000 USDC
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing at $50,000
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			negativeTncSubaccountIds:  []satypes.SubaccountId{constants.Carl_Num0},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs:                    []clobtypes.ClobPair{constants.ClobPair_Btc},
			transferOrWithdrawSubaccount: constants.Dave_Num1,
			isWithdrawal:                 true,
			gatedPerpetualId:             constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1.Params.Id,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails.
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
			},
			expectedWithdrawalsGated: map[uint32]bool{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1.Params.Id: true,
			},
			expectedNegativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1.Params.Id: 4,
			},
			expectedErr: "WithdrawalsAndTransfersBlocked: failed to apply subaccount updates",
		},
		`Can place a liquidation order that is unfilled and cannot be deleveraged due to
		non-overlapping bankruptcy prices for isolated market, withdrawals are gated for isolated subaccount`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1ISO_Short_49USD,
				constants.Dave_Num0_1ISO_Long_50USD_Short,
				constants.Alice_Num0_1ISO_LONG_10_000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.IsoUsd_IsolatedMarket.Params.MarketId: 5_050_000_000, // $50.5 / ISO
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String():              1_000_000_000_000, // $1,000,000 USDC
				constants.IsoCollateralPoolAddress.String(): 1_000_000_000_000, // $1,000,000 USDC
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 ISO short is $49, and closing at $50
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob3_Sell025ISO_Price50_GTB11,
			},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			negativeTncSubaccountIds:  []satypes.SubaccountId{constants.Carl_Num0},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs:                    []clobtypes.ClobPair{constants.ClobPair_3_Iso},
			transferOrWithdrawSubaccount: constants.Alice_Num0,
			isWithdrawal:                 true,
			gatedPerpetualId:             constants.IsoUsd_IsolatedMarket.Params.Id,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails.
				// Dave's bankruptcy price to close 1 BTC long is $50, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				constants.Carl_Num0_1ISO_Short_49USD,
				constants.Dave_Num0_1ISO_Long_50USD_Short,
				constants.Alice_Num0_1ISO_LONG_10_000USD,
			},
			expectedWithdrawalsGated: map[uint32]bool{
				constants.IsoUsd_IsolatedMarket.Params.Id: true,
			},
			expectedNegativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 4,
			},
			expectedErr: "WithdrawalsAndTransfersBlocked: failed to apply subaccount updates",
		},
		`Can place a liquidation order that is unfilled and cannot be deleveraged due to
		non-overlapping bankruptcy prices for isolated market, withdrawals are not gated for non-isolated subaccount`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1ISO_Short_49USD,
				constants.Dave_Num0_1ISO_Long_50USD_Short,
				constants.Alice_Num0_10_000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.IsoUsd_IsolatedMarket.Params.MarketId: 5_050_000_000, // $50.5 / ISO
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String():              1_000_000_000_000, // $1,000,000 USDC
				constants.IsoCollateralPoolAddress.String(): 1_000_000_000_000, // $1,000,000 USDC
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 ISO short is $49, and closing at $50
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob3_Sell025ISO_Price50_GTB11,
			},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			negativeTncSubaccountIds:  []satypes.SubaccountId{constants.Carl_Num0},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs:                    []clobtypes.ClobPair{constants.ClobPair_3_Iso},
			transferOrWithdrawSubaccount: constants.Alice_Num0,
			isWithdrawal:                 true,
			gatedPerpetualId:             constants.IsoUsd_IsolatedMarket.Params.Id,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails.
				// Dave's bankruptcy price to close 1 BTC long is $50, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				constants.Carl_Num0_1ISO_Short_49USD,
				constants.Dave_Num0_1ISO_Long_50USD_Short,
				constants.Alice_Num0_10_000USD,
			},
			expectedSubaccountsAfterWithdrawal: []satypes.Subaccount{
				constants.Carl_Num0_1ISO_Short_49USD,
				constants.Dave_Num0_1ISO_Long_50USD_Short,
				// Alice is not an isolated subaccount, and so can still withdraw.
				Alice_Num0_AfterWithdrawal,
			},
			expectedWithdrawalsGated: map[uint32]bool{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Id: false,
				constants.IsoUsd_IsolatedMarket.Params.Id:                        true,
			},
			expectedNegativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 4,
			},
			expectedErr: "",
		},
		`Can place a liquidation order that is partially-filled filled, deleveraging is skipped but
		its still negative TNC, withdrawals are gated`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_025BTC_Long_50000USD,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 1_000_000_000_000, // $1,000,000 USDC
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10,
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing 0.75 BTC at $50,000
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},

			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			negativeTncSubaccountIds:  []satypes.SubaccountId{constants.Carl_Num0},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs:                    []clobtypes.ClobPair{constants.ClobPair_Btc},
			transferOrWithdrawSubaccount: constants.Dave_Num1,
			isWithdrawal:                 false,
			gatedPerpetualId:             constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1.Params.Id,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails for remaining amount.
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(49_999_000_000-12_499_750_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-75_000_000), // -0.75 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				// Dave_Num0 does not change since deleveraging against this subaccount failed.
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 12_499_750_000),
						},
					},
				},
			},
			expectedWithdrawalsGated: map[uint32]bool{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1.Params.Id: true,
			},
			expectedNegativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1.Params.Id: 4,
			},
			expectedErr: "WithdrawalsAndTransfersBlocked: failed to apply subaccount updates",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *assettypes.GenesisState) {
						genesisState.Assets = []assettypes.Asset{
							*constants.Usdc,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						// Set oracle prices in the genesis.
						pricesGenesis := constants.TestPricesGenesisState

						// Make a copy of the MarketPrices slice to avoid modifying by reference.
						marketPricesCopy := make([]prices.MarketPrice, len(pricesGenesis.MarketPrices))
						copy(marketPricesCopy, pricesGenesis.MarketPrices)

						for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
							exponent, exists := constants.TestMarketIdsToExponents[marketId]
							require.True(t, exists)

							marketPricesCopy[marketId] = prices.MarketPrice{
								Id:       marketId,
								Price:    oraclePrice,
								Exponent: exponent,
							}
						}

						pricesGenesis.MarketPrices = marketPricesCopy
						*genesisState = pricesGenesis
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						// If the collateral pool address is already in bank genesis state, update it.
						foundPools := make(map[string]struct{})
						for i, bal := range genesisState.Balances {
							usdcBal, exists := tc.collateralPoolBalances[bal.Address]
							if exists {
								foundPools[bal.Address] = struct{}{}
								genesisState.Balances[i] = banktypes.Balance{
									Address: bal.Address,
									Coins: sdktypes.Coins{
										sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(usdcBal)),
									},
								}
							}
						}
						// If the collateral pool address is not in the bank genesis state, add it.
						for addr, bal := range tc.collateralPoolBalances {
							_, exists := foundPools[addr]
							if exists {
								continue
							}
							genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
								Address: addr,
								Coins: sdktypes.Coins{
									sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(bal)),
								},
							})
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = tc.liquidityTiers
						genesisState.Perpetuals = tc.perpetuals
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
						genesisState.ClobPairs = tc.clobPairs
						genesisState.LiquidationsConfig = tc.liquidationConfig
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Create all existing orders.
			existingOrderMsgs := make([]clobtypes.MsgPlaceOrder, len(tc.placedMatchableOrders))
			for i, matchableOrder := range tc.placedMatchableOrders {
				existingOrderMsgs[i] = clobtypes.MsgPlaceOrder{Order: matchableOrder.MustGetOrder()}
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, existingOrderMsgs...) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			_, err := tApp.App.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
				BlockHeight:                3,
				LiquidatableSubaccountIds:  tc.liquidatableSubaccountIds,
				NegativeTncSubaccountIds:   tc.negativeTncSubaccountIds,
				SubaccountOpenPositionInfo: clobtest.GetOpenPositionsFromSubaccounts(tc.subaccounts),
			})
			require.NoError(t, err)

			// Verify test expectations.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			for _, expectedSubaccount := range tc.expectedSubaccounts {
				require.Equal(
					t,
					expectedSubaccount,
					tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id),
				)
			}
			for perpetualId, expectedWithdrawalsGated := range tc.expectedWithdrawalsGated {
				negativeTncSubaccountSeenAtBlock, exists, err := tApp.App.SubaccountsKeeper.GetNegativeTncSubaccountSeenAtBlock(
					ctx,
					perpetualId,
				)
				require.NoError(t, err)
				require.Equal(t, expectedWithdrawalsGated, exists)
				require.Equal(t, tc.expectedNegativeTncSubaccountSeenAtBlock[perpetualId], negativeTncSubaccountSeenAtBlock)
			}

			// Verify withdrawals are blocked by trying to create a transfer message that withdraws funds.
			var msg proto.Message
			if tc.isWithdrawal {
				withdrawMsg := sendingtypes.MsgWithdrawFromSubaccount{
					Sender:    tc.transferOrWithdrawSubaccount,
					Recipient: tc.transferOrWithdrawSubaccount.Owner,
					AssetId:   constants.Usdc.Id,
					Quantums:  1,
				}
				msg = &withdrawMsg
			} else {
				transferMsg := sendingtypes.MsgCreateTransfer{
					Transfer: &sendingtypes.Transfer{
						Sender:    tc.transferOrWithdrawSubaccount,
						Recipient: constants.Bob_Num0,
						AssetId:   constants.Usdc.Id,
						Amount:    1,
					},
				}
				msg = &transferMsg
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: tc.transferOrWithdrawSubaccount.Owner,
					Gas:                  1000000,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				msg,
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
			ctx = tApp.AdvanceToBlock(
				5,
				testapp.AdvanceToBlockOptions{
					ValidateFinalizeBlock: func(
						ctx sdktypes.Context,
						request abcitypes.RequestFinalizeBlock,
						response abcitypes.ResponseFinalizeBlock,
					) (haltchain bool) {
						// Note the first TX is MsgProposedOperations, the second is all other TXs.
						execResult := response.TxResults[1]
						if tc.expectedErr != "" {
							require.True(t, execResult.IsErr())
							require.Equal(t, satypes.ErrFailedToUpdateSubaccounts.ABCICode(), execResult.Code)
							require.Contains(t, execResult.Log, tc.expectedErr)
						} else {
							require.False(t, execResult.IsErr())
						}
						return false
					},
				},
			)
			var expectedSubaccountsAfterWithdrawal []satypes.Subaccount
			// If an error was expected for withdrawal / transfer, subaccounts should be the same as
			// the ones expected after chain initialization, otherwise the test-case should have a set
			// of expected subaccounts after a successful withdrawal.
			if tc.expectedErr != "" {
				expectedSubaccountsAfterWithdrawal = tc.expectedSubaccounts
			} else {
				expectedSubaccountsAfterWithdrawal = tc.expectedSubaccountsAfterWithdrawal
			}
			for _, expectedSubaccount := range expectedSubaccountsAfterWithdrawal {
				require.Equal(
					t,
					expectedSubaccount,
					tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id),
				)
			}

			// If an error was expected for the withdrawal / transfer, verify that transfers and withdrawals are
			// unblocked after the withdrawal gating period passes.
			if tc.expectedErr != "" {
				_, err = tApp.App.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
					BlockHeight: tc.expectedNegativeTncSubaccountSeenAtBlock[tc.gatedPerpetualId] +
						satypes.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
					LiquidatableSubaccountIds:  tc.liquidatableSubaccountIds,
					NegativeTncSubaccountIds:   []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: clobtest.GetOpenPositionsFromSubaccounts(tc.subaccounts),
				})
				require.NoError(t, err)
				tApp.AdvanceToBlock(
					tc.expectedNegativeTncSubaccountSeenAtBlock[tc.gatedPerpetualId]+
						satypes.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS+
						1,
					testapp.AdvanceToBlockOptions{},
				)
				for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: tc.transferOrWithdrawSubaccount.Owner,
						Gas:                  1000000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					msg,
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
				tApp.AdvanceToBlock(
					tc.expectedNegativeTncSubaccountSeenAtBlock[tc.gatedPerpetualId]+
						satypes.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS+
						2,
					testapp.AdvanceToBlockOptions{},
				)
			}
		})
	}
}
