package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"

	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// NumMessages processes a query request/response for the NumMessages from state.
func (k Keeper) NumMessages(
	c context.Context,
	req *types.QueryNumMessagesRequest,
) (
	*types.QueryNumMessagesResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	numMessages := k.GetNumMessages(ctx)

	return &types.QueryNumMessagesResponse{
		NumMessages: numMessages,
	}, nil
}

// Message processes a query request/response for the Message from state.
func (k Keeper) Message(
	c context.Context,
	req *types.QueryMessageRequest,
) (
	*types.QueryMessageResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	delayedMessage, found := k.GetMessage(ctx, req.Id)

	if !found {
		return nil, status.Error(codes.NotFound, "message not found")
	}

	return &types.QueryMessageResponse{
		Message: &delayedMessage,
	}, nil
}

// BlockMessageIds processes a query request/response for the BlockMessageIds from state.
func (k Keeper) BlockMessageIds(
	c context.Context,
	req *types.QueryBlockMessageIdsRequest,
) (
	*types.QueryBlockMessageIdsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.BlockHeight < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid block height")
	}

	ctx := sdk.UnwrapSDKContext(c)

	blockMessageIds, found := k.GetBlockMessageIds(ctx, req.BlockHeight)
	if !found {
		return nil, status.Error(codes.NotFound, "block message ids not found")
	}
	return &types.QueryBlockMessageIdsResponse{
		MessageIds: blockMessageIds.Ids,
	}, nil
}
