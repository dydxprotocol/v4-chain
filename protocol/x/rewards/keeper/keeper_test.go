package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
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
		Address: types.TreasuryModuleAddress.String(),
		Coins: []sdk.Coin{{
			Denom:  TestRewardTokenDenom,
			Amount: sdkmath.NewInt(0),
		}},
	}
)

func TestRewardShareStorage_DefaultValue(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
			tApp := testapp.NewTestAppBuilder(t).Build()
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
	makerAddress := TestAddress1
	takerAddress := TestAddress2

	tests := map[string]struct {
		prevTakerRewardShare *types.RewardShare
		prevMakerRewardShare *types.RewardShare
		fill                 clobtypes.FillForProcess
		revSharesForFill     revsharetypes.RevSharesForFill
		feeTiers             []*feetierstypes.PerpetualFeeTier

		expectedTakerShare types.RewardShare
		expectedMakerShare types.RewardShare
	}{
		"positive maker fee, positive taker fees reduced by maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(1_000_000),
				FillQuoteQuantums:                 big.NewInt(800_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares:             []revsharetypes.RevShare{},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{},
				FeeSourceToRevSharePpm:   map[revsharetypes.RevShareFeeSource]uint32{},
				AffiliateRevShare:        nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(200_000), // 2 - 0.1% * 800 -(2 * 0.5)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(1_000_000),
			},
		},
		"negative maker fee, positive taker fees reduced by 0.1% maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(-1_000_000),
				FillQuoteQuantums:                 big.NewInt(750_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares:             []revsharetypes.RevShare{},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{},
				FeeSourceToRevSharePpm:   map[revsharetypes.RevShareFeeSource]uint32{},
				AffiliateRevShare:        nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(250_000), // 2 - 0.1% * 750 - (2 * 0.5)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(0),
			},
		},
		"negative maker fee, positive taker fees reduced by 0.05% maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(-1_000_000),
				FillQuoteQuantums:                 big.NewInt(750_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares:             []revsharetypes.RevShare{},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{},
				FeeSourceToRevSharePpm:   map[revsharetypes.RevShareFeeSource]uint32{},
				AffiliateRevShare:        nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -500,  // -0.05%
					TakerFeePpm: 2_000, // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(625_000), // 2 - 0.05% * 750 - (2 * 0.5)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(0),
			},
		},
		"positive maker fee, positive taker fees offset by maker rebate, no previous share": {
			prevTakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(700_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(500_000),
				FillQuoteQuantums:                 big.NewInt(750_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares:             []revsharetypes.RevShare{},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{},
				FeeSourceToRevSharePpm:   map[revsharetypes.RevShareFeeSource]uint32{},
				AffiliateRevShare:        nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
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
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(700_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(500_000),
				FillQuoteQuantums:                 big.NewInt(750_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares:             []revsharetypes.RevShare{},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{},
				FeeSourceToRevSharePpm:   map[revsharetypes.RevShareFeeSource]uint32{},
				AffiliateRevShare:        nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: 1_000, // 0.1%
					TakerFeePpm: 2_000, // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(350_000), // 0.7 - (0.7 * 0.5)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(500_000),
			},
		},
		"positive maker + taker fees reduced by maker rebate, no previous share with net fee revshare": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(1_000_000),
				FillQuoteQuantums:                 big.NewInt(800_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 9,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(200_000),
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(0),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            0,
				},
				AffiliateRevShare: nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(180_000), // (2 - 0.1% * 800 - 0.5*2) * (1 - 0.1)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(900_000), // 1 * (1 - 0.1)
			},
		},
		"positive maker + taker fees reduced by maker rebate, no previous share with multiple net fee revshare": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(1_000_000),
				FillQuoteQuantums:                 big.NewInt(800_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(400_000),
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(0),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 200_000, // 20%
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            0,
				},
				AffiliateRevShare: nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(160_000), // (2 - 0.1% * 800 - 0.5*2) * (1 - 0.1)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(800_000), // 1 * (1 - 0.2)
			},
		},
		"positive maker + taker fees reduced by maker rebate, no previous share and taker + net fee revshare": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(1_000_000),
				FillQuoteQuantums:                 big.NewInt(800_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         takerAddress,
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(200_000),
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(200_000),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            100_000, // 10%
				},
				AffiliateRevShare: &revsharetypes.RevShare{
					Recipient:         takerAddress,
					RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      revsharetypes.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(200_000),
					RevSharePpm:       100_000, // 10%
				},
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(180_000), // (2 - 0.1% * 800 - 1) * (1 - 0.1)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(900_000), // 1 * (1 - 0.1)
			},
		},
		"positive maker + taker fees reduced by maker rebate, no previous share and taker + order router rev share": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(1_000_000),
				FillQuoteQuantums:                 big.NewInt(800_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 0,
				TakerOrderRouterAddr:              constants.AliceAccAddress.String(),
				MakerOrderRouterAddr:              constants.BobAccAddress.String(),
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         takerAddress,
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER,
						QuoteQuantums:     big.NewInt(50_000),
						RevSharePpm:       50_000, // 5%
					},
					{
						Recipient:         makerAddress,
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER,
						QuoteQuantums:     big.NewInt(100_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(200_000),
					revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE:            big.NewInt(100_000),
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(50_000),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE:            100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            50_000,  // 5%
				},
				AffiliateRevShare: nil,
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(135_000), // ((2 - 0.1% * 800) - (2 * 0.5) - 0.05) * (1 - 0.1)) * 2
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(810_000), // (1 - 0.1) * (1 - 0.1)
			},
		},
		"positive maker + taker fees reduced by maker rebate, taker + net fee revshare,rolling taker volume > 50 mil": {
			prevTakerRewardShare: nil,
			prevMakerRewardShare: nil,
			fill: clobtypes.FillForProcess{
				TakerAddr:                         takerAddress,
				TakerFeeQuoteQuantums:             big.NewInt(2_000_000),
				MakerAddr:                         makerAddress,
				MakerFeeQuoteQuantums:             big.NewInt(1_000_000),
				FillQuoteQuantums:                 big.NewInt(800_000_000),
				ProductId:                         uint32(1),
				MarketId:                          uint32(1),
				MonthlyRollingTakerVolumeQuantums: 60_000_000_000_000,
			},
			revSharesForFill: revsharetypes.RevSharesForFill{
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         takerAddress,
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(200_000),
						RevSharePpm:       100_000, // 10%
					},
				},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(200_000),
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(200_000),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            100_000, // 10%
				},
				AffiliateRevShare: &revsharetypes.RevShare{
					Recipient:         takerAddress,
					RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      revsharetypes.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(200_000),
					RevSharePpm:       100_000, // 10%
				},
			},
			feeTiers: []*feetierstypes.PerpetualFeeTier{
				{
					MakerFeePpm: -1_000, // -0.1%
					TakerFeePpm: 2_000,  // 0.2%
				},
			},
			expectedTakerShare: types.RewardShare{
				Address: takerAddress,
				Weight:  dtypes.NewInt(1_080_000), // (2 - 0.1% * 800 - 0) * (1 - 0.1)
			},
			expectedMakerShare: types.RewardShare{
				Address: makerAddress,
				Weight:  dtypes.NewInt(900_000), // 1 * (1 - 0.1)
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
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
				tc.fill,
				tc.revSharesForFill,
			)

			// Check the new reward shares.
			require.Equal(t, tc.expectedTakerShare, k.GetRewardShare(ctx, takerAddress))
			require.Equal(t, tc.expectedMakerShare, k.GetRewardShare(ctx, makerAddress))
		})
	}
}

func TestProcessRewardsForBlock(t *testing.T) {
	testRewardTokenMarketId := uint32(33)
	testRewardTokenMarket := "test-market"
	TestRewardTokenDenomExp := int32(-18)

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
			rewardShares: []types.RewardShare{},
			tokenPrice:   tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(1000, 18), // 1000 full coins
			),
			feeMultiplierPpm: 1_000_000, // 100%
			// 1$ / 2$ * 100% = 0.5 full coin, all paid to TestAddress1
			expectedBalances: []banktypes.Balance{
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(1000, 18), // 1000 full coins
						),
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
			tokenPrice: tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(1000, 18), // 1000 full coins
			),
			feeMultiplierPpm: 1_000_000, // 100%
			// 1$ / 2$ * 100% = 0.5 full coin, all paid to TestAddress1
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(5e17), // 0.5 full coin
					}},
				},
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(9995, 17), // 999.5 full coins
						),
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
			tokenPrice: tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(1, 20), // 100 full coins
			),
			feeMultiplierPpm: 950_000, // 95%
			// 1$ / 2$ * 95% = 0.475 full coin, all paid to TestAddress1
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(475, 15), // 0.475 full coin
						),
					}},
				},
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(99525, 15), // 99.525 full coin
						),
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
			tokenPrice: tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(2, 17), // 0.2 full coins
			),
			feeMultiplierPpm: 1_000_000, // 100%
			// 1$ / 2$ * 100% = 0.5 full coin > 0.2 full coin. Pay 0.2 full coin to TestAddress1.
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(2, 17), // 0.2 full coins
						),
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
				big_testutil.Int64MulPow10(2_000_123, 18), // ~2_000_123 full coin.
			),
			feeMultiplierPpm: 990_000, // 99%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("507667050000000000000", 10)),
						), // $1025.59 weight / $2 price * 99% ~= 507.667 full coin
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("1000543500000000000000", 10)),
						), // $2021.3 weight / $2 price * 99% ~= 1000 full coin
					}},
				},
				{
					Address: TestAddress3,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("413651700000000000000", 10)),
						), // $835.66 weight / $2 price * 99% ~= 413 full coin
					}},
				},
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("1998201137750000000000000", 10)),
						), // 2_000_123 - 507.667 - 1000.5435 - 413.6517 ~= 1_998_201.1 full coins
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
			tokenPrice: tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(10, 18),
			), // 10 full coins
			feeMultiplierPpm: 1_000_000, // 100%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1_666_666_666_666_666_666), // 1/6 of 10 = 1.666666 full coins
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(3_333_333_333_333_333_333), // 1/3 of 10 = 3.333333 full coins
					}},
				},
				{
					Address: TestAddress3,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(5, 18),
						), // 1/2 of 10 = 5 full coins
					}},
				},
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1), // 1e-18 full coins left due to rounding
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
			tokenPrice: tokenPrice1_18Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(100, 18),
			), // 100 full coins
			feeMultiplierPpm: 990_000, // 99%
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("56722081676906396819", 10)),
						), // 56.722081 full coins
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("225876400433682688", 10)),
						), // 0.225876 full coin
					}},
				},
				{
					Address: TestAddress3,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("43052041922659920491", 10)),
						), // 43.052041 full coins
					}},
				},
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(2), // 2e-18 full coins left due to rounding
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
			tokenPrice: tokenPrice2Usdc,
			treasuryAccountBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(1, 15),
			), // 0.001 full coins
			feeMultiplierPpm: 990_000, // 0.99
			expectedBalances: []banktypes.Balance{
				{
					Address: TestAddress1,
					Coins: []sdk.Coin{{
						Denom: TestRewardTokenDenom,
						Amount: sdkmath.NewIntFromBigInt(
							big_testutil.Int64MulPow10(99_999_999, 7),
						), // 0.001 * 100_000_000 / 100_000_001 = 0.00099999999 full coin
					}},
				},
				{
					Address: TestAddress2,
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(9_999_999), // rounded to 9.9e-12 full coins
					}},
				},
				{
					Address: types.TreasuryModuleAddress.String(),
					Coins: []sdk.Coin{{
						Denom:  TestRewardTokenDenom,
						Amount: sdkmath.NewInt(1), // 0.001 - 0.00099999999 + 9.9e-12 = 1e-18 full coin left due to rounding
					}},
				},
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set up treasury account balance in genesis state
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: types.TreasuryModuleAddress.String(),
							Coins: []sdk.Coin{
								sdk.NewCoin(TestRewardTokenDenom, tc.treasuryAccountBalance),
							},
						})
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RewardsKeeper

			// Set up PricesKeeper
			_, err := keepertest.CreateTestMarket(
				t,
				ctx,
				&tApp.App.PricesKeeper,
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
					"expected balance: %s, got: %s",
					expectedBalance.Coins[0],
					gotBalance,
				)
			}
		})
	}
}
