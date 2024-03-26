package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ListLimitParams(
	ctx context.Context,
	req *types.ListLimitParamsRequest,
) (*types.ListLimitParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return &types.ListLimitParamsResponse{
		LimitParamsList: k.GetAllLimitParams(sdkCtx),
	}, nil
}

func (k Keeper) CapacityByDenom(
	ctx context.Context,
	req *types.QueryCapacityByDenomRequest,
) (*types.QueryCapacityByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	limiterCapacityList, err := k.GetLimiterCapacityListForDenom(sdkCtx, req.Denom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCapacityByDenomResponse{
		LimiterCapacityList: limiterCapacityList,
	}, nil
}

func (k Keeper) AllPendingSendPackets(
	ctx context.Context,
	req *types.QueryAllPendingSendPacketsRequest,
) (*types.QueryAllPendingSendPacketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.PendingSendPacketPrefix))
	pendingPackets := make([]types.PendingSendPacket, 0)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		channelId, sequence, err := types.SplitPendingSendPacketKey(iterator.Key())
		if err != nil {
			log.ErrorLog(sdkCtx, "unexpected PendingSendPacket key format", err)
			return nil, err
		}
		pendingPackets = append(pendingPackets, types.PendingSendPacket{
			ChannelId: channelId,
			Sequence:  sequence,
		})
	}
	return &types.QueryAllPendingSendPacketsResponse{
		PendingSendPackets: pendingPackets,
	}, nil
}
