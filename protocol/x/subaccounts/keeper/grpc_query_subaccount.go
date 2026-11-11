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

// MaxSubaccountAllPageLimit bounds how many subaccounts a single SubaccountAll query can decode and
// buffer into one response. Without a per-handler cap, a very large requested limit forces the node to
// unmarshal and hold that many records in memory for one request regardless of any chain-wide
// query-gas-limit, since gas accounting does not by itself bound response size. Larger limits are
// clamped rather than rejected so key-based paginators (e.g. the liquidation daemon, whose page limit
// is operator-configurable) keep working by following NextKey.
const MaxSubaccountAllPageLimit = 1000

func (k Keeper) SubaccountAll(
	c context.Context,
	req *types.QueryAllSubaccountRequest,
) (*types.QuerySubaccountAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pagination := &query.PageRequest{}
	if req.Pagination != nil {
		*pagination = *req.Pagination
	}
	// CountTotal forces query.Paginate to walk the entire prefix store to completion regardless of
	// Limit, so capping Limit alone does not bound the cost of a count_total=true request.
	if pagination.CountTotal {
		return nil, status.Error(codes.InvalidArgument, "count_total is not supported for this query")
	}
	// query.Paginate turns CountTotal back on whenever Limit is 0, so the default has to be set here.
	if pagination.Limit == 0 {
		pagination.Limit = query.DefaultLimit
	}
	pagination.Limit = min(pagination.Limit, MaxSubaccountAllPageLimit)

	var subaccounts []types.Subaccount
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	store := ctx.KVStore(k.storeKey)
	subaccountStore := prefix.NewStore(store, []byte(types.SubaccountKeyPrefix))

	pageRes, err := query.Paginate(subaccountStore, pagination, func(key []byte, value []byte) error {
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
