package keeper

import (
	"context"
	"errors"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

func (k msgServer) CreatePerpetual(
	goCtx context.Context,
	msg *types.MsgCreatePerpetual,
) (*types.MsgCreatePerpetualResponse, error) {
	return &types.MsgCreatePerpetualResponse{}, errors.New("Not implemented")
}
