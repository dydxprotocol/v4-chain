package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) getProposer(ctx sdk.Context) (stakingtypes.Validator, error) {
	proposerConsAddress := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	proposer, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, proposerConsAddress)
	if !found {
		k.Logger(ctx).Error(
			"Failed to get proposer by consensus address",
			"proposer",
			proposerConsAddress.String(),
		)
		return stakingtypes.Validator{}, errors.New("failed to get proposer by consensus address")
	}

	return proposer, nil
}
