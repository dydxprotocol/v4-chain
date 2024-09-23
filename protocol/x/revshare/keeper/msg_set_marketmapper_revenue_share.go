package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k msgServer) SetMarketMapperRevenueShare(
	goCtx context.Context,
	msg *types.MsgSetMarketMapperRevenueShare,
) (*types.MsgSetMarketMapperRevenueShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if sender is authorized to set revenue share
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	unconditionalRevShareConfig, err := k.GetUnconditionalRevShareConfigParams(ctx)
	if err != nil {
		return nil, err
	}
	affiliateTiers, err := k.affiliatesKeeper.GetAllAffiliateTiers(ctx)
	if err != nil {
		return nil, err
	}
	affiliateWhitelist, err := k.affiliatesKeeper.GetAffiliateWhitelist(ctx)
	if err != nil {
		return nil, err
	}

	if !k.ValidateRevShareSafety(affiliateTiers, unconditionalRevShareConfig, msg.Params, affiliateWhitelist) {
		return nil, errorsmod.Wrapf(
			types.ErrRevShareSafetyViolation,
			"rev share safety violation",
		)
	}

	// Set market mapper revenue share
	if err := k.SetMarketMapperRevenueShareParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgSetMarketMapperRevenueShareResponse{}, nil
}
