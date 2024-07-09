package keeper

import (
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
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
	params := k.GetParams(ctx)
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
			if vault.GetUsdcPosition().Cmp(params.ActivationThresholdQuoteQuantums.BigInt()) == -1 {
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
func (k Keeper) RefreshVaultClobOrders(ctx sdk.Context, vaultId types.VaultId) (err error) {
	// Cancel CLOB orders from last block.
	orderIdsToCancel, err := k.GetVaultClobOrderIds(
		ctx.WithBlockHeight(ctx.BlockHeight()-1),
		vaultId,
	)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get vault clob order IDs to cancel", err, "vaultId", vaultId)
		return err
	}
	orderExpirationSeconds := k.GetParams(ctx).OrderExpirationSeconds
	for _, orderId := range orderIdsToCancel {
		if _, exists := k.clobKeeper.GetLongTermOrderPlacement(ctx, *orderId); exists {
			err := k.clobKeeper.HandleMsgCancelOrder(ctx, clobtypes.NewMsgCancelOrderStateful(
				*orderId,
				uint32(ctx.BlockTime().Unix())+orderExpirationSeconds,
			), true)
			if err != nil {
				log.ErrorLogWithError(ctx, "Failed to cancel order", err, "orderId", orderId, "vaultId", vaultId)
			}
			vaultId.IncrCounterWithLabels(
				metrics.VaultCancelOrder,
				metrics.GetLabelForBoolValue(metrics.Success, err == nil),
			)
		}
	}
	// Place new CLOB orders.
	ordersToPlace, err := k.GetVaultClobOrders(ctx, vaultId)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get vault clob orders to place", err, "vaultId", vaultId)
		return err
	}

	for i, order := range ordersToPlace {
		err := k.PlaceVaultClobOrder(ctx, order)
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to place order", err, "order", order, "vaultId", vaultId)
		}

		vaultId.IncrCounterWithLabels(
			metrics.VaultPlaceOrder,
			metrics.GetLabelForBoolValue(metrics.Success, err == nil),
		)

		// Send indexer messages. We expect orderIdsToCancel and ordersToPlace to have the same length
		// and the order to place at each index to be a replacement of the order to cancel at the same index.
		replacedOrderId := orderIdsToCancel[i]
		if replacedOrderId == nil {
			k.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewLongTermOrderPlacementEvent(
						*order,
					),
				),
			)
		} else {
			k.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewLongTermOrderReplacementEvent(
						*replacedOrderId,
						*order,
					),
				),
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
	clobPair, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists {
		return orders, errorsmod.Wrap(
			types.ErrClobPairNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
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

	// Get parameters.
	params := k.GetParams(ctx)

	// Calculate order size (in base quantums).
	orderSizePctPpm := lib.BigU(params.OrderSizePctPpm)
	orderSize := lib.QuoteToBaseQuantums(
		new(big.Int).Mul(equity, orderSizePctPpm),
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)
	orderSize.Quo(orderSize, lib.BigIntOneMillion())

	// Round (towards-zero) order size to the nearest multiple of step size.
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
		params.SpreadMinPpm,
		params.SpreadBufferPpm+marketParam.MinPriceChangePpm,
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
		GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + params.OrderExpirationSeconds,
	}
	skewFactorPpm := lib.BigU(params.SkewFactorPpm)

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

	orderIds, err := k.GetVaultClobOrderIds(ctx, vaultId)
	if err != nil {
		return orders, err
	}
	orders = make([]*clobtypes.Order, 2*params.Layers)
	for i := uint32(0); i < params.Layers; i++ {
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
) (orderIds []*clobtypes.OrderId, err error) {
	clobPair, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists {
		return orderIds, errorsmod.Wrap(
			types.ErrClobPairNotFound,
			fmt.Sprintf("VaultId: %v", vaultId),
		)
	}
	vault := vaultId.ToSubaccountId()
	constructOrderId := func(
		side clobtypes.Order_Side,
		layer uint32,
	) *clobtypes.OrderId {
		return &clobtypes.OrderId{
			SubaccountId: *vault,
			ClientId:     k.GetVaultClobOrderClientId(ctx, side, uint8(layer)),
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   clobPair.Id,
		}
	}

	layers := k.GetParams(ctx).Layers
	orderIds = make([]*clobtypes.OrderId, 2*layers)
	for i := uint32(0); i < layers; i++ {
		// Construct ask order ID at this layer.
		orderIds[2*i] = constructOrderId(clobtypes.Order_SIDE_SELL, i)

		// Construct bid order ID at this layer.
		orderIds[2*i+1] = constructOrderId(clobtypes.Order_SIDE_BUY, i)
	}

	return orderIds, nil
}

// GetVaultClobOrderClientId returns the client ID for a CLOB order where
// - 1st bit is `side-1` (subtract 1 as buy_side = 1, sell_side = 2)
//
// - 2nd bit is `block height % 2`
//   - block height bit alternates between 0 and 1 to ensure that client IDs
//     are different in two consecutive blocks (otherwise, order placement would
//     fail because the same order IDs are already marked for cancellation)
//
// - next 8 bits are `layer`
func (k Keeper) GetVaultClobOrderClientId(
	ctx sdk.Context,
	side clobtypes.Order_Side,
	layer uint8,
) uint32 {
	sideBit := uint32(side - 1)
	sideBit <<= 31

	blockHeightBit := uint32(ctx.BlockHeight() % 2)
	blockHeightBit <<= 30

	layerBits := uint32(layer) << 22

	return sideBit | blockHeightBit | layerBits
}

// PlaceVaultClobOrder places a vault CLOB order as an order internal to the protocol,
// skipping various logs, metrics, and validations.
func (k Keeper) PlaceVaultClobOrder(
	ctx sdk.Context,
	order *clobtypes.Order,
) error {
	// Place an internal clob order.
	return k.clobKeeper.HandleMsgPlaceOrder(ctx, clobtypes.NewMsgPlaceOrder(*order), true)
}
