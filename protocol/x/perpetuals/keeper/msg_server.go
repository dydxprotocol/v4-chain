package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

type msgServer struct {
	Keeper types.PerpetualsKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.PerpetualsKeeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
