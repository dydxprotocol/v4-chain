package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateBlockRateLimitConfiguration updates the equity tier limit configuration returning an error
// if the configuration is invalid.
func (k msgServer) XOperate(
	goCtx context.Context,
	msg *types.MsgXOperate,
) (resp *types.MsgXOperateResponse, err error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	sid := msg.Sid
	for _, clobid := range msg.CancelAlls {
		k.Keeper.CancelAllOrders(ctx, sid, clobid)
	}

	for _, iid := range msg.Cancels {
		uid := types.XUID{
			Sid: sid,
			Iid: iid,
		}
		k.Keeper.RemoveOrderById(ctx, uid.ToBytes())
	}

	for _, placeOrder := range msg.Orders {
		uid := types.XUID{
			Sid: sid,
			Iid: placeOrder.Iid,
		}
		order := types.FormOrder(uid, placeOrder.Base)
		_, _, err := k.Keeper.ProcessOrder(ctx, order, placeOrder.PlaceFlags)
		if err != nil {
			ctx.Logger().Error("Error processing order", "error", err)
		}
	}

	return &types.MsgXOperateResponse{}, nil
}
