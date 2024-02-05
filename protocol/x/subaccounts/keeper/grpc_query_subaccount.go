package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SubaccountAll(
	c context.Context,
	req *types.QueryAllSubaccountRequest,
) (*types.QuerySubaccountAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var subaccounts []types.Subaccount
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	store := ctx.KVStore(k.storeKey)
	subaccountStore := prefix.NewStore(store, []byte(types.SubaccountKeyPrefix))

	pageRes, err := query.Paginate(subaccountStore, req.Pagination, func(key []byte, value []byte) error {
		var subaccount types.Subaccount
		if err := k.cdc.Unmarshal(value, &subaccount); err != nil {
			return err
		}

		subaccounts = append(subaccounts, subaccount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySubaccountAllResponse{Subaccount: subaccounts, Pagination: pageRes}, nil
}

func (k Keeper) Subaccount(
	c context.Context,
	req *types.QueryGetSubaccountRequest,
) (*types.QuerySubaccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	val := k.GetSubaccount(
		ctx,
		types.SubaccountId{
			Owner:  req.Owner,
			Number: req.Number,
		},
	)

	return &types.QuerySubaccountResponse{Subaccount: val}, nil
}
