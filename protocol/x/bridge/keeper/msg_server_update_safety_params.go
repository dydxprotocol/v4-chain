package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/bridge/types"
)

// UpdateSafetyParams updates the SafetyParams in state.
func (k msgServer) UpdateSafetyParams(
	goCtx context.Context,
	msg *types.MsgUpdateSafetyParams,
) (*types.MsgUpdateSafetyParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateSafetyParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateSafetyParamsResponse{}, nil
}
