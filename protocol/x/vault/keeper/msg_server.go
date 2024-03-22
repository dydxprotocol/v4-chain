package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

type msgServer struct {
	Keeper types.VaultKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.VaultKeeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
