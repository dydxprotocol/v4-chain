package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) AffiliateInfo(c context.Context,
	req *types.AffiliateInfoRequest) (*types.AffiliateInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	addr, err := sdk.AccAddressFromBech32(req.GetAddress())
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidAddress, "address: %s, error: %s",
			req.GetAddress(), err.Error())
	}

	affiliateOverridesMap, err := k.GetAffiliateOverridesMap(ctx)
	if err != nil {
		return nil, err
	}
	tierLevel := uint32(0)
	feeSharePpm := uint32(0)
	isWhitelisted := false
	if _, exists := affiliateOverridesMap[addr.String()]; exists {
		isWhitelisted = true
	}
	tierLevel, feeSharePpm, err = k.GetTierForAffiliate(ctx, addr.String(), affiliateOverridesMap)
	if err != nil {
		return nil, err
	}

	userStats := k.statsKeeper.GetUserStats(ctx, addr.String())
	referredVolume := userStats.Affiliate_30DReferredVolumeQuoteQuantums
	attributedVolume := userStats.Affiliate_30DAttributedVolumeQuoteQuantums
	stakedAmount := k.statsKeeper.GetStakedBaseTokens(ctx, req.GetAddress())

	return &types.AffiliateInfoResponse{
		IsWhitelisted:               isWhitelisted,
		Tier:                        tierLevel,
		FeeSharePpm:                 feeSharePpm,
		StakedAmount:                dtypes.NewIntFromBigInt(stakedAmount),
		ReferredVolume_30DRolling:   dtypes.NewIntFromBigInt(lib.BigU(referredVolume)),
		AttributedVolume_30DRolling: dtypes.NewIntFromBigInt(lib.BigU(attributedVolume)),
	}, nil
}

func (k Keeper) ReferredBy(ctx context.Context,
	req *types.ReferredByRequest) (*types.ReferredByResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check req.Address is a valid bech32 address
	_, err := sdk.AccAddressFromBech32(req.GetAddress())
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidAddress, "address: %s, error: %s",
			req.GetAddress(), err.Error())
	}

	affiliateAddr, exists := k.GetReferredBy(sdkCtx, req.GetAddress())
	if !exists {
		return &types.ReferredByResponse{}, nil
	}

	return &types.ReferredByResponse{
		AffiliateAddress: affiliateAddr,
	}, nil
}

func (k Keeper) AllAffiliateTiers(c context.Context,
	req *types.AllAffiliateTiersRequest) (*types.AllAffiliateTiersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	affiliateTiers, err := k.GetAllAffiliateTiers(ctx)
	if err != nil {
		return nil, err
	}

	return &types.AllAffiliateTiersResponse{
		Tiers: affiliateTiers,
	}, nil
}

func (k Keeper) AffiliateWhitelist(c context.Context,
	req *types.AffiliateWhitelistRequest) (*types.AffiliateWhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	affiliateWhitelist, err := k.GetAffiliateWhitelist(ctx)
	if err != nil {
		return nil, err
	}

	return &types.AffiliateWhitelistResponse{
		Whitelist: affiliateWhitelist,
	}, nil
}

func (k Keeper) AffiliateParameters(c context.Context,
	req *types.AffiliateParametersRequest) (*types.AffiliateParametersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	affiliateParameters, err := k.GetAffiliateParameters(ctx)
	if err != nil {
		return nil, err
	}

	return &types.AffiliateParametersResponse{Parameters: affiliateParameters}, nil
}

func (k Keeper) AffiliateOverrides(c context.Context,
	req *types.AffiliateOverridesRequest) (*types.AffiliateOverridesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	affiliateOverrides, err := k.GetAffiliateOverrides(ctx)
	if err != nil {
		return nil, err
	}

	return &types.AffiliateOverridesResponse{Overrides: affiliateOverrides}, nil
}
