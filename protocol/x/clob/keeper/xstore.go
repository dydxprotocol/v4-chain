package keeper

import (
	"fmt"
	"math"
	"math/big"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func (k Keeper) GetXTimestampStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.XTimestampKeyPrefix),
	)
}
func (k Keeper) GetXOrderStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.XOrderKeyPrefix),
	)
}
func (k Keeper) GetXStopStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.XStopKeyPrefix),
	)
}
func (k Keeper) GetXRestingStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.XRestingKeyPrefix),
	)
}
func (k Keeper) GetXExpiryStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.XExpiryKeyPrefix),
	)
}

func (k Keeper) TriggerOrder(
	ctx sdk.Context,
	uidBytes []byte,
) (
	sizeSum uint64, // filled size
	sizeRem uint64, // remaining size
	err error,
) {
	// Get the order.
	orderStore := k.GetXOrderStore(ctx)
	var order types.XOrder
	k.cdc.MustUnmarshal(orderStore.Get(uidBytes), &order)
	priority := order.GetPriority()

	// Check if the order is triggerable.
	if !order.Base.IsStop() {
		return 0, 0, fmt.Errorf("order is not triggerable")
	}

	// Remove the order from the stop store.
	stopStore := k.GetXStopStore(ctx)
	stopKey := order.ToStopKey(priority)
	if !stopStore.Has(stopKey) {
		return 0, 0, fmt.Errorf("order not found in stop store")
	}
	stopStore.Delete(stopKey)

	// Remove the order from the expiry store.
	expiryKey := order.ToExpiryKeyOrNil(priority)
	if expiryKey != nil {
		expiryStore := k.GetXExpiryStore(ctx)
		if !expiryStore.Has(expiryKey) {
			panic(fmt.Errorf("expiry does not exist with key %s", expiryKey))
		}
		expiryStore.Delete(expiryKey)
	}

	// Remove the order from the order store.
	orderStore.Delete(uidBytes)

	// Process the order.
	return k.ProcessLiveOrder(ctx, order)
}

func (k Keeper) GetOrderById(
	ctx sdk.Context,
	uidBytes []byte,
) (
	order types.XOrder,
	found bool,
) {
	orderStore := k.GetXOrderStore(ctx)
	orderBytes := orderStore.Get(uidBytes)
	if len(orderBytes) == 0 {
		return types.XOrder{}, false
	}
	k.cdc.MustUnmarshal(orderBytes, &order)
	return order, true
}

func (k Keeper) ProcessOrder(
	ctx sdk.Context,
	order types.XOrder,
	placeFlags uint32,
) (
	sizeSum uint64, // filled size
	sizeRem uint64, // remaining size
	err error,
) {
	// Get existing order.
	orderStore := k.GetXOrderStore(ctx)
	uidBytes := order.Uid.ToBytes()
	orderExists := orderStore.Has(uidBytes)

	// Replacement Flags
	replaceFlags := types.GetReplaceFlags(placeFlags)
	switch replaceFlags {
	// New-only orders: Return error if order already exists.
	case types.Order_REPLACE_FLAGS_NEW_ONLY:
		if orderExists {
			return 0, 0, fmt.Errorf("order already exists with key %s", uidBytes)
		}
	// Upsert orders: Remove existing order if it exists.
	case types.Order_REPLACE_FLAGS_UPSERT:
		if orderExists {
			k.RemoveOrderById(ctx, uidBytes)
		}
	// Incremental size orders: If order exists, remove it from state.
	// Place the new order at the incremented size, unless the size overflows, in which case use the size of the new order.
	case types.Order_REPLACE_FLAGS_INC_SIZE:
		if orderExists {
			existingOrder, _ := k.GetOrderById(ctx, uidBytes)
			order.Base.Quantums = lib.Max(order.Base.Quantums, order.Base.Quantums+existingOrder.Base.Quantums)
			k.RemoveOrderById(ctx, uidBytes)
		}
	// Decremental size orders: If order exists, remove it from state.
	// Place the new order at the decremented size, unless the size is zero.
	case types.Order_REPLACE_FLAGS_DEC_SIZE:
		existingOrder, found := k.GetOrderById(ctx, uidBytes)
		if !found {
			return 0, 0, fmt.Errorf("dec size: order does not exist with key %s", uidBytes)
		}
		if existingOrder.Base.Quantums <= order.Base.Quantums {
			k.RemoveOrderById(ctx, uidBytes)
			return 0, 0, fmt.Errorf("dec size: size decremented to zero")
		}
		order.Base.Quantums = existingOrder.Base.Quantums - order.Base.Quantums
		k.RemoveOrderById(ctx, uidBytes)
	default:
		panic(fmt.Errorf("invalid place flags: %d", placeFlags))
	}

	// Stop orders: Put it in state and let it trigger later. Return.
	if order.Base.IsStop() {
		err := k.AddOrderToState(ctx, order, true)
		if err != nil {
			return 0, 0, err
		}
		return 0, order.Base.Quantums, nil
	}

	// All other orders are live orders.
	return k.ProcessLiveOrder(ctx, order)
}

func (k Keeper) ProcessLiveOrder(
	ctx sdk.Context,
	order types.XOrder,
) (
	sizeSum uint64, // filled size
	sizeRem uint64, // remaining size
	err error,
) {
	// A live order is either:
	// - not a stop order
	// - a stop order that has been triggered
	sizeSum = 0
	sizeRem = order.Base.Quantums

	// Resize order if reduce-only.
	if order.Base.IsReduceOnly() {
		subaccountId, err := k.GetSubaccountIdFromSID(ctx, order.Uid.Sid)
		if err != nil {
			return 0, 0, err
		}
		sa := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
		clobPair, found := k.GetClobPair(ctx, types.ClobPairId(order.Uid.Iid.ClobId))
		if !found {
			return 0, 0, fmt.Errorf("clob pair not found")
		}
		perpetualId := clobPair.MustGetPerpetualId()
		pos, posExists := sa.GetPerpetualPositionForId(perpetualId)
		if !posExists ||
			pos.Quantums.Sign() == 0 ||
			(pos.Quantums.Sign() > 0) == (order.Base.Side() == types.Order_SIDE_BUY) {
			return 0, 0, fmt.Errorf("order is reduce-only and position is zero or wrong side")
		}
		absSize := new(big.Int).Abs(pos.Quantums.BigInt())
		if absSize.Cmp(lib.BigU(order.Base.Quantums)) < 0 {
			if !absSize.IsUint64() {
				panic(fmt.Errorf("abs size is not uint64"))
			}
			sizeRem = absSize.Uint64()
		}
	}

	// Start matching the order against the resting orders.
	// Iterate only from the best possible price to the price of the order.
	orderStore := k.GetXOrderStore(ctx)
	isBidBook := order.Base.Side() != types.Order_SIDE_BUY // isBidBook is true for sell orders.
	itStartKey := types.ToRestingKey(
		order.Uid.Iid.ClobId,
		isBidBook,
		getStartSubticks(isBidBook),
		0,
	)
	itEndKey := types.ToRestingKey(
		order.Uid.Iid.ClobId,
		isBidBook,
		order.Base.Subticks,
		math.MaxUint64,
	)
	it := k.GetXRestingStore(ctx).Iterator(itStartKey, itEndKey)
	defer it.Close()
	for ; it.Valid() && sizeRem > 0; it.Next() {
		// Get resting order.
		makerUidBytes := it.Value()
		makerOrder, found := k.GetOrderById(ctx, makerUidBytes)
		if !found {
			panic(fmt.Errorf("maker order not found during iteration"))
		}

		// Self-Trade Prevention.
		if makerOrder.Uid.Sid == order.Uid.Sid {
			stpType := order.Base.GetStp()
			if stpType == types.Order_STP_EXPIRE_MAKER {
				k.RemoveOrderById(ctx, makerUidBytes)
				continue
			} else if stpType == types.Order_STP_EXPIRE_TAKER {
				return sizeSum, 0, fmt.Errorf("stp is expire taker")
			} else if stpType == types.Order_STP_EXPIRE_BOTH {
				k.RemoveOrderById(ctx, makerUidBytes)
				return sizeSum, 0, fmt.Errorf("stp is expire both")
			}
			panic(fmt.Errorf("invalid stp type: %d", stpType))
		}

		// Return early if matches and post-only.
		if order.Base.GetTif() == types.Order_TIME_IN_FORCE_POST_ONLY {
			return sizeSum, 0, fmt.Errorf("order is post-only and matches resting order")
		}

		// Get the fill size.
		fillSize := makerOrder.Base.Quantums
		if fillSize > sizeRem {
			fillSize = sizeRem
		}

		// Attempt to update subaccounts.
		updates, err := k.GetSubaccountUpdatesForMatch(ctx, order, makerOrder, fillSize)
		if err != nil {
			return sizeSum, 0, err
		}
		_, results, err := k.subaccountsKeeper.UpdateSubaccounts(ctx, updates, satypes.Match)
		if err != nil {
			return sizeSum, 0, err
		}
		takerResult := results[0]
		makerResult := results[1]
		if !takerResult.IsSuccess() {
			return sizeSum, 0, nil // TODO: should an error be returned?
		}
		if !makerResult.IsSuccess() {
			k.RemoveOrderById(ctx, makerUidBytes)
			continue
		}

		// Update resting order.
		// If order is fully-filled, remove it from state.
		// If order is partially-filled, update it in state.
		makerOrder.Base.Quantums -= fillSize
		if makerOrder.Base.Quantums == 0 {
			k.RemoveOrderById(ctx, makerUidBytes)
		} else {
			orderStore.Set(makerOrder.Uid.ToBytes(), k.cdc.MustMarshal(&makerOrder))
		}

		// Cancel Maker OCO Order.
		if makerOrder.Base.HasOcoClientId() {
			ocoOrderId := makerOrder.Uid // Copy by value.
			ocoOrderId.Iid.ClientId = makerOrder.Base.MustGetOcoClientId()
			if ocoOrderId == makerOrder.Uid {
				panic(fmt.Errorf("oco order id is the same as the maker order id"))
			}
			k.RemoveOrderById(ctx, ocoOrderId.ToBytes())
		}

		// Cancel Taker OCO Order. Optimization: Only attempt to cancel on the first fill.
		if sizeSum == 0 && order.Base.HasOcoClientId() {
			ocoOrderId := order.Uid // Copy by value.
			ocoOrderId.Iid.ClientId = order.Base.MustGetOcoClientId()
			if ocoOrderId == order.Uid {
				panic(fmt.Errorf("oco order id is the same as the taker order id"))
			}
			k.RemoveOrderById(ctx, ocoOrderId.ToBytes())
		}

		// Update counters.
		sizeRem -= fillSize
		sizeSum += fillSize
	}

	// If the order is IOC, there is no remaining size.
	if order.Base.GetTif() == types.Order_TIME_IN_FORCE_IOC {
		sizeRem = 0
	}

	// If there is remaining size, attempt to place the order as a resting order.
	if sizeRem > 0 {
		order.Base.Quantums = sizeRem
		err := k.AddOrderToState(ctx, order, false)
		if err != nil {
			return sizeSum, 0, err
		}
	}

	return sizeSum, sizeRem, nil
}

func (k Keeper) RemoveOrderById(
	ctx sdk.Context,
	uidBytes []byte,
) (found bool) {
	// Remove order by id.
	orderStore := k.GetXOrderStore(ctx)
	orderBytes := orderStore.Get(uidBytes)
	if len(orderBytes) == 0 {
		return false
	}
	var order types.XOrder
	k.cdc.MustUnmarshal(orderBytes, &order)
	orderStore.Delete(uidBytes)

	// Precompute value which is used multiple times.
	priority := order.GetPriority()

	// Remove from Resting store.
	restingStore := k.GetXRestingStore(ctx)
	restingKey := order.ToRestingKey(priority)
	if restingStore.Has(restingKey) {
		restingStore.Delete(restingKey)
	}

	// Remove from Stop store if triggerable.
	if order.Base.IsStop() {
		stopStore := k.GetXStopStore(ctx)
		stopKey := order.ToStopKey(priority)
		if stopStore.Has(stopKey) {
			stopStore.Delete(stopKey)
		}
	}

	// Remove expiry.
	expiryKey := order.ToExpiryKeyOrNil(priority)
	if expiryKey != nil {
		expiryStore := k.GetXExpiryStore(ctx)
		if !expiryStore.Has(expiryKey) {
			panic(fmt.Errorf("expiry does not exist with key %s", expiryKey))
		}
		expiryStore.Delete(expiryKey)
	}

	return true
}

func (k Keeper) AddOrderToState(
	ctx sdk.Context,
	order types.XOrder,
	asUntriggeredStop bool,
) error {
	if asUntriggeredStop && !order.Base.IsStop() {
		panic(fmt.Errorf("order is not a stop order"))
	}

	// Precompute these values which may be used multiple times.

	priority := order.GetPriority()
	uidBytes := order.Uid.ToBytes()

	// Ensure key does not already exist in the order store or the expiry store (if applicable).
	orderStore := k.GetXOrderStore(ctx)
	orderKey := uidBytes
	if orderStore.Has(orderKey) {
		return fmt.Errorf("order already exists with key %s", orderKey)
	}
	var expiryStore prefix.Store
	expiryKeyOrNil := order.ToExpiryKeyOrNil(priority)
	if expiryKeyOrNil != nil {
		expiryStore = k.GetXExpiryStore(ctx)
		if expiryStore.Has(expiryKeyOrNil) {
			return fmt.Errorf("order already exists with key %s", expiryKeyOrNil)
		}
	}

	// Ensure key does not exist in the stop/resting store (as applicable).
	// If it does not exist in this store, then start adding it to state.
	if asUntriggeredStop {
		stopStore := k.GetXStopStore(ctx)
		stopKey := order.ToStopKey(priority)
		if stopStore.Has(stopKey) {
			return fmt.Errorf("order already exists with key %s", stopKey)
		}
		stopStore.Set(stopKey, uidBytes)
	} else {
		restingStore := k.GetXRestingStore(ctx)
		restingKey := order.ToRestingKey(priority)
		if restingStore.Has(restingKey) {
			return fmt.Errorf("order already exists with key %s", restingKey)
		}
		restingStore.Set(restingKey, uidBytes)
	}

	// Finish add the order to state.
	if expiryKeyOrNil != nil {
		expiryStore.Set(expiryKeyOrNil, uidBytes)
	}
	orderStore.Set(orderKey, k.cdc.MustMarshal(&order))

	return nil
}

func (k Keeper) CancelAllOrders(
	ctx sdk.Context,
	sid uint64,
	clobId uint32,
) {
	// Iterate over the order store for all orders for the given sid/clobId. itEnd is exclusive.
	store := k.GetXOrderStore(ctx)
	itStart := types.XUID{Sid: sid, Iid: types.XIID{ClobId: clobId, ClientId: 0}}
	itEnd := types.XUID{Sid: sid, Iid: types.XIID{ClobId: clobId + 1, ClientId: 0}}
	it := store.Iterator(itStart.ToBytes(), itEnd.ToBytes())
	defer it.Close()
	for ; it.Valid(); it.Next() {
		k.RemoveOrderById(ctx, it.Key())
	}
}

func (k Keeper) ExpireXOrders(ctx sdk.Context) {
	// Iterate over the expiry store til the next second, exclusive.
	store := k.GetXExpiryStore(ctx)
	iterEnd := lib.Uint32ToBytes(uint32(ctx.BlockTime().Unix() + 1))
	it := store.Iterator(nil, iterEnd)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		k.RemoveOrderById(ctx, it.Value())
	}
}

func (k Keeper) TriggerXOrders(ctx sdk.Context) error {
	// Iterate over the stop store for each clob pair.
	store := k.GetXStopStore(ctx)
	for _, clobPair := range k.GetAllClobPairs(ctx) {
		// Get the last traded prices in the block.
		minPrice, maxPrice, found := k.GetTradePricesForPerpetual(ctx, clobPair.MustGetPerpetualId())
		if !found {
			continue
		}

		// Iterate over the traded prices upwards.
		itaEnd := types.ToStopKey(clobPair.Id, false, uint64(maxPrice)+1, 0)
		ita := store.Iterator(nil, itaEnd)
		defer ita.Close()
		for ; ita.Valid(); ita.Next() {
			_, _, err := k.TriggerOrder(ctx, ita.Value())
			if err != nil {
				return err
			}
		}

		// Iterate over the traded prices downwards.
		itbEnd := types.ToStopKey(clobPair.Id, true, uint64(minPrice)-1, 0)
		itb := store.Iterator(nil, itbEnd)
		defer itb.Close()
		for ; itb.Valid(); itb.Next() {
			_, _, err := k.TriggerOrder(ctx, itb.Value())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Stateful Helper Functions

func (k Keeper) GetSubaccountIdFromSID(
	ctx sdk.Context,
	sid uint64,
) (
	satypes.SubaccountId,
	error,
) {
	accountNumber, subaccountNumber := types.SplitSID(sid)
	address, err := k.accountKeeper.GetAccountByNumber(ctx, accountNumber)
	if err != nil {
		return satypes.SubaccountId{}, err
	}
	return satypes.SubaccountId{
		Owner:  address.String(),
		Number: subaccountNumber,
	}, nil
}

func (k Keeper) GetSubaccountUpdatesForMatch(
	ctx sdk.Context,
	takerOrder types.XOrder,
	makerOrder types.XOrder,
	fillAmount uint64,
) (
	[]satypes.Update,
	error,
) {
	clobPair, ok := k.GetClobPair(ctx, types.ClobPairId(takerOrder.Uid.Iid.ClobId))
	if !ok {
		return nil, fmt.Errorf("clob pair not found")
	}
	perpetualId := clobPair.MustGetPerpetualId()

	// Get the subaccount IDs.
	makerSubaccountId, err := k.GetSubaccountIdFromSID(ctx, makerOrder.Uid.Sid)
	if err != nil {
		return nil, err
	}
	takerSubaccountId, err := k.GetSubaccountIdFromSID(ctx, takerOrder.Uid.Sid)
	if err != nil {
		return nil, err
	}

	// Get the fill quote quantums.
	fillQuoteQuantums := types.FillAmountToQuoteQuantums(
		types.Subticks(makerOrder.Base.Subticks),
		satypes.BaseQuantums(fillAmount),
		clobPair.QuantumConversionExponent,
	)

	// Taker fees and maker fees/rebates are rounded towards positive infinity.
	makerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(ctx, makerSubaccountId.Owner, false)
	takerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(ctx, takerSubaccountId.Owner, true)
	takerFee := lib.BigMulPpm(fillQuoteQuantums, lib.BigI(takerFeePpm), true)
	makerFee := lib.BigMulPpm(fillQuoteQuantums, lib.BigI(makerFeePpm), true)

	// Determine the balance deltas.
	takerQuoteDelta := new(big.Int).Set(fillQuoteQuantums)
	makerQuoteDelta := new(big.Int).Set(fillQuoteQuantums)
	takerPerpDelta := lib.BigU(fillAmount)
	makerPerpDelta := lib.BigU(fillAmount)
	if takerOrder.Base.Side() == types.Order_SIDE_BUY {
		takerQuoteDelta.Neg(takerQuoteDelta)
		makerPerpDelta.Neg(makerPerpDelta)
	} else {
		makerQuoteDelta.Neg(makerQuoteDelta)
		takerPerpDelta.Neg(takerPerpDelta)
	}
	takerQuoteDelta.Sub(takerQuoteDelta, takerFee)
	makerQuoteDelta.Sub(makerQuoteDelta, makerFee)

	// Create the subaccount update.
	updates := []satypes.Update{
		{
			SubaccountId: takerSubaccountId,
			AssetUpdates: []satypes.AssetUpdate{{
				AssetId:          assettypes.AssetUsdc.Id,
				BigQuantumsDelta: takerQuoteDelta,
			}},
			PerpetualUpdates: []satypes.PerpetualUpdate{{
				PerpetualId:      perpetualId,
				BigQuantumsDelta: takerPerpDelta,
			}},
		},
		{
			SubaccountId: makerSubaccountId,
			AssetUpdates: []satypes.AssetUpdate{{
				AssetId:          assettypes.AssetUsdc.Id,
				BigQuantumsDelta: makerQuoteDelta,
			}},
			PerpetualUpdates: []satypes.PerpetualUpdate{{
				PerpetualId:      perpetualId,
				BigQuantumsDelta: makerPerpDelta,
			}},
		},
	}

	return updates, nil
}

// Helper functions

func getStartSubticks(isBidBook bool) uint64 {
	if isBidBook {
		return math.MaxUint64
	}
	return 0
}
