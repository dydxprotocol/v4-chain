package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k Keeper) ListingVaultDepositParams(
	ctx sdk.Context,
	req *types.QueryListingVaultDepositParams,
) (*types.QueryListingVaultDepositParamsResponse, error) {
	params := k.GetListingVaultDepositParams(ctx)
	return &types.QueryListingVaultDepositParamsResponse{
		Params: params,
	}, nil
}
