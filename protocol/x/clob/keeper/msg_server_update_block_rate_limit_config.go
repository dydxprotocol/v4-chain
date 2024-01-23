package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateBlockRateLimitConfiguration updates the block rate limit configuration returning an error
// if the configuration is invalid.
func (k msgServer) UpdateBlockRateLimitConfiguration(
	goCtx context.Context,
	msg *types.MsgUpdateBlockRateLimitConfiguration,
) (resp *types.MsgUpdateBlockRateLimitConfigurationResponse, err error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.Keeper.InitializeBlockRateLimit(ctx, msg.BlockRateLimitConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateBlockRateLimitConfigurationResponse{}, nil
}
