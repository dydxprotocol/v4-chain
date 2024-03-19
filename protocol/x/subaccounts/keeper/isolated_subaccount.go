package keeper

import (
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type positionStateTransition uint

const (
	opened positionStateTransition = iota
	closed
)

// Represents a state transition for an isolated perpetual.
type isolatedPerpetualStateTranisition struct {
	perpetualId uint32
	// TODO(DEC-715): Support non-USDC assets.
	// USDC position size of the subaccount that has a state change for an isolated perpetual.
	usdcQuantumsBeforeUpdate *big.Int
	transition               positionStateTransition
}

// checkIsolatedSubaccountConstaints will validate all `updates` to the relevant subaccounts against
// isolated subaccount constraints and computes whether each update leads to a state change
// (open / close) to an isolated perpetual occurred due to the updates if they are valid.
// This function checks each update in isolation, so if multiple updates for the same subaccount id
// are passed in, they are not evaluated separately.
// The input subaccounts must be settled.
//
// Returns a `success` value of `true` if all updates are valid.
// Returns a `successPerUpdates` value, which is a slice of `UpdateResult`.
// These map to the updates and are used to indicate which of the updates
// caused a failure, if any.
// Returns a `isolatedPerpetualStateTransitions` value, which is a slice of
// `isolatedPerpetualStateTransition`.
// These map to the updates and are used to indicate which of the updates opened or closed an
// isolated perpetual position.
func (k Keeper) checkIsolatedSubaccountConstraints(
	ctx sdk.Context,
	settledUpdates []settledUpdate,
	perpetuals []perptypes.Perpetual,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	isolatedPerpetualStateTranisitions []*isolatedPerpetualStateTranisition,
	err error,
) {
	success = true
	successPerUpdate = make([]types.UpdateResult, len(settledUpdates))
	isolatedPerpetualStateTranisitions = make([]*isolatedPerpetualStateTranisition, len(settledUpdates))
	var perpIdToMarketType = make(map[uint32]perptypes.PerpetualMarketType)

	for _, perpetual := range perpetuals {
		perpIdToMarketType[perpetual.GetId()] = perpetual.Params.MarketType
	}

	for i, u := range settledUpdates {
		result, stateTransition, err := processIsolatedPerpetualUpdates(u, perpIdToMarketType)
		if err != nil {
			return false, nil, nil, err
		}
		if result != types.Success {
			success = false
		}

		successPerUpdate[i] = result
		isolatedPerpetualStateTranisitions[i] = stateTransition
	}

	return success, successPerUpdate, isolatedPerpetualStateTranisitions, nil
}

// processIsolatedPerpetualUpdates checks whether the perpetual updates to a settled subaccount
// violates constraints for isolated perpetuals and computes whether the perpetual updates result in
// a state change (open / close) for an isolated perpetual position if the updates are valid.
// This function assumes the settled subaccount is valid and does not violate the constraints.
// The constraint being checked is:
//   - a subaccount with a position in an isolated perpetual cannot have updates for other
//     perpetuals
//   - a subaccount with a position in a non-isolated perpetual cannot have updates for isolated
//     perpetuals
//   - a subaccount with no positions cannot be updated to have positions in multiple isolated
//     perpetuals or a combination of isolated and non-isolated perpetuals
//
// If there is a state change (open / close) from the perpetual updates, it is returned along with
// the perpetual id of the isolated perpetual and the size of the USDC position in the subaccount.
func processIsolatedPerpetualUpdates(
	settledUpdate settledUpdate,
	perpIdToMarketType map[uint32]perptypes.PerpetualMarketType,
) (types.UpdateResult, *isolatedPerpetualStateTranisition, error) {
	// If there are no perpetual updates, then this update does not violate constraints for isolated
	// markets.
	if len(settledUpdate.PerpetualUpdates) == 0 {
		return types.Success, nil, nil
	}

	// Check if the updates contain an update to an isolated perpetual.
	hasIsolatedUpdate := false
	isolatedUpdatePerpetualId := uint32(math.MaxUint32)
	isolatedUpdate := (*types.PerpetualUpdate)(nil)
	for _, perpetualUpdate := range settledUpdate.PerpetualUpdates {
		marketType, exists := perpIdToMarketType[perpetualUpdate.PerpetualId]
		if !exists {
			return types.UpdateCausedError, nil, errorsmod.Wrap(
				perptypes.ErrPerpetualDoesNotExist, lib.UintToString(perpetualUpdate.PerpetualId),
			)
		}

		if marketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
			hasIsolatedUpdate = true
			isolatedUpdatePerpetualId = perpetualUpdate.PerpetualId
			isolatedUpdate = &types.PerpetualUpdate{
				PerpetualId:      perpetualUpdate.PerpetualId,
				BigQuantumsDelta: perpetualUpdate.GetBigQuantums(),
			}
			break
		}
	}

	// Check if the subaccount has a position in an isolated perpetual.
	// Assumes the subaccount itself does not violate the isolated perpetual constraints.
	isIsolatedSubaccount := false
	isolatedPositionPerpetualId := uint32(math.MaxUint32)
	isolatedPerpetualPosition := (*types.PerpetualPosition)(nil)
	hasPerpetualPositions := len(settledUpdate.SettledSubaccount.PerpetualPositions) > 0
	for _, perpetualPosition := range settledUpdate.SettledSubaccount.PerpetualPositions {
		marketType, exists := perpIdToMarketType[perpetualPosition.PerpetualId]
		if !exists {
			return types.UpdateCausedError, nil, errorsmod.Wrap(
				perptypes.ErrPerpetualDoesNotExist, lib.UintToString(perpetualPosition.PerpetualId),
			)
		}

		if marketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
			isIsolatedSubaccount = true
			isolatedPositionPerpetualId = perpetualPosition.PerpetualId
			isolatedPerpetualPosition = perpetualPosition
			break
		}
	}

	// A subaccount with a perpetual position in an isolated perpetual cannot have updates to other
	// non-isolated perpetuals.
	if isIsolatedSubaccount && !hasIsolatedUpdate {
		return types.ViolatesIsolatedSubaccountConstraints, nil, nil
	}

	// A subaccount with perpetual positions in non-isolated perpetuals cannot have an update
	// to an isolated perpetual.
	if !isIsolatedSubaccount && hasPerpetualPositions && hasIsolatedUpdate {
		return types.ViolatesIsolatedSubaccountConstraints, nil, nil
	}

	// There cannot be more than a single perpetual update if an update to an isolated perpetual
	// exists in the slice of perpetual updates.
	if hasIsolatedUpdate && len(settledUpdate.PerpetualUpdates) > 1 {
		return types.ViolatesIsolatedSubaccountConstraints, nil, nil
	}

	// Note we can assume that if `hasIsolatedUpdate` is true, there is only a single perpetual
	// update for the subaccount, given the above check.
	// A subaccount with a perpetual position in an isolated perpetual cannot have an update to
	// another isolated perpetual.
	if isIsolatedSubaccount &&
		hasIsolatedUpdate &&
		isolatedPositionPerpetualId != isolatedUpdatePerpetualId {
		return types.ViolatesIsolatedSubaccountConstraints, nil, nil
	}

	return types.Success,
		getIsolatedPerpetualStateTransition(
			settledUpdate.SettledSubaccount,
			isolatedPerpetualPosition,
			isolatedUpdate,
		),
		nil
}

// getIsolatedPerpetualStateTransition computes whether an isolated perpetual position will be
// opened or closed for a subaccount given an isolated perpetual update for the subaccount.
// Input subaccount account must be settled.
func getIsolatedPerpetualStateTransition(
	settledSubaccount types.Subaccount,
	isolatedPerpetualPosition *types.PerpetualPosition,
	isolatedPerpetualUpdate *types.PerpetualUpdate,
) *isolatedPerpetualStateTranisition {
	// If there is no update to an isolated perpetual position, then no state transitions have
	// happened for the isolated perpetual.
	if isolatedPerpetualUpdate == nil {
		return nil
	}

	perpetualId := isolatedPerpetualUpdate.PerpetualId
	// TODO(DEC-715): Support non-USDC assets.
	usdcQuantumsBeforeUpdate := new(big.Int).Set(settledSubaccount.GetUsdcPosition())

	// If the subaccount has no isolated perpetual position, then this update is opening an isolated
	// perpetual position.
	if isolatedPerpetualPosition == nil {
		return &isolatedPerpetualStateTranisition{
			perpetualId:              perpetualId,
			usdcQuantumsBeforeUpdate: usdcQuantumsBeforeUpdate,
			transition:               opened,
		}
	}

	// If the position size after the update is zero, then this update is closing an isolated
	// perpetual position.
	if finalPositionSize := new(big.Int).Add(
		isolatedPerpetualPosition.GetBigQuantums(),
		isolatedPerpetualUpdate.GetBigQuantums(),
	); finalPositionSize.Cmp(lib.BigInt0()) == 0 {
		return &isolatedPerpetualStateTranisition{
			perpetualId:              perpetualId,
			usdcQuantumsBeforeUpdate: usdcQuantumsBeforeUpdate,
			transition:               closed,
		}
	}

	// If a position was not opened or closed, no state transition happened from the perpetual
	// update.
	return nil
}

// transferCollateralForIsolatedPerpetual transfers collateral between an isolated collateral pool
// and the cross-perpetual collateral pool based on whether an isolated perpetual position was
// opened or closed in a subaccount.
// Note: This uses the `x/bank` keeper and modifies `x/bank` state.
func (k *Keeper) transferCollateralForIsolatedPerpetual(
	ctx sdk.Context,
	updatedSubaccount types.Subaccount,
	stateTransition *isolatedPerpetualStateTranisition,
) error {
	// No collateral to transfer if no state transition.
	if stateTransition == nil {
		return nil
	}

	isolatedCollateralPoolAddr, err := k.GetCollateralPoolFromPerpetualId(ctx, stateTransition.perpetualId)
	if err != nil {
		return err
	}

	// If an isolated perpetual position was opened in the subaccount, then move collateral equivalent
	// to the USDC asset position size of the subaccount before the update from the
	// cross-perpetual collateral pool to the isolated perpetual collateral pool.
	if stateTransition.transition == opened {
		_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
			ctx,
			// TODO(DEC-715): Support non-USDC assets.
			assettypes.AssetUsdc.Id,
			stateTransition.usdcQuantumsBeforeUpdate,
		)
		if err != nil {
			return err
		}

		if err := k.bankKeeper.SendCoins(
			ctx,
			types.ModuleAddress,
			isolatedCollateralPoolAddr,
			[]sdk.Coin{coinToTransfer},
		); err != nil {
			return err
		}
		return nil
		// If the isolated perpetual position was closed, then move collateral equivalent to the USDC
		// asset position size of the subaccount after the update from the isolated perpetual collateral
		// pool to the cross-perpetual collateral pool.
	} else if stateTransition.transition == closed {
		_, coinToTransfer, err := k.assetsKeeper.ConvertAssetToCoin(
			ctx,
			// TODO(DEC-715): Support non-USDC assets.
			assettypes.AssetUsdc.Id,
			updatedSubaccount.GetUsdcPosition(),
		)
		if err != nil {
			return err
		}

		if err := k.bankKeeper.SendCoins(
			ctx,
			isolatedCollateralPoolAddr,
			types.ModuleAddress,
			[]sdk.Coin{coinToTransfer},
		); err != nil {
			return err
		}
		return nil
	}

	return nil
}
