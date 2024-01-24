package keeper

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// NextDelayedMessageId processes a query request/response for the NextDelayedMessageId from state.
func (k Keeper) NextDelayedMessageId(
	c context.Context,
	req *types.QueryNextDelayedMessageIdRequest,
) (
	*types.QueryNextDelayedMessageIdResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	nextDelayedMessageId := k.GetNextDelayedMessageId(ctx)

	return &types.QueryNextDelayedMessageIdResponse{
		NextDelayedMessageId: nextDelayedMessageId,
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
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
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

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	blockMessageIds, found := k.GetBlockMessageIds(ctx, req.BlockHeight)
	if !found {
		return nil, status.Error(codes.NotFound, "block message ids not found")
	}
	return &types.QueryBlockMessageIdsResponse{
		MessageIds: blockMessageIds.Ids,
	}, nil
}
