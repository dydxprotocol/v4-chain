package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateBlockRateLimitConfiguration updates the equity tier limit configuration returning an error
// if the configuration is invalid.
func (k msgServer) UpdateEquityTierLimitConfiguration(
	goCtx context.Context,
	configuration *types.MsgUpdateEquityTierLimitConfiguration,
) (*types.MsgUpdateEquityTierLimitConfigurationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Keeper.InitializeEquityTierLimit(ctx, configuration.EquityTierLimitConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateEquityTierLimitConfigurationResponse{}, nil
}
