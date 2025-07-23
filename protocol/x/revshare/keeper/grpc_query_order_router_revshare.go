package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k Keeper) OrderRouterRevShare(
	ctx context.Context,
	req *types.QueryOrderRouterRevShare,
) (*types.QueryOrderRouterRevShareResponse, error) {
	if req == nil {
		return nil, types.ErrEmptyRequest
	}
	revSharePpm, err := k.GetOrderRouterRevShare(sdk.UnwrapSDKContext(ctx), req.Address)
	if err != nil {
		return nil, err
	}
	return &types.QueryOrderRouterRevShareResponse{
		OrderRouterRevShare: types.OrderRouterRevShare{
			Address:  req.Address,
			SharePpm: revSharePpm,
		},
	}, nil
}
