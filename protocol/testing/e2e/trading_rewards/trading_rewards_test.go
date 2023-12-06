package trading_rewards_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricefeed_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	vesttypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

const (
	BlockTimeDuration             = 2 * time.Second
	NumBlocksPerMinute            = int64(time.Minute / BlockTimeDuration) // 30
	BlockHeightAtFirstFundingTick = 1000
	TestRewardsTokenMarketId      = 30
	TestRewardsTokenPriceExponent = -8
	TestBtcMarketId               = 0
	TestEthMarketId               = 1
	TestBtcPriceExponent          = -5
	TestEthPriceExponent          = -6
	GTBLimit                      = 20
)

var (
	TestRewardsVestEntry = vesttypes.VestEntry{
		VesterAccount:   rewardstypes.VesterAccountName,
		TreasuryAccount: rewardstypes.TreasuryAccountName,
		Denom:           lib.DefaultBaseDenom,
		StartTime:       TestRewardsVestStartTime,
		EndTime:         TestRewardsVestEndTime,
	}
	OrderTemplate_ShortTerm_Btc = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side: clobtypes.Order_SIDE_BUY,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: 0,
		},
	}
	// Genesis time of the chain
	GenesisTime               = time.Unix(1696132500, 0) // Sun Oct 01 2023 03:55:00 GMT+0000
	TestRewardsVestStartTime  = time.Unix(1696132800, 0) // Sun Oct 01 2023 04:00:00 GMT+0000
	TestRewardsVestEndTime    = time.Unix(1853985600, 0) // Sun Oct 01 2028 04:00:00 GMT+0000
	RewardsVesterAccAddress   = authtypes.NewModuleAddress(rewardstypes.VesterAccountName)
	RewardsTreasuryAccAddress = authtypes.NewModuleAddress(rewardstypes.TreasuryAccountName)
	HeightAtVestStart         = testapp.EstimatedHeightForBlockTime(
		GenesisTime,
		TestRewardsVestStartTime,
		BlockTimeDuration,
	)
	TestAccountStartingTokenBalance = big_testutil.Int64MulPow10(5, 23)
)

type expectedBalancesAtBlock struct {
	Height           int64
	ExpectedBalances []expectedBalance
}

type expectedBalance struct {
	AccAddress sdk.AccAddress
	Balance    *big.Int
}

type TestHumanOraclePrice struct {
	MarketId      uint32
	PriceExponent int32
	HumanPrice    string
}

func TestTradingRewards(t *testing.T) {
	tests := map[string]struct {
		// nth block after vesting starts -> orders placed during that block, in human readable form
		testHumanOrders          map[int64][]clobtest.TestHumanOrder
		vestEntries              []vesttypes.VestEntry
		rewardsParams            rewardstypes.Params
		humanOraclePrices        []TestHumanOraclePrice
		initRewardsVesterBalance *big.Int
		// nth block after vesting starts -> expectedBalance
		expectedBalances []expectedBalancesAtBlock
		// // nth block after vesting starts -> expected balance of rewards treasury
		// expectedRewardsTreasuryBalances map[int]*big.Int
	}{
		"Every block, only one taker account gets rewards": {
			testHumanOrders: map[int64][]clobtest.TestHumanOrder{
				2: {
					// Bob BTC maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Bob_Num0),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+2),
						),
						HumanPrice: "28003",
						HumanSize:  "1",
					},
					// Alice BTC taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Alice_Num0),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+2),
						), HumanPrice: "28003",
						HumanSize: "1",
					},
				},
				13: {
					// Alice BTC maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Alice_Num0),
							clobtest.WithClientId(1),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+13),
						), HumanPrice: "28003",
						HumanSize: "1",
					},
					// Bob BTC taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Bob_Num0),
							clobtest.WithClientId(1),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+13),
						),
						HumanPrice: "28003",
						HumanSize:  "1",
					},
				},
			},
			vestEntries: []vesttypes.VestEntry{TestRewardsVestEntry},
			humanOraclePrices: []TestHumanOraclePrice{
				{
					MarketId:      TestRewardsTokenMarketId,
					PriceExponent: TestRewardsTokenPriceExponent,
					HumanPrice:    "1.95",
				},
				{
					MarketId:      TestBtcMarketId,
					PriceExponent: TestBtcPriceExponent,
					HumanPrice:    "28003",
				},
			},
			rewardsParams: rewardstypes.Params{
				TreasuryAccount:  rewardstypes.TreasuryAccountName,
				Denom:            lib.DefaultBaseDenom,
				DenomExponent:    lib.BaseDenomExponent,
				MarketId:         TestRewardsTokenMarketId,
				FeeMultiplierPpm: 990_000, // 99%
			},
			expectedBalances: []expectedBalancesAtBlock{
				{
					Height: 0,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// 200 million full coins
							Balance: big_testutil.Int64MulPow10(200_000_000, 18),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// 0 full coins
							Balance: big.NewInt(0),
						},
					},
				},
				{
					Height: 1,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// ~199999997.47 full coins
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"199999997465993634576010055",
								10,
							)),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// ~2.53 full coins
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"2534006365423989945",
								10,
							)),
						},
					},
				},
				{
					Height: 2,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// ~199999994.93 full coins
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"199999994931987269152020110",
								10,
							)),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// Total of ~5.06 full coins have vested, which is less than calculated
							// rewards (~5.5 full coins). So all reward tokens were distributed.
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"0",
								10,
							)),
						},
						{
							AccAddress: constants.AliceAccAddress,
							// starting balance + ~5.06 full coins rewards
							Balance: new(big.Int).Add(
								TestAccountStartingTokenBalance,
								big_testutil.MustFirst(new(big.Int).SetString(
									"5068012730847979890",
									10,
								)),
							),
						},
						{
							AccAddress: constants.BobAccAddress,
							// starting balance, no rewards
							Balance: TestAccountStartingTokenBalance,
						},
					},
				},
				{
					Height: 12,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// ~199999969.59 full coins
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"199999969591923614912120660",
								10,
							)),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// ~25.34 full coins. Note this is exactly 10x the amount vested per block,
							// since 10 blocks has passed since the last check.
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"25340063654239899450",
								10,
							)),
						},
					},
				},
				{
					Height: 13,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// ~199999967.05 full coins
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"199999967057917249488130715",
								10,
							)),
						},
						{
							AccAddress: constants.BobAccAddress,
							// Starting balance: 500000000000000000000000
							// Total rewards = (TakerFee - TakerVolume * MaxMakerRebate) * 0.99
							//               = ($28003 * 0.05% - $28003 * 0.011%) * 0.99
							//               = ($14.0015 - $3.08033) 0.99 = $10.8119583
							// Reward tokens = $10.8119583 / $1.95 = 5.544594 full coins
							Balance: new(big.Int).Add(
								TestAccountStartingTokenBalance,
								big_testutil.MustFirst(new(big.Int).SetString(
									"5544594000000000000",
									10,
								)),
							),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// 25.34 + 2.534 - 5.544594 ~= 22.329 full coins
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"22329476019663889395",
								10,
							)),
						},
					},
				},
			},
			initRewardsVesterBalance: big_testutil.Int64MulPow10(200_000_000, 18), // 200 million full coins
		},
		"Multiple accounts gets rewards": {
			testHumanOrders: map[int64][]clobtest.TestHumanOrder{
				10: {
					// Bob BTC maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Bob_Num0),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "28003",
						HumanSize:  "2",
					},
					// Alice BTC taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Alice_Num0),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						), HumanPrice: "28003",
						HumanSize: "2",
					},
					// Alice BTC maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Alice_Num0),
							clobtest.WithClientId(1),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						), HumanPrice: "28003",
						HumanSize: "2",
					},
					// Bob BTC taker order
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Bob_Num0),
							clobtest.WithClientId(1),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "28003",
						HumanSize:  "2",
					},
					// Carl ETH maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Carl_Num0),
							clobtest.WithClobPairid(TestEthMarketId),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "1605",
						HumanSize:  "20",
					},
					// Dave ETH taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Dave_Num0),
							clobtest.WithClobPairid(TestEthMarketId),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "1605",
						HumanSize:  "20",
					},
					// Dave ETH maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Dave_Num0),
							clobtest.WithClobPairid(TestEthMarketId),
							clobtest.WithClientId(1),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "1605",
						HumanSize:  "20",
					},
					// Carl ETH taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Carl_Num0),
							clobtest.WithClobPairid(TestEthMarketId),
							clobtest.WithClientId(1),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "1605",
						HumanSize:  "20",
					},
				},
			},
			vestEntries: []vesttypes.VestEntry{TestRewardsVestEntry},
			humanOraclePrices: []TestHumanOraclePrice{
				{
					MarketId:      TestRewardsTokenMarketId,
					PriceExponent: TestRewardsTokenPriceExponent,
					HumanPrice:    "1.95",
				},
				{
					MarketId:      TestBtcMarketId,
					PriceExponent: TestBtcPriceExponent,
					HumanPrice:    "28003",
				},
				{
					MarketId:      TestEthMarketId,
					PriceExponent: TestEthPriceExponent,
					HumanPrice:    "1605",
				},
			},
			rewardsParams: rewardstypes.Params{
				TreasuryAccount:  rewardstypes.TreasuryAccountName,
				Denom:            lib.DefaultBaseDenom,
				DenomExponent:    lib.BaseDenomExponent,
				MarketId:         TestRewardsTokenMarketId,
				FeeMultiplierPpm: 990_000, // 99%
			},
			expectedBalances: []expectedBalancesAtBlock{
				{
					Height: 0,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// 200 million full coins
							Balance: big_testutil.Int64MulPow10(200_000_000, 18),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// 0 full coins
							Balance: big.NewInt(0),
						},
					},
				},
				{
					Height: 10,
					// Net fees:
					// - Alice and Bob: $21.842_340
					// - Carl and Dave: $12.519
					// Total rewards tokens distributed: ~25.34 (less than the value of net fees)
					// Entitled reward tokens:
					// - Alice and Bob: 8.0539
					// - Carl and Dave: 4.616
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// ~199999974.659 full coins, since ~25.34 full coins have vested
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"199999974659936345760100550",
								10,
							)),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// All vested rewards were distributed, only rounding dusts left.
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"2",
								10,
							)),
						},
						{
							AccAddress: constants.AliceAccAddress,
							Balance: new(big.Int).Add(
								TestAccountStartingTokenBalance,
								big_testutil.MustFirst(new(big.Int).SetString(
									"8053910091363583686",
									10,
								)),
							),
						},
						{
							AccAddress: constants.BobAccAddress,
							Balance: new(big.Int).Add(
								TestAccountStartingTokenBalance,
								big_testutil.MustFirst(new(big.Int).SetString(
									"8053910091363583686",
									10,
								)),
							),
						},
						{
							AccAddress: constants.CarlAccAddress,
							Balance: new(big.Int).Add(
								TestAccountStartingTokenBalance,
								big_testutil.MustFirst(new(big.Int).SetString(
									"4616121735756366038",
									10,
								)),
							),
						},
						{
							AccAddress: constants.DaveAccAddress,
							Balance: new(big.Int).Add(
								TestAccountStartingTokenBalance,
								big_testutil.MustFirst(new(big.Int).SetString(
									"4616121735756366038",
									10,
								)),
							),
						},
					},
				},
			},
			initRewardsVesterBalance: big_testutil.Int64MulPow10(200_000_000, 18), // 200 million full coins
		},
		"rewards fee multiplier = 0, no rewards are distributed": {
			testHumanOrders: map[int64][]clobtest.TestHumanOrder{
				10: {
					// Bob BTC maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Bob_Num0),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "28003",
						HumanSize:  "2",
					},
					// Alice BTC taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Alice_Num0),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						), HumanPrice: "28003",
						HumanSize: "2",
					},
					// Carl ETH maker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_BUY),
							clobtest.WithSubaccountId(constants.Carl_Num0),
							clobtest.WithClobPairid(TestEthMarketId),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "1605",
						HumanSize:  "20",
					},
					// Dave ETH taker order.
					{
						Order: clobtest.GenerateOrderUsingTemplate(
							OrderTemplate_ShortTerm_Btc,
							clobtest.WithSide(clobtypes.Order_SIDE_SELL),
							clobtest.WithSubaccountId(constants.Dave_Num0),
							clobtest.WithClobPairid(TestEthMarketId),
							clobtest.WithGTB(HeightAtVestStart+GTBLimit+10),
						),
						HumanPrice: "1605",
						HumanSize:  "20",
					},
				},
			},
			vestEntries: []vesttypes.VestEntry{TestRewardsVestEntry},
			humanOraclePrices: []TestHumanOraclePrice{
				{
					MarketId:      TestRewardsTokenMarketId,
					PriceExponent: TestRewardsTokenPriceExponent,
					HumanPrice:    "1.95",
				},
				{
					MarketId:      TestBtcMarketId,
					PriceExponent: TestBtcPriceExponent,
					HumanPrice:    "28003",
				},
				{
					MarketId:      TestEthMarketId,
					PriceExponent: TestEthPriceExponent,
					HumanPrice:    "1605",
				},
			},
			rewardsParams: rewardstypes.Params{
				TreasuryAccount:  rewardstypes.TreasuryAccountName,
				Denom:            lib.DefaultBaseDenom,
				DenomExponent:    lib.BaseDenomExponent,
				MarketId:         TestRewardsTokenMarketId,
				FeeMultiplierPpm: 0, // 0%
			},
			expectedBalances: []expectedBalancesAtBlock{
				{
					Height: 0,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// 200 million full coins
							Balance: big_testutil.Int64MulPow10(200_000_000, 18),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// 0 full coins
							Balance: big.NewInt(0),
						},
					},
				},
				{
					Height: 10,
					ExpectedBalances: []expectedBalance{
						{
							AccAddress: RewardsVesterAccAddress,
							// ~199999974.659 full coins, since ~25.34 full coins have vested
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"199999974659936345760100550",
								10,
							)),
						},
						{
							AccAddress: RewardsTreasuryAccAddress,
							// No rewards were distributed. ~25.34 full coins have vested.
							Balance: big_testutil.MustFirst(new(big.Int).SetString(
								"25340063654239899450",
								10,
							)),
						},
						{
							AccAddress: constants.AliceAccAddress,
							Balance:    TestAccountStartingTokenBalance,
						},
						{
							AccAddress: constants.BobAccAddress,
							Balance:    TestAccountStartingTokenBalance,
						},
						{
							AccAddress: constants.CarlAccAddress,
							Balance:    TestAccountStartingTokenBalance,
						},
						{
							AccAddress: constants.DaveAccAddress,
							Balance:    TestAccountStartingTokenBalance,
						},
					},
				},
			},
			initRewardsVesterBalance: big_testutil.Int64MulPow10(200_000_000, 18), // 200 million full coins
		},
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
					// Initialize sender module with its initial balance.
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *banktypes.GenesisState) {
							genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
								Address: RewardsTreasuryAccAddress.String(),
								Coins: sdk.Coins{
									sdk.NewCoin(lib.DefaultBaseDenom, sdk.NewInt(0)),
								},
							})
							genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
								Address: RewardsVesterAccAddress.String(),
								Coins: sdk.Coins{
									sdk.NewCoin(lib.DefaultBaseDenom, sdk.NewIntFromBigInt(
										tc.initRewardsVesterBalance,
									)),
								},
							})
						},
					)
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *vesttypes.GenesisState) {
							genesisState.VestEntries = tc.vestEntries
						},
					)
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *rewardstypes.GenesisState) {
							genesisState.Params = tc.rewardsParams
						},
					)
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			for _, humanOraclePrice := range tc.humanOraclePrices {
				pricefeed_testutil.UpdateIndexPrice(
					t,
					ctx,
					tApp.App,
					humanOraclePrice.MarketId,
					pricestest.MustHumanPriceToMarketPrice(
						humanOraclePrice.HumanPrice,
						humanOraclePrice.PriceExponent,
					),
					// Only index price past a certain threshold is used for premium calculation.
					// Use additional buffer here to ensure `test-race` passes.
					time.Now().Add(1*time.Hour),
				)
			}

			// Iterate through blocks that have expected states.
			for _, expectedBalancesAtBlock := range tc.expectedBalances {
				nthBlockAfterVest := expectedBalancesAtBlock.Height
				// If there are orders for this block, place them.
				if orders, exists := tc.testHumanOrders[nthBlockAfterVest]; exists {
					// Advance to the block before the block we want to place orders on,
					// to make sure orders at placed at the correct block.
					targetHeight := int64(HeightAtVestStart) + nthBlockAfterVest - 1
					if ctx.BlockHeight() < targetHeight {
						ctx = tApp.AdvanceToBlock(uint32(targetHeight), testapp.AdvanceToBlockOptions{
							BlockTime: TestRewardsVestStartTime.Add(
								BlockTimeDuration * time.Duration(nthBlockAfterVest-1),
							),
							LinearBlockTimeInterpolation: true,
						})
					}
					// Place orders on the book.
					for _, testHumanOrder := range orders {
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
				}

				// Advance to target `nthBlockSinceVest` after vesting starts.
				ctx = tApp.AdvanceToBlock(testapp.EstimatedHeightForBlockTime(
					GenesisTime,
					TestRewardsVestStartTime.Add(
						BlockTimeDuration*time.Duration(nthBlockAfterVest),
					),
					BlockTimeDuration,
				), testapp.AdvanceToBlockOptions{
					BlockTime: TestRewardsVestStartTime.Add(
						BlockTimeDuration * time.Duration(nthBlockAfterVest),
					),
					LinearBlockTimeInterpolation: true,
				})

				for _, expectedBalance := range expectedBalancesAtBlock.ExpectedBalances {
					gotBalance := tApp.App.BankKeeper.GetBalance(
						ctx,
						expectedBalance.AccAddress,
						lib.DefaultBaseDenom,
					).Amount.BigInt()
					require.Equal(t,
						expectedBalance.Balance,
						gotBalance,
						"unexpected balance for address %s at %d block since vest; expected %s, got %s",
						expectedBalance.AccAddress.String(),
						nthBlockAfterVest,
						expectedBalance.Balance,
						gotBalance,
					)
				}
			}
		})
	}
}
