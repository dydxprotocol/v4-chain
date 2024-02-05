package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
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

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
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

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
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

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetSafetyParams(ctx)
	return &types.QuerySafetyParamsResponse{
		Params: params,
	}, nil
}

// AcknowledgedEventInfo processes a query request/response for `AcknowledgedEventInfo` from state.
func (k Keeper) AcknowledgedEventInfo(
	c context.Context,
	req *types.QueryAcknowledgedEventInfoRequest,
) (
	*types.QueryAcknowledgedEventInfoResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	acknowledgedEventInfo := k.GetAcknowledgedEventInfo(ctx)
	return &types.QueryAcknowledgedEventInfoResponse{
		Info: acknowledgedEventInfo,
	}, nil
}

// RecognizedEventInfo processes a query request/response for the following
// that has a greater `NextId`:
// - the `AcknowledgedEventInfo` from state
// - the `RecognizedEventInfo` from memory
// Since RecognizedEventInfo is from memory, the value is not deterministic based on state
// and therefore may be different between nodes.
func (k Keeper) RecognizedEventInfo(
	c context.Context,
	req *types.QueryRecognizedEventInfoRequest,
) (
	*types.QueryRecognizedEventInfoResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	acknowledgedEventInfo := k.GetAcknowledgedEventInfo(ctx)
	recognizedEventInfo := k.bridgeEventManager.GetRecognizedEventInfo()

	// If `AcknowledgedEventInfo` from state has a greater `NextId`, use that in response.
	// This implies that the EventInfo that has a greater `NextId` also has a equal-or-higher
	// value of `EthBlockHeight`.
	if acknowledgedEventInfo.NextId > recognizedEventInfo.NextId {
		recognizedEventInfo = acknowledgedEventInfo
	}
	return &types.QueryRecognizedEventInfoResponse{
		Info: recognizedEventInfo,
	}, nil
}

func (k Keeper) DelayedCompleteBridgeMessages(
	c context.Context,
	req *types.QueryDelayedCompleteBridgeMessagesRequest,
) (
	*types.QueryDelayedCompleteBridgeMessagesResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	return &types.QueryDelayedCompleteBridgeMessagesResponse{
		Messages: k.GetDelayedCompleteBridgeMessages(ctx, req.Address),
	}, nil
}
