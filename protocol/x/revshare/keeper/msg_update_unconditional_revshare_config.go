package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k msgServer) UpdateUnconditionalRevShareConfig(
	goCtx context.Context,
	msg *types.MsgUpdateUnconditionalRevShareConfig,
) (*types.MsgUpdateUnconditionalRevShareConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if sender is authorized to set revenue share
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}
	err := msg.ValidateBasic()
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
	marketMapperRevShareParams := k.GetMarketMapperRevenueShareParams(ctx)
	if !k.ValidateRevShareSafety(affiliateTiers, msg.Config, marketMapperRevShareParams, affiliateWhitelist) {
		return nil, errorsmod.Wrapf(
			types.ErrRevShareSafetyViolation,
			"rev share safety violation",
		)
	}
	k.SetUnconditionalRevShareConfigParams(ctx, msg.Config)
	return &types.MsgUpdateUnconditionalRevShareConfigResponse{}, nil
}
