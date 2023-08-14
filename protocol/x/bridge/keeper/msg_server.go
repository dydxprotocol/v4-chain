package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

type msgServer struct {
	Keeper types.BridgeKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.BridgeKeeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
