package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k msgServer) SetMarketMapperRevShareDetailsForMarket(
	goCtx context.Context,
	msg *types.MsgSetMarketMapperRevShareDetailsForMarket,
) (*types.MsgSetMarketMapperRevShareDetailsForMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if sender is authorized to set revenue share
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// We do not validate if the market exists in x/prices here because that
	// creates a circular dependency between x/prices and x/revshare. We can
	// assume that governance provides a strong safety check for this

	// Set market mapper revenue share details
	k.SetMarketMapperRevShareDetails(ctx, msg.MarketId, msg.Params)

	return &types.MsgSetMarketMapperRevShareDetailsForMarketResponse{}, nil
}
