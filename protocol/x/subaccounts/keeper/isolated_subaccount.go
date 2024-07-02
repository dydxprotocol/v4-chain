package keeper

import (
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// checkIsolatedSubaccountConstaints will validate all `updates` to the relevant subaccounts against
// isolated subaccount constraints.
// This function checks each update in isolation, so if multiple updates for the same subaccount id
// are passed in, they are not evaluated separately.
// The input subaccounts must be settled.
//
// Returns a `success` value of `true` if all updates are valid.
// Returns a `successPerUpdates` value, which is a slice of `UpdateResult`.
// These map to the updates and are used to indicate which of the updates
// caused a failure, if any.
func (k Keeper) checkIsolatedSubaccountConstraints(
	ctx sdk.Context,
	settledUpdates []types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
) {
	success = true
	successPerUpdate = make([]types.UpdateResult, len(settledUpdates))

	for i, u := range settledUpdates {
		result := isValidIsolatedPerpetualUpdates(u, perpInfos)
		if result != types.Success {
			success = false
		}

		successPerUpdate[i] = result
	}

	return success, successPerUpdate
}

// Checks whether the perpetual updates to a settled subaccount violates constraints for isolated
// perpetuals. This function assumes the settled subaccount is valid and does not violate the
// the constraints.
// The constraint being checked is:
//   - a subaccount with a position in an isolated perpetual cannot have updates for other
//     perpetuals
//   - a subaccount with a position in a non-isolated perpetual cannot have updates for isolated
//     perpetuals
//   - a subaccount with no positions cannot be updated to have positions in multiple isolated
//     perpetuals or a combination of isolated and non-isolated perpetuals
func isValidIsolatedPerpetualUpdates(
	settledUpdate types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) types.UpdateResult {
	// If there are no perpetual updates, then this update does not violate constraints for isolated
	// markets.
	if len(settledUpdate.PerpetualUpdates) == 0 {
		return types.Success
	}

	// Check if the updates contain an update to an isolated perpetual.
	hasIsolatedUpdate := false
	isolatedUpdatePerpetualId := uint32(math.MaxUint32)
	for _, perpetualUpdate := range settledUpdate.PerpetualUpdates {
		perpInfo := perpInfos.MustGet(perpetualUpdate.PerpetualId)

		if perpInfo.Perpetual.Params.MarketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
			hasIsolatedUpdate = true
			isolatedUpdatePerpetualId = perpetualUpdate.PerpetualId
			break
		}
	}

	// Check if the subaccount has a position in an isolated perpetual.
	// Assumes the subaccount itself does not violate the isolated perpetual constraints.
	isIsolatedSubaccount := false
	isolatedPositionPerpetualId := uint32(math.MaxUint32)
	hasPerpetualPositions := len(settledUpdate.SettledSubaccount.PerpetualPositions) > 0
	for _, perpetualPosition := range settledUpdate.SettledSubaccount.PerpetualPositions {
		perpInfo := perpInfos.MustGet(perpetualPosition.PerpetualId)

		if perpInfo.Perpetual.Params.MarketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
			isIsolatedSubaccount = true
			isolatedPositionPerpetualId = perpetualPosition.PerpetualId
			break
		}
	}

	// A subaccount with a perpetual position in an isolated perpetual cannot have updates to other
	// non-isolated perpetuals.
	if isIsolatedSubaccount && !hasIsolatedUpdate {
		return types.ViolatesIsolatedSubaccountConstraints
	}

	// A subaccount with perpetual positions in non-isolated perpetuals cannot have an update
	// to an isolated perpetual.
	if !isIsolatedSubaccount && hasPerpetualPositions && hasIsolatedUpdate {
		return types.ViolatesIsolatedSubaccountConstraints
	}

	// There cannot be more than a single perpetual update if an update to an isolated perpetual
	// exists in the slice of perpetual updates.
	if hasIsolatedUpdate && len(settledUpdate.PerpetualUpdates) > 1 {
		return types.ViolatesIsolatedSubaccountConstraints
	}

	// Note we can assume that if `hasIsolatedUpdate` is true, there is only a single perpetual
	// update for the subaccount, given the above check.
	// A subaccount with a perpetual position in an isolated perpetual cannot have an update to
	// another isolated perpetual.
	if isIsolatedSubaccount &&
		hasIsolatedUpdate &&
		isolatedPositionPerpetualId != isolatedUpdatePerpetualId {
		return types.ViolatesIsolatedSubaccountConstraints
	}

	return types.Success
}

// GetIsolatedPerpetualStateTransition computes whether an isolated perpetual position will be
// opened or closed for a subaccount.
// This function assumes that the subaccount is valid under isolated perpetual constraints.
// The input `settledUpdate` must have an updated subaccount (`settledUpdate.SettledSubaccount`),
// so all the updates must have been applied already to the subaccount.
func GetIsolatedPerpetualStateTransition(
	settledUpdateWithUpdatedSubaccount types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) (*types.IsolatedPerpetualPositionStateTransition, error) {
	// This subaccount needs to have had the updates in the `settledUpdate` already applied to it.
	updatedSubaccount := settledUpdateWithUpdatedSubaccount.SettledSubaccount
	// If there are no perpetual updates, then no perpetual position could have been opened or closed
	// on the subaccount.
	if len(settledUpdateWithUpdatedSubaccount.PerpetualUpdates) == 0 {
		return nil, nil
	}

	// If there are more than 1 valid perpetual update, or more than 1 valid perpetual position on the
	// subaccount, it is not isolated to an isolated perpetual, and so no isolated perpetual position
	// could have been opened or closed.
	if len(settledUpdateWithUpdatedSubaccount.PerpetualUpdates) > 1 ||
		len(updatedSubaccount.PerpetualPositions) > 1 {
		return nil, nil
	}

	// Now, from the above checks, we know there is only a single perpetual update and 0 or 1 perpetual
	// positions.
	perpetualUpdate := settledUpdateWithUpdatedSubaccount.PerpetualUpdates[0]
	perpInfo := perpInfos.MustGet(perpetualUpdate.PerpetualId)
	// If the perpetual update is not for an isolated perpetual, no isolated perpetual position is
	// being opened or closed.
	if perpInfo.Perpetual.Params.MarketType != perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
		return nil, nil
	}

	// If the updated subaccount does not have any perpetual positions, then an isolated perpetual
	// position must have been closed due to the perpetual update.
	if len(updatedSubaccount.PerpetualPositions) == 0 {
		return &types.IsolatedPerpetualPositionStateTransition{
			SubaccountId:  updatedSubaccount.Id,
			PerpetualId:   perpetualUpdate.PerpetualId,
			QuoteQuantums: updatedSubaccount.GetUsdcPosition(),
			Transition:    types.Closed,
		}, nil
	}

	// After the above checks, the subaccount must have only a single perpetual position, which is for
	// the same isolated perpetual as the perpetual update.
	perpetualPosition := updatedSubaccount.PerpetualPositions[0]
	// If the size of the update and the position are the same, the perpetual update must have opened
	// the position.
	if perpetualUpdate.GetBigQuantums().Cmp(perpetualPosition.GetBigQuantums()) == 0 {
		if len(settledUpdateWithUpdatedSubaccount.AssetUpdates) != 1 {
			return nil, errorsmod.Wrapf(
				types.ErrFailedToUpdateSubaccounts,
				"Subaccount with id %v opened perpteual position with perpetual id %d with invalid number of"+
					" changes to asset positions (%d), should only be 1 asset update",
				updatedSubaccount.Id,
				perpetualUpdate.PerpetualId,
				len(settledUpdateWithUpdatedSubaccount.AssetUpdates),
			)
		}
		if settledUpdateWithUpdatedSubaccount.AssetUpdates[0].AssetId != assettypes.AssetUsdc.Id {
			return nil, errorsmod.Wrapf(
				types.ErrFailedToUpdateSubaccounts,
				"Subaccount with id %v opened perpteual position with perpetual id %d without a change to the"+
					" quote currency's asset position.",
				updatedSubaccount.Id,
				perpetualUpdate.PerpetualId,
			)
		}
		// Collateral equal to the quote currency asset position before the update was applied needs to be transferred.
		// Subtract the delta from the updated subaccount's quote currency asset position size to get the size
		// of the quote currency asset position.
		quoteQuantumsBeforeUpdate := new(big.Int).Sub(
			updatedSubaccount.GetUsdcPosition(),
			settledUpdateWithUpdatedSubaccount.AssetUpdates[0].GetBigQuantums(),
		)
		return &types.IsolatedPerpetualPositionStateTransition{
			SubaccountId:  updatedSubaccount.Id,
			PerpetualId:   perpetualUpdate.PerpetualId,
			QuoteQuantums: quoteQuantumsBeforeUpdate,
			Transition:    types.Opened,
		}, nil
	}

	// The isolated perpetual position changed size but was not opened or closed.
	return nil, nil
}

// transferCollateralForIsolatedPerpetual transfers collateral between an isolated collateral pool
// and the cross-perpetual collateral pool based on whether an isolated perpetual position was
// opened or closed in a subaccount.
// Note: This uses the `x/bank` keeper and modifies `x/bank` state.
func (k *Keeper) transferCollateralForIsolatedPerpetual(
	ctx sdk.Context,
	stateTransition *types.IsolatedPerpetualPositionStateTransition,
) error {
	// No collateral to transfer if no state transition.
	if stateTransition == nil {
		return nil
	}

	// If there are zero quantums to transfer, don't transfer collateral.
	if stateTransition.QuoteQuantums.Sign() == 0 {
		return nil
	}

	isolatedCollateralPoolAddr, err := k.GetCollateralPoolFromPerpetualId(ctx, stateTransition.PerpetualId)
	if err != nil {
		return err
	}
	var toModuleAddr sdk.AccAddress
	var fromModuleAddr sdk.AccAddress

	// If an isolated perpetual position was opened in the subaccount, then move collateral from the
	// cross-perpetual collateral pool to the isolated perpetual collateral pool.
	if stateTransition.Transition == types.Opened {
		toModuleAddr = isolatedCollateralPoolAddr
		fromModuleAddr = types.ModuleAddress
		// If the isolated perpetual position was closed, then move collateral from the isolated
		// perpetual collateral pool to the cross-perpetual collateral pool.
	} else if stateTransition.Transition == types.Closed {
		toModuleAddr = types.ModuleAddress
		fromModuleAddr = isolatedCollateralPoolAddr
	} else {
		// Should never hit this.
		return errorsmod.Wrapf(
			types.ErrFailedToUpdateSubaccounts,
			"Invalid state transition %v for isolated perpetual with id %d in subaccount with id %v",
			stateTransition,
			stateTransition.PerpetualId,
			stateTransition.SubaccountId,
		)
	}

	// Invalid to transfer negative quantums. This should already be caught by collateralization
	// checks as well.
	if stateTransition.QuoteQuantums.Sign() == -1 {
		return errorsmod.Wrapf(
			types.ErrFailedToUpdateSubaccounts,
			"Subaccount with id %v %s perpteual position with perpetual id %d with negative collateral %s to transfer",
			stateTransition.SubaccountId,
			stateTransition.Transition.String(),
			stateTransition.PerpetualId,
			stateTransition.QuoteQuantums.String(),
		)
	}

	// Transfer collateral between collateral pools.
	_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
		ctx,
		// TODO(DEC-715): Support non-USDC assets.
		assettypes.AssetUsdc.Id,
		stateTransition.QuoteQuantums,
	)
	if err != nil {
		return err
	}

	if err = k.bankKeeper.SendCoins(
		ctx,
		fromModuleAddr,
		toModuleAddr,
		[]sdk.Coin{coinToTransfer},
	); err != nil {
		return err
	}

	return nil
}

// computeAndExecuteCollateralTransfer computes collateral transfers resulting from updates to
// a subaccount and executes the collateral transfer using `x/bank`.`
// The input `settledUpdate` must have an updated subaccount (`settledUpdate.SettledSubaccount`),
// so all the updates must have been applied already to the subaccount.
// Note: This uses the `x/bank` keeper and modifies `x/bank` state.
func (k *Keeper) computeAndExecuteCollateralTransfer(
	ctx sdk.Context,
	settledUpdateWithUpdatedSubaccount types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) error {
	// The subaccount in `settledUpdateWithUpdatedSubaccount` already has the perpetual updates
	// and asset updates applied to it.
	stateTransition, err := GetIsolatedPerpetualStateTransition(
		settledUpdateWithUpdatedSubaccount,
		perpInfos,
	)
	if err != nil {
		return err
	}
	if err := k.transferCollateralForIsolatedPerpetual(
		ctx,
		stateTransition,
	); err != nil {
		return err
	}

	return nil
}
