package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO(CLOB-863) Add tests for these endpoints.
func (k Keeper) AllAssets(
	c context.Context,
	req *types.QueryAllAssetsRequest,
) (*types.QueryAllAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var assets []types.Asset
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	store := ctx.KVStore(k.storeKey)
	assetStore := prefix.NewStore(store, []byte(types.AssetKeyPrefix))

	pageRes, err := query.Paginate(assetStore, req.Pagination, func(key []byte, value []byte) error {
		var asset types.Asset
		if err := k.cdc.Unmarshal(value, &asset); err != nil {
			return err
		}

		assets = append(assets, asset)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAssetsResponse{Asset: assets, Pagination: pageRes}, nil
}

func (k Keeper) Asset(c context.Context, req *types.QueryAssetRequest) (*types.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	val, exists := k.GetAsset(
		ctx,
		req.Id,
	)

	if !exists {
		return nil,
			status.Error(
				codes.NotFound,
				fmt.Sprintf(
					"Asset id %+v not found.",
					req.Id,
				),
			)
	}

	return &types.QueryAssetResponse{Asset: val}, nil
}
