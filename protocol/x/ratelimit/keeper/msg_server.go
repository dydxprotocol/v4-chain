package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetLimitParams(
	goCtx context.Context,
	msg *types.MsgSetLimitParams,
) (*types.MsgSetLimitParamsResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// TODO(CORE-825): Implement messages.

	return &types.MsgSetLimitParamsResponse{}, nil
}

func (k msgServer) DeleteLimitParams(
	goCtx context.Context,
	msg *types.MsgDeleteLimitParams,
) (*types.MsgDeleteLimitParamsResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// TODO(CORE-825): Implement messages.

	return &types.MsgDeleteLimitParamsResponse{}, nil
}
