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

func (k Keeper) PendingSettlements(
	ctx context.Context,
	request *types.QueryPendingSettlementsRequest,
) (*types.QueryPendingSettlementsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	pending, pageResponse, err := k.GetPendingPage(sdk.UnwrapSDKContext(ctx), request.Pagination)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.QueryPendingSettlementsResponse{PendingSettlements: pending, Pagination: pageResponse}, nil
}

func (k Keeper) PendingSettlement(
	ctx context.Context,
	request *types.QueryPendingSettlementRequest,
) (*types.QueryPendingSettlementResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	pending, found := k.GetPendingForPacket(sdk.UnwrapSDKContext(ctx), request.SourceChannel, request.Sequence)
	if !found {
		return nil, status.Error(codes.NotFound, types.ErrSettlementNotFound.Error())
	}
	return &types.QueryPendingSettlementResponse{PendingSettlement: pending}, nil
}
