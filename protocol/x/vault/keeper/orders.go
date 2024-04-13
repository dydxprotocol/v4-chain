package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
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
		if totalShares.NumShares.BigInt().Sign() <= 0 {
			continue
		}

		// Skip if vault has no perpetual positions and strictly less than `activation_threshold_quote_quantums` USDC.
		vault := k.subaccountsKeeper.GetSubaccount(ctx, *vaultId.ToSubaccountId())
		if vault.PerpetualPositions == nil || len(vault.PerpetualPositions) == 0 {
			if vault.GetUsdcPosition().Cmp(params.ActivationThresholdQuoteQuantums.BigInt()) == -1 {
				continue
			}
		}

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
}

// RefreshVaultClobOrders refreshes orders of a CLOB vault.
func (k Keeper) RefreshVaultClobOrders(ctx sdk.Context, vaultId types.VaultId) (err error) {
	// Cancel CLOB orders from last block.
	ordersToCancel, err := k.GetVaultClobOrders(
		ctx.WithBlockHeight(ctx.BlockHeight()-1),
		vaultId,
	)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get vault clob orders to cancel", err, "vaultId", vaultId)
		return err
	}
	orderExpirationSeconds := k.GetParams(ctx).OrderExpirationSeconds
	for _, order := range ordersToCancel {
		if _, exists := k.clobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId); exists {
			err := k.clobKeeper.HandleMsgCancelOrder(ctx, clobtypes.NewMsgCancelOrderStateful(
				order.OrderId,
				uint32(ctx.BlockTime().Unix())+orderExpirationSeconds,
			))
			if err != nil {
				log.ErrorLogWithError(ctx, "Failed to cancel order", err, "order", order, "vaultId", vaultId)
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
	for _, order := range ordersToPlace {
		err := k.PlaceVaultClobOrder(ctx, order)
		if err != nil {
			log.ErrorLogWithError(ctx, "Failed to place order", err, "order", order, "vaultId", vaultId)
		}
		vaultId.IncrCounterWithLabels(
			metrics.VaultPlaceOrder,
			metrics.GetLabelForBoolValue(metrics.Success, err == nil),
		)
	}

	return nil
}

// GetVaultClobOrders returns a list of long term orders for a given CLOB vault.
// Let n be number of layers, then the function returns orders at [a_0, b_0, a_1, b_1, ..., a_{n-1}, b_{n-1}]
// where a_i and b_i are the ask price and bid price at i-th layer. To compute a_i and b_i:
// - a_i = oraclePrice * (1 + skew_i) * (1 + spread)^{i+1}
// - b_i = oraclePrice * (1 + skew_i) / (1 + spread)^{i+1}
// - skew_i = -leverage_i * spread * skew_factor
// - leverage_i = leverage +/- i * order_size_pct\ (- for ask and + for bid)
// - leverage = open notional / equity
// - spread = max(spread_min, spread_buffer + min_price_change)
// and size of each order is calculated as `order_size * equity / oraclePrice`.
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

	// Get vault (subaccount 0 of corresponding module account).
	vault := vaultId.ToSubaccountId()
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
		int256.MustFromBig(inventory),
		perpetual.Params.AtomicResolution,
		marketPrice.GetPrice(),
		marketPrice.GetExponent(),
	).ToBig()
	leverage := new(big.Rat).Quo(
		new(big.Rat).SetInt(openNotional),
		new(big.Rat).SetInt(equity),
	)
	// Get parameters.
	params := k.GetParams(ctx)
	// Calculate order size (in base quantums).
	// order_size = order_size_pct * equity / oracle_price
	// = order_size_pct * equity / (price * 10^exponent / 10^quote_atomic_resolution) / 10^base_atomic_resolution
	// = order_size_pct * equity / (price * 10^(exponent - quote_atomic_resolution + base_atomic_resolution))
	orderSizeBaseQuantums := lib.BigRatMulPpm(
		new(big.Rat).SetInt(equity),
		params.OrderSizePctPpm,
	)
	orderSizeBaseQuantums = orderSizeBaseQuantums.Quo(
		orderSizeBaseQuantums,
		lib.BigMulPow10(
			new(big.Int).SetUint64(marketPrice.Price),
			marketPrice.Exponent-lib.QuoteCurrencyAtomicResolution+perpetual.Params.AtomicResolution,
		),
	)
	orderSizeBaseQuantumsRounded := lib.BigRatRoundToNearestMultiple(
		orderSizeBaseQuantums,
		uint32(clobPair.StepBaseQuantums),
		false,
	)
	// If order size is non-positive, return empty orders.
	if orderSizeBaseQuantumsRounded <= 0 {
		return []*clobtypes.Order{}, nil
	}
	// Calculate spread.
	spreadPpm := lib.Max(
		params.SpreadMinPpm,
		params.SpreadBufferPpm+marketParam.MinPriceChangePpm,
	)
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
	// Initialize spreadBaseMultiplier and spreadMultiplier as `1 + spread`.
	spreadBaseMultiplier := new(big.Rat).SetFrac(
		new(big.Int).SetUint64(uint64(lib.OneMillion+spreadPpm)),
		lib.BigIntOneMillion(),
	)
	spreadMultiplier := new(big.Rat).Set(spreadBaseMultiplier)
	// Construct one ask and one bid for each layer.
	constructOrder := func(
		side clobtypes.Order_Side,
		layer uint32,
		spreadMultipler *big.Rat,
	) *clobtypes.Order {
		// Calculate size that will have been filled before this layer is matched.
		sizeFilledByThisLayer := new(big.Rat).SetFrac(
			new(big.Int).SetUint64(uint64(params.OrderSizePctPpm*layer)),
			lib.BigIntOneMillion(),
		)

		// Ask: leverage_i = leverage - i * order_size_pct
		// Bid: leverage_i = leverage + i * order_size_pct
		var leverageI *big.Rat
		if side == clobtypes.Order_SIDE_SELL {
			leverageI = new(big.Rat).Sub(
				leverage,
				sizeFilledByThisLayer,
			)
		} else {
			leverageI = new(big.Rat).Add(
				leverage,
				sizeFilledByThisLayer,
			)
		}

		// skew_i = -leverage_i * spread * skew_factor
		skewI := lib.BigRatMulPpm(leverageI, spreadPpm)
		skewI = lib.BigRatMulPpm(skewI, params.SkewFactorPpm)
		skewI = skewI.Neg(skewI)

		// Ask: price = oracle price * (1 + skew_i) * (1 + spread)^(i+1)
		// Bid: price = oracle price * (1 + skew_i) / (1 + spread)^(i+1)
		orderSubticks := new(big.Rat).Add(skewI, new(big.Rat).SetUint64(1))
		orderSubticks = orderSubticks.Mul(orderSubticks, oracleSubticks)
		if side == clobtypes.Order_SIDE_SELL {
			orderSubticks = orderSubticks.Mul(
				orderSubticks,
				spreadMultipler,
			)
		} else {
			orderSubticks = orderSubticks.Quo(
				orderSubticks,
				spreadMultipler,
			)
		}

		return &clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *vault,
				ClientId:     k.GetVaultClobOrderClientId(ctx, side, uint8(layer)),
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   clobPair.Id,
			},
			Side:     side,
			Quantums: orderSizeBaseQuantumsRounded,
			Subticks: lib.BigRatRoundToNearestMultiple(
				orderSubticks,
				clobPair.SubticksPerTick,
				side == clobtypes.Order_SIDE_SELL, // round up for asks and down for bids.
			),
			GoodTilOneof: goodTilBlockTime,
		}
	}
	orders = make([]*clobtypes.Order, 2*params.Layers)
	for i := uint32(0); i < params.Layers; i++ {
		// Construct ask at this layer.
		orders[2*i] = constructOrder(
			clobtypes.Order_SIDE_SELL,
			i,
			spreadMultiplier,
		)

		// Construct bid at this layer.
		orders[2*i+1] = constructOrder(
			clobtypes.Order_SIDE_BUY,
			i,
			spreadMultiplier,
		)

		// Update spreadMultiplier for next layer.
		spreadMultiplier = spreadMultiplier.Mul(spreadMultiplier, spreadBaseMultiplier)
	}

	return orders, nil
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
