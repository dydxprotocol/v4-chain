package lib

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	perpInfos map[uint32]perptypes.PerpInfo,
) (
	settledSubaccount types.Subaccount,
	fundingPayments map[uint32]dtypes.SerializableInt,
	err error,
) {
	totalNetSettlementPpm := big.NewInt(0)

	newPerpetualPositions := []*types.PerpetualPosition{}
	fundingPayments = make(map[uint32]dtypes.SerializableInt)

	// Iterate through and settle all perpetual positions.
	for _, p := range subaccount.PerpetualPositions {
		perpInfo, found := perpInfos[p.PerpetualId]
		if !found {
			return types.Subaccount{}, nil, errorsmod.Wrapf(types.ErrPerpetualInfoDoesNotExist, "%d", p.PerpetualId)
		}

		// Call the stateless utility function to get the net settlement and new funding index.
		bigNetSettlementPpm, newFundingIndex := perplib.GetSettlementPpmWithPerpetual(
			perpInfo.Perpetual,
			p.GetBigQuantums(),
			p.FundingIndex.BigInt(),
		)
		// Record non-zero funding payment (to be later emitted in SubaccountUpdateEvent to indexer).
		// Note: Funding payment is the negative of settlement, i.e. positive settlement is equivalent
		// to a negative funding payment (position received funding payment) and vice versa.
		if bigNetSettlementPpm.Cmp(lib.BigInt0()) != 0 {
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
	return newSubaccount, fundingPayments, nil
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
	bigCurNetCollateral *big.Int,
	bigCurInitialMargin *big.Int,
	bigCurMaintenanceMargin *big.Int,
	bigNewNetCollateral *big.Int,
	bigNewMaintenanceMargin *big.Int,
) types.UpdateResult {
	// Determine whether the subaccount was previously undercollateralized before the update.
	var underCollateralizationResult = types.StillUndercollateralized
	if bigCurInitialMargin.Cmp(bigCurNetCollateral) <= 0 {
		underCollateralizationResult = types.NewlyUndercollateralized
	}

	// If the maintenance margin is increasing, then the subaccount is undercollateralized.
	if bigNewMaintenanceMargin.Cmp(bigCurMaintenanceMargin) > 0 {
		return underCollateralizationResult
	}

	// If the maintenance margin is zero, it means the subaccount must have no open positions, and negative net
	// collateral. If the net collateral is not improving then this transition is not valid.
	if bigNewMaintenanceMargin.BitLen() == 0 || bigCurMaintenanceMargin.BitLen() == 0 {
		if bigNewMaintenanceMargin.BitLen() == 0 &&
			bigCurMaintenanceMargin.BitLen() == 0 &&
			bigNewNetCollateral.Cmp(bigCurNetCollateral) > 0 {
			return types.Success
		}

		return underCollateralizationResult
	}

	// Note that here we are effectively checking that
	// `newNetCollateral / newMaintenanceMargin >= curNetCollateral / curMaintenanceMargin`.
	// However, to avoid rounding errors, we factor this as
	// `newNetCollateral * curMaintenanceMargin >= curNetCollateral * newMaintenanceMargin`.
	bigCurRisk := new(big.Int).Mul(bigNewNetCollateral, bigCurMaintenanceMargin)
	bigNewRisk := new(big.Int).Mul(bigCurNetCollateral, bigNewMaintenanceMargin)

	// The subaccount is not well-collateralized, and the state transition leaves the subaccount in a
	// "more-risky" state (collateral relative to margin requirements is decreasing).
	if bigNewRisk.Cmp(bigCurRisk) > 0 {
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
