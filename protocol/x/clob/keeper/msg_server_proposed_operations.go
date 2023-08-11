package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/x/clob/types"
)

func (k msgServer) ProposedOperations(
	goCtx context.Context,
	msg *types.MsgProposedOperations,
) (*types.MsgProposedOperationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addToOrderbookCollatCheckOrderHashesSet := msg.GetAddToOrderbookCollatCheckOrderHashesSet()

	if err := k.Keeper.ProcessProposerOperations(
		ctx,
		msg.OperationsQueue,
		addToOrderbookCollatCheckOrderHashesSet,
	); err != nil {
		panic(
			sdkerrors.Wrapf(
				err,
				"Block height: %d",
				ctx.BlockHeight(),
			),
		)
	}

	return &types.MsgProposedOperationsResponse{}, nil
}
