package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k Keeper) PermissionlessMarketListingStatus(
	ctx context.Context,
	req *types.QueryPermissionlessMarketListingStatus,
) (*types.QueryPermissionlessMarketListingStatusResponse, error) {

	enabled, err := k.IsPermissionlessListingEnabled(sdk.UnwrapSDKContext(ctx))
	return &types.QueryPermissionlessMarketListingStatusResponse{
		Enabled: enabled,
	}, err
}
