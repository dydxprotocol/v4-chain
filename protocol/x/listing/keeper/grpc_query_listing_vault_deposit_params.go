package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k Keeper) ListingVaultDepositParams(
	ctx context.Context,
	req *types.QueryListingVaultDepositParams,
) (*types.QueryListingVaultDepositParamsResponse, error) {
	params := k.GetListingVaultDepositParams(sdk.UnwrapSDKContext(ctx))
	return &types.QueryListingVaultDepositParamsResponse{
		Params: params,
	}, nil
}
