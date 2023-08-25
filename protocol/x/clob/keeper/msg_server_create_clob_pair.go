package keeper

import (
	"context"
	"errors"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k msgServer) CreateClobPair(
	goCtx context.Context,
	msg *types.MsgCreateClobPair,
) (*types.MsgCreateClobPairResponse, error) {
	// TODO(CORE-502): Implement message handler.
	return &types.MsgCreateClobPairResponse{}, errors.New("Not implemented")
}
