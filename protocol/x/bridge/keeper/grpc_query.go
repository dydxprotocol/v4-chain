package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/bridge/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// EventParams processes a query request/response for the EventParams from state.
func (k Keeper) EventParams(
	c context.Context,
	req *types.QueryEventParamsRequest,
) (
	*types.QueryEventParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetEventParams(ctx)
	return &types.QueryEventParamsResponse{
		Params: params,
	}, nil
}

// ProposeParams processes a query request/response for the ProposeParams from state.
func (k Keeper) ProposeParams(
	c context.Context,
	req *types.QueryProposeParamsRequest,
) (
	*types.QueryProposeParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetProposeParams(ctx)
	return &types.QueryProposeParamsResponse{
		Params: params,
	}, nil
}

// SafetyParams processes a query request/response for the SafetyParams from state.
func (k Keeper) SafetyParams(
	c context.Context,
	req *types.QuerySafetyParamsRequest,
) (
	*types.QuerySafetyParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetSafetyParams(ctx)
	return &types.QuerySafetyParamsResponse{
		Params: params,
	}, nil
}

// NextAcknowledgedEventId processes a query request/response for the NextAcknowledgedEventId from state.
func (k Keeper) NextAcknowledgedEventId(
	c context.Context,
	req *types.QueryNextAcknowledgedEventIdRequest,
) (
	*types.QueryNextAcknowledgedEventIdResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryNextAcknowledgedEventIdResponse{
		Id: k.GetNextAcknowledgedEventId(ctx),
	}, nil
}

// NextRecognizedEventId processes a query request/response for
// the greater of:
// - the NextAcknowledgedEventId from state
// - the NextRecognizedEventId from memory
// Since NextRecognizedEventId is from memory, the value is not deterministic based on state
// and therefore may be different between nodes.
func (k Keeper) NextRecognizedEventId(
	c context.Context,
	req *types.QueryNextRecognizedEventIdRequest,
) (
	*types.QueryNextRecognizedEventIdResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	ackEventId := k.GetNextAcknowledgedEventId(ctx)
	recEventId := uint32(0) // TODO(CORE-323): get the next recognized event id from bridgeEventManager
	return &types.QueryNextRecognizedEventIdResponse{
		Id: lib.MaxUint32(ackEventId, recEventId),
	}, nil
}
