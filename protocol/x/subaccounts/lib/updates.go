package lib

import (
	"fmt"
	"math/big"
	"sort"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	assetslib "github.com/dydxprotocol/v4-chain/protocol/x/assets/lib"
	perplib "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetSettledSubaccountWithPerpetuals returns 1. a new settled subaccount given an unsettled subaccount,
// updating the USDC AssetPosition, FundingIndex, and LastFundingPayment fields accordingly
// (does not persist any changes) and 2. a map with perpetual ID as key and last funding
// payment as value (for emitting funding payments to indexer).
func GetSettledSubaccountWithPerpetuals(
	subaccount types.Subaccount,
	perpInfos perptypes.PerpInfos,
) (
	settledSubaccount types.Subaccount,
	fundingPayments map[uint32]dtypes.SerializableInt,
) {
	totalNetSettlementPpm := big.NewInt(0)

	newPerpetualPositions := []*types.PerpetualPosition{}
	fundingPayments = make(map[uint32]dtypes.SerializableInt)

	// Iterate through and settle all perpetual positions.
	for _, p := range subaccount.PerpetualPositions {
		perpInfo := perpInfos.MustGet(p.PerpetualId)

		// Call the stateless utility function to get the net settlement and new funding index.
		bigNetSettlementPpm, newFundingIndex := perplib.GetSettlementPpmWithPerpetual(
			perpInfo.Perpetual,
			p.GetBigQuantums(),
			p.FundingIndex.BigInt(),
		)
		// Record non-zero funding payment (to be later emitted in SubaccountUpdateEvent to indexer).
		// Note: Funding payment is the negative of settlement, i.e. positive settlement is equivalent
		// to a negative funding payment (position received funding payment) and vice versa.
		if bigNetSettlementPpm.BitLen() != 0 {
			fundingPayments[p.PerpetualId] = dtypes.NewIntFromBigInt(
				new(big.Int).Neg(
					new(big.Int).Div(bigNetSettlementPpm, lib.BigIntOneMillion()),
				),
			)
		}

		// Aggregate all net settlements.
		totalNetSettlementPpm.Add(totalNetSettlementPpm, bigNetSettlementPpm)

		// Update cached funding index of the perpetual position.
		newPerpetualPositions = append(
			newPerpetualPositions, &types.PerpetualPosition{
				PerpetualId:  p.PerpetualId,
				Quantums:     p.Quantums,
				FundingIndex: dtypes.NewIntFromBigInt(newFundingIndex),
			},
		)
	}

	newSubaccount := types.Subaccount{
		Id:                 subaccount.Id,
		AssetPositions:     subaccount.AssetPositions,
		PerpetualPositions: newPerpetualPositions,
		MarginEnabled:      subaccount.MarginEnabled,
	}
	newUsdcPosition := new(big.Int).Add(
		subaccount.GetUsdcPosition(),
		// `Div` implements Euclidean division (unlike Go). When the diviser is positive,
		// division result always rounds towards negative infinity.
		totalNetSettlementPpm.Div(totalNetSettlementPpm, lib.BigIntOneMillion()),
	)
	// TODO(CLOB-993): Remove this function and use `UpdateAssetPositions` instead.
	newSubaccount.SetUsdcAssetPosition(newUsdcPosition)
	return newSubaccount, fundingPayments
}

// IsValidStateTransitionForUndercollateralizedSubaccount returns an `UpdateResult`
// denoting whether this state transition is valid. This function accepts the collateral and
// margin requirements of a subaccount before and after an update ("cur" and
// "new", respectively).
//
// This function should only be called if the account is undercollateralized after the update.
//
// A state transition is valid if the subaccount enters a
// "less-or-equally-risky" state after an update.
// i.e.`newNetCollateral / newMaintenanceMargin >= curNetCollateral / curMaintenanceMargin`.
//
// Otherwise, the state transition is invalid. If the account was previously undercollateralized,
// `types.StillUndercollateralized` is returned. If the account was previously
// collateralized and is now undercollateralized, `types.NewlyUndercollateralized` is
// returned.
//
// Note that the inequality `newNetCollateral / newMaintenanceMargin >= curNetCollateral / curMaintenanceMargin`
// has divide-by-zero issue when margin requirements are zero. To make sure the state
// transition is valid, we special case this scenario and only allow state transition that improves net collateral.
func IsValidStateTransitionForUndercollateralizedSubaccount(
	riskCur margin.Risk,
	riskNew margin.Risk,
) types.UpdateResult {
	// Determine whether the subaccount was previously undercollateralized before the update.
	var underCollateralizationResult = types.StillUndercollateralized
	if riskCur.IMR.Cmp(riskCur.NC) <= 0 {
		underCollateralizationResult = types.NewlyUndercollateralized
	}

	// If the maintenance margin is increasing, then the subaccount is undercollateralized.
	if riskNew.MMR.Cmp(riskCur.MMR) > 0 {
		return underCollateralizationResult
	}

	// If the maintenance margin is zero, it means the subaccount must have no open positions, and negative net
	// collateral. If the net collateral is not improving then this transition is not valid.
	if riskNew.MMR.BitLen() == 0 || riskCur.MMR.BitLen() == 0 {
		if riskNew.MMR.BitLen() == 0 &&
			riskCur.MMR.BitLen() == 0 &&
			riskNew.NC.Cmp(riskCur.NC) > 0 {
			return types.Success
		}

		return underCollateralizationResult
	}

	// Note that here we are effectively checking that
	// `newNetCollateral / newMaintenanceMargin >= curNetCollateral / curMaintenanceMargin`.
	// However, to avoid rounding errors, we factor this as
	// `newNetCollateral * curMaintenanceMargin >= curNetCollateral * newMaintenanceMargin`.
	newNcOldMmr := new(big.Int).Mul(riskNew.NC, riskCur.MMR)
	oldNcNewMmr := new(big.Int).Mul(riskCur.NC, riskNew.MMR)

	// The subaccount is not well-collateralized, and the state transition leaves the subaccount in a
	// "more-risky" state (collateral relative to margin requirements is decreasing).
	if oldNcNewMmr.Cmp(newNcOldMmr) > 0 {
		return underCollateralizationResult
	}

	// The subaccount is in a "less-or-equally-risky" state (margin requirements are decreasing or unchanged,
	// collateral relative to margin requirements is decreasing or unchanged).
	// This subaccount is undercollateralized in this state, but we still consider this state transition valid.
	return types.Success
}

// ApplyUpdatesToPositions merges a slice of `types.UpdatablePositions` and `types.PositionSize`
// (i.e. concrete types *types.AssetPosition` and `types.AssetUpdate`) into a slice of `types.PositionSize`.
// If a given `PositionSize` shares an ID with an `UpdatablePositionSize`, the update and position are merged
// into a single `PositionSize`.
//
// An error is returned if two updates share the same position id.
//
// Note: There are probably performance implications here for allocating a new slice of PositionSize,
// and for allocating new slices when converting the concrete types to interfaces. However, without doing
// this there would be a lot of duplicate code for calculating changes for both `Assets` and `Perpetuals`.
func ApplyUpdatesToPositions[
	P types.PositionSize,
	U types.PositionSize,
](positions []P, updates []U) ([]types.PositionSize, error) {
	var result []types.PositionSize = make([]types.PositionSize, 0, len(positions)+len(updates))

	updateMap := make(map[uint32]types.PositionSize, len(updates))
	updateIndexMap := make(map[uint32]int, len(updates))
	for i, update := range updates {
		// Check for non-unique updates (two updates to the same position).
		id := update.GetId()
		_, exists := updateMap[id]
		if exists {
			errMsg := fmt.Sprintf("Multiple updates exist for position %v", update.GetId())
			return nil, errorsmod.Wrap(types.ErrNonUniqueUpdatesPosition, errMsg)
		}

		updateMap[id] = update
		updateIndexMap[id] = i
		result = append(result, update)
	}

	// Iterate over each position, if the position shares an ID with
	// an update, then we "merge" the update and the position into a new `PositionUpdate`.
	for _, pos := range positions {
		id := pos.GetId()
		update, exists := updateMap[id]
		if !exists {
			result = append(result, pos)
		} else {
			var newPos = types.NewPositionUpdate(id)

			// Add the position size and update together to get the new size.
			var bigNewPositionSize = new(big.Int).Add(
				pos.GetBigQuantums(),
				update.GetBigQuantums(),
			)

			newPos.SetBigQuantums(bigNewPositionSize)

			// Replace update with `PositionUpdate`
			index := updateIndexMap[id]
			result[index] = newPos
		}
	}

	return result, nil
}

// GetUpdatedAssetPositions filters out all the asset positions on a subaccount that have
// been updated. This will include any asset postions that were closed due to an update.
// TODO(DEC-1295): look into reducing code duplication here using Generics+Reflect.
func GetUpdatedAssetPositions(
	update types.SettledUpdate,
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

// GetUpdatedPerpetualPositions filters out all the perpetual positions on a subaccount that have
// been updated. This will include any perpetual postions that were closed due to an update or that
// received / paid out funding payments..
func GetUpdatedPerpetualPositions(
	update types.SettledUpdate,
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
	settledUpdates []types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
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
				perpInfo := perpInfos.MustGet(pu.PerpetualId)
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
	settledUpdates []types.SettledUpdate,
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

// GetRiskForSubaccount returns the risk value of the `Subaccount` after updates are applied.
// It is used to get information about speculative changes to the `Subaccount`.
// The input subaccount must be settled.
//
// The provided update can also be "zeroed" in order to get information about
// the current state of the subaccount (i.e. with no changes).
//
// If two position updates reference the same position, an error is returned.
func GetRiskForSubaccount(
	settledUpdate types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) (
	risk margin.Risk,
	err error,
) {
	// Initialize return values.
	risk = margin.ZeroRisk()

	// Merge updates and assets.
	assetSizes, err := ApplyUpdatesToPositions(
		settledUpdate.SettledSubaccount.AssetPositions,
		settledUpdate.AssetUpdates,
	)
	if err != nil {
		return risk, err
	}

	// Merge updates and perpetuals.
	perpetualSizes, err := ApplyUpdatesToPositions(
		settledUpdate.SettledSubaccount.PerpetualPositions,
		settledUpdate.PerpetualUpdates,
	)
	if err != nil {
		return risk, err
	}

	// Iterate over all assets and updates and calculate change to net collateral and margin requirements.
	for _, size := range assetSizes {
		r, err := assetslib.GetNetCollateralAndMarginRequirements(
			size.GetId(),
			size.GetBigQuantums(),
		)
		if err != nil {
			return risk, err
		}
		risk.AddInPlace(r)
	}

	// Iterate over all perpetuals and updates and calculate change to net collateral and margin requirements.
	for _, size := range perpetualSizes {
		perpInfo := perpInfos.MustGet(size.GetId())
		r := perplib.GetNetCollateralAndMarginRequirements(
			perpInfo.Perpetual,
			perpInfo.Price,
			perpInfo.LiquidityTier,
			size.GetBigQuantums(),
		)
		risk.AddInPlace(r)
	}

	return risk, nil
}
