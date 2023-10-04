package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorlib "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k msgServer) ProposedOperations(
	goCtx context.Context,
	msg *types.MsgProposedOperations,
) (resp *types.MsgProposedOperationsResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defer func() {
		if err != nil {
			errorlib.LogErrorWithBlockHeight(ctx, err)
		}
	}()

	if err := k.Keeper.ProcessProposerOperations(
		ctx,
		msg.GetOperationsQueue(),
	); err != nil {
		return nil, err
	}

	return &types.MsgProposedOperationsResponse{}, nil
}
