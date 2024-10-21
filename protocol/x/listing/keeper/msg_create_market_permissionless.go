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
		k.Logger(ctx).Error("failed to create PML market", "error", err)
		return nil, err
	}

	perpetualId, err := k.Keeper.CreatePerpetual(ctx, marketId, msg.Ticker)
	if err != nil {
		k.Logger(ctx).Error("failed to create perpetual for PML market", "error", err)
		return nil, err
	}

	clobPairId, err := k.Keeper.CreateClobPair(ctx, perpetualId)
	if err != nil {
		k.Logger(ctx).Error("failed to create clob pair for PML market", "error", err)
		return nil, err
	}

	err = k.Keeper.DepositToMegavaultforPML(ctx, *msg.SubaccountId, clobPairId)
	if err != nil {
		k.Logger(ctx).Error("failed to deposit to megavault for PML market", "error", err)
		return nil, err
	}

	return &types.MsgCreateMarketPermissionlessResponse{}, nil
}
