package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateBlockRateLimitConfiguration updates the block rate limit configuration returning an error
// if the configuration is invalid.
func (k msgServer) UpdateBlockRateLimitConfiguration(goCtx context.Context, configuration *types.MsgUpdateBlockRateLimitConfiguration) (*types.MsgUpdateBlockRateLimitConfigurationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Keeper.InitializeBlockRateLimit(ctx, configuration.BlockRateLimitConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateBlockRateLimitConfigurationResponse{}, nil
}
