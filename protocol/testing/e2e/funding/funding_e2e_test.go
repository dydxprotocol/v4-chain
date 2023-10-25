package funding_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricefeed_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

type TestHumanOrder struct {
	Order      clobtypes.Order
	HumanPrice string
	HumanSize  string
}

const (
	BlockTimeDuration             = 2 * time.Second
	NumBlocksPerMinute            = int64(time.Minute / BlockTimeDuration) // 30
	BlockHeightAtFirstFundingTick = 1000
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
	// Genesis time of the chain
	GenesisTime                               = time.Unix(1690000000, 0)
	FirstFundingSampleTick                    = time.Unix(1690000050, 0)
	FirstFundingTick                          = time.Unix(1690002000, 0)
	LastFundingSampleOfSecondFundingTickEpoch = time.Unix(1690005570, 0)
	SecondFundingTick                         = time.Unix(
		1690002000+int64(epochstypes.FundingTickEpochDuration),
		0,
	)
)

func TestFunding(t *testing.T) {
	tests := map[string]struct {
		testHumanOrders   []TestHumanOrder
		initialIndexPrice map[uint32]string
		// index price to be used in premium calculation
		indexPriceForPremium map[uint32]string
		// oracle price for funding index calculation
		oracelPriceForFundingIndex map[uint32]string
		expectedFundingPremium     int32
		expectedFundingIndex       int64
	}{
		"Test funding": {
			testHumanOrders: []TestHumanOrder{
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
				// Matched orders to set up Alice and Bob's positions.
				{
					Order:      OrderTemplate_Bob_Num0_Id1_Clob0_Sell_LongTerm,
					HumanPrice: "28003",
					HumanSize:  "1",
				},
				{
					Order:      OrderTemplate_Alice_Num0_Id1_Clob0_Buy_LongTerm,
					HumanPrice: "28003",
					HumanSize:  "1",
				},
			},
			initialIndexPrice: map[uint32]string{
				0: "28002",
			},
			indexPriceForPremium: map[uint32]string{
				0: "27960",
			},
			oracelPriceForFundingIndex: map[uint32]string{
				0: "27000",
			},
			expectedFundingPremium: 1430, // 28_000 / 27_960 - 1 ~= 0.001430
			// 1430 / 8 * 27000 * 10^(btc_atomic_resolution - quote_atomic_resolution) ~= 483
			expectedFundingIndex: 483,
		},
		// TODO(CORE-712): Add more test cases
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
				pricestest.MustHumanPriceToMarketPrice(tc.initialIndexPrice[0], -5),
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
				pricestest.MustHumanPriceToMarketPrice(tc.indexPriceForPremium[0], -5),
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
				pricestest.MustHumanPriceToMarketPrice(tc.oracelPriceForFundingIndex[0], -5),
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
			// TODO(CORE-703): Settle Alice and Bob's positions, so we can measure that the funding payment as processed.
		})
	}
}
