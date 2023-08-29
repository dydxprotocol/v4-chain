package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// UpdateProposeParams updates the ProposeParams in state.
func (k msgServer) UpdateProposeParams(
	goCtx context.Context,
	msg *types.MsgUpdateProposeParams,
) (*types.MsgUpdateProposeParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.GetAuthority()) {
		return nil, errors.Wrapf(
			types.ErrInvalidAuthority,
			"message authority %s is not valid for sending update propose params messages",
			msg.GetAuthority(),
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateProposeParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateProposeParamsResponse{}, nil
}
