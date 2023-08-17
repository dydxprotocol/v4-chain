package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// UpdateProposeParams updates the ProposeParams in state.
func (k msgServer) UpdateProposeParams(
	goCtx context.Context,
	msg *types.MsgUpdateProposeParams,
) (*types.MsgUpdateProposeParamsResponse, error) {
	if k.Keeper.GetGovAuthority() != msg.Authority {
		return nil, errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority: expected %s, got %s",
			k.Keeper.GetGovAuthority(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateProposeParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateProposeParamsResponse{}, nil
}
