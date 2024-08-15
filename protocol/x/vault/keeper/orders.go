package keeper

import (
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// RefreshAllVaultOrders refreshes all orders for all vaults by
// 1. Cancelling all existing orders.
// 2. Placing new orders.
func (k Keeper) RefreshAllVaultOrders(ctx sdk.Context) {
	// Iterate through all vaults.
	numActiveVaults := 0
	totalSharesIterator := k.getTotalSharesIterator(ctx)
	defer totalSharesIterator.Close()
	for ; totalSharesIterator.Valid(); totalSharesIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(totalSharesIterator.Key())
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to get vault ID from state key", err)
			continue
		}
		var totalShares types.NumShares
		k.cdc.MustUnmarshal(totalSharesIterator.Value(), &totalShares)

		// Skip if TotalShares is non-positive.
		if totalShares.NumShares.Sign() <= 0 {
			continue
		}

		// Skip if vault has no perpetual positions and strictly less than `activation_threshold_quote_quantums` USDC.
		vault := k.subaccountsKeeper.GetSubaccount(ctx, *vaultId.ToSubaccountId())
		if vault.PerpetualPositions == nil || len(vault.PerpetualPositions) == 0 {
			quotingParams, exists := k.GetVaultQuotingParams(ctx, *vaultId)
			if !exists {
				log.ErrorLog(ctx, "Non-existent vault params", "vaultId", *vaultId)
				continue
			}
			if vault.GetUsdcPosition().Cmp(quotingParams.ActivationThresholdQuoteQuantums.BigInt()) == -1 {
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
				oldOrderPlacement.Order.Subticks != orderToPlace.Subticks {
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
	clobPair, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists || clobPair.Status == clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT {
		return []*clobtypes.Order{}, nil
	}
	perpId := clobPair.Metadata.(*clobtypes.ClobPair_PerpetualClobMetadata).PerpetualClobMetadata.PerpetualId
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpId)
	if err != nil {
		return orders, errorsmod.Wrap(
			err,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	marketParam, exists := k.pricesKeeper.GetMarketParam(ctx, perpetual.Params.MarketId)
	if !exists {
		return orders, errorsmod.Wrap(
			types.ErrMarketParamNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	marketPrice, err := k.pricesKeeper.GetMarketPrice(ctx, perpetual.Params.MarketId)
	if err != nil {
		return orders, errorsmod.Wrap(
			err,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	} else if marketPrice.Price == 0 {
		// Market price can be zero upon market initialization or due to invalid exchange config.
		return orders, errorsmod.Wrap(
			types.ErrZeroMarketPrice,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}

	// Calculate leverage = open notional / equity.
	equity, err := k.GetVaultEquity(ctx, vaultId)
	if err != nil {
		return orders, err
	}
	if equity.Sign() <= 0 {
		return orders, errorsmod.Wrap(
			types.ErrNonPositiveEquity,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	inventory := k.GetVaultInventoryInPerpetual(ctx, vaultId, perpId)
	openNotional := lib.BaseToQuoteQuantums(
		inventory,
		perpetual.Params.AtomicResolution,
		marketPrice.GetPrice(),
		marketPrice.GetExponent(),
	)
	leveragePpm := new(big.Int).Mul(openNotional, lib.BigIntOneMillion())
	leveragePpm.Quo(leveragePpm, equity)

	// Get vault parameters.
	quotingParams, exists := k.GetVaultQuotingParams(ctx, vaultId)
	if !exists {
		return orders, errorsmod.Wrap(
			types.ErrVaultParamsNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
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
	spreadPpm := lib.BigU(lib.Max(
		quotingParams.SpreadMinPpm,
		quotingParams.SpreadBufferPpm+marketParam.MinPriceChangePpm,
	))
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

		return &clobtypes.Order{
			OrderId:      *orderId,
			Side:         side,
			Quantums:     orderSize.Uint64(), // Validated to be a uint64 above.
			Subticks:     subticksRounded,
			GoodTilOneof: goodTilBlockTime,
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

	quotingParams, exists := k.GetVaultQuotingParams(ctx, vaultId)
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

// ReplaceVaultClobOrder replaces a vault CLOB order internal to the protocol and
// emits order replacement indexer event.
func (k Keeper) ReplaceVaultClobOrder(
	ctx sdk.Context,
	vaultId types.VaultId,
	oldOrderId *clobtypes.OrderId,
	newOrder *clobtypes.Order,
) error {
	quotingParams, exists := k.GetVaultQuotingParams(ctx, vaultId)
	if !exists {
		return errorsmod.Wrap(
			types.ErrVaultParamsNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	// Cancel old order.
	err := k.clobKeeper.HandleMsgCancelOrder(ctx, clobtypes.NewMsgCancelOrderStateful(
		*oldOrderId,
		uint32(ctx.BlockTime().Unix())+quotingParams.OrderExpirationSeconds,
	))
	vaultId.IncrCounterWithLabels(
		metrics.VaultCancelOrder,
		metrics.GetLabelForBoolValue(metrics.Success, err == nil),
	)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to cancel order", err, "orderId", oldOrderId, "vaultId", vaultId)
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
