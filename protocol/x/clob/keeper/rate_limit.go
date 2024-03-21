package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// RateLimitCancelOrder passes order cancellations with valid clob pairs to `cancelOrderRateLimiter`.
// The rate limiting is only performed during `CheckTx` and `ReCheckTx`.
func (k *Keeper) RateLimitCancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) error {
	// Only rate limit during `CheckTx` and `ReCheckTx`.
	if lib.IsDeliverTxMode(ctx) {
		return nil
	}

	_, found := k.GetClobPair(ctx, types.ClobPairId(msg.OrderId.GetClobPairId()))
	// If the clob pair isn't found then we expect order cancellation validation to fail the order cancellation as
	// being invalid.
	if !found {
		return nil
	}

	// Ensure that the GTB is valid before we attempt to rate limit. This is to prevent a replay attack
	// where short-term order cancellations with GTBs in the past or the far future could be replayed by an adversary.
	// Normally transaction replay attacks rely on sequence numbers being part of the signature and being incremented
	// for each transaction but sequence number verification is skipped for short-term orders.
	if msg.OrderId.IsShortTermOrder() {
		nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)
		if err := k.validateGoodTilBlock(msg.GetGoodTilBlock(), nextBlockHeight); err != nil {
			return err
		}
	}
	return k.placeCancelOrderRateLimiter.RateLimit(ctx, msg)
}

// RateLimitPlaceOrder passes orders with valid clob pairs to `placeOrderRateLimiter`.
// The rate limiting is only performed during `CheckTx` and `ReCheckTx`.
func (k *Keeper) RateLimitPlaceOrder(ctx sdk.Context, msg *types.MsgPlaceOrder) error {
	// Only rate limit during `CheckTx` and `ReCheckTx`.
	if lib.IsDeliverTxMode(ctx) {
		return nil
	}

	_, found := k.GetClobPair(ctx, msg.Order.GetClobPairId())
	// If the clob pair isn't found then we expect order validation to fail the order as being invalid.
	if !found {
		return nil
	}

	// Ensure that the GTB is valid before we attempt to rate limit. This is to prevent a replay attack
	// where short-term order placements with GTBs in the past or the far future could be replayed by an adversary.
	// Normally transaction replay attacks rely on sequence numbers being part of the signature and being incremented
	// for each transaction but sequence number verification is skipped for short-term orders.
	if msg.Order.IsShortTermOrder() {
		nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)
		if err := k.validateGoodTilBlock(msg.Order.GetGoodTilBlock(), nextBlockHeight); err != nil {
			return err
		}
	}

	return k.placeCancelOrderRateLimiter.RateLimit(ctx, msg)
}

// RateLimitBatchCancel passes orders with valid clob pairs to `placeOrderRateLimiter`.
// The rate limiting is only performed during `CheckTx` and `ReCheckTx`.
func (k *Keeper) RateLimitBatchCancel(ctx sdk.Context, msg *types.MsgBatchCancel) error {
	// Only rate limit during `CheckTx` and `ReCheckTx`.
	if lib.IsDeliverTxMode(ctx) {
		return nil
	}

	for _, batch := range msg.ShortTermCancels {
		_, found := k.GetClobPair(ctx, types.ClobPairId(batch.GetClobPairId()))
		// If the clob pair isn't found then we expect order validation to fail the order as being invalid.
		if !found {
			return nil
		}
	}

	// Ensure that the GTB is valid before we attempt to rate limit. This is to prevent a replay attack
	// where short-term order placements with GTBs in the past or the far future could be replayed by an adversary.
	// Normally transaction replay attacks rely on sequence numbers being part of the signature and being incremented
	// for each transaction but sequence number verification is skipped for short-term orders.
	nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)
	if err := k.validateGoodTilBlock(msg.GetGoodTilBlock(), nextBlockHeight); err != nil {
		return err
	}

	return k.placeCancelOrderRateLimiter.RateLimit(ctx, msg)
}

func (k *Keeper) PruneRateLimits(ctx sdk.Context) {
	k.placeCancelOrderRateLimiter.PruneRateLimits(ctx)
}
