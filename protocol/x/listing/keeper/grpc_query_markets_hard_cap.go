package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k Keeper) MarketsHardCap(
	ctx context.Context,
	req *types.QueryMarketsHardCap,
) (*types.QueryMarketsHardCapResponse, error) {
	hardCap, err := k.GetMarketsHardCap(sdk.UnwrapSDKContext(ctx))
	return &types.QueryMarketsHardCapResponse{
		HardCap: hardCap,
	}, err
}
