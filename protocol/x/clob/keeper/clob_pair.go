package keeper

import (
	"fmt"
	"sort"

	gogotypes "github.com/cosmos/gogoproto/types"

	storetypes "cosmossdk.io/store/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	dydxlog "github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// getClobPairStore returns a prefix store where the ClobPair objects are stored.
func (k Keeper) getClobPairStore(
	ctx sdk.Context,
) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.ClobPairKeyPrefix))
}

// clobPairKey returns the store key to retrieve a ClobPair by id.
func clobPairKey(
	id types.ClobPairId,
) []byte {
	return lib.Uint32ToKey(id.ToUint32())
}

// (Same behavior as old `CreatePerpetualClobPair`) creates the
// objects in state and in-memory.
// This function should ONLY be used to initialize genesis state
// or test setups.
// Regular CLOB logic should use `CreatePerpetualClobPair` instead.
func (k Keeper) CreatePerpetualClobPairAndMemStructs(
	ctx sdk.Context,
	clobPairId uint32,
	perpetualId uint32,
	stepSizeBaseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	status types.ClobPair_Status,
) (types.ClobPair, error) {
	clobPair, err := k.createPerpetualClobPair(
		ctx,
		clobPairId,
		perpetualId,
		stepSizeBaseQuantums,
		quantumConversionExponent,
		subticksPerTick,
		status,
	)
	if err != nil {
		return types.ClobPair{}, err
	}

	k.MemClob.CreateOrderbook(clobPair)
	k.SetClobPairIdForPerpetual(clobPair)

	return clobPair, nil
}

// CreatePerpetualClobPair creates a new perpetual CLOB pair in the store.
// Additionally, it creates an order book matching the ID of the newly created CLOB pair.
//
// An error will occur if any of the fields fail validation (see ValidateClobPair for details),
// or if the `perpetualId` cannot be found.
// In the event of an error, the store will not be updated nor will a matching order book be created.
//
// Returns the newly created CLOB pair and an error if one occurs.
func (k Keeper) CreatePerpetualClobPair(
	ctx sdk.Context,
	clobPairId uint32,
	perpetualId uint32,
	stepSizeBaseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	status types.ClobPair_Status,
) (types.ClobPair, error) {
	clobPair, err := k.createPerpetualClobPair(
		ctx,
		clobPairId,
		perpetualId,
		stepSizeBaseQuantums,
		quantumConversionExponent,
		subticksPerTick,
		status,
	)
	if err != nil {
		return types.ClobPair{}, err
	}

	// Don't stage events for the genesis block.
	if lib.IsDeliverTxMode(ctx) {
		if err := k.StageNewClobPairSideEffects(ctx, clobPair); err != nil {
			return clobPair, err
		}
	}

	return clobPair, nil
}

func (k Keeper) createPerpetualClobPair(
	ctx sdk.Context,
	clobPairId uint32,
	perpetualId uint32,
	stepSizeBaseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	status types.ClobPair_Status,
) (types.ClobPair, error) {
	clobPair := types.ClobPair{
		Metadata: &types.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &types.PerpetualClobMetadata{
				PerpetualId: perpetualId,
			},
		},
		Id:                        clobPairId,
		StepBaseQuantums:          stepSizeBaseQuantums.ToUint64(),
		QuantumConversionExponent: quantumConversionExponent,
		SubticksPerTick:           subticksPerTick,
		Status:                    status,
	}
	if err := k.ValidateClobPairCreation(ctx, &clobPair); err != nil {
		return clobPair, err
	}

	// Write the `ClobPair` to state.
	k.SetClobPair(ctx, clobPair)

	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		panic(err)
	}
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return types.ClobPair{}, err
	}

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypePerpetualMarket,
		indexerevents.PerpetualMarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewPerpetualMarketCreateEvent(
				perpetualId,
				clobPair.Id,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				clobPair.Status,
				clobPair.QuantumConversionExponent,
				perpetual.Params.AtomicResolution,
				clobPair.SubticksPerTick,
				clobPair.StepBaseQuantums,
				perpetual.Params.LiquidityTier,
				perpetual.Params.MarketType,
				perpetual.Params.DefaultFundingPpm,
			),
		),
	)

	return clobPair, nil
}

// ValidateClobPairCreation validates a CLOB pair's fields are suitable for CLOB pair creation
// and that the perpetual ID is associated with an existing perpetual.
func (k Keeper) ValidateClobPairCreation(ctx sdk.Context, clobPair *types.ClobPair) error {
	// If the desired CLOB pair ID is already in use, return an error.
	if clobPair, exists := k.GetClobPair(ctx, clobPair.GetClobPairId()); exists {
		return errorsmod.Wrapf(
			types.ErrClobPairAlreadyExists,
			"id=%v, existing clob pair=%v",
			clobPair.Id,
			clobPair,
		)
	}

	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		return errorsmod.Wrap(
			types.ErrInvalidClobPairParameter,
			err.Error(),
		)
	}

	// Verify the perpetual ID is not already associated with an existing CLOB pair.
	if clobPairId, found := k.PerpetualIdToClobPairId[perpetualId]; found {
		return errorsmod.Wrapf(
			types.ErrPerpetualAssociatedWithExistingClobPair,
			"perpetual id=%v, existing clob pair id=%v",
			perpetualId,
			clobPairId,
		)
	}

	return k.validateClobPair(ctx, clobPair)
}

// validateClobPair validates a CLOB pair's fields are suitable for CLOB pair creation.
//
// Stateful Validation:
//   - Must be a perpetual CLOB pair with a perpetualId matching a perpetual in the store.
//
// Stateless Validation
//   - `clobPair.Validate()` returns no error.
func (k Keeper) validateClobPair(ctx sdk.Context, clobPair *types.ClobPair) error {
	if err := clobPair.Validate(); err != nil {
		return err
	}

	// TODO(DEC-1535): update this validation when we implement "spot"/"asset" clob pairs.
	switch clobPair.Metadata.(type) {
	case *types.ClobPair_PerpetualClobMetadata:
		perpetualId, err := clobPair.GetPerpetualId()
		if err != nil {
			return errorsmod.Wrapf(
				err,
				"CLOB pair (%+v) has invalid perpetual.",
				clobPair,
			)
		}
		// Validate the perpetual referenced by the CLOB pair exists.
		if _, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId); err != nil {
			return errorsmod.Wrapf(
				err,
				"CLOB pair (%+v) has invalid perpetual.",
				clobPair,
			)
		}
	default:
		return errorsmod.Wrapf(
			types.ErrInvalidClobPairParameter,
			// TODO(DEC-1535): update this error message when we implement "spot"/"asset" clob pairs.
			"CLOB pair (%+v) is not a perpetual CLOB.",
			clobPair,
		)
	}
	return nil
}

func (k Keeper) ApplySideEffectsForNewClobPair(
	ctx sdk.Context,
	clobPair types.ClobPair,
) {
	if created := k.MemClob.MaybeCreateOrderbook(clobPair); !created {
		dydxlog.ErrorLog(
			ctx,
			"ApplySideEffectsForNewClobPair: Orderbook already exists for CLOB pair",
			"clob_pair", clobPair,
		)
		return
	}
	k.SetClobPairIdForPerpetual(clobPair)
}

// StageNewClobPairSideEffects stages a ClobPair creation event, so that any in-memory side effects
// can happen later when the transaction and block is committed.
// Note the staged event will be processed only if below are both true:
//   - The current transaction is committed.
//   - The current block is agreed upon and committed by consensus.
func (k Keeper) StageNewClobPairSideEffects(ctx sdk.Context, clobPair types.ClobPair) error {
	lib.AssertDeliverTxMode(ctx)

	k.finalizeBlockEventStager.StageFinalizeBlockEvent(
		ctx,
		&types.ClobStagedFinalizeBlockEvent{
			Event: &types.ClobStagedFinalizeBlockEvent_CreateClobPair{
				CreateClobPair: &clobPair,
			},
		},
	)

	return nil
}

// SetClobPair sets a specific `ClobPair` in the store from its index.
func (k Keeper) SetClobPair(ctx sdk.Context, clobPair types.ClobPair) {
	store := k.getClobPairStore(ctx)
	b := k.cdc.MustMarshal(&clobPair)
	store.Set(clobPairKey(clobPair.GetClobPairId()), b)
}

// InitMemClobOrderbooks initializes the memclob with `ClobPair`s from state.
func (k Keeper) InitMemClobOrderbooks(ctx sdk.Context) {
	clobPairs := k.GetAllClobPairs(ctx)
	for _, clobPair := range clobPairs {
		// Create the corresponding orderbook in the memclob.
		k.MemClob.MaybeCreateOrderbook(
			clobPair,
		)
	}
}

// HydrateClobPairAndPerpetualMapping hydrates the in-memory mapping between clob pair and perpetual.
func (k Keeper) HydrateClobPairAndPerpetualMapping(ctx sdk.Context) {
	clobPairs := k.GetAllClobPairs(ctx)
	for _, clobPair := range clobPairs {
		// Create the corresponding mapping between clob pair and perpetual.
		k.SetClobPairIdForPerpetual(
			clobPair,
		)
	}
}

// SetClobPairIdForPerpetual sets the mapping between clob pair and perpetual.
func (k Keeper) SetClobPairIdForPerpetual(clobPair types.ClobPair) {
	// If this `ClobPair` is for a perpetual, add the `clobPairId` to the list of CLOB pair IDs
	// that facilitate trading of this perpetual.
	if perpetualClobMetadata := clobPair.GetPerpetualClobMetadata(); perpetualClobMetadata != nil {
		perpetualId := perpetualClobMetadata.PerpetualId
		clobPairIds, exists := k.PerpetualIdToClobPairId[perpetualId]
		if !exists {
			clobPairIds = make([]types.ClobPairId, 0)
		}

		for _, clobPairId := range clobPairIds {
			if clobPairId == clobPair.GetClobPairId() {
				// The mapping already exists, so return.
				return
			}
		}

		k.PerpetualIdToClobPairId[perpetualId] = append(
			clobPairIds,
			clobPair.GetClobPairId(),
		)
	}
}

// GetClobPairIdForPerpetual gets the first CLOB pair ID associated with the provided perpetual ID.
// It returns an error if there are no CLOB pair IDs associated with the perpetual ID.
func (k Keeper) GetClobPairIdForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
) (
	clobPairId types.ClobPairId,
	err error,
) {
	clobPairIds, exists := k.PerpetualIdToClobPairId[perpetualId]
	if !exists {
		return 0, errorsmod.Wrapf(
			types.ErrNoClobPairForPerpetual,
			"Perpetual ID %d has no associated CLOB pairs",
			perpetualId,
		)
	}

	if len(clobPairIds) == 0 {
		panic("GetClobPairIdForPerpetual: Perpetual ID was created without a CLOB pair ID.")
	}

	if len(clobPairIds) > 1 {
		panic("GetClobPairIdForPerpetual: Perpetual ID was created with multiple CLOB pair IDs.")
	}

	return clobPairIds[0], nil
}

// GetClobPair returns a clobPair from its index
func (k Keeper) GetClobPair(
	ctx sdk.Context,
	id types.ClobPairId,
) (val types.ClobPair, found bool) {
	store := k.getClobPairStore(ctx)

	b := store.Get(clobPairKey(id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveClobPair removes a clobPair from the store
func (k Keeper) RemoveClobPair(
	ctx sdk.Context,
	id types.ClobPairId,
) {
	store := k.getClobPairStore(ctx)
	store.Delete(clobPairKey(id))
}

// GetAllClobPairs returns all clobPair, sorted by ClobPair id.
func (k Keeper) GetAllClobPairs(ctx sdk.Context) (list []types.ClobPair) {
	store := k.getClobPairStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClobPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})

	return
}

// validateLiquidationAgainstClobPairStatus returns an error if placing the provided
// liquidation order would conflict with the clob pair's current status.
func (k Keeper) validateLiquidationAgainstClobPairStatus(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
) error {
	clobPair, found := k.GetClobPair(ctx, liquidationOrder.GetClobPairId())
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidClob,
			"Clob %v is not a valid clob",
			liquidationOrder.GetClobPairId(),
		)
	}

	if !types.IsSupportedClobPairStatus(clobPair.Status) {
		panic(
			fmt.Sprintf(
				"validateLiquidationAgainstClobPairStatus: clob pair status %v is not supported",
				clobPair.Status,
			),
		)
	}

	if clobPair.Status != types.ClobPair_STATUS_ACTIVE {
		return errorsmod.Wrapf(
			types.ErrLiquidationConflictsWithClobPairStatus,
			"Liquidation order %+v cannot be placed for clob pair with status %+v",
			liquidationOrder,
			clobPair.Status,
		)
	}

	return nil
}

// validateOrderAgainstClobPairStatus returns an error if placing the provided
// order would conflict with the clob pair's current status.
func (k Keeper) validateOrderAgainstClobPairStatus(
	ctx sdk.Context,
	order types.Order,
	clobPair types.ClobPair,
) error {
	if !types.IsSupportedClobPairStatus(clobPair.Status) {
		// Validation should only be called against ClobPairs in state, implying we have a ClobPair with
		// an unsupported status in state.
		panic(
			fmt.Sprintf(
				"validateOrderAgainstClobPairStatus: clob pair status %v is not supported",
				clobPair.Status,
			),
		)
	}

	switch clobPair.Status {
	case types.ClobPair_STATUS_INITIALIZING:
		// Reject stateful orders. Short-term orders expire within the ShortBlockWindow, ensuring
		// stale short term order expiration. Stateful orders are rejected to prevent long-term orders
		// from controlling the price at which new orders can be placed. This is necessary when all orders
		// are post-only.
		if order.IsStatefulOrder() {
			return errorsmod.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order %+v must not be stateful for clob pair with status %+v",
				order,
				clobPair.Status,
			)
		}

		// Reject non-post-only orders. During the initializing phase we only allow post-only orders.
		// This allows liquidity to build around the oracle price without any real trading happening.
		if order.TimeInForce != types.Order_TIME_IN_FORCE_POST_ONLY {
			return errorsmod.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order %+v must be post-only for clob pair with status %+v",
				order,
				clobPair.Status,
			)
		}

		// Reject orders on the wrong side of the market. This is to ensure liquidity is building around
		// the oracle price. For instance without this check a user could place an ask far below the oracle
		// price, thereby preventing any bids at or above the specified price of the ask.
		currentOraclePriceSubticksRat := k.GetOraclePriceSubticksRat(ctx, clobPair)
		currentOraclePriceSubticks := lib.BigRatRound(currentOraclePriceSubticksRat, false).Uint64()
		// Throw error if order is a buy and order subticks is greater than oracle price subticks
		if order.IsBuy() && order.Subticks > currentOraclePriceSubticks {
			return errorsmod.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order subticks %+v must be less than or equal to oracle price subticks %+v for clob pair with status %+v",
				order.Subticks,
				currentOraclePriceSubticks,
				clobPair.Status,
			)
		}
		// Throw error if order is a sell and order subticks is less than oracle price subticks
		if !order.IsBuy() && order.Subticks < currentOraclePriceSubticks {
			return errorsmod.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order subticks %+v must be greater than or equal to oracle price subticks %+v for clob pair with status %+v",
				order.Subticks,
				currentOraclePriceSubticks,
				clobPair.Status,
			)
		}
	case types.ClobPair_STATUS_FINAL_SETTLEMENT:
		return errorsmod.Wrapf(
			types.ErrOrderConflictsWithClobPairStatus,
			"Order %+v disallowed, trading is disabled for clob pair with status %+v",
			order,
			clobPair.Status,
		)
	}

	return nil
}

// mustGetClobPair fetches a ClobPair from state given its id.
// This function panics if the ClobPair is not found.
func (k Keeper) mustGetClobPair(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
) types.ClobPair {
	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		panic(
			fmt.Sprintf(
				"mustGetClobPair: ClobPair with id %+v not found",
				clobPairId,
			),
		)
	}
	return clobPair
}

// mustGetClobPairForPerpetualId fetches a ClobPair from state given a perpetual id.
// This function panics if the ClobPair is not found.
func (k Keeper) mustGetClobPairForPerpetualId(
	ctx sdk.Context,
	perpetualId uint32,
) types.ClobPair {
	clobPairId, err := k.GetClobPairIdForPerpetual(ctx, perpetualId)
	if err != nil {
		panic(err)
	}
	return k.mustGetClobPair(ctx, clobPairId)
}

// UpdateClobPair overwrites a ClobPair in state and sends an update to the indexer.
// This function returns an error if the update includes an unsupported transition
// for the ClobPair's status.
func (k Keeper) UpdateClobPair(
	ctx sdk.Context,
	clobPair types.ClobPair,
) error {
	oldClobPair, found := k.GetClobPair(ctx, types.ClobPairId(clobPair.Id))
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidClobPairUpdate,
			"UpdateClobPair: ClobPair with id %d not found in state",
			clobPair.Id,
		)
	}

	// Note, only perpetual clob pairs are currently supported. Neither the old nor the
	// new clob pair should be spot.
	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		return errorsmod.Wrap(
			types.ErrInvalidClobPairUpdate,
			err.Error(),
		)
	}
	oldPerpetualId, err := oldClobPair.GetPerpetualId()
	if err != nil {
		return errorsmod.Wrap(
			types.ErrInvalidClobPairUpdate,
			err.Error(),
		)
	}
	if perpetualId != oldPerpetualId {
		return errorsmod.Wrap(
			types.ErrInvalidClobPairUpdate,
			"UpdateClobPair: cannot update ClobPair perpetual id",
		)
	}
	if clobPair.StepBaseQuantums != oldClobPair.StepBaseQuantums {
		return errorsmod.Wrapf(
			types.ErrInvalidClobPairUpdate,
			"UpdateClobPair: cannot update ClobPair step base quantums",
		)
	}
	if clobPair.SubticksPerTick != oldClobPair.SubticksPerTick {
		return errorsmod.Wrapf(
			types.ErrInvalidClobPairUpdate,
			"UpdateClobPair: cannot update ClobPair subticks per tick",
		)
	}
	if clobPair.QuantumConversionExponent != oldClobPair.QuantumConversionExponent {
		return errorsmod.Wrapf(
			types.ErrInvalidClobPairUpdate,
			"UpdateClobPair: cannot update ClobPair quantum conversion exponent",
		)
	}

	oldStatus := oldClobPair.Status
	newStatus := clobPair.Status
	if !types.IsSupportedClobPairStatusTransition(oldStatus, newStatus) {
		return errorsmod.Wrapf(
			types.ErrInvalidClobPairStatusTransition,
			"Cannot transition from status %+v to status %+v",
			oldStatus,
			newStatus,
		)
	}

	if err := k.validateClobPair(ctx, &clobPair); err != nil {
		return err
	}

	k.SetClobPair(ctx, clobPair)

	// Send UpdateClobPair to indexer.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeUpdateClobPair,
		indexerevents.UpdateClobPairEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewUpdateClobPairEvent(
				clobPair.GetClobPairId(),
				clobPair.Status,
				clobPair.QuantumConversionExponent,
				types.SubticksPerTick(clobPair.GetSubticksPerTick()),
				satypes.BaseQuantums(clobPair.GetStepBaseQuantums()),
			),
		),
	)

	// If newly transitioning to final settlement, enter final settlement.
	if newStatus == types.ClobPair_STATUS_FINAL_SETTLEMENT && oldStatus != newStatus {
		k.mustTransitionToFinalSettlement(ctx, clobPair.GetClobPairId())
	}

	return nil
}

// getInternalOperationClobPairId returns the ClobPairId associated with the operation. This function
// will panic if called for a PreexistingStatefulOrder internal operation since this operation type
// should never be included in MsgProposedOperations.
func (k Keeper) getInternalOperationClobPairId(
	ctx sdk.Context,
	internalOperation types.InternalOperation,
) (
	clobPairId types.ClobPairId,
	err error,
) {
	switch castedOperation := internalOperation.Operation.(type) {
	case *types.InternalOperation_Match:
		switch castedMatch := castedOperation.Match.Match.(type) {
		case *types.ClobMatch_MatchOrders:
			clobPairId = types.ClobPairId(castedMatch.MatchOrders.TakerOrderId.ClobPairId)
		case *types.ClobMatch_MatchPerpetualLiquidation:
			clobPairId = types.ClobPairId(castedMatch.MatchPerpetualLiquidation.ClobPairId)
		case *types.ClobMatch_MatchPerpetualDeleveraging:
			clobPairId, err = k.GetClobPairIdForPerpetual(
				ctx,
				castedMatch.MatchPerpetualDeleveraging.PerpetualId,
			)
		}
	case *types.InternalOperation_ShortTermOrderPlacement:
		clobPairId = types.ClobPairId(castedOperation.ShortTermOrderPlacement.Order.OrderId.ClobPairId)
	case *types.InternalOperation_OrderRemoval:
		clobPairId = types.ClobPairId(castedOperation.OrderRemoval.OrderId.ClobPairId)
	case *types.InternalOperation_PreexistingStatefulOrder:
		// this helper is only used in ProcessOperations (DeliverTx) which should not contain
		// this operation type, so panic.
		panic(
			"getInternalOperationClobPairId: should never be called for preexisting stateful order " +
				"internal operations",
		)
	default:
		panic(
			fmt.Sprintf(
				"getInternalOperationClobPairId: Unrecognized operation type for operation: %+v",
				internalOperation.GetInternalOperationTextString(),
			),
		)
	}

	return clobPairId, err
}

// validateInternalOperationAgainstClobPairStatus validates that an internal
// operation is valid for its associated ClobPair's current status. This function will panic if the
// ClobPair cannot be found or if the ClobPair's status is not supported.
// Returns an error if one is encountered during validation.
func (k Keeper) validateInternalOperationAgainstClobPairStatus(
	ctx sdk.Context,
	internalOperation types.InternalOperation,
) error {
	clobPairId, err := k.getInternalOperationClobPairId(ctx, internalOperation)
	if err != nil {
		return err
	}

	// Fail if the ClobPair cannot be found.
	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidClob,
			"CLOB pair ID %d not found in state",
			clobPairId,
		)
	}

	// Verify ClobPair fetched from state has a supported status.
	if !types.IsSupportedClobPairStatus(clobPair.Status) {
		panic(
			"validateInternalOperationAgainstClobPairStatus: ClobPair's status is not supported",
		)
	}

	// Branch validation logic for supported statuses requiring validation.
	switch clobPair.Status {
	case types.ClobPair_STATUS_INITIALIZING:
		// All operations are invalid for initializing clob pairs.
		return errorsmod.Wrapf(
			types.ErrOperationConflictsWithClobPairStatus,
			"Operation %s invalid for ClobPair with id %d with status %s",
			internalOperation.GetInternalOperationTextString(),
			clobPairId,
			types.ClobPair_STATUS_INITIALIZING,
		)
	case types.ClobPair_STATUS_FINAL_SETTLEMENT:
		// Only allow deleveraging events. This allows the protocol to close out open
		// positions in the market. All other operations are not allowed.
		if match := internalOperation.GetMatch(); match != nil && match.GetMatchPerpetualDeleveraging() != nil {
			return nil
		}
		return errorsmod.Wrapf(
			types.ErrOperationConflictsWithClobPairStatus,
			"Operation %s invalid for ClobPair with id %d with status %s",
			internalOperation.GetInternalOperationTextString(),
			clobPairId,
			types.ClobPair_STATUS_FINAL_SETTLEMENT,
		)
	}

	return nil
}

// IsPerpetualClobPairActive returns true if the ClobPair associated with the provided perpetual id
// has the active status. Returns an error if the ClobPair cannot be found.
func (k Keeper) IsPerpetualClobPairActive(
	ctx sdk.Context,
	perpetualId uint32,
) (bool, error) {
	clobPairId, err := k.GetClobPairIdForPerpetual(ctx, perpetualId)
	if err != nil {
		return false, err
	}

	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		return false, errorsmod.Wrapf(
			types.ErrInvalidClob,
			"IsPerpetualClobPairActive: did not find clob pair with id = %d",
			clobPairId,
		)
	}

	return clobPair.Status == types.ClobPair_STATUS_ACTIVE, nil
}

// GetNextClobPairID returns the next clob pair id to be used from the module store
func (k Keeper) GetNextClobPairID(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.NextClobPairIDKey))
	var result gogotypes.UInt32Value
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// SetNextClobPairID sets the next clob pair id to be used
func (k Keeper) SetNextClobPairID(ctx sdk.Context, nextID uint32) {
	store := ctx.KVStore(k.storeKey)
	value := gogotypes.UInt32Value{Value: nextID}
	store.Set([]byte(types.NextClobPairIDKey), k.cdc.MustMarshal(&value))
}

// AcquireNextClobPairID returns the next clob pair id to be used and increments
// the next clob pair id
func (k Keeper) AcquireNextClobPairID(ctx sdk.Context) uint32 {
	nextID := k.GetNextClobPairID(ctx)
	// if clob pair id already exists, increment until we find one that doesn't
	maxAttempts, attempts := 1000, 0
	for {
		_, found := k.GetClobPair(ctx, types.ClobPairId(nextID))
		if !found {
			break
		}
		nextID++

		// panic if we've tried too many times and are stuck in a loop
		attempts++
		if attempts >= maxAttempts {
			panic("Exceeded maximum attempts to find a unique clob pair id")
		}
	}

	k.SetNextClobPairID(ctx, nextID+1)
	return nextID
}
