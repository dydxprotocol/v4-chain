package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) State(ctx context.Context, request *types.QueryStateRequest) (*types.QueryStateResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return &types.QueryStateResponse{Controls: k.GetControls(sdkCtx), Ledger: k.GetLedger(sdkCtx)}, nil
}

func (k Keeper) Participants(
	ctx context.Context,
	request *types.QueryParticipantsRequest,
) (*types.QueryParticipantsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	return &types.QueryParticipantsResponse{
		Participants: k.GetAllParticipants(sdk.UnwrapSDKContext(ctx)),
	}, nil
}

func (k Keeper) PendingSettlements(
	ctx context.Context,
	request *types.QueryPendingSettlementsRequest,
) (*types.QueryPendingSettlementsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	return &types.QueryPendingSettlementsResponse{
		PendingSettlements: k.GetAllPending(sdk.UnwrapSDKContext(ctx)),
	}, nil
}
