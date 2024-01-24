package keeper

import (
	"context"
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// AddPremiumVotes add new premium sample votes from a proposer to the application.
func (k msgServer) AddPremiumVotes(
	goCtx context.Context,
	msg *types.MsgAddPremiumVotes,
) (*types.MsgAddPremiumVotesResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Validate.
	if err := k.Keeper.PerformStatefulPremiumVotesValidation(ctx, msg); err != nil {
		panic(fmt.Sprintf(
			"PerformStatefulPremiumVotesValidation failed, err = %v",
			err,
		))
	}

	err := k.Keeper.AddPremiumVotes(
		ctx,
		msg.Votes,
	)

	if err != nil {
		panic(fmt.Sprintf(
			"AddPremiumVotes failed, err = %v",
			err,
		))
	}

	return &types.MsgAddPremiumVotesResponse{}, nil
}
