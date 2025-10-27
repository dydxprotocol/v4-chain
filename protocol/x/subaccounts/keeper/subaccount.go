package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	streamingtypes "github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	salib "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	gometrics "github.com/hashicorp/go-metrics"
)

// SetSubaccount set a specific subaccount in the store from its index.
// Note that empty subaccounts are removed from state.
func (k Keeper) SetSubaccount(ctx sdk.Context, subaccount types.Subaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SubaccountKeyPrefix))
	key := subaccount.Id.ToStateKey()

	if len(subaccount.PerpetualPositions) == 0 && len(subaccount.AssetPositions) == 0 {
		if store.Has(key) {
			store.Delete(key)
		}
	} else {
		if !store.Has(key) {
			metrics.IncrCounterWithLabels(
				metrics.SubaccountCreatedCount,
				1,
				metrics.GetLabelForStringValue(
					metrics.Callback,
					metrics.GetCallbackMetricFromCtx(ctx),
				),
			)
		}
		b := k.cdc.MustMarshal(&subaccount)
		store.Set(key, b)
	}
}

// GetCollateralPoolForSubaccount returns the collateral pool address for a subaccount
// based on the subaccount's perpetual positions. If the subaccount holds a position in an isolated
// market, the collateral pool address will be the isolated market's pool address. Otherwise, the
// collateral pool address will be the module's pool address.
func (k Keeper) GetCollateralPoolForSubaccount(ctx sdk.Context, subaccountId types.SubaccountId) (
	sdk.AccAddress,
	error,
) {
	subaccount := k.GetSubaccount(ctx, subaccountId)
	return k.getCollateralPoolForSubaccount(ctx, subaccount)
}

func (k Keeper) getCollateralPoolForSubaccount(ctx sdk.Context, subaccount types.Subaccount) (
	sdk.AccAddress,
	error,
) {
	// Use the default collateral pool if the subaccount has no perpetual positions.
	if len(subaccount.PerpetualPositions) == 0 {
		return types.ModuleAddress, nil
	}

	return k.GetCollateralPoolFromPerpetualId(ctx, subaccount.PerpetualPositions[0].PerpetualId)
}

// GetCollateralPoolForSubaccountWithPerpetuals returns the collateral pool address based on the
// perpetual passed in as an argument.
func (k Keeper) GetCollateralPoolFromPerpetualId(ctx sdk.Context, perpetualId uint32) (sdk.AccAddress, error) {
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return nil, err
	}

	if perpetual.Params.MarketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
		return authtypes.NewModuleAddress(types.ModuleName + ":" + lib.UintToString(perpetual.GetId())), nil
	}

	return authtypes.NewModuleAddress(types.ModuleName), nil
}

// GetSubaccount returns a subaccount from its index.
//
// Note that this function is getting called very frequently; metrics in this function
// should be sampled to reduce CPU time.
func (k Keeper) GetSubaccount(
	ctx sdk.Context,
	id types.SubaccountId,
) (val types.Subaccount) {
	if rand.Float64() < metrics.LatencyMetricSampleRate {
		defer metrics.ModuleMeasureSinceWithLabels(
			types.ModuleName,
			[]string{metrics.GetSubaccount, metrics.Latency},
			time.Now(),
			[]gometrics.Label{
				metrics.GetLabelForStringValue(
					metrics.SampleRate,
					fmt.Sprintf("%f", metrics.LatencyMetricSampleRate),
				),
			},
		)
	}

	// Check state for the subaccount.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SubaccountKeyPrefix))
	b := store.Get(id.ToStateKey())

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

func (k Keeper) GetStreamSubaccountUpdate(
	ctx sdk.Context,
	id types.SubaccountId,
	snapshot bool,
) (val types.StreamSubaccountUpdate) {
	subaccount := k.GetSubaccount(ctx, id)
	assetPositions := make([]*types.SubaccountAssetPosition, len(subaccount.AssetPositions))
	for i, ap := range subaccount.AssetPositions {
		assetPositions[i] = &types.SubaccountAssetPosition{
			AssetId:  ap.AssetId,
			Quantums: ap.Quantums.BigInt().Uint64(),
		}
	}
	perpetualPositions := make([]*types.SubaccountPerpetualPosition, len(subaccount.PerpetualPositions))
	for i, pp := range subaccount.PerpetualPositions {
		perpetualPositions[i] = &types.SubaccountPerpetualPosition{
			PerpetualId: pp.PerpetualId,
			Quantums:    pp.Quantums.BigInt().Int64(),
		}
	}

	return types.StreamSubaccountUpdate{
		SubaccountId:              &id,
		UpdatedAssetPositions:     assetPositions,
		UpdatedPerpetualPositions: perpetualPositions,
		Snapshot:                  snapshot,
	}
}

// GetAllSubaccount returns all subaccount.
// For more performant searching and iteration, use `ForEachSubaccount`.
func (k Keeper) GetAllSubaccount(ctx sdk.Context) (list []types.Subaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SubaccountKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SubaccountKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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

// GetRandomSubaccount returns a random subaccount. Will return an error if there are no subaccounts.
func (k Keeper) GetRandomSubaccount(ctx sdk.Context, rand *rand.Rand) (types.Subaccount, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SubaccountKeyPrefix))

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SubaccountKeyPrefix))

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
	perpInfos perptypes.PerpInfos,
	requireUniqueSubaccount bool,
) (
	settledUpdates []types.SettledUpdate,
	subaccountIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt,
	err error,
) {
	var idToSettledSubaccount = make(map[types.SubaccountId]types.Subaccount)
	var idToLeverageMap = make(map[types.SubaccountId]map[uint32]uint32)
	settledUpdates = make([]types.SettledUpdate, len(updates))
	subaccountIdToFundingPayments = make(map[types.SubaccountId]map[uint32]dtypes.SerializableInt)

	// Iterate over all updates and query the relevant `Subaccounts`.
	for i, u := range updates {
		settledSubaccount, exists := idToSettledSubaccount[u.SubaccountId]
		var fundingPayments map[uint32]dtypes.SerializableInt
		var leverageMap map[uint32]uint32

		if exists && requireUniqueSubaccount {
			return nil, nil, types.ErrNonUniqueUpdatesSubaccount
		}

		// Get and store the settledSubaccount and leverage if SubaccountId doesn't exist in maps.
		if !exists {
			subaccount := k.GetSubaccount(ctx, u.SubaccountId)
			settledSubaccount, fundingPayments = salib.GetSettledSubaccountWithPerpetuals(subaccount, perpInfos)

			// Only fetch leverage if there are perpetual updates or perpetual positions
			// to avoid unnecessary gas consumption
			if len(u.PerpetualUpdates) > 0 || len(settledSubaccount.PerpetualPositions) > 0 {
				if leverage, found := k.GetLeverage(ctx, &u.SubaccountId); found {
					leverageMap = leverage
				}
			}

			idToSettledSubaccount[u.SubaccountId] = settledSubaccount
			idToLeverageMap[u.SubaccountId] = leverageMap
			subaccountIdToFundingPayments[u.SubaccountId] = fundingPayments
		} else {
			// Reuse cached leverage map if there are perpetual updates
			// or perpetual positions
			if len(u.PerpetualUpdates) > 0 || len(settledSubaccount.PerpetualPositions) > 0 {
				leverageMap = idToLeverageMap[u.SubaccountId]
			}
		}

		settledUpdate := types.SettledUpdate{
			SettledSubaccount: settledSubaccount,
			AssetUpdates:      u.AssetUpdates,
			PerpetualUpdates:  u.PerpetualUpdates,
			LeverageMap:       leverageMap,
		}

		settledUpdates[i] = settledUpdate
	}

	return settledUpdates, subaccountIdToFundingPayments, nil
}

func GenerateStreamSubaccountUpdate(
	settledUpdate types.SettledUpdate,
	fundingPayments map[uint32]dtypes.SerializableInt,
) types.StreamSubaccountUpdate {
	// Get updated perpetual positions
	updatedPerpetualPositions := salib.GetUpdatedPerpetualPositions(
		settledUpdate,
		fundingPayments,
	)
	// Convert updated perpetual positions to SubaccountPerpetualPosition type
	perpetualPositions := make([]*types.SubaccountPerpetualPosition, len(updatedPerpetualPositions))
	for i, pp := range updatedPerpetualPositions {
		perpetualPositions[i] = &types.SubaccountPerpetualPosition{
			PerpetualId: pp.PerpetualId,
			Quantums:    pp.Quantums.BigInt().Int64(),
		}
	}

	updatedAssetPositions := salib.GetUpdatedAssetPositions(settledUpdate)
	assetPositionsWithQuoteBalance := indexerevents.AddQuoteBalanceFromPerpetualPositions(
		updatedPerpetualPositions,
		updatedAssetPositions,
	)

	// Convert updated asset positions to SubaccountAssetPosition type
	assetPositions := make([]*types.SubaccountAssetPosition, len(assetPositionsWithQuoteBalance))
	for i, ap := range assetPositionsWithQuoteBalance {
		assetPositions[i] = &types.SubaccountAssetPosition{
			AssetId:  ap.AssetId,
			Quantums: ap.Quantums.BigInt().Uint64(),
		}
	}

	return types.StreamSubaccountUpdate{
		SubaccountId:              settledUpdate.SettledSubaccount.Id,
		UpdatedAssetPositions:     assetPositions,
		UpdatedPerpetualPositions: perpetualPositions,
		Snapshot:                  false,
	}
}

// UpdateSubaccounts validates and applies all `updates` to the relevant subaccounts as long as this is a
// valid state-transition for all subaccounts involved. All `updates` are made atomically, meaning that
// all state-changes will either succeed or all will fail.
//
// Returns a boolean indicating whether the update was successfully applied or not. If `false`, then no
// updates to any subaccount were made. A second return value returns an array of `UpdateResult` which map
// to the `updates` to indicate which of the updates caused a failure, if any.
// This function also transfers collateral between the cross-perpetual collateral pool and isolated
// perpetual collateral pools if any of the updates led to an isolated perpetual posititon to be opened
// or closed. This is done using the `x/bank` keeper and updates `x/bank` state.
//
// Each `SubaccountId` in the `updates` must be unique or an error is returned.
func (k Keeper) UpdateSubaccounts(
	ctx sdk.Context,
	updates []types.Update,
	updateType types.UpdateType,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	defer metrics.ModuleMeasureSinceWithLabels(
		types.ModuleName,
		[]string{metrics.UpdateSubaccounts, metrics.Latency},
		time.Now(),
		[]gometrics.Label{
			metrics.GetLabelForStringValue(metrics.UpdateType, updateType.String()),
		},
	)

	perpInfos, err := k.GetAllRelevantPerpetuals(ctx, updates)
	if err != nil {
		return false, nil, err
	}

	settledUpdates, subaccountIdToFundingPayments, err := k.getSettledUpdates(ctx, updates, perpInfos, true)
	if err != nil {
		return false, nil, err
	}

	success, successPerUpdate, err = k.internalCanUpdateSubaccountsWithLeverage(
		ctx,
		settledUpdates,
		updateType,
		perpInfos,
	)

	if !success || err != nil {
		return success, successPerUpdate, err
	}

	// Get OpenInterestDelta from the updates, and persist the OI change if any.
	perpOpenInterestDelta := salib.GetDeltaOpenInterestFromUpdates(settledUpdates, updateType)
	if perpOpenInterestDelta != nil {
		if err := k.perpetualsKeeper.ModifyOpenInterest(
			ctx,
			perpOpenInterestDelta.PerpetualId,
			perpOpenInterestDelta.BaseQuantums,
		); err != nil {
			return false, nil, errorsmod.Wrapf(
				types.ErrCannotModifyPerpOpenInterestForOIMF,
				"perpId = %v, delta = %v, settledUpdates = %+v, err = %v",
				perpOpenInterestDelta.PerpetualId,
				perpOpenInterestDelta.BaseQuantums,
				settledUpdates,
				err,
			)
		}
	}

	// Apply the updates to asset positions and perpetual positions.
	for i := range settledUpdates {
		settledUpdates[i].SettledSubaccount = salib.CalculateUpdatedSubaccount(
			settledUpdates[i],
			perpInfos,
		)
	}

	// Transfer collateral between collateral pools for any isolated perpetual positions that changed
	// state due to an update.
	for _, settledUpdateWithUpdatedSubaccount := range settledUpdates {
		if err := k.computeAndExecuteCollateralTransfer(
			ctx,
			// The subaccount in `settledUpdateWithUpdatedSubaccount` already has the perpetual updates
			// and asset updates applied to it.
			settledUpdateWithUpdatedSubaccount,
			perpInfos,
		); err != nil {
			return false, nil, err
		}
	}

	// Apply all updates, including a subaccount update event in the Indexer block message
	// per update and emit a cometbft event for each settled funding payment.
	for _, u := range settledUpdates {
		k.SetSubaccount(ctx, u.SettledSubaccount)
		// Below access is safe because for all updated subaccounts' IDs, this map
		// is populated as GetSettledSubaccountWithPerpetuals() is called in getSettledUpdates().
		fundingPayments := subaccountIdToFundingPayments[*u.SettledSubaccount.Id]
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeSubaccountUpdate,
			indexerevents.SubaccountUpdateEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewSubaccountUpdateEvent(
					u.SettledSubaccount.Id,
					salib.GetUpdatedPerpetualPositions(
						u,
						fundingPayments,
					),
					salib.GetUpdatedAssetPositions(u),
					fundingPayments,
				),
			),
		)

		// If DeliverTx and GRPC streaming is on, emit a generated subaccount update to stream.
		if lib.IsDeliverTxMode(ctx) && k.GetFullNodeStreamingManager().Enabled() {
			if k.GetFullNodeStreamingManager().TracksSubaccountId(*u.SettledSubaccount.Id) {
				subaccountUpdate := GenerateStreamSubaccountUpdate(u, fundingPayments)
				k.GetFullNodeStreamingManager().SendSubaccountUpdate(
					ctx,
					subaccountUpdate,
				)
			}
		}

		// Emit an event indicating a funding payment was paid / received for each settled funding
		// payment. Note that `fundingPaid` is positive if the subaccount paid funding,
		// and negative if the subaccount received funding.
		// Note the perpetual IDs are sorted first to ensure event emission determinism.
		sortedPerpIds := lib.GetSortedKeys[lib.Sortable[uint32]](fundingPayments)
		for _, perpetualId := range sortedPerpIds {
			fundingPaid := fundingPayments[perpetualId]
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
// This method automatically fetches leverage configuration for all subaccounts
// being updated and applies leverage-aware margin requirements.
//
// Returns a `success` value of `true` if all updates are valid.
// Returns a `successPerUpdates` value, which is a slice of `UpdateResult`.
// These map to the updates and are used to indicate which of the updates
// caused a failure, if any.
func (k Keeper) CanUpdateSubaccounts(
	ctx sdk.Context,
	updates []types.Update,
	updateType types.UpdateType,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	defer metrics.ModuleMeasureSinceWithLabels(
		types.ModuleName,
		[]string{metrics.CanUpdateSubaccounts, metrics.Latency},
		time.Now(),
		[]gometrics.Label{
			metrics.GetLabelForStringValue(metrics.UpdateType, updateType.String()),
		},
	)

	perpInfos, err := k.GetAllRelevantPerpetuals(ctx, updates)
	if err != nil {
		return false, nil, err
	}

	settledUpdates, _, err := k.getSettledUpdates(ctx, updates, perpInfos, false)
	if err != nil {
		return false, nil, err
	}

	success, successPerUpdate, err = k.internalCanUpdateSubaccountsWithLeverage(ctx, settledUpdates, updateType, perpInfos)
	return success, successPerUpdate, err
}

func (k Keeper) internalCanUpdateSubaccountsWithLeverage(
	ctx sdk.Context,
	settledUpdates []types.SettledUpdate,
	updateType types.UpdateType,
	perpInfos perptypes.PerpInfos,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	// TODO(TRA-99): Add integration / E2E tests on order placement / matching with this new
	// constraint.
	// Check if the updates satisfy the isolated perpetual constraints.
	success, successPerUpdate = k.checkIsolatedSubaccountConstraints(
		ctx,
		settledUpdates,
		perpInfos,
	)
	if !success {
		return success, successPerUpdate, nil
	}

	// Block all withdrawals and transfers if either of the following is true within the last
	// `WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`:
	// - There was a negative TNC subaccount seen for any of the collateral pools of subaccounts being updated
	// - There was a chain outage that lasted at least five minutes.
	if updateType == types.Withdrawal || updateType == types.Transfer {
		lastBlockNegativeTncSubaccountSeen, negativeTncSubaccountExists, err := k.getLastBlockNegativeSubaccountSeen(
			ctx,
			settledUpdates,
		)
		if err != nil {
			return false, nil, err
		}
		currentBlock := uint32(ctx.BlockHeight())

		// Panic if the current block is less than the last block a negative TNC subaccount was seen.
		if negativeTncSubaccountExists && currentBlock < lastBlockNegativeTncSubaccountSeen {
			panic(
				fmt.Sprintf(
					"internalCanUpdateSubaccounts: current block (%d) is less than the last "+
						"block a negative TNC subaccount was seen (%d)",
					currentBlock,
					lastBlockNegativeTncSubaccountSeen,
				),
			)
		}

		// Panic if the current block is less than the last block a chain outage was seen.
		downtimeInfo := k.blocktimeKeeper.GetDowntimeInfoFor(
			ctx,
			types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_CHAIN_OUTAGE_DURATION,
		)
		chainOutageExists := downtimeInfo.BlockInfo.Height > 0 && downtimeInfo.Duration > 0
		if chainOutageExists && currentBlock < downtimeInfo.BlockInfo.Height {
			panic(
				fmt.Sprintf(
					"internalCanUpdateSubaccounts: current block (%d) is less than the last "+
						"block a chain outage was seen (%d)",
					currentBlock,
					downtimeInfo.BlockInfo.Height,
				),
			)
		}

		negativeTncSubaccountSeen := negativeTncSubaccountExists && currentBlock-lastBlockNegativeTncSubaccountSeen <
			types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS
		chainOutageSeen := chainOutageExists && currentBlock-downtimeInfo.BlockInfo.Height <
			types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS

		if negativeTncSubaccountSeen || chainOutageSeen {
			success = false
			for i := range settledUpdates {
				successPerUpdate[i] = types.WithdrawalsAndTransfersBlocked
			}
			metrics.IncrCounterWithLabels(
				metrics.SubaccountWithdrawalsAndTransfersBlocked,
				1,
				metrics.GetLabelForStringValue(metrics.UpdateType, updateType.String()),
				metrics.GetLabelForBoolValue(metrics.SubaccountsNegativeTncSubaccountSeen, negativeTncSubaccountSeen),
				metrics.GetLabelForBoolValue(metrics.ChainOutageSeen, chainOutageSeen),
			)
			return success, successPerUpdate, nil
		}
	}

	// Get delta open interest from the updates.
	// `perpOpenInterestDelta` is nil if the update type is not `Match` or if the updates
	// do not result in OI changes.
	perpOpenInterestDelta := salib.GetDeltaOpenInterestFromUpdates(settledUpdates, updateType)

	// Temporily apply open interest delta to perpetuals, so IMF is calculated based on open interest after the update.
	// `perpOpenInterestDeltas` is only present for `Match` update type.
	if perpOpenInterestDelta != nil {
		perpInfo := perpInfos.MustGet(perpOpenInterestDelta.PerpetualId)
		existingValue := big.NewInt(0)
		if !perpInfo.Perpetual.OpenInterest.IsNil() {
			existingValue.Set(perpInfo.Perpetual.OpenInterest.BigInt())
		}
		perpInfo.Perpetual.OpenInterest = dtypes.NewIntFromBigInt(
			new(big.Int).Add(existingValue, perpOpenInterestDelta.BaseQuantums),
		)
		perpInfos[perpOpenInterestDelta.PerpetualId] = perpInfo

		// Reset the OpenInterest to the original value.
		defer func() {
			perpInfo.Perpetual.OpenInterest = dtypes.NewIntFromBigInt(existingValue)
			perpInfos[perpOpenInterestDelta.PerpetualId] = perpInfo
		}()
	}

	riskCurMap := make(map[string]margin.Risk)

	// Iterate over all updates.
	for i, u := range settledUpdates {
		// Check all updated perps are updatable.
		for _, perpUpdate := range u.PerpetualUpdates {
			updatable, err := k.perpetualsKeeper.IsPositionUpdatable(ctx, perpUpdate.GetId())
			if err != nil {
				return false, nil, err
			}
			if !updatable {
				return false, nil, errorsmod.Wrapf(
					types.ErrProductPositionNotUpdatable,
					"type: perpetual, id: %d",
					perpUpdate.GetId(),
				)
			}
		}

		// Check all updated assets are updatable.
		for _, assetUpdate := range u.AssetUpdates {
			updatable, err := k.assetsKeeper.IsPositionUpdatable(ctx, assetUpdate.GetId())
			if err != nil {
				return false, nil, err
			}
			if !updatable {
				return false, nil, errorsmod.Wrapf(
					types.ErrProductPositionNotUpdatable,
					"type: asset, id: %d",
					assetUpdate.GetId(),
				)
			}
		}

		// Get the new collateralization and margin requirements with the update applied.
		riskNew, err := salib.GetRiskForSettledUpdate(u, perpInfos)
		if err != nil {
			return false, nil, err
		}

		var result = types.Success

		// The subaccount is not well-collateralized after the update.
		// We must now check if the state transition is valid.
		if !riskNew.IsInitialCollateralized() {
			// Get the current collateralization and margin requirements without the update applied.
			bytes, err := proto.Marshal(u.SettledSubaccount.Id)
			if err != nil {
				return false, nil, err
			}
			saKey := string(bytes)

			// Cache the current collateralization and margin requirements for the subaccount.
			if _, ok := riskCurMap[saKey]; !ok {
				riskCurMap[saKey], err = salib.GetRiskForSubaccount(
					u.SettledSubaccount,
					perpInfos,
					u.LeverageMap,
				)
				if err != nil {
					return false, nil, err
				}
			}

			// Determine whether the state transition is valid.
			result = salib.IsValidStateTransitionForUndercollateralizedSubaccount(
				riskCurMap[saKey],
				riskNew,
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

func (k Keeper) GetNetCollateralAndMarginRequirements(
	ctx sdk.Context,
	update types.Update,
) (
	risk margin.Risk,
	err error,
) {
	// Get leverage configuration for this subaccount
	var leverageMap map[uint32]uint32
	if leverage, found := k.GetLeverage(ctx, &update.SubaccountId); found {
		leverageMap = leverage
	}

	return k.GetNetCollateralAndMarginRequirementsWithLeverage(ctx, update, leverageMap)
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
func (k Keeper) GetNetCollateralAndMarginRequirementsWithLeverage(
	ctx sdk.Context,
	update types.Update,
	leverageMap map[uint32]uint32,
) (
	risk margin.Risk,
	err error,
) {
	subaccount := k.GetSubaccount(ctx, update.SubaccountId)

	perpInfos, err := k.GetAllRelevantPerpetuals(ctx, []types.Update{update})
	if err != nil {
		return risk, err
	}
	settledSubaccount, _ := salib.GetSettledSubaccountWithPerpetuals(subaccount, perpInfos)

	settledUpdate := types.SettledUpdate{
		SettledSubaccount: settledSubaccount,
		AssetUpdates:      update.AssetUpdates,
		PerpetualUpdates:  update.PerpetualUpdates,
	}
	updatedSubaccount := salib.CalculateUpdatedSubaccount(settledUpdate, perpInfos)

	return salib.GetRiskForSubaccount(
		updatedSubaccount,
		perpInfos,
		leverageMap,
	)
}

// GetAllRelevantPerpetuals returns all relevant perpetual information for a given set of updates.
// This includes all perpetuals that exist on the accounts already and all perpetuals that are
// being updated in the input updates.
func (k Keeper) GetAllRelevantPerpetuals(
	ctx sdk.Context,
	updates []types.Update,
) (
	perptypes.PerpInfos,
	error,
) {
	subaccountIdsSet := make(map[types.SubaccountId]struct{})
	perpIdsSet := make(map[uint32]struct{})

	// Add all relevant perpetuals in every update.
	for _, update := range updates {
		// If this subaccount has not been processed already, get all of its existing perpetuals.
		if _, exists := subaccountIdsSet[update.SubaccountId]; !exists {
			sa := k.GetSubaccount(ctx, update.SubaccountId)
			for _, postition := range sa.PerpetualPositions {
				perpIdsSet[postition.PerpetualId] = struct{}{}
			}
			subaccountIdsSet[update.SubaccountId] = struct{}{}
		}

		// Add all perpetuals in the update.
		for _, perpUpdate := range update.PerpetualUpdates {
			perpIdsSet[perpUpdate.GetId()] = struct{}{}
		}
	}

	// Important: Sort the perpIds to ensure determinism!
	sortedPerpIds := lib.GetSortedKeys[lib.Sortable[uint32]](perpIdsSet)

	// Get all perpetual information from state.
	ltCache := make(map[uint32]perptypes.LiquidityTier)
	perpInfos := make(perptypes.PerpInfos, len(sortedPerpIds))
	for _, perpId := range sortedPerpIds {
		perpetual, price, err := k.perpetualsKeeper.GetPerpetualAndMarketPrice(ctx, perpId)
		if err != nil {
			return nil, err
		}

		ltId := perpetual.Params.LiquidityTier
		if _, ok := ltCache[ltId]; !ok {
			liquidityTierFromState, err := k.perpetualsKeeper.GetLiquidityTier(ctx, ltId)
			if err != nil {
				return nil, err
			}
			ltCache[ltId] = liquidityTierFromState
		}
		liquidityTier := ltCache[ltId]

		perpInfos[perpId] = perptypes.PerpInfo{
			Perpetual:     perpetual,
			Price:         price,
			LiquidityTier: liquidityTier,
		}
	}

	return perpInfos, nil
}

func (k Keeper) GetFullNodeStreamingManager() streamingtypes.FullNodeStreamingManager {
	return k.streamingManager
}

// GetInsuranceFundBalance returns the current balance of the specific insurance fund based on the
// perpetual (in quote quantums).
// This calls the Bank Keeper’s GetBalance() function for the Module Address of the insurance fund.
func (k Keeper) GetInsuranceFundBalance(ctx sdk.Context, perpetualId uint32) (balance *big.Int) {
	usdcAsset, exists := k.assetsKeeper.GetAsset(ctx, assettypes.AssetUsdc.Id)
	if !exists {
		panic("GetInsuranceFundBalance: Usdc asset not found in state")
	}
	insuranceFundAddr, err := k.perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, perpetualId)
	if err != nil {
		return nil
	}
	insuranceFundBalance := k.bankKeeper.GetBalance(
		ctx,
		insuranceFundAddr,
		usdcAsset.Denom,
	)

	// Return as big.Int.
	return insuranceFundBalance.Amount.BigInt()
}

func (k Keeper) GetCrossInsuranceFundBalance(ctx sdk.Context) (balance *big.Int) {
	usdcAsset, exists := k.assetsKeeper.GetAsset(ctx, assettypes.AssetUsdc.Id)
	if !exists {
		panic("GetCrossInsuranceFundBalance: Usdc asset not found in state")
	}
	insuranceFundBalance := k.bankKeeper.GetBalance(
		ctx,
		perptypes.InsuranceFundModuleAddress,
		usdcAsset.Denom,
	)

	// Return as big.Int.
	return insuranceFundBalance.Amount.BigInt()
}
