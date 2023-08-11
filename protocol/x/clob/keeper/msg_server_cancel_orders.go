package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

func (k msgServer) CancelOrder(
	goCtx context.Context,
	msg *types.MsgCancelOrder,
) (*types.MsgCancelOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	return &types.MsgCancelOrderResponse{}, nil
}
