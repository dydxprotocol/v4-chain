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

	// Set market mapper revenue share details
	if err := k.SetMarketMapperRevShareDetails(ctx, msg.Params.MarketId, *msg.Params.MarketMapperRevShareDetails); err != nil {
		return nil, err
	}

	return &types.MsgSetMarketMapperRevShareDetailsForMarketResponse{}, nil
}
