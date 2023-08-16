package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// UpdateSafetyParams updates the SafetyParams in state.
func (k msgServer) UpdateSafetyParams(
	goCtx context.Context,
	msg *types.MsgUpdateSafetyParams,
) (*types.MsgUpdateSafetyParamsResponse, error) {
	if k.Keeper.GetAuthority() != msg.Authority {
		return nil, errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority: expected %s, got %s",
			k.Keeper.GetAuthority(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateSafetyParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateSafetyParamsResponse{}, nil
}
