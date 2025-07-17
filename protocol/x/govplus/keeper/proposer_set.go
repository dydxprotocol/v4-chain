package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BlockProposerSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	fmt.Println("tian", "BlockProposerSetUpdates")
	powerReduction := k.stakingKeeper.PowerReduction(ctx)
	allValidators, err := k.stakingKeeper.GetAllValidators(ctx)
	fmt.Println("tian", "BlockProposerSetUpdates",
		"power reduction", powerReduction,
		"num validators", len(allValidators),
	)
	if err != nil {
		return nil, err
	}

	test := true
	for _, validator := range allValidators {
		power := math.NewInt(validator.ConsensusPower(powerReduction))
		fmt.Println("tian", "BlockProposerSetUpdates",
			"validator op addr", validator.OperatorAddress,
			// "pub key", validator.ConsensusPubkey,
			"power", power)
		if test {
			updates = append(updates, validator.ABCIValidatorUpdateCanPropose(
				powerReduction,
				true,
			))
			test = false
		} else {
			updates = append(updates, validator.ABCIValidatorUpdateCanPropose(
				powerReduction,
				false,
			))
		}
	}
	fmt.Println("tian", "BlockProposerSetUpdates",
		"num updates", len(updates),
		"updates", updates,
	)

	return updates, nil
}
