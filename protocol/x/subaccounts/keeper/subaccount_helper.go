package keeper

import (
	"sort"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// getUpdatedAssetPositions filters out all the asset positions on a subaccount that have
// been updated. This will include any asset postions that were closed due to an update.
// TODO(DEC-1295): look into reducing code duplication here using Generics+Reflect.
func getUpdatedAssetPositions(
	update SettledUpdate,
) []*types.AssetPosition {
	assetIdToPositionMap := make(map[uint32]*types.AssetPosition)
	for _, assetPosition := range update.SettledSubaccount.AssetPositions {
		assetIdToPositionMap[assetPosition.AssetId] = assetPosition
	}

	updatedAssetIds := make(map[uint32]struct{})
	for _, assetUpdate := range update.AssetUpdates {
		updatedAssetIds[assetUpdate.AssetId] = struct{}{}
	}

	updatedAssetPositions := make([]*types.AssetPosition, 0, len(updatedAssetIds))
	for updatedId := range updatedAssetIds {
		assetPosition, exists := assetIdToPositionMap[updatedId]
		// If a position does not exist on the subaccount with the asset id of an update, it must
		// have been deleted due to quantums becoming 0. This needs to be included in the event, so we
		// construct a position with the AssetId of the update and a Quantums value of 0. The other
		// properties are left as the default values as a 0-sized position indicates the position is
		// closed.
		if !exists {
			assetPosition = &types.AssetPosition{
				AssetId:  updatedId,
				Quantums: dtypes.ZeroInt(),
			}
		}
		updatedAssetPositions = append(updatedAssetPositions, assetPosition)
	}

	// Sort the asset positions in ascending order by asset id.
	sort.Slice(updatedAssetPositions, func(i, j int) bool {
		return updatedAssetPositions[i].GetId() < updatedAssetPositions[j].GetId()
	})

	return updatedAssetPositions
}

// getUpdatedPerpetualPositions filters out all the perpetual positions on a subaccount that have
// been updated. This will include any perpetual postions that were closed due to an update or that
// received / paid out funding payments..
func getUpdatedPerpetualPositions(
	update SettledUpdate,
	fundingPayments map[uint32]dtypes.SerializableInt,
) []*types.PerpetualPosition {
	perpetualIdToPositionMap := make(map[uint32]*types.PerpetualPosition)
	for _, perpetualPosition := range update.SettledSubaccount.PerpetualPositions {
		perpetualIdToPositionMap[perpetualPosition.PerpetualId] = perpetualPosition
	}

	// `updatedPerpetualIds` indicates which perpetuals were either explicitly updated
	// (through update.PerpetualUpdates) or implicitly updated (had non-zero last funding
	// payment).
	updatedPerpetualIds := make(map[uint32]struct{})
	for _, perpetualUpdate := range update.PerpetualUpdates {
		updatedPerpetualIds[perpetualUpdate.PerpetualId] = struct{}{}
	}
	// Mark perpetuals with non-zero funding payment also as updated.
	for perpetualIdWithNonZeroLastFunding := range fundingPayments {
		updatedPerpetualIds[perpetualIdWithNonZeroLastFunding] = struct{}{}
	}

	updatedPerpetualPositions := make([]*types.PerpetualPosition, 0, len(updatedPerpetualIds))
	for updatedId := range updatedPerpetualIds {
		perpetualPosition, exists := perpetualIdToPositionMap[updatedId]
		// If a position does not exist on the subaccount with the perpetual id of an update, it must
		// have been deleted due to quantums becoming 0. This needs to be included in the event, so we
		// construct a position with the PerpetualId of the update and a Quantums value of 0. The other
		// properties are left as the default values as a 0-sized position indicates the position is
		// closed and thus the funding index and the side of the position does not matter.
		if !exists {
			perpetualPosition = &types.PerpetualPosition{
				PerpetualId: updatedId,
				Quantums:    dtypes.ZeroInt(),
			}
		}
		updatedPerpetualPositions = append(updatedPerpetualPositions, perpetualPosition)
	}

	// Sort the perpetual positions in ascending order by perpetual id.
	sort.Slice(updatedPerpetualPositions, func(i, j int) bool {
		return updatedPerpetualPositions[i].GetId() < updatedPerpetualPositions[j].GetId()
	})

	return updatedPerpetualPositions
}

// For each settledUpdate in settledUpdates, updates its SettledSubaccount.PerpetualPositions
// to reflect settledUpdate.PerpetualUpdates.
// For newly created positions, use `perpIdToFundingIndex` map to populate the `FundingIndex` field.
func UpdatePerpetualPositions(
	settledUpdates []SettledUpdate,
	perpInfos map[uint32]types.PerpInfo,
) {
	// Apply the updates.
	for i, u := range settledUpdates {
		// Build a map of all the Subaccount's Perpetual Positions by id.
		perpetualPositionsMap := make(map[uint32]*types.PerpetualPosition)
		for _, pp := range u.SettledSubaccount.PerpetualPositions {
			perpetualPositionsMap[pp.PerpetualId] = pp
		}

		// Update the perpetual positions.
		for _, pu := range u.PerpetualUpdates {
			// Check if the `Subaccount` already has a position with the same id.
			// If so – we update the size of the existing position, otherwise
			// we create a new position.
			if pp, exists := perpetualPositionsMap[pu.PerpetualId]; exists {
				curQuantums := pp.GetBigQuantums()
				updateQuantums := pu.GetBigQuantums()
				newQuantums := curQuantums.Add(curQuantums, updateQuantums)

				// Handle the case where the position is now closed.
				if newQuantums.Sign() == 0 {
					delete(perpetualPositionsMap, pu.PerpetualId)
				}
				pp.Quantums = dtypes.NewIntFromBigInt(newQuantums)
			} else {
				// This subaccount does not have a matching position for this update.
				// Create the new position.
				perpInfo, exists := perpInfos[pu.PerpetualId]
				if !exists {
					// Invariant: `perpInfos` contains all existing perpetauls,
					// and perpetual position update must refer to an existing perpetual.
					panic(errorsmod.Wrapf(types.ErrPerpetualInfoDoesNotExist, "%d", pu.PerpetualId))
				}
				perpetualPosition := &types.PerpetualPosition{
					PerpetualId:  pu.PerpetualId,
					Quantums:     dtypes.NewIntFromBigInt(pu.GetBigQuantums()),
					FundingIndex: perpInfo.Perpetual.FundingIndex,
				}

				// Add the new position to the map.
				perpetualPositionsMap[pu.PerpetualId] = perpetualPosition
			}
		}

		// Convert the new PerpetualPostiion values back into a slice.
		perpetualPositions := make([]*types.PerpetualPosition, 0, len(perpetualPositionsMap))
		for _, value := range perpetualPositionsMap {
			perpetualPositions = append(perpetualPositions, value)
		}

		// Sort the new PerpetualPositions in ascending order by Id.
		sort.Slice(perpetualPositions, func(i, j int) bool {
			return perpetualPositions[i].GetId() < perpetualPositions[j].GetId()
		})

		settledUpdates[i].SettledSubaccount.PerpetualPositions = perpetualPositions
	}
}

// For each settledUpdate in settledUpdates, updates its SettledSubaccount.AssetPositions
// to reflect settledUpdate.AssetUpdates.
func UpdateAssetPositions(
	settledUpdates []SettledUpdate,
) {
	// Apply the updates.
	for i, u := range settledUpdates {
		// Build a map of all the Subaccount's Asset Positions by id.
		assetPositionsMap := make(map[uint32]*types.AssetPosition)
		for _, ap := range u.SettledSubaccount.AssetPositions {
			assetPositionsMap[ap.AssetId] = ap
		}

		// Update the asset positions.
		for _, au := range u.AssetUpdates {
			// Check if the `Subaccount` already has a position with the same id.
			// If so - we update the size of the existing position, otherwise
			// we create a new position.
			if ap, exists := assetPositionsMap[au.AssetId]; exists {
				curQuantums := ap.GetBigQuantums()
				updateQuantums := au.GetBigQuantums()
				newQuantums := curQuantums.Add(curQuantums, updateQuantums)

				ap.Quantums = dtypes.NewIntFromBigInt(newQuantums)

				// Handle the case where the position is now closed.
				if ap.Quantums.Sign() == 0 {
					delete(assetPositionsMap, au.AssetId)
				}
			} else {
				// This subaccount does not have a matching asset position for this update.

				// Create the new asset position.
				assetPosition := &types.AssetPosition{
					AssetId:  au.AssetId,
					Quantums: dtypes.NewIntFromBigInt(au.GetBigQuantums()),
				}

				// Add the new asset position to the map.
				assetPositionsMap[au.AssetId] = assetPosition
			}
		}

		// Convert the new AssetPostiion values back into a slice.
		assetPositions := make([]*types.AssetPosition, 0, len(assetPositionsMap))
		for _, value := range assetPositionsMap {
			assetPositions = append(assetPositions, value)
		}

		// Sort the new AssetPositions in ascending order by AssetId.
		sort.Slice(assetPositions, func(i, j int) bool {
			return assetPositions[i].GetId() < assetPositions[j].GetId()
		})

		settledUpdates[i].SettledSubaccount.AssetPositions = assetPositions
	}
}
