package keeper

import (
	"context"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
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

func (k msgServer) ClaimYieldForSubaccount(
	goCtx context.Context,
	msg *types.MsgClaimYieldForSubaccount,
) (
	response *types.MsgClaimYieldForSubaccountResponse,
	err error,
) {

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	err = k.ClaimYieldForSubaccountFromIdAndSetNewState(ctx, msg.Id)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimYieldForSubaccountResponse{}, nil
}
