package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// UpdateOperatorParams updates the operator parameters of megavault.
func (k msgServer) UpdateOperatorParams(
	goCtx context.Context,
	msg *types.MsgUpdateOperatorParams,
) (*types.MsgUpdateOperatorParamsResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	if err := k.SetOperatorParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateOperatorParamsResponse{}, nil
}
