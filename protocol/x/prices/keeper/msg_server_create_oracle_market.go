package keeper

import (
	"context"
	"errors"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func (k msgServer) CreateOracleMarket(
	goCtx context.Context,
	msg *types.MsgCreateOracleMarket,
) (*types.MsgCreateOracleMarketResponse, error) {
	return &types.MsgCreateOracleMarketResponse{}, errors.New("Not implemented")
}
