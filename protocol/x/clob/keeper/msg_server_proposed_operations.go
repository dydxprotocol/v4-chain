package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k msgServer) ProposedOperations(
	goCtx context.Context,
	msg *types.MsgProposedOperations,
) (*types.MsgProposedOperationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.ProcessProposerOperations(
		ctx,
		msg.GetOperationsQueue(),
	); err != nil {
		err = errorsmod.Wrapf(
			err,
			"Block height: %d",
			ctx.BlockHeight(),
		)
		ctx.Logger().Error(err.Error())
		return nil, err
	}

	return &types.MsgProposedOperationsResponse{}, nil
}
