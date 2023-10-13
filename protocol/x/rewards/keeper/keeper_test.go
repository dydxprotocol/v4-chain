package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

const (
	TestAddress1         = "dydx16h7p7f4dysrgtzptxx2gtpt5d8t834g9dj830z"
	TestAddress2         = "dydx168pjt8rkru35239fsqvz7rzgeclakp49zx3aum"
	TestAddress3         = "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70"
	TestRewardTokenDenom = "test-denom"
)

var (
	ZeroTreasuryAccountBalance = banktypes.Balance{
		Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
		Coins: []sdk.Coin{{
			Denom:  TestRewardTokenDenom,
			Amount: sdkmath.NewInt(0),
		}},
	}
)

func TestRewardShareStorage_DefaultValue(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	require.Equal(t,
		types.RewardShare{
			Address: TestAddress1,
			Weight:  dtypes.NewInt(0),
		},
		k.GetRewardShare(ctx, TestAddress1),
	)
}

func TestRewardShareStorage_Exists(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	val := types.RewardShare{
		Address: TestAddress1,
		Weight:  dtypes.NewInt(12_345_678),
	}

	err := k.SetRewardShare(ctx, val)
	require.NoError(t, err)
	require.Equal(t, val, k.GetRewardShare(ctx, TestAddress1))
}

func TestSetRewardShare_FailsWithNonpositiveWeight(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	val := types.RewardShare{
		Address: TestAddress1,
		Weight:  dtypes.NewInt(0),
	}

	err := k.SetRewardShare(ctx, val)
	require.ErrorContains(t, err, "Invalid weight 0: weight must be positive")
}

func TestAddRewardShareToAddress(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()

	tests := map[string]struct {
		prevRewardShare     *types.RewardShare // nil if no previous share
		newWeight           *big.Int
		expectedRewardShare types.RewardShare
		expectedErr         error
	}{
		"no previous share": {
			prevRewardShare: nil,
			newWeight:       big.NewInt(12_345_678),
			expectedRewardShare: types.RewardShare{
				Address: TestAddress1,
				Weight:  dtypes.NewInt(12_345_678),
			},
		},
		"with previous share": {
			prevRewardShare: &types.RewardShare{
				Address: TestAddress1,
				Weight:  dtypes.NewInt(100_000),
			},
			newWeight: big.NewInt(500),
			expectedRewardShare: types.RewardShare{
				Address: TestAddress1,
				Weight:  dtypes.NewInt(100_500),
			},
		},
		"fails with non-positive weight": {
			newWeight:   big.NewInt(0),
			expectedErr: fmt.Errorf("Invalid weight 0: weight must be positive"),
			expectedRewardShare: types.RewardShare{
				Address: TestAddress1,
				Weight:  dtypes.NewInt(0),
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp.Reset()
			ctx := tApp.InitChain()
			k := tApp.App.RewardsKeeper

			if tc.prevRewardShare != nil {
				err := k.SetRewardShare(ctx, *tc.prevRewardShare)
				require.NoError(t, err)
			}

			err := k.AddRewardShareToAddress(ctx, TestAddress1, tc.newWeight)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

			// Check the new reward share.
			require.Equal(t, tc.expectedRewardShare, k.GetRewardShare(ctx, TestAddress1))
		})
	}
}

func TestAddRewardSharesForFill(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	makerAddress := TestAddress1
	takerAdderss := TestAddress2

	tests := map[string]struct {
		prevTakerRewardShare *types.RewardShare
		prevMakerRewardShare *types.RewardShare
		fillQuoteQuantums    *big.Int
		takerFeeQuantums     *big.Int
		makerFeeQuantums     *big.Int
		feeTiers             []*feetierstypes.PerpetualFeeTier

		expectedTakerShare types.RewardShare
		expectedMakerShare types.RewardShare
	}{
		"positive maker fee, positive taker fees reduced by maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fillQuoteQuantums:    big.NewInt(800_000_000), // $800
			takerFeeQuantums:     big.NewInt(2_000_000),   // $2
			makerFeeQuantums:     big.NewInt(1_000_000),   // $1
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAdderss,
				Weight:  dtypes.NewInt(1_200_000), // 2 - 0.1% * 800
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(1_000_000),
			},
		},
		"negative maker fee, positive taker fees reduced by 0.1% maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fillQuoteQuantums:    big.NewInt(750_000_000), // $750
			takerFeeQuantums:     big.NewInt(2_000_000),   // $2
			makerFeeQuantums:     big.NewInt(-1_000_000),  // $1
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAdderss,
				Weight:  dtypes.NewInt(1_250_000), // 2 - 0.1% * 750
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(0),
			},
		},
		"negative maker fee, positive taker fees reduced by 0.05% maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fillQuoteQuantums:    big.NewInt(750_000_000), // $750
			takerFeeQuantums:     big.NewInt(2_000_000),   // $2
			makerFeeQuantums:     big.NewInt(-1_000_000),  // $1
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -500,  // -0.05%
					TakerFeePpm: 2_000, // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAdderss,
				Weight:  dtypes.NewInt(1_625_000), // 2 - 0.05% * 750
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(0),
			},
		},
		"positive maker fee, positive taker fees offset by maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fillQuoteQuantums:    big.NewInt(750_000_000), // $750
			takerFeeQuantums:     big.NewInt(700_000),     // $0.7
			makerFeeQuantums:     big.NewInt(500_000),     // $1
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAdderss,
				Weight:  dtypes.NewInt(0), // $0.7 - $750 * 0.1% < 0
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(500_000),
			},
		},
		"positive maker fee, positive taker fees, no maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fillQuoteQuantums:    big.NewInt(750_000_000), // $750
			takerFeeQuantums:     big.NewInt(700_000),     // $0.7
			makerFeeQuantums:     big.NewInt(500_000),     // $1
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: 1_000, // 0.1%
					TakerFeePpm: 2_000, // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAdderss,
				Weight:  dtypes.NewInt(700_000),
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(500_000),
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp.Reset()
			ctx := tApp.InitChain()
			k := tApp.App.RewardsKeeper

			feeTiersKeeper := tApp.App.FeeTiersKeeper
			err := feeTiersKeeper.SetPerpetualFeeParams(ctx, feetierstypes.PerpetualFeeParams{
				Tiers: tc.feeTiers,
			})
			require.NoError(t, err)

			if tc.prevTakerRewardShare != nil {
				err := k.SetRewardShare(ctx, *tc.prevTakerRewardShare)
				require.NoError(t, err)
			}
			if tc.prevMakerRewardShare != nil {
				err := k.SetRewardShare(ctx, *tc.prevMakerRewardShare)
				require.NoError(t, err)
			}

			k.AddRewardSharesForFill(
				ctx,
				takerAdderss,
				makerAddress,
				tc.fillQuoteQuantums,
				tc.takerFeeQuantums,
				tc.makerFeeQuantums,
			)

			// Check the new reward shares.
			require.Equal(t, tc.expectedTakerShare, k.GetRewardShare(ctx, takerAdderss))
			require.Equal(t, tc.expectedMakerShare, k.GetRewardShare(ctx, makerAddress))
		})
	}
}

func TestProcessRewardsForBlock(t *testing.T) {
	testRewardTokenMarketId := uint32(33)
	testRewardTokenMarket := "test-market"
	// TODO(CORE-645): Update test to -18 denom for consistency with prod.
	TestRewardTokenDenomExp := int32(-6)

	tokenPrice2Usdc := pricestypes.MarketPrice{
		Id:       testRewardTokenMarketId,
		Price:    200_000_000, // 2$ per full coin.
		Exponent: -8,
	}

	tokenPrice1_18Usdc := pricestypes.MarketPrice{
		Id:       testRewardTokenMarketId,
		Price:    118_000_000, // 1.18$ per full coin.
		Exponent: -8,
	}

	tests := map[string]struct {
		rewardShares           []types.RewardShare
		tokenPrice             pricestypes.MarketPrice
		treasuryAccountBalance sdkmath.Int
		feeMultiplierPpm       uint32
		expectedBalances       []banktypes.Balance
	}{
		"zero reward share, no change in treasury balance": {
			rewardShares:           []types.RewardShare{},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(1_000_000_000), // 1000 full coins
			feeMultiplierPpm:       1_000_000,                     // 100%
			// 1$ / 2$ * 100% = 0.5 full coin, all paid to TestAddress1
			expectedBalances: []banktypes.Balance{
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1_000_000_000),
					}},
				},
			},
		},
		"one reward share, enough treasury balance": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(1_000_000), // $1 weight of fee
				},
			},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(1_000_000_000), // 1000 full coins
			feeMultiplierPpm:       1_000_000,                     // 100%
			// 1$ / 2$ * 100% = 0.5 full coin, all paid to TestAddress1
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(500_000),
					}},
				},
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(999_500_000), // 999.5 full coins
					}},
				},
			},
		},
		"one reward share, enough treasury balance, 0.99 fee multiplier": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(1_000_000), // $1 weight of fee
				},
			},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(1_000_000_000), // 1000 full coins
			feeMultiplierPpm:       950_000,                       // 95%
			// 1$ / 2$ * 95% = 0.475 full coin, all paid to TestAddress1
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(475_000),
					}},
				},
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(999_525_000), // 999.525 full coins
					}},
				},
			},
		},
		"one reward share, not enough treasury balance": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(1_000_000), // $1 weight of fee
				},
			},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(200_000), // 0.2 full coin
			feeMultiplierPpm:       1_000_000,               // 100%
			// 1$ / 2$ * 100% = 0.5 full coin > 0.2 full coin. Pay 0.2 full coin to TestAddress1.
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(200_000),
					}},
				},
				ZeroTreasuryAccountBalance, // No balance left in treasury.
			},
		},
		"one reward share, zero treasury balance": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(1_000_000), // $1 weight of fee
				},
			},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(0),
			feeMultiplierPpm:       1_000_000, // 100%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(0),
					}}, // No balance to pay out to TestAddress1.
				},
				ZeroTreasuryAccountBalance,
			},
		},
		"three reward shares, enough treasury balance, fee multipler = 0.99, realistic numbers": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(1_025_590_000), // $1025.59 weight of fee
				},
				{
					Address: TestAddress2,
					Weight:  dtypes.NewInt(2_021_300_000), // $2021.3 weight of fee
				},
				{
					Address: TestAddress3,
					Weight:  dtypes.NewInt(835_660_000), // $835.66 weight of fee
				},
			},
			tokenPrice: tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(2_000_123, 18), //~2_000_123 full coin.
			), // 1000 full coins
			feeMultiplierPpm: 990_000, // 99%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(507_667_050), // $1 weight / $2 price * 99% = 0.495 full coin
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1_000_543_500), // $2021.3 weight / $2 price * 99% ~= 1000 full coin
					}},
				},
				{
					Address: TestAddress3,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(413_651_700), // $835.66 weight / $2 price * 99% ~= 413 full coin
					}},
				},
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("2000122999999998078137750", 10)),
						), // ~2_000_122.9 full coins
					}},
				},
			},
		},
		"three reward shares, not enough treasury balance": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(10_000_000), // $10 weight of fee
				},
				{
					Address: TestAddress2,
					Weight:  dtypes.NewInt(20_000_000), // $20 weight of fee
				},
				{
					Address: TestAddress3,
					Weight:  dtypes.NewInt(30_000_000), // $30 weight of fee
				},
			},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(10_000_000), // 10 full coins
			feeMultiplierPpm:       1_000_000,                  // 100%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1_666_666), // 1/6 of 10 = 1.666666 full coins
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(3_333_333), // 1/3 of 10 = 3.333333 full coins
					}},
				},
				{
					Address: TestAddress3,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(5_000_000), // 1/2 of 10 = 5 full coins
					}},
				},
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1), // 0.000001 full coins left due to rounding
					}},
				},
			},
		},
		"three reward shares, not enough treasury balance, $1.18 token price, 0.99 fee multiplier": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(125_560_000), // $125.56 weight of fee (~56.72% of total weight)
				},
				{
					Address: TestAddress2,
					Weight:  dtypes.NewInt(500_000), // $0.5 weight of fee (~0.23% of total weight)
				},
				{
					Address: TestAddress3,
					Weight:  dtypes.NewInt(95_300_000), // $95.3 weight of fee (~43.05% of total weight)
				},
			},
			tokenPrice:             tokenPrice1_18Usdc,
			treasuryAccountBalance: sdkmath.NewInt(100_000_000), // 100 full coins
			feeMultiplierPpm:       990_000,                     // 99%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(56_722081), // 56.722081 full coins
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(225_876), // 0.225876 full coin
					}},
				},
				{
					Address: TestAddress3,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(43_052_041), // 43.052041 full coins
					}},
				},
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(2), // 0.000002 full coins left due to rounding
					}},
				},
			},
		},
		"2 reward shares, one address reward was rounded to 0, fee multipler = 0.99": {
			rewardShares: []types.RewardShare{
				{
					Address: TestAddress1,
					Weight:  dtypes.NewInt(100_000_000), // $100 weight of fee
				},
				{
					Address: TestAddress2,
					Weight:  dtypes.NewInt(1), // $0.000001 weight of fee
				},
			},
			tokenPrice:             tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewInt(1_000), // 0.001 full coins
			feeMultiplierPpm:       990_000,               // 0.99
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(999),
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(0), // rounded to 0
					}},
				},
				{
					Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1),
					}},
				},
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set up treasury account balance in genesis state
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: authtypes.NewModuleAddress(types.TreasuryAccountName).String(),
							Coins: []sdk.Coin{
								sdk.NewCoin(TestRewardTokenDenom, tc.treasuryAccountBalance),
							},
						})
					},
				)
				return genesis
			}).WithTesting(t).Build()

			tApp.Reset()
			ctx := tApp.InitChain()
			k := tApp.App.RewardsKeeper

			// Set up PricesKeeper
			_, err := tApp.App.PricesKeeper.CreateMarket(
				ctx,
				pricestypes.MarketParam{
					Id:                 testRewardTokenMarketId,
					Pair:               testRewardTokenMarket,
					Exponent:           tc.tokenPrice.Exponent,
					MinExchanges:       uint32(1),
					MinPriceChangePpm:  uint32(50),
					ExchangeConfigJson: "{}",
				},
				tc.tokenPrice,
			)
			require.NoError(t, err)

			// Set up RewardsKeeper
			err = k.SetParams(
				ctx,
				types.Params{
					TreasuryAccount:  types.TreasuryAccountName,
					Denom:            TestRewardTokenDenom,
					DenomExponent:    TestRewardTokenDenomExp,
					MarketId:         testRewardTokenMarketId,
					FeeMultiplierPpm: tc.feeMultiplierPpm,
				},
			)
			require.NoError(t, err)

			for _, rewardShare := range tc.rewardShares {
				err := k.AddRewardShareToAddress(ctx, rewardShare.Address, rewardShare.Weight.BigInt())
				require.NoError(t, err)
			}

			err = k.ProcessRewardsForBlock(ctx)
			require.NoError(t, err)

			for _, expectedBalance := range tc.expectedBalances {
				gotBalance := tApp.App.BankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(expectedBalance.Address),
					TestRewardTokenDenom,
				)
				require.Equal(t,
					expectedBalance.Coins[0], // Only checking reward token balance in `expectedBalances`.
					gotBalance,
				)
			}
		})
	}
}
