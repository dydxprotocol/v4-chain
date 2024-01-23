package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type msgServer struct {
	types.DelayMsgKeeper
}

func NewMsgServerImpl(keeper types.DelayMsgKeeper) types.MsgServer {
	return &msgServer{keeper}
}

// DelayMessage delays execution of a message by a given number of blocks.
func (k msgServer) DelayMessage(
	goCtx context.Context,
	msg *types.MsgDelayMessage,
) (*types.MsgDelayMessageResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	// x/delaymsg accepts messages that may have been created by other modules. In this case, the
	// ValidateBasic method of the message will not have been called. We call it here to ensure
	// that the message is valid before continuing.
	if err := msg.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"msg.ValidateBasic failed, err = %v",
			err,
		)
	}

	if !k.HasAuthority(msg.GetAuthority()) {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"%v is not recognized as a valid authority for sending messages",
			msg.GetAuthority(),
		)
	}

	sdkMsg, err := msg.GetMessage()
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"GetMessage for MsgDelayedMessage failed, err = %v",
			err,
		)
	}

	id, err := k.DelayMessageByBlocks(ctx, sdkMsg, msg.DelayBlocks)

	if err != nil {
		return nil, fmt.Errorf("DelayMessageByBlocks failed, err = %w", err)
	}

	return &types.MsgDelayMessageResponse{
		Id: uint64(id),
	}, nil
}
