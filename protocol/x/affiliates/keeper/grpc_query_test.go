package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	constants "github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
)

func TestAffiliateInfo(t *testing.T) {
	testCases := map[string]struct {
		req         *types.AffiliateInfoRequest
		res         *types.AffiliateInfoResponse
		setup       func(ctx sdk.Context, k keeper.Keeper, tApp *testapp.TestApp)
		expectError error
	}{

		"Success": {
			req: &types.AffiliateInfoRequest{
				Address: constants.AliceAccAddress.String(),
			},
			res: &types.AffiliateInfoResponse{
				IsWhitelisted:  false,
				Tier:           0,
				FeeSharePpm:    types.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm,
				ReferredVolume: dtypes.NewIntFromUint64(types.DefaultAffiliateTiers.Tiers[0].ReqReferredVolumeQuoteQuantums),
				StakedAmount:   dtypes.NewIntFromUint64(uint64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins) * 1e18),
			},
			setup: func(ctx sdk.Context, k keeper.Keeper, tApp *testapp.TestApp) {
				err := k.RegisterAffiliate(ctx, constants.BobAccAddress.String(), constants.AliceAccAddress.String())
				require.NoError(t, err)

				stakingKeeper := tApp.App.StakingKeeper

				err = stakingKeeper.SetDelegation(ctx,
					stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
						constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
							new(big.Int).Mul(
								big.NewInt(int64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
								big.NewInt(1e18),
							),
						),
					),
				)
				require.NoError(t, err)
			},
		},
		"NonExistentAddress": {
			req: &types.AffiliateInfoRequest{
				Address: constants.AliceAccAddress.String(),
			},
			res: &types.AffiliateInfoResponse{
				IsWhitelisted:  false,
				Tier:           0,
				FeeSharePpm:    types.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm,
				ReferredVolume: dtypes.NewIntFromUint64(types.DefaultAffiliateTiers.Tiers[0].ReqReferredVolumeQuoteQuantums),
				StakedAmount:   dtypes.NewIntFromUint64(uint64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins) * 1e18),
			},
			setup: func(ctx sdk.Context, k keeper.Keeper, tApp *testapp.TestApp) {
				stakingKeeper := tApp.App.StakingKeeper
				err := stakingKeeper.SetDelegation(ctx,
					stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
						constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
							big.NewInt(int64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
						),
					),
				)
				require.NoError(t, err)
			},
		},
		"InvalidAddress": {
			req: &types.AffiliateInfoRequest{
				Address: "invalid_address",
			},
			res:         nil,
			setup:       func(ctx sdk.Context, k keeper.Keeper, tApp *testapp.TestApp) {},
			expectError: types.ErrInvalidAddress,
		},
		"EmptyRequest": {
			req:         &types.AffiliateInfoRequest{},
			res:         nil,
			setup:       func(ctx sdk.Context, k keeper.Keeper, tApp *testapp.TestApp) {},
			expectError: types.ErrInvalidAddress,
		},
		"Whitelisted": {
			req: &types.AffiliateInfoRequest{
				Address: constants.AliceAccAddress.String(),
			},
			res: &types.AffiliateInfoResponse{
				IsWhitelisted:  true,
				Tier:           0,
				FeeSharePpm:    120_000,
				ReferredVolume: dtypes.NewIntFromUint64(0),
				StakedAmount:   dtypes.NewIntFromUint64(0),
			},
			setup: func(ctx sdk.Context, k keeper.Keeper, tApp *testapp.TestApp) {
				err := k.RegisterAffiliate(ctx, constants.BobAccAddress.String(), constants.AliceAccAddress.String())
				require.NoError(t, err)

				stakingKeeper := tApp.App.StakingKeeper

				err = stakingKeeper.SetDelegation(ctx,
					stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
						constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
							big.NewInt(0),
						),
					),
				)
				require.NoError(t, err)

				affiliatesWhitelist := types.AffiliateWhitelist{
					Tiers: []types.AffiliateWhitelist_Tier{
						{
							Addresses:        []string{constants.AliceAccAddress.String()},
							TakerFeeSharePpm: 120_000, // 12%
						},
					},
				}
				err = k.SetAffiliateWhitelist(ctx, affiliatesWhitelist)
				require.NoError(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.AffiliatesKeeper

			// Set up affiliate tiers
			tiers := types.DefaultAffiliateTiers
			err := k.UpdateAffiliateTiers(ctx, tiers)
			require.NoError(t, err)

			// Run the setup function
			tc.setup(ctx, k, tApp)

			// Call the AffiliateInfo method
			res, err := k.AffiliateInfo(ctx, tc.req)

			// Check the result
			if tc.res == nil {
				require.ErrorIs(t, err, tc.expectError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestReferredBy(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper
	testCases := map[string]struct {
		req         *types.ReferredByRequest
		setup       func(ctx sdk.Context, k keeper.Keeper)
		expected    *types.ReferredByResponse
		expectError error
	}{
		"Success": {
			req: &types.ReferredByRequest{
				Address: constants.AliceAccAddress.String(),
			},
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				err := k.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
			expected: &types.ReferredByResponse{
				AffiliateAddress: constants.BobAccAddress.String(),
			},
		},
		"Affiliate not registered": {
			req: &types.ReferredByRequest{
				Address: constants.DaveAccAddress.String(),
			},
			setup:       func(ctx sdk.Context, k keeper.Keeper) {},
			expected:    nil,
			expectError: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
			require.NoError(t, err)
			tc.setup(ctx, k)
			res, err := k.ReferredBy(ctx, tc.req)

			if tc.expected == nil {
				require.ErrorIs(t, err, tc.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, res)
			}
		})
	}
}

func TestAllAffiliateTiers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	req := &types.AllAffiliateTiersRequest{}

	tiers := types.DefaultAffiliateTiers
	err := k.UpdateAffiliateTiers(ctx, tiers)
	require.NoError(t, err)

	res, err := k.AllAffiliateTiers(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, &types.AllAffiliateTiersResponse{Tiers: tiers}, res)
}

func TestAffiliateWhitelist(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	req := &types.AffiliateWhitelistRequest{}
	whitelist := types.AffiliateWhitelist{
		Tiers: []types.AffiliateWhitelist_Tier{
			{
				Addresses:        []string{constants.AliceAccAddress.String()},
				TakerFeeSharePpm: 100_000,
			},
		},
	}
	err := k.SetAffiliateWhitelist(ctx, whitelist)
	require.NoError(t, err)

	res, err := k.AffiliateWhitelist(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, &types.AffiliateWhitelistResponse{Whitelist: whitelist}, res)
}
