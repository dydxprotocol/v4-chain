package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"

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
	// x/delaymsg accepts messages that may have been created by other modules. In this case, the
	// ValidateBasic method of the message will not have been called. We call it here to ensure
	// that the message is valid before continuing.
	if err := msg.ValidateBasic(); err != nil {
		k.Logger(ctx).Error(
			"DelayMessage failed because msg.ValidateBasic failed",
			constants.ErrorLogKey,
			err,
		)
		return nil, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"msg.ValidateBasic failed, err = %v",
			err,
		)
	}

	if !k.HasAuthority(msg.GetAuthority()) {
		k.Logger(ctx).Error(
			"DelayMessage failed because msg.Authority is not recognized as a valid authority for sending messages",
			"authority",
			msg.GetAuthority(),
		)
		return nil, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"%v is not recognized as a valid authority for sending messages",
			msg.GetAuthority(),
		)
	}

	sdkMsg, err := msg.GetMessage()
	if err != nil {
		k.Logger(ctx).Error(
			"GetMessage for MsgDelayMessage failed",
			constants.ErrorLogKey,
			err,
		)
		return nil, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"GetMessage for MsgDelayedMessage failed, err = %v",
			err,
		)
	}

	id, err := k.DelayMessageByBlocks(ctx, sdkMsg, msg.DelayBlocks)

	if err != nil {
		k.Logger(ctx).Error(
			"DelayMessageByBlocks failed",
			constants.ErrorLogKey,
			err,
		)
		return nil, fmt.Errorf("DelayMessageByBlocks failed, err  = %w", err)
	}

	return &types.MsgDelayMessageResponse{
		Id: uint64(id),
	}, nil
}
