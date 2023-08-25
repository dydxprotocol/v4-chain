package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type msgServer struct {
	types.DelayMsgKeeper
}

func NewMsgServerImpl(keeper types.DelayMsgKeeper) types.MsgServer {
	return &msgServer{keeper}
}

func (k msgServer) DelayMessage(
	goCtx context.Context,
	msg *types.MsgDelayMessage,
) (*types.MsgDelayMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authorities := k.GetAuthorities()
	if _, ok := authorities[msg.GetAuthority()]; !ok {
		panic(fmt.Errorf(
			"%v is not recognized as a valid authority for sending delayed messages",
			msg.GetAuthority(),
		))
	}

	var sdkMsg sdk.Msg

	if err := k.DecodeMessage(msg.Msg, &sdkMsg); err != nil {
		panic(fmt.Errorf("UnmarshalInterface for DelayedMessage failed, err = %w", err))
	}

	id, err := k.DelayMessageByBlocks(ctx, sdkMsg, msg.DelayBlocks)

	if err != nil {
		panic(fmt.Errorf("DelayMessageByBlocks failed, err  = %w", err))
	}

	return &types.MsgDelayMessageResponse{
		Id: uint64(id),
	}, nil
}
