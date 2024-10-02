package keeper

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib/vault"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// RefreshAllVaultOrders refreshes all orders for all vaults by
// 1. Cancelling all existing orders.
// 2. Placing new orders.
func (k Keeper) RefreshAllVaultOrders(ctx sdk.Context) {
	// Iterate through all vaults.
	numActiveVaults := 0
	vaultParamsIterator := k.getVaultParamsIterator(ctx)
	defer vaultParamsIterator.Close()
	for ; vaultParamsIterator.Valid(); vaultParamsIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(vaultParamsIterator.Key())
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get vault ID from state key", err)
			continue
		}
		var vaultParams types.VaultParams
		k.cdc.MustUnmarshal(vaultParamsIterator.Value(), &vaultParams)

		if vaultParams.Status != types.VaultStatus_VAULT_STATUS_QUOTING &&
			vaultParams.Status != types.VaultStatus_VAULT_STATUS_CLOSE_ONLY {
			continue
		}

		// Skip if vault has no perpetual positions and strictly less than `activation_threshold_quote_quantums` USDC.
		if vaultParams.QuotingParams == nil {
			defaultQuotingParams := k.GetDefaultQuotingParams(ctx)
			vaultParams.QuotingParams = &defaultQuotingParams
		}
		vault := k.subaccountsKeeper.GetSubaccount(ctx, *vaultId.ToSubaccountId())
		if len(vault.PerpetualPositions) == 0 {
			if vault.GetUsdcPosition().Cmp(vaultParams.QuotingParams.ActivationThresholdQuoteQuantums.BigInt()) == -1 {
				continue
			}
		}

		// Count current vault as active.
		numActiveVaults++

		// Refresh orders depending on vault type.
		// Currently only supported vault type is CLOB.
		switch vaultId.Type {
		case types.VaultType_VAULT_TYPE_CLOB:
			err := k.RefreshVaultClobOrders(ctx, *vaultId)
			if err != nil {
				log.ErrorLogWithError(ctx, "Failed to refresh vault clob orders", err, "vaultId", *vaultId)
			}
		default:
			log.ErrorLog(ctx, "Failed to refresh vault orders: unknown vault type", "vaultId", *vaultId)
		}
	}

	// Emit metric on number of active vaults.
	metrics.SetGauge(
		metrics.NumActiveVaults,
		float32(numActiveVaults),
	)
}

// RefreshVaultClobOrders refreshes orders of a CLOB vault.
// Note: Client IDs are deterministically constructed based on layer and side. A client ID has its
// last bit flipped only upon order replacement.
func (k Keeper) RefreshVaultClobOrders(ctx sdk.Context, vaultId types.VaultId) (err error) {
	// Get client IDs of most recently placed orders, if any.
	mostRecentClientIds := k.GetMostRecentClientIds(ctx, vaultId)
	// Get orders to place.
	ordersToPlace, err := k.GetVaultClobOrders(ctx, vaultId)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get vault clob orders to place", err, "vaultId", vaultId)
		return err
	}

	clientIds := make([]uint32, len(ordersToPlace))
	for i, orderToPlace := range ordersToPlace {
		if i >= len(mostRecentClientIds) { // when a vault first starts quoting or when `layers` increases.
			// Place order.
			err = k.PlaceVaultClobOrder(ctx, vaultId, orderToPlace)
		} else {
			oldClientId := mostRecentClientIds[i]
			oldOrderId := vaultId.GetClobOrderId(oldClientId)
			oldOrderPlacement, exists := k.clobKeeper.GetLongTermOrderPlacement(ctx, *oldOrderId)
			if !exists { // when order expires / fully fills.
				// Flip client ID because
				// - for an expired order: order expiration event is a block event and order placement
				//   is a tx event. As block events are processed after tx events, indexer will set
				//   new order to expired status if same order ID is used.
				// - for a fully filled order: a fully filled order is added to `RemovedStatefulOrderIds`
				//   in x/clob, which is checked against when placing an order. Order placement fails
				//   if same order ID is used.
				orderToPlace.OrderId.ClientId = oldClientId ^ 1
				err = k.PlaceVaultClobOrder(ctx, vaultId, orderToPlace)
			} else if oldOrderPlacement.Order.Quantums != orderToPlace.Quantums ||
				oldOrderPlacement.Order.Subticks != orderToPlace.Subticks ||
				oldOrderPlacement.Order.Side != orderToPlace.Side {
				// Replace old order with new order.
				// Flip last bit of old client ID to get new client ID to make sure they are different
				// as order placement fails if the same order ID is already marked for cancellation.
				orderToPlace.OrderId.ClientId = oldClientId ^ 1
				err = k.ReplaceVaultClobOrder(ctx, vaultId, oldOrderId, orderToPlace)
			} else {
				// No need to place/replace as existing order is already as desired.
				clientIds[i] = oldClientId
				continue
			}
		}
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to place/replace vault clob order", err, "vaultId", vaultId)
		}
		clientIds[i] = orderToPlace.OrderId.ClientId
	}
	k.SetMostRecentClientIds(ctx, vaultId, clientIds)

	// Cancel any orders that are no longer needed.
	_, quotingParams, exists := k.GetVaultAndQuotingParams(ctx, vaultId)
	if !exists {
		return types.ErrVaultParamsNotFound
	}
	for i := len(ordersToPlace); i < len(mostRecentClientIds); i++ {
		_, err = k.TryToCancelVaultClobOrder(
			ctx,
			vaultId,
			mostRecentClientIds[i],
			quotingParams.OrderExpirationSeconds,
		)
		if err != nil {
			log.ErrorLogWithError(
				ctx,
				"Failed to cancel no longer needed vault clob order",
				err,
				"vaultId",
				vaultId,
			)
		}
	}

	return nil
}

// GetVaultClobOrders returns a list of long term orders for a given CLOB vault.
// Let n be number of layers, then the function returns orders at [a_0, b_0, a_1, b_1, ..., a_{n-1}, b_{n-1}]
// where a_i and b_i are the ask price and bid price at i-th layer. To compute a_i and b_i:
// - a_i = oraclePrice * (1 + ask_spread_i)
// - b_i = oraclePrice * (1 - bid_spread_i)
//
// - ask_spread_i = (1 + skew_i) * spread
// - bid_spread_i = (1 - skew_i) * spread
//
// - skew_i is computed differently based on order side and leverage_i:
//   - ask, leverage_i < 0: (skew_factor * leverage_i - 1)^2 - 1
//   - bid, leverage_i < 0: -skew_factor * leverage_i
//   - ask, leverage_i >= 0: -skew_factor * leverage_i
//   - bid, leverage_i >= 0: -((skew_factor * leverage_i + 1)^2 - 1)
//
// - leverage_i = leverage +/- i * order_size_pct\ (- for ask and + for bid)
// - leverage = open notional / equity
//
// - spread = max(spread_min, spread_buffer + min_price_change)
//
// size of each order is calculated as `order_size_pct * equity / oraclePrice`.
func (k Keeper) GetVaultClobOrders(
	ctx sdk.Context,
	vaultId types.VaultId,
) (orders []*clobtypes.Order, err error) {
	// Get clob pair, perpetual, market parameter, and market price that correspond to this vault.
	clobPair, perpetual, marketParam, marketPrice, err := k.GetVaultClobPerpAndMarket(ctx, vaultId)
	if errors.Is(err, types.ErrClobPairNotFound) || clobPair.Status == clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT {
		return []*clobtypes.Order{}, nil
	} else if err != nil {
		return orders, err
	} else if marketPrice.Price == 0 {
		// Market price can be zero upon market initialization or due to invalid exchange config.
		return orders, errorsmod.Wrap(
			types.ErrZeroMarketPrice,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	// Get vault leverage and equity.
	leverage, equity, err := k.GetVaultLeverageAndEquity(ctx, vaultId, &perpetual, &marketPrice)
	if err != nil {
		return orders, err
	}

	leveragePpm := new(big.Int).Mul(leverage.Num(), lib.BigIntOneMillion())
	leveragePpm = lib.BigDivCeil(leveragePpm, leverage.Denom())

	// Get vault parameters.
	vaultParams, quotingParams, exists := k.GetVaultAndQuotingParams(ctx, vaultId)
	if !exists {
		return orders, errorsmod.Wrap(
			types.ErrVaultParamsNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	// No orders if vault is deactivated, stand-by, or close-only with zero leverage.
	if vaultParams.Status == types.VaultStatus_VAULT_STATUS_DEACTIVATED ||
		vaultParams.Status == types.VaultStatus_VAULT_STATUS_STAND_BY ||
		vaultParams.Status == types.VaultStatus_VAULT_STATUS_CLOSE_ONLY && leverage.Sign() == 0 {
		return []*clobtypes.Order{}, nil
	}

	// Calculate order size (in base quantums).
	orderSizePctPpm := lib.BigU(quotingParams.OrderSizePctPpm)
	orderSize := lib.QuoteToBaseQuantums(
		new(big.Int).Mul(equity, orderSizePctPpm),
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)
	orderSize.Quo(orderSize, lib.BigIntOneMillion())

	// Round (towards-zero) order size to the nearest multiple of step size.
	// Note: below division by StepBaseQuantums is safe as x/clob disallows
	// a clob pair's StepBaseQuantums to be zero.
	stepSize := lib.BigU(clobPair.StepBaseQuantums)
	orderSize.Quo(orderSize, stepSize).Mul(orderSize, stepSize)

	// If order size is zero, return empty orders.
	if orderSize.Sign() == 0 {
		return []*clobtypes.Order{}, nil
	}

	// If order size is not a valid uint64, return error.
	if !orderSize.IsUint64() {
		return []*clobtypes.Order{}, errorsmod.Wrap(
			types.ErrInvalidOrderSize,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	// Calculate spread.
	spreadPpm := lib.BigU(vault.SpreadPpm(&quotingParams, &marketParam))
	// Get oracle price in subticks.
	oracleSubticks := clobtypes.PriceToSubticks(
		marketPrice,
		clobPair,
		perpetual.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	// Get order expiration time.
	goodTilBlockTime := &clobtypes.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + quotingParams.OrderExpirationSeconds,
	}
	skewFactorPpm := lib.BigU(quotingParams.SkewFactorPpm)

	// Construct one ask and one bid for each layer.
	constructOrder := func(
		side clobtypes.Order_Side,
		layer uint32,
		orderId *clobtypes.OrderId,
	) *clobtypes.Order {
		// Ask: leverage_i = leverage - i * order_size_pct
		// Bid: leverage_i = leverage + i * order_size_pct
		leveragePpmI := lib.BigU(layer)
		leveragePpmI.Mul(leveragePpmI, orderSizePctPpm)
		if side == clobtypes.Order_SIDE_SELL {
			leveragePpmI.Neg(leveragePpmI)
		}
		leveragePpmI.Add(leveragePpmI, leveragePpm)

		// Calculate skew.
		skewPpmI := lib.BigMulPpm(leveragePpmI, skewFactorPpm, true)
		if leveragePpm.Sign() < 0 {
			if side == clobtypes.Order_SIDE_SELL {
				// ask when short: skew_i = (skew_factor * leverage_i - 1)^2 - 1
				skewPpmI.Sub(skewPpmI, lib.BigIntOneMillion())
				skewPpmI = lib.BigMulPpm(skewPpmI, skewPpmI, true)
				skewPpmI.Sub(skewPpmI, lib.BigIntOneMillion())
			} else {
				// bid when short: skew_i = -skew_factor * leverage_i
				skewPpmI.Neg(skewPpmI)
			}
		} else {
			if side == clobtypes.Order_SIDE_SELL {
				// ask when long: skew_i = -skew_factor * leverage_i
				skewPpmI.Neg(skewPpmI)
			} else {
				// bid when long: skew_i = -((skew_factor * leverage_i + 1)^2 - 1)
				skewPpmI.Add(skewPpmI, lib.BigIntOneMillion())
				skewPpmI = lib.BigMulPpm(skewPpmI, skewPpmI, true)
				skewPpmI.Sub(skewPpmI, lib.BigIntOneMillion())
				skewPpmI.Neg(skewPpmI)
			}
		}

		// Calculate skewed spread.
		skewedSpreadPpmI := lib.BigIntOneMillion()
		if side == clobtypes.Order_SIDE_SELL {
			skewedSpreadPpmI.Add(skewedSpreadPpmI, skewPpmI)
		} else {
			skewedSpreadPpmI.Sub(skewPpmI, skewedSpreadPpmI)
		}
		// To maintain precision, delay division by 1 million until after multiplying skewed spread with oracle price.
		skewedSpreadPpmI.Mul(skewedSpreadPpmI, spreadPpm)
		skewedSpreadPpmI.Add(skewedSpreadPpmI, lib.BigIntOneTrillion())

		// Determine order subticks.
		subticks := lib.BigMulPpm(
			oracleSubticks.Num(),
			skewedSpreadPpmI,
			side == clobtypes.Order_SIDE_SELL,
		)
		divisor := lib.BigIntOneMillion() // delayed division by 1 million as noted above.
		divisor.Mul(divisor, oracleSubticks.Denom())
		if side == clobtypes.Order_SIDE_SELL {
			subticks = lib.BigDivCeil(subticks, divisor)
		} else {
			subticks = new(big.Int).Quo(subticks, divisor)
		}

		// Bound subticks between the minimum and maximum subticks.
		// Note: below division by SubticksPerTick is safe as x/clob disallows
		// a clob pair's SubticksPerTick to be zero.
		subticksPerTick := lib.BigU(clobPair.SubticksPerTick)
		subticks = lib.BigIntRoundToMultiple(
			subticks,
			subticksPerTick,
			side == clobtypes.Order_SIDE_SELL,
		)

		minSubticks := uint64(clobPair.SubticksPerTick)
		maxSubticks := uint64(math.MaxUint64 - (uint64(math.MaxUint64) % uint64(clobPair.SubticksPerTick)))
		subticksRounded := lib.BigUint64Clamp(
			subticks,
			minSubticks,
			maxSubticks,
		)

		// If the side would increase the vault's inventory, make the order post-only.
		timeInForceType := clobtypes.Order_TIME_IN_FORCE_UNSPECIFIED
		if (side == clobtypes.Order_SIDE_SELL && leveragePpm.Sign() <= 0) ||
			(side == clobtypes.Order_SIDE_BUY && leveragePpm.Sign() >= 0) {
			timeInForceType = clobtypes.Order_TIME_IN_FORCE_POST_ONLY
		}

		return &clobtypes.Order{
			OrderId:      *orderId,
			Side:         side,
			Quantums:     orderSize.Uint64(), // Validated to be a uint64 above.
			Subticks:     subticksRounded,
			GoodTilOneof: goodTilBlockTime,
			TimeInForce:  timeInForceType,
		}
	}

	orderIds := k.GetVaultClobOrderIds(ctx, vaultId)
	orders = make([]*clobtypes.Order, 2*quotingParams.Layers)
	if len(orders) != len(orderIds) { // sanity check
		return orders, errorsmod.Wrap(
			types.ErrOrdersAndOrderIdsDiffLen,
			fmt.Sprintf(
				"VaultId: %v, len(Orders):%d, len(OrderIds):%d",
				vaultId,
				len(orders),
				len(orderIds),
			),
		)
	}
	for i := uint32(0); i < quotingParams.Layers; i++ {
		// Construct ask at this layer.
		orders[2*i] = constructOrder(clobtypes.Order_SIDE_SELL, i, orderIds[2*i])

		// Construct bid at this layer.
		orders[2*i+1] = constructOrder(clobtypes.Order_SIDE_BUY, i, orderIds[2*i+1])
	}

	if vaultParams.Status == types.VaultStatus_VAULT_STATUS_CLOSE_ONLY {
		// In close-only mode with non-zero leverage.
		reduceOnlyTotalOrderSize := k.GetVaultInventoryInPerpetual(ctx, vaultId, perpetual.Params.Id)
		stepSize := lib.BigU(clobPair.StepBaseQuantums)
		reduceOnlyTotalOrderSize.Quo(reduceOnlyTotalOrderSize, stepSize)
		reduceOnlyTotalOrderSize.Mul(reduceOnlyTotalOrderSize, stepSize)
		if reduceOnlyTotalOrderSize.Sign() == 0 {
			return []*clobtypes.Order{}, nil
		}

		// If vault is long, only need sell orders.
		reduceOnlySide := clobtypes.Order_SIDE_SELL
		if leverage.Sign() < 0 {
			// If vault is short, only need buy orders.
			reduceOnlySide = clobtypes.Order_SIDE_BUY
		}
		reduceOnlyOrders := make([]*clobtypes.Order, 0, len(orders))
		totalOrderSize := reduceOnlyTotalOrderSize.Uint64()
		for _, order := range orders {
			if order.Side == reduceOnlySide {
				if totalOrderSize == 0 {
					break
				}

				order.Quantums = lib.Min(order.Quantums, totalOrderSize)
				totalOrderSize -= order.Quantums
				reduceOnlyOrders = append(reduceOnlyOrders, order)
			}
		}
		return reduceOnlyOrders, nil
	}

	return orders, nil
}

// GetVaultClobOrderIds returns a list of order IDs for a given CLOB vault.
// Let n be number of layers, then the function returns order IDs
// [a_0, b_0, a_1, b_1, ..., a_{n-1}, b_{n-1}] where a_i and b_i are respectively
// ask and bid order IDs at the i-th layer.
func (k Keeper) GetVaultClobOrderIds(
	ctx sdk.Context,
	vaultId types.VaultId,
) (orderIds []*clobtypes.OrderId) {
	vault := vaultId.ToSubaccountId()
	constructOrderId := func(
		side clobtypes.Order_Side,
		layer uint32,
	) *clobtypes.OrderId {
		return &clobtypes.OrderId{
			SubaccountId: *vault,
			ClientId:     types.GetVaultClobOrderClientId(side, uint8(layer)),
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   clobtypes.ClobPairId(vaultId.Number).ToUint32(),
		}
	}

	_, quotingParams, exists := k.GetVaultAndQuotingParams(ctx, vaultId)
	if !exists {
		return []*clobtypes.OrderId{}
	}
	layers := quotingParams.Layers
	orderIds = make([]*clobtypes.OrderId, 2*layers)
	for i := uint32(0); i < layers; i++ {
		// Construct ask order ID at this layer.
		orderIds[2*i] = constructOrderId(clobtypes.Order_SIDE_SELL, i)

		// Construct bid order ID at this layer.
		orderIds[2*i+1] = constructOrderId(clobtypes.Order_SIDE_BUY, i)
	}

	return orderIds
}

// PlaceVaultClobOrder places a vault CLOB order internal to the protocol, skipping various
// logs, metrics, and validations
func (k Keeper) PlaceVaultClobOrder(
	ctx sdk.Context,
	vaultId types.VaultId,
	order *clobtypes.Order,
) error {
	err := k.clobKeeper.HandleMsgPlaceOrder(ctx, clobtypes.NewMsgPlaceOrder(*order), true)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to place order", err, "order", order, "vaultId", vaultId)
	}
	vaultId.IncrCounterWithLabels(
		metrics.VaultPlaceOrder,
		metrics.GetLabelForBoolValue(metrics.Success, err == nil),
	)
	return err
}

// CancelVaultClobOrder cancels a vault CLOB order.
func (k Keeper) CancelVaultClobOrder(
	ctx sdk.Context,
	vaultId types.VaultId,
	orderId *clobtypes.OrderId,
	orderExpirationSeconds uint32,
) error {
	err := k.clobKeeper.HandleMsgCancelOrder(ctx, clobtypes.NewMsgCancelOrderStateful(
		*orderId,
		uint32(ctx.BlockTime().Unix())+orderExpirationSeconds,
	))
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to cancel order", err, "orderId", orderId, "vaultId", vaultId)
	}
	vaultId.IncrCounterWithLabels(
		metrics.VaultCancelOrder,
		metrics.GetLabelForBoolValue(metrics.Success, err == nil),
	)
	return err
}

// TryToCancelVaultClobOrder tries to cancel a vault CLOB order. Returns whether the order exists
// and whether cancellation errors.
func (k Keeper) TryToCancelVaultClobOrder(
	ctx sdk.Context,
	vaultId types.VaultId,
	clientId uint32,
	orderExpirationSeconds uint32,
) (
	orderExists bool,
	err error,
) {
	orderId := vaultId.GetClobOrderId(clientId)
	_, exists := k.clobKeeper.GetLongTermOrderPlacement(ctx, *orderId)
	if exists {
		err = k.CancelVaultClobOrder(ctx, vaultId, orderId, orderExpirationSeconds)
		return true, err
	}
	return false, nil
}

// ReplaceVaultClobOrder replaces a vault CLOB order internal to the protocol and
// emits order replacement indexer event.
func (k Keeper) ReplaceVaultClobOrder(
	ctx sdk.Context,
	vaultId types.VaultId,
	oldOrderId *clobtypes.OrderId,
	newOrder *clobtypes.Order,
) error {
	_, quotingParams, exists := k.GetVaultAndQuotingParams(ctx, vaultId)
	if !exists {
		return errorsmod.Wrap(
			types.ErrVaultParamsNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	// Cancel old order.
	err := k.CancelVaultClobOrder(ctx, vaultId, oldOrderId, quotingParams.OrderExpirationSeconds)
	if err != nil {
		return err
	}

	// Place new order.
	err = k.PlaceVaultClobOrder(ctx, vaultId, newOrder)
	return err
}

// GetMostRecentClientIds returns the most recent client IDs for a vault.
func (k Keeper) GetMostRecentClientIds(
	ctx sdk.Context,
	vaultId types.VaultId,
) []uint32 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MostRecentClientIdsKeyPrefix))
	bytes := store.Get(vaultId.ToStateKey())
	if bytes == nil {
		return []uint32{}
	}
	return lib.BytesToUint32Array(bytes)
}

// SetMostRecentClientIds sets the most recent client IDs for a vault.
func (k Keeper) SetMostRecentClientIds(
	ctx sdk.Context,
	vaultId types.VaultId,
	clientIds []uint32,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MostRecentClientIdsKeyPrefix))
	store.Set(vaultId.ToStateKey(), lib.Uint32ArrayToBytes(clientIds))
}
