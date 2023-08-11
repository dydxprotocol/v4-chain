package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

func (k msgServer) PlaceOrder(goCtx context.Context, msg *types.MsgPlaceOrder) (*types.MsgPlaceOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	return &types.MsgPlaceOrderResponse{}, nil
}
