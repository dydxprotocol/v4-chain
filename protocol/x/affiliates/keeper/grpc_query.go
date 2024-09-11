package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

	tierLevel, feeSharePpm, err := k.GetTierForAffiliate(ctx, addr.String())
	if err != nil {
		return nil, err
	}

	referredVolume, err := k.GetReferredVolume(ctx, req.GetAddress())
	if err != nil {
		return nil, err
	}

	stakedAmount := k.statsKeeper.GetStakedAmount(ctx, req.GetAddress())

	return &types.AffiliateInfoResponse{
		Tier:           tierLevel,
		FeeSharePpm:    feeSharePpm,
		ReferredVolume: dtypes.NewIntFromBigInt(referredVolume),
		StakedAmount:   dtypes.NewIntFromBigInt(stakedAmount),
	}, nil
}

func (k Keeper) ReferredBy(ctx context.Context,
	req *types.ReferredByRequest) (*types.ReferredByResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	affiliateAddr, exists := k.GetReferredBy(sdkCtx, req.GetAddress())
	if !exists {
		return &types.ReferredByResponse{}, errorsmod.Wrapf(
			types.ErrAffiliateNotFound,
			"affiliate not found for address: %s",
			req.GetAddress(),
		)
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
	// TODO(OTE-791): Implement `AffiliateWhitelist` RPC method.
	return nil, errors.New("not implemented")
}
