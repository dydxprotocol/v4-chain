package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) DelayMessage(
	goCtx context.Context,
	msg *types.MsgDelayMessage,
) (*types.MsgDelayMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO(CORE-437): Filter out non-authorized messages.

	var sdkMsg sdk.Msg
	if err := k.cdc.UnmarshalInterface(msg.Msg, &sdkMsg); err != nil {
		return nil, err
	}

	id, err := k.DelayMessageByBlocks(ctx, sdkMsg, msg.DelayBlocks)

	if err != nil {
		return nil, err
	}

	return &types.MsgDelayMessageResponse{
		Id: uint64(id),
	}, nil
}
