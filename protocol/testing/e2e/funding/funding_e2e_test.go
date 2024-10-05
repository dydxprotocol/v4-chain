package funding_test

import (
	"testing"
	"time"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	clobtest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	pricefeed_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/pricefeed"
	pricestest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	epochstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	sendingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cometbft/cometbft/types"

	"github.com/stretchr/testify/require"
)

const (
	BlockTimeDuration             = 2 * time.Second
	NumBlocksPerMinute            = int64(time.Minute / BlockTimeDuration) // 30
	BlockHeightAtFirstFundingTick = 1000
	TestTransferTDaiForSettlement = 10_000_000_000_000
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
				AssetId:   assettypes.AssetTDai.Id,
				Amount:    TestTransferTDaiForSettlement,
			},
		},
		{
			Transfer: &sendingtypes.Transfer{
				Sender:    constants.Dave_Num0,
				Recipient: constants.Bob_Num0,
				AssetId:   assettypes.AssetTDai.Id,
				Amount:    TestTransferTDaiForSettlement,
			},
		},
		{
			Transfer: &sendingtypes.Transfer{
				Sender:    constants.Dave_Num0,
				Recipient: constants.Carl_Num0,
				AssetId:   assettypes.AssetTDai.Id,
				Amount:    TestTransferTDaiForSettlement,
			},
		},
	}
)

type expectedSettlements struct {
	SubaccountId satypes.SubaccountId
	Settlement   int64
}

func getSubaccountTDaiBalance(subaccount satypes.Subaccount) int64 {
	return subaccount.AssetPositions[0].Quantums.BigInt().Int64()
}

func TestFunding(t *testing.T) {
	tests := map[string]struct {
		testHumanOrders    []clobtest.TestHumanOrder
		initialDaemonPrice map[uint32]string
		// daemon price to be used in premium calculation
		daemonPriceForPremium map[uint32]string
		// oracle price for funding index calculation
		oracelPriceForFundingIndex map[uint32]string
		// address -> funding
		expectedSubaccountSettlements []expectedSettlements
		expectedFundingPremiums       []perptypes.MarketPremiums
		expectedFundingIndex          int64
	}{
		"daemon price below impact bid, positive funding, longs pay shorts": {
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
			initialDaemonPrice: map[uint32]string{
				TestMarketId: "28002",
			},
			daemonPriceForPremium: map[uint32]string{
				TestMarketId: "27960",
			},
			oracelPriceForFundingIndex: map[uint32]string{
				TestMarketId: "27000",
			},
			expectedFundingPremiums: []perptypes.MarketPremiums{
				{
					PerpetualId: 0,
					// 28_000 / 27_960 - 1 ~= 0.001430
					Premiums: constants.GenerateConstantFundingPremiums(1430, 60),
				},
			},
			// 1430 / 8 * 27000 * 10^(btc_atomic_resolution - quote_atomic_resolution) ~= 482.625
			expectedFundingIndex: 482,
			expectedSubaccountSettlements: []expectedSettlements{
				{
					SubaccountId: constants.Alice_Num0,
					// Alice is long 0.8 BTC, pays funding
					// Theoretical (from funding rate): 0.00143 / 8 * 27_000 * 0.8 ~= $3.864
					// Actual (from funding index): 482 * 8e9 (base quantums) / 1e6 (quote atomic) = $3.856
					Settlement: -3_856_000,
				},
				{
					SubaccountId: constants.Bob_Num0,
					// Bob is short 1 BTC, receives funding
					// Theoretical (from funding rate): 0.00143 / 8 * 27_000 * 1 ~= $4.82625
					// Actual (from funding index): 482 * 1e10 (base quantums) / 1e6 (quote atomic) = $4.82
					Settlement: 4_820_000,
				},
				{
					SubaccountId: constants.Carl_Num0,
					// Carl is long 0.2 BTC, pays funding
					// Theoretical (from funding rate): 0.00143 / 8 * 27_000 * 0.2 ~= $0.96525
					// Actual (from funding index): 482 * 2e9 (base quantums) / 1e6 (quote atomic) = $0.964
					Settlement: -964_000,
				},
			},
		},
		"daemon price above impact ask, negative funding, final funding rate clamped": {
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
			initialDaemonPrice: map[uint32]string{
				0: "28002",
			},
			daemonPriceForPremium: map[uint32]string{
				0: "34000",
			},
			oracelPriceForFundingIndex: map[uint32]string{
				0: "33500",
			},
			expectedFundingPremiums: []perptypes.MarketPremiums{
				{
					PerpetualId: 0,
					// 28005  / 34000 - 1 ~= -0.176323 (-17.6%)
					Premiums: constants.GenerateConstantFundingPremiums(-176_323, 60),
				},
			},
			// Funding rate clamp = premium_rate_clamp_factor * (initial_margin - maintenance_margin)
			//                    = 600% * (0.05 - 0.03) = 12% = 120_000 ppm
			// 120_000 / 8 * 33500 * 10^(btc_atomic_resolution - quote_atomic_resolution) = 50_250
			expectedFundingIndex: -50_250,
			expectedSubaccountSettlements: []expectedSettlements{
				{
					SubaccountId: constants.Alice_Num0,
					// Alice is long 0.8 BTC, receives funding
					// Theoretical (from funding rate): 0.12 / 8 * 33_500 * 0.8 = $402
					// Actual (from funding index): 50_250 * 8e9 (base quantums) / 1e6 (quote atomic) = $402
					Settlement: 402_000_000,
				},
				{
					SubaccountId: constants.Bob_Num0,
					// Bob is short 1 BTC, pays funding
					// Theoretical (from funding rate): 0.12 / 8 * 33_500 * 1 = $502.5
					// Actual (from funding index): 50_250 * 1e10 (base quantums) / 1e6 (quote atomic) = $502.5
					Settlement: -502_500_000,
				},
				{
					SubaccountId: constants.Carl_Num0,
					// Carl is long 0.2 BTC, receives funding
					// Theoretical (from funding rate): 0.12 / 8 * 33_500 * 0.2 = $100.5
					// Actual (from funding index): 50_250 * 2e9 (base quantums) / 1e6 (quote atomic) = $100.5
					Settlement: 100_500_000,
				},
			},
		},
		"daemon price between impact bid and ask, zero funding": {
			testHumanOrders: []clobtest.TestHumanOrder{
				// Unmatched orders to generate funding premiums.
				{
					Order:      OrderTemplate_Bob_Num0_Id0_Clob0_Sell_LongTerm,
					HumanPrice: "28005", // Impact ask price
					HumanSize:  "2",
				},
				{
					Order:      OrderTemplate_Alice_Num0_Id0_Clob0_Buy_LongTerm,
					HumanPrice: "28000", // Impact bid price
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
			initialDaemonPrice: map[uint32]string{
				0: "28002",
			},
			daemonPriceForPremium: map[uint32]string{
				0: "28003", // Between impact bid and ask
			},
			oracelPriceForFundingIndex: map[uint32]string{
				0: "27500",
			},
			expectedFundingPremiums: nil,
			expectedFundingIndex:    0,
			expectedSubaccountSettlements: []expectedSettlements{
				{
					SubaccountId: constants.Alice_Num0,
					Settlement:   0,
				},
				{
					SubaccountId: constants.Bob_Num0,
					Settlement:   0,
				},
				{
					SubaccountId: constants.Carl_Num0,
					Settlement:   0,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				// UpdateDaemonPrice only contacts the tApp.App.Server causing non-determinism in the
				// other App instances in TestApp used for non-determinism checking.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					genesis.GenesisTime = GenesisTime
					return genesis
				}).Build()

			// rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			// rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			// require.NoError(t, conversionErr)

			// tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			// tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			// tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			// tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			// tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			// tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			// tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			// tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			ctx := tApp.InitChain()

			rate := sdaiservertypes.TestSDAIEventRequest.ConversionRate

			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{},
				rate,
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz},
			})

			// ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

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

			// Update initial daemon price. This price is meant to be within the impact price range,
			// leading to zero sampled premiums.
			pricefeed_testutil.UpdateDaemonPrice(
				t,
				ctx,
				tApp.App,
				TestMarketId,
				pricestest.MustHumanPriceToMarketPrice(tc.initialDaemonPrice[TestMarketId], -5),
				// Only daemon price past a certain threshold is used for premium calculation.
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

			// Update daemon price for each validator so they use this price for premium calculation.
			pricefeed_testutil.UpdateDaemonPrice(
				t,
				ctx,
				tApp.App,
				TestMarketId,
				pricestest.MustHumanPriceToMarketPrice(tc.daemonPriceForPremium[TestMarketId], -5),
				// Only daemon price past a certain threshold is used for premium calculation.
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
			require.Equal(t, tc.expectedFundingPremiums, premiumSamples.AllMarketPremiums)

			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight()+NumBlocksPerMinute-1), testapp.AdvanceToBlockOptions{
				BlockTime: SecondFundingTick.Add(-BlockTimeDuration),
			})

			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{
					0: {
						SpotPrice: pricestest.MustHumanPriceToMarketPrice(
							tc.oracelPriceForFundingIndex[0],
							-5,
						),
						PnlPrice: pricestest.MustHumanPriceToMarketPrice(
							tc.oracelPriceForFundingIndex[0],
							-5,
						),
					},
				},
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight())+1, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz},
				BlockTime:          SecondFundingTick,
			})

			premiumSamples = tApp.App.PerpetualsKeeper.GetPremiumSamples(ctx)
			require.Equal(t, uint32(0), premiumSamples.NumPremiums)
			require.Len(t, premiumSamples.AllMarketPremiums, 0)

			// Check that the funding index is correctly updated.
			btcPerp, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, 0)

			require.NoError(t, err)
			require.Equal(t, tc.expectedFundingIndex, btcPerp.FundingIndex.BigInt().Int64())

			subaccsBeforeSettlement := []satypes.Subaccount{}
			totalTDaiBalanceBeforeSettlement := int64(0)
			for _, expectedSettlements := range tc.expectedSubaccountSettlements {
				subaccBeforeSettlement := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, expectedSettlements.SubaccountId)
				// Before settlement, each perpetual position should have zero funding index, since these positions
				// were opened when BTC perpetual has zero funding idnex.
				require.Equal(t, int64(0), subaccBeforeSettlement.PerpetualPositions[0].FundingIndex.BigInt().Int64())
				subaccsBeforeSettlement = append(subaccsBeforeSettlement, subaccBeforeSettlement)
				totalTDaiBalanceBeforeSettlement += getSubaccountTDaiBalance(subaccBeforeSettlement)
			}

			// Send transfers from Dave to subaccounts that has positions, so that funding is settled for these accounts.
			for _, transfer := range TestTranfers {
				// Invoke CheckTx.
				CheckTx_MsgDepositToSubaccount := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: transfer.Transfer.Sender.Owner,
						Gas:                  150_000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					&transfer,
				)
				resp := tApp.CheckTx(CheckTx_MsgDepositToSubaccount)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			ctx = tApp.AdvanceToBlock(uint32(ctx.BlockHeight()+1), testapp.AdvanceToBlockOptions{})

			totalTDaiBalanceAfterSettlement := int64(0)
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
				totalTDaiBalanceAfterSettlement += getSubaccountTDaiBalance(subaccAfterSettlement)

				require.Equal(t,
					getSubaccountTDaiBalance(subaccsBeforeSettlement[i])+expectedSettlements.Settlement,
					getSubaccountTDaiBalance(subaccAfterSettlement)-TestTransferTDaiForSettlement,
					"subaccount id: %v, expected settlement: %v, got settlement: %v,"+
						"balance before settlement: %v, balance after (minus test transfer): %v",
					expectedSettlements.SubaccountId,
					expectedSettlements.Settlement,
					getSubaccountTDaiBalance(subaccAfterSettlement)-TestTransferTDaiForSettlement-
						getSubaccountTDaiBalance(subaccsBeforeSettlement[i]),
					getSubaccountTDaiBalance(subaccsBeforeSettlement[i]),
					getSubaccountTDaiBalance(subaccAfterSettlement)-TestTransferTDaiForSettlement,
				)
			}

			// Check that the involved subaccounts has the same total balance before and after the transfer
			// (besides transfers from Dave).
			require.Equal(t,
				totalTDaiBalanceBeforeSettlement,
				totalTDaiBalanceAfterSettlement-TestTransferTDaiForSettlement*3,
			)
		})
	}
}
