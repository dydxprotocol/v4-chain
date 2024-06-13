package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	perplib "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
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
	requireUniqueSubaccount bool,
) (
	settledUpdates []SettledUpdate,
	subaccountIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt,
	err error,
) {
	var idToSettledSubaccount = make(map[types.SubaccountId]types.Subaccount)
	settledUpdates = make([]SettledUpdate, len(updates))
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

		settledUpdate := SettledUpdate{
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

	settledUpdates, subaccountIdToFundingPayments, err := k.getSettledUpdates(ctx, updates, true)
	if err != nil {
		return false, nil, err
	}

	allPerps := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	success, successPerUpdate, err = k.internalCanUpdateSubaccounts(
		ctx,
		settledUpdates,
		updateType,
		allPerps,
	)

	if !success || err != nil {
		return success, successPerUpdate, err
	}

	// Get a mapping from perpetual Id to current perpetual funding index.
	perpIdToFundingIndex := make(map[uint32]dtypes.SerializableInt, len(allPerps))
	for _, perp := range allPerps {
		perpIdToFundingIndex[perp.Params.Id] = perp.FundingIndex
	}

	// Get OpenInterestDelta from the updates, and persist the OI change if any.
	perpOpenInterestDelta := GetDeltaOpenInterestFromUpdates(settledUpdates, updateType)
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

	// Apply the updates to perpetual positions.
	UpdatePerpetualPositions(
		settledUpdates,
		perpIdToFundingIndex,
	)

	// Apply the updates to asset positions.
	UpdateAssetPositions(settledUpdates)

	// Transfer collateral between collateral pools for any isolated perpetual positions that changed
	// state due to an update.
	for _, settledUpdateWithUpdatedSubaccount := range settledUpdates {
		if err := k.computeAndExecuteCollateralTransfer(
			ctx,
			// The subaccount in `settledUpdateWithUpdatedSubaccount` already has the perpetual updates
			// and asset updates applied to it.
			settledUpdateWithUpdatedSubaccount,
			allPerps,
		); err != nil {
			return false, nil, err
		}
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
			indexerevents.SubaccountUpdateEventVersion,
			indexer_manager.GetBytes(
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

	settledUpdates, _, err := k.getSettledUpdates(ctx, updates, false)
	if err != nil {
		return false, nil, err
	}

	allPerps := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	success, successPerUpdate, err = k.internalCanUpdateSubaccounts(ctx, settledUpdates, updateType, allPerps)
	return success, successPerUpdate, err
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
	// Fetch all relevant perpetuals.
	perpetuals := make(map[uint32]perptypes.Perpetual, len(subaccount.PerpetualPositions))
	for _, p := range subaccount.PerpetualPositions {
		perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, p.PerpetualId)
		if err != nil {
			return types.Subaccount{}, nil, err
		}
		perpetuals[p.PerpetualId] = perpetual
	}

	return GetSettledSubaccountWithPerpetuals(subaccount, perpetuals)
}

// GetSettledSubaccountWithPerpetuals returns 1. a new settled subaccount given an unsettled subaccount,
// updating the USDC AssetPosition, FundingIndex, and LastFundingPayment fields accordingly
// (does not persist any changes) and 2. a map with perpetual ID as key and last funding
// payment as value (for emitting funding payments to indexer).
//
// Note that this is a stateless utility function.
func GetSettledSubaccountWithPerpetuals(
	subaccount types.Subaccount,
	perpetuals map[uint32]perptypes.Perpetual,
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
		perpetual, found := perpetuals[p.PerpetualId]
		if !found {
			return types.Subaccount{},
				nil,
				errorsmod.Wrap(
					perptypes.ErrPerpetualDoesNotExist, lib.UintToString(p.PerpetualId),
				)
		}

		// Call the stateless utility function to get the net settlement and new funding index.
		bigNetSettlementPpm, newFundingIndex := perplib.GetSettlementPpmWithPerpetual(
			perpetual,
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

func checkPositionUpdatable(
	ctx sdk.Context,
	pk types.ProductKeeper,
	p types.PositionSize,
) (
	err error,
) {
	updatable, err := pk.IsPositionUpdatable(
		ctx,
		p.GetId(),
	)
	if err != nil {
		return err
	}

	if !updatable {
		return errorsmod.Wrapf(
			types.ErrProductPositionNotUpdatable,
			"type: %v, id: %d",
			p.GetProductType(),
			p.GetId(),
		)
	}
	return nil
}

// internalCanUpdateSubaccounts will validate all `updates` to the relevant subaccounts and compute
// if any of the updates led to an isolated perpetual position being opened or closed.
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
	settledUpdates []SettledUpdate,
	updateType types.UpdateType,
	perpetuals []perptypes.Perpetual,
) (
	success bool,
	successPerUpdate []types.UpdateResult,
	err error,
) {
	// TODO(TRA-99): Add integration / E2E tests on order placement / matching with this new
	// constraint.
	// Check if the updates satisfy the isolated perpetual constraints.
	success, successPerUpdate, err = k.checkIsolatedSubaccountConstraints(
		ctx,
		settledUpdates,
		perpetuals,
	)
	if err != nil {
		return false, nil, err
	}
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
	perpOpenInterestDelta := GetDeltaOpenInterestFromUpdates(settledUpdates, updateType)

	bigCurNetCollateral := make(map[string]*big.Int)
	bigCurInitialMargin := make(map[string]*big.Int)
	bigCurMaintenanceMargin := make(map[string]*big.Int)

	// Iterate over all updates.
	for i, u := range settledUpdates {
		// Check all updated perps are updatable.
		for _, perpUpdate := range u.PerpetualUpdates {
			err := checkPositionUpdatable(ctx, k.perpetualsKeeper, perpUpdate)
			if err != nil {
				return false, nil, err
			}
		}

		// Check all updated assets are updatable.
		for _, assetUpdate := range u.AssetUpdates {
			err := checkPositionUpdatable(ctx, k.assetsKeeper, assetUpdate)
			if err != nil {
				return false, nil, err
			}
		}

		// Branch the state to calculate the new OIMF after OI increase.
		// The branched state is only needed for this purpose and is always discarded.
		branchedContext, _ := ctx.CacheContext()

		// Temporily apply open interest delta to perpetuals, so IMF is calculated based on open interest after the update.
		// `perpOpenInterestDeltas` is only present for `Match` update type.
		if perpOpenInterestDelta != nil {
			if err := k.perpetualsKeeper.ModifyOpenInterest(
				branchedContext,
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
		// Get the new collateralization and margin requirements with the update applied.
		bigNewNetCollateral,
			bigNewInitialMargin,
			bigNewMaintenanceMargin,
			err := k.internalGetNetCollateralAndMarginRequirements(
			branchedContext,
			u,
		)

		// if `internalGetNetCollateralAndMarginRequirements`, returns error.
		if err != nil {
			return false, nil, err
		}

		var result = types.Success

		// The subaccount is not well-collateralized after the update.
		// We must now check if the state transition is valid.
		if bigNewInitialMargin.Cmp(bigNewNetCollateral) > 0 {
			// Get the current collateralization and margin requirements without the update applied.
			emptyUpdate := SettledUpdate{
				SettledSubaccount: u.SettledSubaccount,
			}

			bytes, err := proto.Marshal(u.SettledSubaccount.Id)
			if err != nil {
				return false, nil, err
			}
			saKey := string(bytes)

			// Cache the current collateralization and margin requirements for the subaccount.
			if _, ok := bigCurNetCollateral[saKey]; !ok {
				bigCurNetCollateral[saKey],
					bigCurInitialMargin[saKey],
					bigCurMaintenanceMargin[saKey],
					err = k.internalGetNetCollateralAndMarginRequirements(
					ctx,
					emptyUpdate,
				)
				if err != nil {
					return false, nil, err
				}
			}

			// Determine whether the state transition is valid.
			result = IsValidStateTransitionForUndercollateralizedSubaccount(
				bigCurNetCollateral[saKey],
				bigCurInitialMargin[saKey],
				bigCurMaintenanceMargin[saKey],
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

	settledUpdate := SettledUpdate{
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
	settledUpdate SettledUpdate,
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

	// Iterate over all assets and updates and calculate change to net collateral and margin requirements.
	for _, size := range assetSizes {
		id := size.GetId()
		bigQuantums := size.GetBigQuantums()

		nc, err := k.assetsKeeper.GetNetCollateral(ctx, id, bigQuantums)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
		}

		imr, mmr, err := k.assetsKeeper.GetMarginRequirements(
			ctx,
			id,
			bigQuantums,
		)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
		}
		bigNetCollateral.Add(bigNetCollateral, nc)
		bigInitialMargin.Add(bigInitialMargin, imr)
		bigMaintenanceMargin.Add(bigMaintenanceMargin, mmr)
	}

	// Iterate over all perpetuals and updates and calculate change to net collateral and margin requirements.
	for _, size := range perpetualSizes {
		perpetual,
			marketPrice,
			liquidityTier,
			err := k.perpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(ctx, size.GetId())
		if err != nil {
			return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
		}
		nc, imr, mmr := perplib.GetNetCollateralAndMarginRequirements(
			perpetual,
			marketPrice,
			liquidityTier,
			size.GetBigQuantums(),
		)
		bigNetCollateral.Add(bigNetCollateral, nc)
		bigInitialMargin.Add(bigInitialMargin, imr)
		bigMaintenanceMargin.Add(bigMaintenanceMargin, mmr)
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
