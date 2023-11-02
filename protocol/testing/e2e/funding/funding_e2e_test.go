package funding_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricefeed_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

const (
	BlockTimeDuration             = 2 * time.Second
	NumBlocksPerMinute            = int64(time.Minute / BlockTimeDuration) // 30
	BlockHeightAtFirstFundingTick = 1000
	TestTransferUsdcForSettlement = 10_000_000_000_000
	TestMarketId                  = 0
)

var (
	OrderTemplate_Alice_Num0_Id0_Clob0_Buy_LongTerm = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side: clobtypes.Order_SIDE_BUY,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(GenesisTime.Add(24 * time.Hour).Unix()),
		},
	}
	OrderTemplate_Alice_Num0_Id1_Clob0_Buy_LongTerm = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side: clobtypes.Order_SIDE_BUY,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(GenesisTime.Add(24 * time.Hour).Unix()),
		},
	}
	OrderTemplate_Bob_Num0_Id0_Clob0_Sell_LongTerm = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side: clobtypes.Order_SIDE_SELL,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(GenesisTime.Add(24 * time.Hour).Unix()),
		},
	}
	OrderTemplate_Bob_Num0_Id1_Clob0_Sell_LongTerm = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side: clobtypes.Order_SIDE_SELL,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(GenesisTime.Add(24 * time.Hour).Unix()),
		},
	}
	OrderTemplate_Carl_Num0_Id0_Clob0_Buy_LongTerm = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side: clobtypes.Order_SIDE_BUY,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(GenesisTime.Add(24 * time.Hour).Unix()),
		},
	}
	// Genesis time of the chain
	GenesisTime                               = time.Unix(1690000000, 0)
	FirstFundingSampleTick                    = time.Unix(1690000050, 0)
	FirstFundingTick                          = time.Unix(1690002000, 0)
	LastFundingSampleOfSecondFundingTickEpoch = time.Unix(1690005570, 0)
	SecondFundingTick                         = time.Unix(
		1690002000+int64(epochstypes.FundingTickEpochDuration),
		0,
	)
	// Test transfers to settle Alice and Bob's subaccount.
	TestTranfers = []sendingtypes.MsgCreateTransfer{
		{
			Transfer: &sendingtypes.Transfer{
				Sender:    constants.Dave_Num0,
				Recipient: constants.Alice_Num0,
				AssetId:   assettypes.AssetUsdc.Id,
				Amount:    TestTransferUsdcForSettlement,
			},
		},
		{
			Transfer: &sendingtypes.Transfer{
				Sender:    constants.Dave_Num0,
				Recipient: constants.Bob_Num0,
				AssetId:   assettypes.AssetUsdc.Id,
				Amount:    TestTransferUsdcForSettlement,
			},
		},
		{
			Transfer: &sendingtypes.Transfer{
				Sender:    constants.Dave_Num0,
				Recipient: constants.Carl_Num0,
				AssetId:   assettypes.AssetUsdc.Id,
				Amount:    TestTransferUsdcForSettlement,
			},
		},
	}
)

type expectedSettlements struct {
	SubaccountId satypes.SubaccountId
	Settlement   int64
}

func getSubaccountUsdcBalance(subaccount satypes.Subaccount) int64 {
	return subaccount.AssetPositions[0].Quantums.BigInt().Int64()
}

func TestFunding(t *testing.T) {
	tests := map[string]struct {
		testHumanOrders   []clobtest.TestHumanOrder
		initialIndexPrice map[uint32]string
		// index price to be used in premium calculation
		indexPriceForPremium map[uint32]string
		// oracle price for funding index calculation
		oracelPriceForFundingIndex map[uint32]string
		// address -> funding
		expectedSubaccountSettlements []expectedSettlements
		expectedFundingPremium        int32
		expectedFundingIndex          int64
	}{
		"Test funding": {
			testHumanOrders: []clobtest.TestHumanOrder{
				// Unmatched orders to generate funding premiums.
				{
					Order:      OrderTemplate_Bob_Num0_Id0_Clob0_Sell_LongTerm,
					HumanPrice: "28005",
					HumanSize:  "2",
				},
				{
					Order:      OrderTemplate_Alice_Num0_Id0_Clob0_Buy_LongTerm,
					HumanPrice: "28000",
					HumanSize:  "2",
				},
				// Matched orders to set up positions for Alice, Bob and Carl
				{
					Order:      OrderTemplate_Bob_Num0_Id1_Clob0_Sell_LongTerm,
					HumanPrice: "28003",
					HumanSize:  "1",
				},
				{
					Order:      OrderTemplate_Alice_Num0_Id1_Clob0_Buy_LongTerm,
					HumanPrice: "28003",
					HumanSize:  "0.8",
				},
				{
					Order:      OrderTemplate_Carl_Num0_Id0_Clob0_Buy_LongTerm,
					HumanPrice: "28003",
					HumanSize:  "0.2",
				},
			},
			initialIndexPrice: map[uint32]string{
				TestMarketId: "28002",
			},
			indexPriceForPremium: map[uint32]string{
				TestMarketId: "27960",
			},
			oracelPriceForFundingIndex: map[uint32]string{
				TestMarketId: "27000",
			},
			expectedFundingPremium: 1430, // 28_000 / 27_960 - 1 ~= 0.001430
			// 1430 / 8 * 27000 * 10^(btc_atomic_resolution - quote_atomic_resolution) ~= 483
			expectedFundingIndex: 483,
			expectedSubaccountSettlements: []expectedSettlements{
				{
					SubaccountId: constants.Alice_Num0,
					// Alice is long 0.8 BTC, pays funding
					// 0.00143 / 8 * 27_000 * 0.8 ~= $3.864
					Settlement: -3_864_000,
				},
				{
					SubaccountId: constants.Bob_Num0,
					// Bob is short 1 BTC, receives funding
					// 0.00143 / 8 * 27_000 * 1 ~= $4.83
					Settlement: 4_830_000,
				},
				{
					SubaccountId: constants.Carl_Num0,
					// Carl is long 0.2 BTC, pays funding
					// 0.00143 / 8 * 27_000 * 0.2 ~= $0.966
					Settlement: -966_000,
				},
			},
		},
		// TODO(CORE-712): Add more test cases
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				// UpdateIndexPrice only contacts the tApp.App.Server causing non-determinism in the
				// other App instances in TestApp used for non-determinism checking.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					genesis.GenesisTime = GenesisTime
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			// Place orders on the book.
			for _, testHumanOrder := range tc.testHumanOrders {
				order := testapp.MustMakeOrderFromHumanInput(
					ctx,
					tApp.App,
					testHumanOrder.Order,
					testHumanOrder.HumanPrice,
					testHumanOrder.HumanSize,
				)

				checkTx := testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(order))
				resp := tApp.CheckTx(checkTx[0])
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			// Update initial index price. This price is meant to be within the impact price range,
			// leading to zero sampled premiums.
			pricefeed_testutil.UpdateIndexPrice(
				t,
				ctx,
				tApp.App,
				TestMarketId,
				pricestest.MustHumanPriceToMarketPrice(tc.initialIndexPrice[TestMarketId], -5),
				// Only index price past a certain threshold is used for premium calculation.
				// Use additional buffer here to ensure `test-race` passes.
				time.Now().Add(1*time.Hour),
			)

			ctx = tApp.AdvanceToBlock(BlockHeightAtFirstFundingTick, testapp.AdvanceToBlockOptions{
				BlockTime:                    FirstFundingTick,
				LinearBlockTimeInterpolation: true,
			})

			premiumSamples := tApp.App.PerpetualsKeeper.GetPremiumSamples(ctx)
			// No non-zero premium samples yet.
			require.Len(t, premiumSamples.AllMarketPremiums, 0)
			// Zero premium samples since we just entered a new `funding-tick` epoch.
			require.Equal(t, uint32(0), premiumSamples.NumPremiums)

			// Update index price for each validator so they use this price for premium calculation.
			pricefeed_testutil.UpdateIndexPrice(
				t,
				ctx,
				tApp.App,
				TestMarketId,
				pricestest.MustHumanPriceToMarketPrice(tc.indexPriceForPremium[TestMarketId], -5),
				// Only index price past a certain threshold is used for premium calculation.
				// Use additional buffer here to ensure `test-race` passes.
				time.Now().Add(1*time.Hour),
			)

			// We just entered a new `funding-tick` epoch, there should be 0 funding premium samples.
			require.Equal(t, tApp.App.PerpetualsKeeper.GetPremiumSamples(ctx).NumPremiums, uint32(0))

			// Advance to the end of the last funding-sample epoch during the second funding-tick epoch.
			// At this point, 60 funding-sample epochs have passed, so we should expect 60 premium samples.
			ctx = tApp.AdvanceToBlock(
				testapp.EstimatedHeightForBlockTime(GenesisTime, LastFundingSampleOfSecondFundingTickEpoch, BlockTimeDuration),
				testapp.AdvanceToBlockOptions{
					BlockTime:                    LastFundingSampleOfSecondFundingTickEpoch,
					LinearBlockTimeInterpolation: true,
				})

			premiumSamples = tApp.App.PerpetualsKeeper.GetPremiumSamples(ctx)
			require.Equal(t, 60, int(premiumSamples.NumPremiums))
			expectedAllMarketPremiums := []perptypes.MarketPremiums{
				{
					PerpetualId: 0,
					Premiums:    constants.GenerateConstantFundingPremiums(tc.expectedFundingPremium, 60),
				},
			}
			require.Equal(t, expectedAllMarketPremiums, premiumSamples.AllMarketPremiums)

			// Update index price for each validator so they propose this price as the new oracle price.
			// This price will be used for calculating the funding index at the end of `funding-tick`.
			pricefeed_testutil.UpdateIndexPrice(
				t,
				ctx,
				tApp.App,
				TestMarketId,
				pricestest.MustHumanPriceToMarketPrice(tc.oracelPriceForFundingIndex[TestMarketId], -5),
				// Only index price past a certain threshold is used for premium calculation.
				// Use additional buffer here to ensure `test-race` passes.
				time.Now().Add(1*time.Hour),
			)
			// Advance another 30 seconds to the end of the second funding-tick epoch. This will trigger processing
			// of `funding-tick`, which calculates the final funding rate and updates the funding index.
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight()+NumBlocksPerMinute-1), testapp.AdvanceToBlockOptions{
				BlockTime: SecondFundingTick.Add(-BlockTimeDuration),
			})
			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{
				BlockTime: SecondFundingTick,
			})

			premiumSamples = tApp.App.PerpetualsKeeper.GetPremiumSamples(ctx)
			require.Equal(t, uint32(0), premiumSamples.NumPremiums)
			require.Len(t, premiumSamples.AllMarketPremiums, 0)

			// Check that the funding index is correctly updated.
			btcPerp, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, 0)
			require.NoError(t, err)
			require.Equal(t, tc.expectedFundingIndex, btcPerp.FundingIndex.BigInt().Int64())

			subaccsBeforeSettlement := []satypes.Subaccount{}
			totalUsdcBalanceBeforeSettlement := int64(0)
			for _, expectedSettlements := range tc.expectedSubaccountSettlements {
				subaccBeforeSettlement := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, expectedSettlements.SubaccountId)
				// Before settlement, each perpetual position should have zero funding index, since these positions
				// were opened when BTC perpetual has zero funding idnex.
				require.Equal(t, int64(0), subaccBeforeSettlement.PerpetualPositions[0].FundingIndex.BigInt().Int64())
				subaccsBeforeSettlement = append(subaccsBeforeSettlement, subaccBeforeSettlement)
				totalUsdcBalanceBeforeSettlement += getSubaccountUsdcBalance(subaccBeforeSettlement)
			}

			// Send transfers from Dave to subaccounts that has positions, so that funding is settled for these accounts.
			for _, transfer := range TestTranfers {
				// Invoke CheckTx.
				CheckTx_MsgDepositToSubaccount := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: testtx.MustGetOnlySignerAddress(&transfer),
						Gas:                  100_000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					&transfer,
				)
				resp := tApp.CheckTx(CheckTx_MsgDepositToSubaccount)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight()+1), testapp.AdvanceToBlockOptions{})

			totalUsdcBalanceAfterSettlement := int64(0)
			for i, expectedSettlements := range tc.expectedSubaccountSettlements {
				subaccAfterSettlement := tApp.App.SubaccountsKeeper.GetSubaccount(
					ctx,
					expectedSettlements.SubaccountId,
				)

				// Before settlement, each perpetual position should have zero funding index, since these positions
				// were opened when BTC perpetual has zero funding idnex.
				// TODO(CORE-723): Start with non-zero funding index on the perpetual.
				require.Equal(t,
					tc.expectedFundingIndex,
					subaccAfterSettlement.PerpetualPositions[0].FundingIndex.BigInt().Int64(),
				)
				totalUsdcBalanceAfterSettlement += getSubaccountUsdcBalance(subaccAfterSettlement)

				require.Equal(t,
					getSubaccountUsdcBalance(subaccsBeforeSettlement[i])+expectedSettlements.Settlement,
					getSubaccountUsdcBalance(subaccAfterSettlement)-TestTransferUsdcForSettlement,
					"subaccount id: %v, expected settlement: %v, balance before settlement: %v, "+
						"balance after (minus test transfer): %v",
					expectedSettlements.SubaccountId,
					expectedSettlements.Settlement,
					getSubaccountUsdcBalance(subaccsBeforeSettlement[i]),
					getSubaccountUsdcBalance(subaccAfterSettlement)-TestTransferUsdcForSettlement,
				)
			}

			// Check that the involved subaccounts has the same total balance before and after the transfer
			// (besides transfers from Dave).
			require.Equal(t,
				totalUsdcBalanceBeforeSettlement,
				totalUsdcBalanceAfterSettlement-TestTransferUsdcForSettlement*3,
			)
		})
	}
}
