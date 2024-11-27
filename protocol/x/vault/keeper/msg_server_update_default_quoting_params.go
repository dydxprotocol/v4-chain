package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// UpdateDefaultQuotingParams updates the default quoting parameters of the vault module.
func (k msgServer) UpdateDefaultQuotingParams(
	goCtx context.Context,
	msg *types.MsgUpdateDefaultQuotingParams,
) (*types.MsgUpdateDefaultQuotingParamsResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	operator := k.GetOperatorParams(ctx).Operator

	// Check if authority is valid (must be a module authority or operator).
	if !k.HasAuthority(msg.Authority) && msg.Authority != operator {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidAuthority,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.SetDefaultQuotingParams(ctx, msg.DefaultQuotingParams); err != nil {
		return nil, err
	}

	return &types.MsgUpdateDefaultQuotingParamsResponse{}, nil
}
