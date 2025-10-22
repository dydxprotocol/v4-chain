package lib

import (
	"math/big"
	"sort"

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
				QuoteBalance: p.QuoteBalance,
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

	// The subaccount is not well-collateralized, and the state transition leaves the subaccount in a
	// "more-risky" state (collateral relative to margin requirements is decreasing).
	if riskNew.Cmp(riskCur) > 0 {
		return underCollateralizationResult
	}

	// The subaccount is in a "less-or-equally-risky" state (margin requirements are decreasing or unchanged,
	// collateral relative to margin requirements is decreasing or unchanged).
	// This subaccount is undercollateralized in this state, but we still consider this state transition valid.
	return types.Success
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

// CalculateUpdatedAssetPositions returns a deep-copy of the asset positions with the updates applied.
func CalculateUpdatedAssetPositions(
	assetPositions []*types.AssetPosition,
	updates []types.AssetUpdate,
) []*types.AssetPosition {
	// Create a map of asset positions by ID.
	positionsMap := make(map[uint32]*types.AssetPosition)
	for _, pos := range assetPositions {
		copy := pos.DeepCopy()
		positionsMap[pos.AssetId] = &copy
	}

	// Iterate over each update and apply it to the positions.
	for _, update := range updates {
		// Check if the position already exists.
		pos, exists := positionsMap[update.AssetId]
		if exists {
			// Update the position.
			quantums := pos.GetBigQuantums()
			quantums.Add(quantums, update.GetBigQuantums())
			if quantums.BitLen() == 0 {
				// The position is now closed.
				delete(positionsMap, update.AssetId)
			} else {
				pos.Quantums = dtypes.NewIntFromBigInt(quantums)
			}
		} else {
			// Create a new position.
			positionsMap[update.AssetId] = &types.AssetPosition{
				AssetId:  update.AssetId,
				Quantums: dtypes.NewIntFromBigInt(update.GetBigQuantums()),
			}
		}
	}

	return lib.MapToSortedSlice[lib.Sortable[uint32]](positionsMap)
}

// CalculateUpdatedPerpetualPositions returns a deep-copy of the perpetual positions with the updates applied.
func CalculateUpdatedPerpetualPositions(
	positions []*types.PerpetualPosition,
	updates []types.PerpetualUpdate,
	perpInfos perptypes.PerpInfos,
) []*types.PerpetualPosition {
	// Create a map of perpetual positions by ID.
	positionsMap := make(map[uint32]*types.PerpetualPosition)
	for _, pos := range positions {
		copy := pos.DeepCopy()
		positionsMap[pos.PerpetualId] = &copy
	}

	// Iterate over each update and apply it to the positions.
	for _, update := range updates {
		// Check if the position already exists.
		pos, exists := positionsMap[update.PerpetualId]
		if exists {
			// Update existing position.
			quantums := pos.GetBigQuantums()
			quantums.Add(quantums, update.GetBigQuantums())

			quoteBalance := pos.GetQuoteBalance()
			quoteBalance.Add(quoteBalance, update.GetBigQuoteBalance())

			if quantums.BitLen() == 0 && quoteBalance.BitLen() == 0 {
				// The position is now closed.
				delete(positionsMap, update.PerpetualId)
			} else {
				pos.Quantums = dtypes.NewIntFromBigInt(quantums)
				pos.QuoteBalance = dtypes.NewIntFromBigInt(quoteBalance)
			}
		} else {
			// Create a new position.
			perpInfo := perpInfos.MustGet(update.PerpetualId)
			positionsMap[update.PerpetualId] = &types.PerpetualPosition{
				PerpetualId:  update.PerpetualId,
				Quantums:     dtypes.NewIntFromBigInt(update.GetBigQuantums()),
				QuoteBalance: dtypes.NewIntFromBigInt(update.GetBigQuoteBalance()),
				FundingIndex: perpInfo.Perpetual.FundingIndex,
			}
		}
	}

	return lib.MapToSortedSlice[lib.Sortable[uint32]](positionsMap)
}

// CalculateUpdatedSubaccount returns a copy of the settled subaccount with the updates applied.
func CalculateUpdatedSubaccount(
	settledUpdate types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) types.Subaccount {
	result := settledUpdate.SettledSubaccount.DeepCopy()
	result.AssetPositions = CalculateUpdatedAssetPositions(
		result.AssetPositions,
		settledUpdate.AssetUpdates,
	)
	result.PerpetualPositions = CalculateUpdatedPerpetualPositions(
		result.PerpetualPositions,
		settledUpdate.PerpetualUpdates,
		perpInfos,
	)
	return result
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
	subaccount types.Subaccount,
	perpInfos perptypes.PerpInfos,
	leverageMap map[uint32]uint32, // leverage per perpetual, nil means no leverage configured
) (
	risk margin.Risk,
	err error,
) {
	// Initialize return values.
	risk = margin.ZeroRisk()

	// Iterate over all assets and updates and calculate change to net collateral and margin requirements.
	for _, pos := range subaccount.AssetPositions {
		r, err := assetslib.GetNetCollateralAndMarginRequirements(
			pos.AssetId,
			pos.GetBigQuantums(),
		)
		if err != nil {
			return risk, err
		}
		risk.AddInPlace(r)
	}

	// Iterate over all perpetuals and updates and calculate change to net collateral and margin requirements.
	for _, pos := range subaccount.PerpetualPositions {
		perpInfo := perpInfos.MustGet(pos.PerpetualId)

		// Get the configured imf for this perpetual (0 if not configured)
		custom_imf_ppm := uint32(0)
		if leverageMap != nil {
			custom_imf_ppm = leverageMap[pos.PerpetualId]
		}

		r := perplib.GetNetCollateralAndMarginRequirements(
			perpInfo.Perpetual,
			perpInfo.Price,
			perpInfo.LiquidityTier,
			pos.GetBigQuantums(),
			pos.GetQuoteBalance(),
			custom_imf_ppm,
		)
		risk.AddInPlace(r)
	}

	return risk, nil
}

// GetRiskForSettledUpdate returns the risk value for a SettledUpdate with embedded leverage.
// This is a convenience function that extracts the leverage from the SettledUpdate and
// calls GetRiskForSubaccount with the updated subaccount.
func GetRiskForSettledUpdate(
	settledUpdate types.SettledUpdate,
	perpInfos perptypes.PerpInfos,
) (
	risk margin.Risk,
	err error,
) {
	updatedSubaccount := CalculateUpdatedSubaccount(settledUpdate, perpInfos)
	return GetRiskForSubaccount(updatedSubaccount, perpInfos, settledUpdate.LeverageMap)
}
