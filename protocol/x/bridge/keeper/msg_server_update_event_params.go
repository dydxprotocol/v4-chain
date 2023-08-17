package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// UpdateEventParams updates the EventParams in state.
func (k msgServer) UpdateEventParams(
	goCtx context.Context,
	msg *types.MsgUpdateEventParams,
) (*types.MsgUpdateEventParamsResponse, error) {
	if k.Keeper.GetGovAuthority() != msg.Authority {
		return nil, errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority: expected %s, got %s",
			k.Keeper.GetGovAuthority(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateEventParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateEventParamsResponse{}, nil
}
