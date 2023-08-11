package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// RateLimitCancelOrder passes order cancellations with valid clob pairs to `cancelOrderRateLimiter`.
func (k *Keeper) RateLimitCancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) error {
	_, found := k.GetClobPair(ctx, types.ClobPairId(msg.OrderId.GetClobPairId()))
	// If the clob pair isn't found then we expect order cancellation validation to fail the order cancellation as
	// being invalid.
	if !found {
		return nil
	}

	return k.cancelOrderRateLimiter.RateLimit(ctx, msg)
}

// RateLimitPlaceOrder passes orders with valid clob pairs to `placeOrderRateLimiter`.
func (k *Keeper) RateLimitPlaceOrder(ctx sdk.Context, msg *types.MsgPlaceOrder) error {
	_, found := k.GetClobPair(ctx, msg.Order.GetClobPairId())
	// If the clob pair isn't found then we expect order validation to fail the order as being invalid.
	if !found {
		return nil
	}

	return k.placeOrderRateLimiter.RateLimit(ctx, msg)
}

func (k *Keeper) PruneRateLimits(ctx sdk.Context) {
	k.placeOrderRateLimiter.PruneRateLimits(ctx)
	k.cancelOrderRateLimiter.PruneRateLimits(ctx)
}
