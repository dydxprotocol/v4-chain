package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	constants "github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, k)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestMsgServer_RegisterAffiliate(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgRegisterAffiliate
		expectErr error
		setup     func(ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "valid registration",
			msg: &types.MsgRegisterAffiliate{
				Referee:   constants.BobAccAddress.String(),
				Affiliate: constants.AliceAccAddress.String(),
			},
			expectErr: nil,
		},
		{
			name: "invalid referee address",
			msg: &types.MsgRegisterAffiliate{
				Referee:   "invalid_address",
				Affiliate: constants.AliceAccAddress.String(),
			},
			expectErr: types.ErrInvalidAddress,
		},
		{
			name: "invalid affiliate address",
			msg: &types.MsgRegisterAffiliate{
				Referee:   constants.BobAccAddress.String(),
				Affiliate: "invalid_address",
			},
			expectErr: types.ErrInvalidAddress,
		},
		{
			name: "referee already has an affiliate",
			msg: &types.MsgRegisterAffiliate{
				Referee:   constants.BobAccAddress.String(),
				Affiliate: constants.AliceAccAddress.String(),
			},
			expectErr: types.ErrAffiliateAlreadyExistsForReferee,
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				err := k.RegisterAffiliate(ctx, constants.BobAccAddress.String(), constants.AliceAccAddress.String())
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, ms, ctx := setupMsgServer(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			err := k.UpdateAffiliateTiers(sdkCtx, types.DefaultAffiliateTiers)
			require.NoError(t, err)
			if tc.setup != nil {
				tc.setup(sdkCtx, k)
			}
			_, err = ms.RegisterAffiliate(ctx, tc.msg)
			if tc.expectErr != nil {
				require.ErrorIs(t, err, tc.expectErr)
			} else {
				require.NoError(t, err)
				affiliate, found := k.GetReferredBy(sdkCtx, tc.msg.Referee)
				require.True(t, found)
				require.Equal(t, tc.msg.Affiliate, affiliate)
			}
		})
	}
}

func TestMsgServer_UpdateAffiliateTiers(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateAffiliateTiers
		expectErr bool
	}{
		{
			name: "Gov module updates tiers",
			msg: &types.MsgUpdateAffiliateTiers{
				Authority: lib.GovModuleAddress.String(),
				Tiers:     types.DefaultAffiliateTiers,
			},
		},
		{
			name: "non-gov module updates tiers",
			msg: &types.MsgUpdateAffiliateTiers{
				Authority: constants.BobAccAddress.String(),
				Tiers:     types.DefaultAffiliateTiers,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, ms, ctx := setupMsgServer(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			_, err := ms.UpdateAffiliateTiers(ctx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tiers, err := k.GetAllAffiliateTiers(sdkCtx)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Tiers, tiers)
			}
		})
	}
}

func TestMsgServer_UpdateAffiliateWhitelist(t *testing.T) {
	whitelist := types.AffiliateWhitelist{
		Tiers: []types.AffiliateWhitelist_Tier{
			{
				Addresses:        []string{constants.AliceAccAddress.String()},
				TakerFeeSharePpm: 100_000, // 10%
			},
		},
	}
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateAffiliateWhitelist
		expectErr bool
	}{
		{
			name: "Gov module updates whitelist",
			msg: &types.MsgUpdateAffiliateWhitelist{
				Authority: lib.GovModuleAddress.String(),
				Whitelist: whitelist,
			},
		},
		{
			name: "non-gov module updates whitelist",
			msg: &types.MsgUpdateAffiliateWhitelist{
				Authority: constants.BobAccAddress.String(),
				Whitelist: whitelist,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, ms, ctx := setupMsgServer(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			_, err := ms.UpdateAffiliateWhitelist(ctx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				whitelist, err := k.GetAffiliateWhitelist(sdkCtx)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Whitelist, whitelist)
			}
		})
	}
}
