package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k Keeper) UnconditionalRevShareConfig(
	ctx context.Context,
	req *types.QueryUnconditionalRevShareConfig,
) (*types.QueryUnconditionalRevShareConfigResponse, error) {
	config, err := k.GetUnconditionalRevShareConfigParams(sdk.UnwrapSDKContext(ctx))
	if err != nil {
		return nil, err
	}
	return &types.QueryUnconditionalRevShareConfigResponse{
		Config: config,
	}, nil
}
