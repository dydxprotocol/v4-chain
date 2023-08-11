package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4/dtypes"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	feetierstypes "github.com/dydxprotocol/v4/x/feetiers/types"
	"github.com/dydxprotocol/v4/x/rewards/types"
	"github.com/stretchr/testify/require"
)

const (
	TestAddress1 = "dydx16h7p7f4dysrgtzptxx2gtpt5d8t834g9dj830z"
	TestAddress2 = "dydx168pjt8rkru35239fsqvz7rzgeclakp49zx3aum"
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

	k.SetRewardShare(ctx, val)
	require.Equal(t, val, k.GetRewardShare(ctx, TestAddress1))
}

func TestAddRewardShareToAddress(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()

	tests := map[string]struct {
		prevRewardShare     *types.RewardShare // nil if no previous share
		newWeight           *big.Int
		expectedRewardShare types.RewardShare
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
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp.Reset()
			ctx := tApp.InitChain()
			k := tApp.App.RewardsKeeper

			if tc.prevRewardShare != nil {
				k.SetRewardShare(ctx, *tc.prevRewardShare)
			}

			k.AddRewardShareToAddress(ctx, TestAddress1, tc.newWeight)

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
				k.SetRewardShare(ctx, *tc.prevTakerRewardShare)
			}
			if tc.prevMakerRewardShare != nil {
				k.SetRewardShare(ctx, *tc.prevMakerRewardShare)
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
