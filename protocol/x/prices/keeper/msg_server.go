package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type msgServer struct {
	Keeper types.PricesKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.PricesKeeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
