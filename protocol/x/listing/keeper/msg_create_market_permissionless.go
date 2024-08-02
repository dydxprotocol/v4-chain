package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k msgServer) CreateMarketPermissionless(
	goCtx context.Context,
	msg *types.MsgCreateMarketPermissionless,
) (*types.MsgCreateMarketPermissionlessResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Check if the number of listed markets is above the hard cap
	numPerpetuals := len(k.PerpetualsKeeper.GetAllPerpetuals(ctx))
	if uint32(numPerpetuals) > k.Keeper.GetMarketsHardCap(ctx) {
		return nil, types.ErrMarketsHardCapReached
	}

	marketId, err := k.Keeper.CreateMarket(ctx, msg.Ticker)
	if err != nil {
		return nil, err
	}

	perpetualId, err := k.Keeper.CreatePerpetual(ctx, marketId, msg.Ticker)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.CreateClobPair(ctx, perpetualId)
	if err != nil {
		return nil, err
	}

	// TODO: vault deposit for PML

	return &types.MsgCreateMarketPermissionlessResponse{}, nil
}
