package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// UpdateSafetyParams updates the SafetyParams in state.
func (k msgServer) UpdateSafetyParams(
	goCtx context.Context,
	msg *types.MsgUpdateSafetyParams,
) (*types.MsgUpdateSafetyParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.GetAuthority()) {
		return nil, errors.Wrapf(
			types.ErrInvalidAuthority,
			"message authority %s is not valid for sending update safety params messages",
			msg.GetAuthority(),
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateSafetyParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateSafetyParamsResponse{}, nil
}
