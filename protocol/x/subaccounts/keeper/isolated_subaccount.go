package keeper

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	settledUpdates []settledUpdate,
	perpetuals []perptypes.Perpetual,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	success = true
	successPerUpdate = make([]types.UpdateResult, len(settledUpdates))
	var perpIdToMarketType = make(map[uint32]perptypes.PerpetualMarketType)

	for _, perpetual := range perpetuals {
		perpIdToMarketType[perpetual.GetId()] = perpetual.Params.MarketType
	}

	for i, u := range settledUpdates {
		result, err := isValidIsolatedPerpetualUpdates(u, perpIdToMarketType)
		if err != nil {
			return false, nil, err
		}
		if result != types.Success {
			success = false
		}

		successPerUpdate[i] = result
	}

	return success, successPerUpdate, nil
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
	settledUpdate settledUpdate,
	perpIdToMarketType map[uint32]perptypes.PerpetualMarketType,
) (types.UpdateResult, error) {
	// If there are no perpetual updates, then this update does not violate constraints for isolated
	// markets.
	if len(settledUpdate.PerpetualUpdates) == 0 {
		return types.Success, nil
	}

	// Check if the updates contain an update to an isolated perpetual.
	hasIsolatedUpdate := false
	isolatedUpdatePerpetualId := uint32(math.MaxUint32)
	for _, perpetualUpdate := range settledUpdate.PerpetualUpdates {
		marketType, exists := perpIdToMarketType[perpetualUpdate.PerpetualId]
		if !exists {
			return types.UpdateCausedError, errorsmod.Wrap(
				perptypes.ErrPerpetualDoesNotExist, lib.UintToString(perpetualUpdate.PerpetualId),
			)
		}

		if marketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
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
		marketType, exists := perpIdToMarketType[perpetualPosition.PerpetualId]
		if !exists {
			return types.UpdateCausedError, errorsmod.Wrap(
				perptypes.ErrPerpetualDoesNotExist, lib.UintToString(perpetualPosition.PerpetualId),
			)
		}

		if marketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
			isIsolatedSubaccount = true
			isolatedPositionPerpetualId = perpetualPosition.PerpetualId
			break
		}
	}

	// A subaccount with a perpetual position in an isolated perpetual cannot have updates to other
	// non-isolated perpetuals.
	if isIsolatedSubaccount && !hasIsolatedUpdate {
		return types.ViolatesIsolatedSubaccountConstraints, nil
	}

	// A subaccount with perpetual positions in non-isolated perpetuals cannot have an update
	// to an isolated perpetual.
	if !isIsolatedSubaccount && hasPerpetualPositions && hasIsolatedUpdate {
		return types.ViolatesIsolatedSubaccountConstraints, nil
	}

	// There cannot be more than a single perpetual update if an update to an isolated perpetual
	// exists in the slice of perpetual updates.
	if hasIsolatedUpdate && len(settledUpdate.PerpetualUpdates) > 1 {
		return types.ViolatesIsolatedSubaccountConstraints, nil
	}

	// Note we can assume that if `hasIsolatedUpdate` is true, there is only a single perpetual
	// update for the subaccount, given the above check.
	// A subaccount with a perpetual position in an isolated perpetual cannot have an update to
	// another isolated perpetual.
	if isIsolatedSubaccount &&
		hasIsolatedUpdate &&
		isolatedPositionPerpetualId != isolatedUpdatePerpetualId {
		return types.ViolatesIsolatedSubaccountConstraints, nil
	}

	return types.Success, nil
}
