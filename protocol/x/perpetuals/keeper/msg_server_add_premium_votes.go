package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// AddPremiumVotes add new premium sample votes from a proposer to the application.
func (k msgServer) AddPremiumVotes(
	goCtx context.Context,
	msg *types.MsgAddPremiumVotes,
) (*types.MsgAddPremiumVotesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
