package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CreatePerpetualClobPair creates a new perpetual CLOB pair in the store.
// Additionally, it creates an order book matching the ID of the newly created CLOB pair.
//
// An error will occur if any of the fields fail validation (see validateClobPair for details),
// or if the `perpetualId` cannot be found.
// In the event of an error, the store will not be updated nor will a matching order book be created.
//
// Returns the newly created CLOB pair and an error if one occurs.
func (k Keeper) CreatePerpetualClobPair(
	ctx sdk.Context,
	clobPairId uint32,
	perpetualId uint32,
	minOrderBaseQuantums satypes.BaseQuantums,
	stepSizeBaseQuantums satypes.BaseQuantums,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	status types.ClobPair_Status,
) (types.ClobPair, error) {
	// If the desired CLOB pair ID is already in use, return an error.
	if clobPair, exists := k.GetClobPair(ctx, types.ClobPairId(clobPairId)); exists {
		return types.ClobPair{}, sdkerrors.Wrapf(
			types.ErrClobPairAlreadyExists,
			"id=%v, existing clob pair=%v",
			clobPairId,
			clobPair,
		)
	}

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
	if err := k.validateClobPair(ctx, &clobPair); err != nil {
		return clobPair, err
	}
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return clobPair, err
	}

	k.createClobPair(ctx, clobPair)
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypePerpetualMarket,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewPerpetualMarketCreateEvent(
				perpetualId,
				clobPairId,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				status,
				quantumConversionExponent,
				perpetual.Params.AtomicResolution,
				subticksPerTick,
				minOrderBaseQuantums.ToUint64(),
				stepSizeBaseQuantums.ToUint64(),
				perpetual.Params.LiquidityTier,
			),
		),
	)

	return clobPair, nil
}

// validateClobPair validates a CLOB pair's fields are suitable for CLOB pair creation.
//
// - Metadata:
//   - Must be a perpetual CLOB pair with a perpetualId matching a perpetual in the store.
//
// - Status:
//   - Must be a supported status.
//
// - StepBaseQuantums:
//   - Must be greater than zero.
//
// - SubticksPerTick:
//   - Must be greater than zero.
func (k Keeper) validateClobPair(ctx sdk.Context, clobPair *types.ClobPair) error {
	if !types.IsSupportedClobPairStatus(clobPair.Status) {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"CLOB pair (%+v) has unsupported status %+v",
			clobPair,
			clobPair.Status,
		)
	}

	// TODO(DEC-1535): update this validation when we implement "spot"/"asset" clob pairs.
	switch clobPair.Metadata.(type) {
	case *types.ClobPair_PerpetualClobMetadata:
		perpetualId, err := clobPair.GetPerpetualId()
		if err != nil {
			return sdkerrors.Wrapf(
				err,
				"CLOB pair (%+v) has invalid perpetual.",
				clobPair,
			)
		}
		// Validate the perpetual referenced by the CLOB pair exists.
		if _, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId); err != nil {
			return sdkerrors.Wrapf(
				err,
				"CLOB pair (%+v) has invalid perpetual.",
				clobPair,
			)
		}
	default:
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			// TODO(DEC-1535): update this error message when we implement "spot"/"asset" clob pairs.
			"CLOB pair (%+v) is not a perpetual CLOB.",
			clobPair,
		)
	}

	if clobPair.StepBaseQuantums <= 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"invalid ClobPair parameter: StepBaseQuantums must be > 0. Got %v",
			clobPair.StepBaseQuantums,
		)
	}

	// Since a subtick will be calculated as (1 tick/SubticksPerTick), the denominator cannot be 0
	// and negative numbers do not make sense.
	if clobPair.SubticksPerTick <= 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairParameter,
			"invalid ClobPair parameter: SubticksPerTick must be > 0. Got %v",
			clobPair.SubticksPerTick,
		)
	}

	return nil
}

// createOrderbook creates a new orderbook in the memclob and stores the perpetualId to clobPairId mapping
// in memory on the keeper.
func (k Keeper) createOrderbook(ctx sdk.Context, clobPair types.ClobPair) {
	// Create the corresponding orderbook in the memclob.
	k.MemClob.CreateOrderbook(ctx, clobPair)

	// If this `ClobPair` is for a perpetual, add the `clobPairId` to the list of CLOB pair IDs
	// that facilitate trading of this perpetual.
	if perpetualClobMetadata := clobPair.GetPerpetualClobMetadata(); perpetualClobMetadata != nil {
		perpetualId := perpetualClobMetadata.PerpetualId
		clobPairIds, exists := k.PerpetualIdToClobPairId[perpetualId]
		if !exists {
			clobPairIds = make([]types.ClobPairId, 0)
		}
		k.PerpetualIdToClobPairId[perpetualId] = append(
			clobPairIds,
			clobPair.GetClobPairId(),
		)
	}
}

// createClobPair creates a new `ClobPair` in the store and creates the corresponding orderbook in the memclob.
// This function returns an error if a value for the ClobPair's id already exists in state.
func (k Keeper) createClobPair(ctx sdk.Context, clobPair types.ClobPair) {
	// Validate the given clob pair id is not already in use.
	if _, exists := k.GetClobPair(ctx, clobPair.GetClobPairId()); exists {
		panic(
			fmt.Sprintf(
				"ClobPair with id %+v already exists in state",
				clobPair.GetClobPairId(),
			),
		)
	}

	// Write the `ClobPair` to state.
	k.setClobPair(ctx, clobPair)

	// Create the corresponding orderbook in the memclob.
	k.createOrderbook(ctx, clobPair)
}

// setClobPair sets a specific `ClobPair` in the store from its index.
func (k Keeper) setClobPair(ctx sdk.Context, clobPair types.ClobPair) {
	b := k.cdc.MustMarshal(&clobPair)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	// Write the `ClobPair` to state.
	store.Set(types.ClobPairKey(clobPair.GetClobPairId()), b)
}

// InitMemClobOrderbooks initializes the memclob with `ClobPair`s from state.
// This is called during app initialization in `app.go`, before any ABCI calls are received.
func (k Keeper) InitMemClobOrderbooks(ctx sdk.Context) {
	clobPairs := k.GetAllClobPair(ctx)
	for _, clobPair := range clobPairs {
		// Create the corresponding orderbook in the memclob.
		k.createOrderbook(
			ctx,
			clobPair,
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
		return 0, sdkerrors.Wrapf(
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))

	b := store.Get(types.ClobPairKey(
		id,
	))
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	store.Delete(types.ClobPairKey(
		id,
	))
}

// GetAllClobPair returns all clobPair
func (k Keeper) GetAllClobPair(ctx sdk.Context) (list []types.ClobPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ClobPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
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
			return sdkerrors.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order %+v must not be stateful for clob pair with status %+v",
				order,
				clobPair.Status,
			)
		}

		// Reject non-post-only orders. During the initializing phase we only allow post-only orders.
		// This allows liquidity to build around the oracle price without any real trading happening.
		if order.TimeInForce != types.Order_TIME_IN_FORCE_POST_ONLY {
			return sdkerrors.Wrapf(
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
			return sdkerrors.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order subticks %+v must be less than or equal to oracle price subticks %+v for clob pair with status %+v",
				order.Subticks,
				currentOraclePriceSubticks,
				clobPair.Status,
			)
		}
		// Throw error if order is a sell and order subticks is less than oracle price subticks
		if !order.IsBuy() && order.Subticks < currentOraclePriceSubticks {
			return sdkerrors.Wrapf(
				types.ErrOrderConflictsWithClobPairStatus,
				"Order subticks %+v must be greater than or equal to oracle price subticks %+v for clob pair with status %+v",
				order.Subticks,
				currentOraclePriceSubticks,
				clobPair.Status,
			)
		}
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

// SetClobPairStatus fetches a ClobPair by id and sets its
// Status property equal to the provided ClobPair_Status. This function returns
// an error if the proposed status transition is not supported.
func (k Keeper) SetClobPairStatus(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	clobPairStatus types.ClobPair_Status,
) error {
	clobPair := k.mustGetClobPair(ctx, clobPairId)

	if !types.IsSupportedClobPairStatusTransition(clobPair.Status, clobPairStatus) {
		return sdkerrors.Wrapf(
			types.ErrInvalidClobPairStatusTransition,
			"Cannot transition from status %+v to status %+v",
			clobPair.Status,
			clobPairStatus,
		)
	}

	clobPair.Status = clobPairStatus
	if err := k.validateClobPair(ctx, &clobPair); err != nil {
		return err
	}

	k.setClobPair(ctx, clobPair)

	return nil
}

// IsPerpetualClobPairInitializing returns true if the ClobPair associated with the provided perpetual id is
// has the initializing status.
func (k Keeper) IsPerpetualClobPairInitializing(
	ctx sdk.Context,
	perpetualId uint32,
) (bool, error) {
	clobPairId, err := k.GetClobPairIdForPerpetual(ctx, perpetualId)
	if err != nil {
		return false, err
	}

	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		return false, sdkerrors.Wrapf(
			types.ErrInvalidClob,
			"GetPerpetualClobPairStatus: did not find clob pair with id = %d",
			clobPairId,
		)
	}

	return clobPair.Status == types.ClobPair_STATUS_INITIALIZING, nil
}
