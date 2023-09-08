package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// SetSubaccount set a specific subaccount in the store from its index.
func (k Keeper) SetSubaccount(ctx sdk.Context, subaccount types.Subaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))
	b := k.cdc.MustMarshal(&subaccount)
	store.Set(types.SubaccountKey(
		*subaccount.Id,
	), b)
}

// GetSubaccount returns a subaccount from its index.
func (k Keeper) GetSubaccount(
	ctx sdk.Context,
	id types.SubaccountId,
) (val types.Subaccount) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetSubaccount,
		metrics.Latency,
	)

	// Check state for the subaccount.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))
	b := store.Get(types.SubaccountKey(id))

	// If subaccount does not exist in state, return a default value.
	if b == nil {
		return types.Subaccount{
			Id: &id,
		}
	}

	// If subaccount does exist in state, unmarshall and return the value.
	k.cdc.MustUnmarshal(b, &val)
	return val
}

// GetAllSubaccount returns all subaccount.
// For more performant searching and iteration, use `ForEachSubaccount`.
func (k Keeper) GetAllSubaccount(ctx sdk.Context) (list []types.Subaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Subaccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ForEachSubaccount performs a callback across all subaccounts.
// The callback function should return a boolean if we should end iteration or not.
// This is more performant than GetAllSubaccount because it does not fetch all at once.
// and you do not need to iterate through all the subaccounts.
func (k Keeper) ForEachSubaccount(ctx sdk.Context, callback func(types.Subaccount) (finished bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var subaccount types.Subaccount
		k.cdc.MustUnmarshal(iterator.Value(), &subaccount)
		done := callback(subaccount)
		if done {
			break
		}
	}
}

// ForEachSubaccountRandomStart performs a callback across all subaccounts.
// The callback function should return a boolean if we should end iteration or not.
// Note that this function starts at a random subaccount using the passed in `rand`
// and iterates from there. `rand` should be seeded for determinism if used in ways
// that affect consensus.
// TODO(CLOB-823): improve how random bytes are selected since bytes distribution
// might not be uniform.
func (k Keeper) ForEachSubaccountRandomStart(
	ctx sdk.Context,
	callback func(types.Subaccount) (finished bool),
	rand *rand.Rand,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))
	prefix, err := k.getRandomBytes(ctx, rand)
	if err != nil {
		return
	}

	// Iterate over subaccounts from the random prefix (inclusive) to the end.
	prefixStartIterator := store.Iterator(prefix, nil)
	defer prefixStartIterator.Close()
	for ; prefixStartIterator.Valid(); prefixStartIterator.Next() {
		var subaccount types.Subaccount
		k.cdc.MustUnmarshal(prefixStartIterator.Value(), &subaccount)
		done := callback(subaccount)
		if done {
			return
		}
	}

	// Iterator over subaccounts from the start to the random prefix (exclusive).
	prefixEndIterator := store.Iterator(nil, prefix)
	defer prefixEndIterator.Close()
	for ; prefixEndIterator.Valid(); prefixEndIterator.Next() {
		var subaccount types.Subaccount
		k.cdc.MustUnmarshal(prefixEndIterator.Value(), &subaccount)
		done := callback(subaccount)
		if done {
			return
		}
	}
}

// GetRandomSubaccount returns a random subaccount. Will return an error if there are no subaccounts.
func (k Keeper) GetRandomSubaccount(ctx sdk.Context, rand *rand.Rand) (types.Subaccount, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))

	prefix, err := k.getRandomBytes(ctx, rand)
	if err != nil {
		return types.Subaccount{}, err
	}
	prefixItr := store.Iterator(prefix, nil)
	defer prefixItr.Close()

	var val types.Subaccount
	k.cdc.MustUnmarshal(prefixItr.Value(), &val)
	return val, nil
}

func (k Keeper) getRandomBytes(ctx sdk.Context, rand *rand.Rand) ([]byte, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SubaccountKeyPrefix))

	// Use the forward iterator to get the first valid key.
	forwardItr := store.Iterator(nil, nil)
	defer forwardItr.Close()
	if !forwardItr.Valid() {
		return nil, errors.New("No subaccounts")
	}

	// Use the reverse iterator to get the last valid key.
	backwardsItr := store.ReverseIterator(nil, nil)
	defer backwardsItr.Close()

	firstKey := forwardItr.Key()
	lastKey := backwardsItr.Key()
	return lib.RandomBytesBetween(firstKey, lastKey, rand), nil
}

// getSettledUpdates takes in a list of updates and for each update, retrieves
// the updated subaccount in its settled form, and returns a list of settledUpdate
// structs and a map that indicates for each subaccount which perpetuals had funding
// updates. If requireUniqueSubaccount is true, the SubaccountIds in the input updates
// must be unique.
func (k Keeper) getSettledUpdates(
	ctx sdk.Context,
	updates []types.Update,
	requireUniqueSubaccount bool,
) (
	settledUpdates []settledUpdate,
	subaccountIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt,
	err error,
) {
	var idToSettledSubaccount = make(map[types.SubaccountId]types.Subaccount)
	settledUpdates = make([]settledUpdate, len(updates))
	subaccountIdToFundingPayments = make(map[types.SubaccountId]map[uint32]dtypes.SerializableInt)

	// Iterate over all updates and query the relevant `Subaccounts`.
	for i, u := range updates {
		settledSubaccount, exists := idToSettledSubaccount[u.SubaccountId]
		var fundingPayments map[uint32]dtypes.SerializableInt

		if exists && requireUniqueSubaccount {
			return nil, nil, types.ErrNonUniqueUpdatesSubaccount
		}

		// Get and store the settledSubaccount if SubaccountId doesn't exist in
		// idToSettledSubaccount map.
		if !exists {
			subaccount := k.GetSubaccount(ctx, u.SubaccountId)
			settledSubaccount, fundingPayments, err = k.getSettledSubaccount(ctx, subaccount)
			if err != nil {
				return nil, nil, err
			}

			idToSettledSubaccount[u.SubaccountId] = settledSubaccount
			subaccountIdToFundingPayments[u.SubaccountId] = fundingPayments
		}

		settledUpdate := settledUpdate{
			SettledSubaccount: settledSubaccount,
			AssetUpdates:      u.AssetUpdates,
			PerpetualUpdates:  u.PerpetualUpdates,
		}

		settledUpdates[i] = settledUpdate
	}

	return settledUpdates, subaccountIdToFundingPayments, nil
}

// UpdateSubaccounts validates and applies all `updates` to the relevant subaccounts as long as this is a
// valid state-transition for all subaccounts involved. All `updates` are made atomically, meaning that
// all state-changes will either succeed or all will fail.
//
// Returns a boolean indicating whether the update was successfully applied or not. If `false`, then no
// updates to any subaccount were made. A second return value returns an array of `UpdateResult` which map
// to the `updates` to indicate which of the updates caused a failure, if any.
//
// Each `SubaccountId` in the `updates` must be unique or an error is returned.
func (k Keeper) UpdateSubaccounts(
	ctx sdk.Context,
	updates []types.Update,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.UpdateSubaccounts,
		metrics.Latency,
	)

	settledUpdates, subaccountIdToFundingPayments, err := k.getSettledUpdates(ctx, updates, true)
	if err != nil {
		return false, nil, err
	}

	success, successPerUpdate, err = k.internalCanUpdateSubaccounts(ctx, settledUpdates)
	if !success || err != nil {
		return success, successPerUpdate, err
	}

	// Get a mapping from perpetual Id to current perpetual funding index.
	allPerps := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	perpIdToFundingIndex := make(map[uint32]dtypes.SerializableInt)
	for _, perp := range allPerps {
		perpIdToFundingIndex[perp.Params.Id] = perp.FundingIndex
	}

	// Apply the updates to perpetual positions.
	success, err = UpdatePerpetualPositions(
		settledUpdates,
		perpIdToFundingIndex,
	)
	if !success || err != nil {
		return success, successPerUpdate, err
	}

	// Apply the updates to asset positions.
	success, err = UpdateAssetPositions(settledUpdates)
	if !success || err != nil {
		return success, successPerUpdate, err
	}

	// Apply all updates, including a subaccount update event in the Indexer block message
	// per update and emit a cometbft event for each settled funding payment.
	for _, u := range settledUpdates {
		k.SetSubaccount(ctx, u.SettledSubaccount)
		// Below access is safe because for all updated subaccounts' IDs, this map
		// is populated as getSettledSubaccount() is called in getSettledUpdates().
		fundingPayments := subaccountIdToFundingPayments[*u.SettledSubaccount.Id]
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeSubaccountUpdate,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewSubaccountUpdateEvent(
					u.SettledSubaccount.Id,
					getUpdatedPerpetualPositions(
						u,
						fundingPayments,
					),
					getUpdatedAssetPositions(u),
					fundingPayments,
				),
			),
		)

		// Emit an event indicating a funding payment was paid / received for each settled funding
		// payment. Note that `fundingPaid` is positive if the subaccount paid funding,
		// and negative if the subaccount received funding.
		for perpetualId, fundingPaid := range fundingPayments {
			ctx.EventManager().EmitEvent(
				types.NewCreateSettledFundingEvent(
					*u.SettledSubaccount.Id,
					perpetualId,
					fundingPaid.BigInt(),
				),
			)
		}
	}

	return success, successPerUpdate, err
}

// CanUpdateSubaccounts will validate all `updates` to the relevant subaccounts.
// The `updates` do not have to contain unique `SubaccountIds`.
// Each update is considered in isolation. Thus if two updates are provided
// with the same `SubaccountId`, they are validated without respect to each
// other.
//
// Returns a `success` value of `true` if all updates are valid.
// Returns a `successPerUpdates` value, which is a slice of `UpdateResult`.
// These map to the updates and are used to indicate which of the updates
// caused a failure, if any.
func (k Keeper) CanUpdateSubaccounts(
	ctx sdk.Context,
	updates []types.Update,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.CanUpdateSubaccounts,
		metrics.Latency,
	)

	settledUpdates, _, err := k.getSettledUpdates(ctx, updates, false)
	if err != nil {
		return false, nil, err
	}

	return k.internalCanUpdateSubaccounts(ctx, settledUpdates)
}

// getSettledSubaccount returns 1. a new settled subaccount given an unsettled subaccount,
// updating the USDC AssetPosition, FundingIndex, and LastFundingPayment fields accordingly
// (does not persist any changes) and 2. a map with perpetual ID as key and last funding
// payment as value (for emitting funding payments to indexer).
func (k Keeper) getSettledSubaccount(
	ctx sdk.Context,
	subaccount types.Subaccount,
) (
	settledSubaccount types.Subaccount,
	fundingPayments map[uint32]dtypes.SerializableInt,
	err error,
) {
	totalNetSettlement := big.NewInt(0)

	newPerpetualPositions := []*types.PerpetualPosition{}
	fundingPayments = make(map[uint32]dtypes.SerializableInt)

	// Iterate through and settle all perpetual positions.
	for _, p := range subaccount.PerpetualPositions {
		bigNetSettlement, newFundingIndex, err := k.perpetualsKeeper.GetSettlement(
			ctx,
			p.PerpetualId,
			p.GetBigQuantums(),
			p.FundingIndex.BigInt(),
		)
		if err != nil {
			return types.Subaccount{}, nil, err
		}
		// Record non-zero funding payment (to be later emitted in SubaccountUpdateEvent to indexer).
		// Note: Funding payment is the negative of settlement, i.e. positive settlement is equivalent
		// to a negative funding payment (position received funding payment) and vice versa.
		if bigNetSettlement.Cmp(lib.BigInt0()) != 0 {
			fundingPayments[p.PerpetualId] = dtypes.NewIntFromBigInt(
				new(big.Int).Neg(bigNetSettlement),
			)
		}

		// Aggregate all net settlements.
		// TODO(DEC-657): For less rounding error, divide net settlement by
		// SomeLargeNumber only after summing all net settlements.
		totalNetSettlement.Add(totalNetSettlement, bigNetSettlement)

		// Update cached funding index of the perpetual position.
		newPerpetualPositions = append(newPerpetualPositions, &types.PerpetualPosition{
			PerpetualId:  p.PerpetualId,
			Quantums:     p.Quantums,
			FundingIndex: dtypes.NewIntFromBigInt(newFundingIndex),
		})
	}

	newSubaccount := types.Subaccount{
		Id:                 subaccount.Id,
		AssetPositions:     subaccount.AssetPositions,
		PerpetualPositions: newPerpetualPositions,
		MarginEnabled:      subaccount.MarginEnabled,
	}
	newUsdcPosition := new(big.Int).Add(subaccount.GetUsdcPosition(), totalNetSettlement)
	err = newSubaccount.SetUsdcAssetPosition(newUsdcPosition)
	if err != nil {
		return types.Subaccount{}, nil, err
	}
	return newSubaccount, fundingPayments, nil
}

// internalCanUpdateSubaccounts will validate all `updates` to the relevant subaccounts.
// The `updates` do not have to contain `Subaccounts` with unique `SubaccountIds`.
// Each update is considered in isolation. Thus if two updates are provided
// with the same `Subaccount`, they are validated without respect to each
// other.
// The input subaccounts must be settled.
//
// Returns a `success` value of `true` if all updates are valid.
// Returns a `successPerUpdates` value, which is a slice of `UpdateResult`.
// These map to the updates and are used to indicate which of the updates
// caused a failure, if any.
func (k Keeper) internalCanUpdateSubaccounts(
	ctx sdk.Context,
	settledUpdates []settledUpdate,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	success = true
	successPerUpdate = make([]types.UpdateResult, len(settledUpdates))

	// Iterate over all updates.
	for i, u := range settledUpdates {
		// Get the new collateralization and margin requirements with the update applied.
		bigNewNetCollateral,
			bigNewInitialMargin,
			bigNewMaintenanceMargin,
			err := k.internalGetNetCollateralAndMarginRequirements(ctx, u)
		if err != nil {
			return false, nil, err
		}

		var result = types.Success

		// The subaccount is not well-collateralized after the update.
		// We must now check if the state transition is valid.
		if bigNewInitialMargin.Cmp(bigNewNetCollateral) > 0 {
			// Get the current collateralization and margin requirements without the update applied.
			emptyUpdate := settledUpdate{
				SettledSubaccount: u.SettledSubaccount,
			}

			bigCurNetCollateral,
				bigCurInitialMargin,
				bigCurMaintenanceMargin,
				err := k.internalGetNetCollateralAndMarginRequirements(
				ctx,
				emptyUpdate,
			)
			if err != nil {
				return false, nil, err
			}

			// Determine whether the state transition is valid.
			result = IsValidStateTransitionForUndercollateralizedSubaccount(
				bigCurNetCollateral,
				bigCurInitialMargin,
				bigCurMaintenanceMargin,
				bigNewNetCollateral,
				bigNewMaintenanceMargin,
			)
		}

		// If this state transition is not valid, the overall success is now false.
		if !result.IsSuccess() {
			success = false
		}

		successPerUpdate[i] = result
	}

	return success, successPerUpdate, nil
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

// GetNetCollateralAndMarginRequirements returns the total net collateral, total initial margin requirement,
// and total maintenance margin requirement for the subaccount as if the `update` was applied.
// It is used to get information about speculative changes to the subaccount.
//
// The provided update can also be "zeroed" in order to get information about
// the current state of the subaccount (i.e. with no changes).
//
// If two position updates reference the same position, an error is returned.
//
// All return values are denoted in quote quantums.
func (k Keeper) GetNetCollateralAndMarginRequirements(
	ctx sdk.Context,
	update types.Update,
) (
	bigNetCollateral *big.Int,
	bigInitialMargin *big.Int,
	bigMaintenanceMargin *big.Int,
	err error,
) {
	subaccount := k.GetSubaccount(ctx, update.SubaccountId)

	settledSubaccount, _, err := k.getSettledSubaccount(ctx, subaccount)
	if err != nil {
		return nil, nil, nil, err
	}

	settledUpdate := settledUpdate{
		SettledSubaccount: settledSubaccount,
		AssetUpdates:      update.AssetUpdates,
		PerpetualUpdates:  update.PerpetualUpdates,
	}

	return k.internalGetNetCollateralAndMarginRequirements(
		ctx,
		settledUpdate,
	)
}

// internalGetNetCollateralAndMarginRequirements returns the total net collateral, total initial margin
// requirement, and total maintenance margin requirement for the `Subaccount` as if unsettled funding
// of existing positions were settled, and the `bigQuoteBalanceDeltaQuantums`, `assetUpdates`, and
// `perpetualUpdates` were applied. It is used to get information about speculative changes to the
// `Subaccount`.
// The input subaccounts must be settled.
//
// The provided update can also be "zeroed" in order to get information about
// the current state of the subaccount (i.e. with no changes).
//
// If two position updates reference the same position, an error is returned.
func (k Keeper) internalGetNetCollateralAndMarginRequirements(
	ctx sdk.Context,
	settledUpdate settledUpdate,
) (
	bigNetCollateral *big.Int,
	bigInitialMargin *big.Int,
	bigMaintenanceMargin *big.Int,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetNetCollateralAndMarginRequirements,
		metrics.Latency,
	)

	// Initialize return values.
	bigNetCollateral = big.NewInt(0)
	bigInitialMargin = big.NewInt(0)
	bigMaintenanceMargin = big.NewInt(0)

	// Merge updates and assets.
	assetSizes, err := applyUpdatesToPositions(
		settledUpdate.SettledSubaccount.AssetPositions,
		settledUpdate.AssetUpdates,
	)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
	}

	// Merge updates and perpetuals.
	perpetualSizes, err := applyUpdatesToPositions(
		settledUpdate.SettledSubaccount.PerpetualPositions,
		settledUpdate.PerpetualUpdates,
	)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
	}

	// The calculate function increments `netCollateral`, `initialMargin`, and `maintenanceMargin`
	// given a `ProductKeeper` and a `PositionSize`.
	calculate := func(pk types.ProductKeeper, size types.PositionSize) error {
		id := size.GetId()
		bigQuantums := size.GetBigQuantums()

		bigNetCollateralQuoteQuantums, err := pk.GetNetCollateral(ctx, id, bigQuantums)
		if err != nil {
			return err
		}

		bigNetCollateral.Add(bigNetCollateral, bigNetCollateralQuoteQuantums)

		bigInitialMarginRequirements,
			bigMaintenanceMarginRequirements,
			err := pk.GetMarginRequirements(ctx, id, bigQuantums)
		if err != nil {
			return err
		}

		bigInitialMargin.Add(bigInitialMargin, bigInitialMarginRequirements)
		bigMaintenanceMargin.Add(bigMaintenanceMargin, bigMaintenanceMarginRequirements)

		return nil
	}

	// Iterate over all assets and updates and calculate change to net collateral and margin requirements.
	for _, size := range assetSizes {
		err := calculate(k.assetsKeeper, size)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
		}
	}

	// Iterate over all perpetuals and updates and calculate change to net collateral and margin requirements.
	// TODO(DEC-110): `perp.GetSettlement()`, factor in unsettled funding.
	for _, size := range perpetualSizes {
		err := calculate(k.perpetualsKeeper, size)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
		}
	}

	return bigNetCollateral, bigInitialMargin, bigMaintenanceMargin, nil
}

// applyUpdatesToPositions merges a slice of `types.UpdatablePositions` and `types.PositionSize`
// (i.e. concrete types *types.AssetPosition` and `types.AssetUpdate`) into a slice of `types.PositionSize`.
// If a given `PositionSize` shares an ID with an `UpdatablePositionSize`, the update and position are merged
// into a single `PositionSize`.
//
// An error is returned if two updates share the same position id.
//
// Note: There are probably performance implications here for allocating a new slice of PositionSize,
// and for allocating new slices when converting the concrete types to interfaces. However, without doing
// this there would be a lot of duplicate code for calculating changes for both `Assets` and `Perpetuals`.
func applyUpdatesToPositions[
	P types.PositionSize,
	U types.PositionSize,
](positions []P, updates []U) ([]types.PositionSize, error) {
	var result []types.PositionSize = make([]types.PositionSize, 0, len(positions)+len(updates))

	updateMap := make(map[uint32]types.PositionSize)
	updateIndexMap := make(map[uint32]int)
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
