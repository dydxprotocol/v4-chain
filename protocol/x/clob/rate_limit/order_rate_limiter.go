package rate_limit

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// A RateLimiter which rate limits types.MsgPlaceOrder.
//
// The rate limiting keeps track of short term and stateful orders placed during
// CheckTx.
type placeOrderRateLimiter struct {
	checkStateShortTermOrderRateLimiter RateLimiter[string]
	checkStateStatefulOrderRateLimiter  RateLimiter[string]
	// The set of rate limited accounts is only stored for telemetry purposes.
	rateLimitedAccounts map[string]bool
}

var _ RateLimiter[*types.MsgPlaceOrder] = (*placeOrderRateLimiter)(nil)

// NewPlaceOrderRateLimiter returns a RateLimiter which rate limits types.MsgPlaceOrder based upon the provided
// types.BlockRateLimitConfiguration. The rate limiter currently supports limiting based upon:
//   - how many short term orders per account (by using string).
//   - how many stateful order per account (by using string).
//
// The rate limiting must only be used during `CheckTx` because the rate limiting information is not recovered
// on application restart preventing it from being deterministic during `DeliverTx`.
//
// Depending upon the provided types.BlockRateLimitConfiguration, the returned RateLimiter may rely on:
//   - `ctx.BlockHeight()` in RateLimit to track which block the rate limit should apply to.
//   - `ctx.BlockHeight()` in PruneRateLimits and should be invoked during `EndBlocker`. If invoked
//     during `PrepareCheckState` one must supply a `ctx` with the previous block height via
//     `ctx.WithBlockHeight(ctx.BlockHeight()-1)`.
func NewPlaceOrderRateLimiter(config types.BlockRateLimitConfiguration) RateLimiter[*types.MsgPlaceOrder] {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	// Return the no-op rate limiter if the configuration is empty.
	if len(config.MaxShortTermOrdersPerNBlocks)+len(config.MaxStatefulOrdersPerNBlocks) == 0 {
		return noOpRateLimiter[*types.MsgPlaceOrder]{}
	}

	r := placeOrderRateLimiter{
		rateLimitedAccounts: make(map[string]bool, 0),
	}
	if len(config.MaxShortTermOrdersPerNBlocks) == 0 {
		r.checkStateShortTermOrderRateLimiter = NewNoOpRateLimiter[string]()
	} else if len(config.MaxShortTermOrdersPerNBlocks) == 1 &&
		config.MaxShortTermOrdersPerNBlocks[0].NumBlocks == 1 {
		r.checkStateShortTermOrderRateLimiter = NewSingleBlockRateLimiter[string](
			"MaxShortTermOrdersPerNBlocks",
			config.MaxShortTermOrdersPerNBlocks[0],
		)
	} else {
		r.checkStateShortTermOrderRateLimiter = NewMultiBlockRateLimiter[string](
			"MaxShortTermOrdersPerNBlocks",
			config.MaxShortTermOrdersPerNBlocks,
		)
	}
	if len(config.MaxStatefulOrdersPerNBlocks) == 0 {
		r.checkStateStatefulOrderRateLimiter = NewNoOpRateLimiter[string]()
	} else if len(config.MaxStatefulOrdersPerNBlocks) == 1 &&
		config.MaxStatefulOrdersPerNBlocks[0].NumBlocks == 1 {
		r.checkStateStatefulOrderRateLimiter = NewSingleBlockRateLimiter[string](
			"MaxStatefulOrdersPerNBlocks",
			config.MaxStatefulOrdersPerNBlocks[0],
		)
	} else {
		r.checkStateStatefulOrderRateLimiter = NewMultiBlockRateLimiter[string](
			"MaxStatefulOrdersPerNBlocks",
			config.MaxStatefulOrdersPerNBlocks,
		)
	}

	return &r
}

func (r *placeOrderRateLimiter) RateLimit(ctx sdk.Context, msg *types.MsgPlaceOrder) (err error) {
	lib.AssertCheckTxMode(ctx)

	if msg.Order.IsShortTermOrder() {
		err = r.checkStateShortTermOrderRateLimiter.RateLimit(
			ctx,
			msg.Order.OrderId.SubaccountId.Owner,
		)
	} else {
		msg.Order.MustBeStatefulOrder()
		err = r.checkStateStatefulOrderRateLimiter.RateLimit(ctx, msg.Order.OrderId.SubaccountId.Owner)
	}

	if err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.RateLimit, metrics.PlaceOrder, metrics.Count},
			1,
			msg.Order.GetOrderLabels(),
		)
		r.rateLimitedAccounts[msg.Order.OrderId.SubaccountId.Owner] = true
	}
	return err
}

func (r *placeOrderRateLimiter) PruneRateLimits(ctx sdk.Context) {
	telemetry.IncrCounter(
		float32(len(r.rateLimitedAccounts)),
		types.ModuleName,
		metrics.RateLimit,
		metrics.PlaceOrderAccounts,
		metrics.Count,
	)
	// Note that this method for clearing the map is optimized by the go compiler significantly
	// and will leave the relative size of the map the same so that it doesn't need to be resized
	// often.
	for key := range r.rateLimitedAccounts {
		delete(r.rateLimitedAccounts, key)
	}
	r.checkStateShortTermOrderRateLimiter.PruneRateLimits(ctx)
	r.checkStateStatefulOrderRateLimiter.PruneRateLimits(ctx)
}

// A RateLimiter which rate limits types.MsgCancelOrder.
//
// The rate limiting keeps track of short term order cancellations during CheckTx.
type cancelOrderRateLimiter struct {
	checkStateShortTermRateLimiter RateLimiter[string]
	// The set of rate limited accounts is only stored for telemetry purposes.
	rateLimitedAccounts map[string]bool
}

var _ RateLimiter[*types.MsgCancelOrder] = (*cancelOrderRateLimiter)(nil)

// NewCancelOrderRateLimiter returns a RateLimiter which rate limits types.MsgCancelOrder based upon the provided
// types.BlockRateLimitConfiguration. The rate limiter currently supports limiting based upon:
//   - how many short term order cancellations per account (by using string).
//
// The rate limiting must only be used during `CheckTx` because the rate limiting information is not recovered
// on application restart preventing it from being deterministic during `DeliverTx`.
//
// Depending upon the provided types.BlockRateLimitConfiguration, the returned RateLimiter may rely on:
//   - `ctx.BlockHeight()` in RateLimit to track which block the rate limit should apply to.
//   - `ctx.BlockHeight()` in PruneRateLimits and should be invoked during `EndBlocker`. If invoked
//     during `PrepareCheckState` one must supply a `ctx` with the previous block height via
//     `ctx.WithBlockHeight(ctx.BlockHeight()-1)`.
func NewCancelOrderRateLimiter(config types.BlockRateLimitConfiguration) RateLimiter[*types.MsgCancelOrder] {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	// Return the no-op rate limiter if the configuration is empty.
	if len(config.MaxShortTermOrderCancellationsPerNBlocks) == 0 {
		return noOpRateLimiter[*types.MsgCancelOrder]{}
	}

	rateLimiter := cancelOrderRateLimiter{
		rateLimitedAccounts: make(map[string]bool, 0),
	}
	if len(config.MaxShortTermOrderCancellationsPerNBlocks) == 1 &&
		config.MaxShortTermOrderCancellationsPerNBlocks[0].NumBlocks == 1 {
		rateLimiter.checkStateShortTermRateLimiter = NewSingleBlockRateLimiter[string](
			"MaxShortTermOrdersPerNBlocks",
			config.MaxShortTermOrderCancellationsPerNBlocks[0],
		)
		return &rateLimiter
	} else {
		rateLimiter.checkStateShortTermRateLimiter = NewMultiBlockRateLimiter[string](
			"MaxShortTermOrdersPerNBlocks",
			config.MaxShortTermOrderCancellationsPerNBlocks,
		)
		return &rateLimiter
	}
}

func (r *cancelOrderRateLimiter) RateLimit(ctx sdk.Context, msg *types.MsgCancelOrder) (err error) {
	lib.AssertCheckTxMode(ctx)

	if msg.OrderId.IsShortTermOrder() {
		err = r.checkStateShortTermRateLimiter.RateLimit(
			ctx,
			msg.OrderId.SubaccountId.Owner,
		)
	}
	if err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.RateLimit, metrics.CancelOrder, metrics.Count},
			1,
			msg.OrderId.GetOrderIdLabels(),
		)
		r.rateLimitedAccounts[msg.OrderId.SubaccountId.Owner] = true
	}
	return err
}

func (r *cancelOrderRateLimiter) PruneRateLimits(ctx sdk.Context) {
	telemetry.IncrCounter(
		float32(len(r.rateLimitedAccounts)),
		types.ModuleName,
		metrics.RateLimit,
		metrics.CancelOrderAccounts,
		metrics.Count,
	)
	// Note that this method for clearing the map is optimized by the go compiler significantly
	// and will leave the relative size of the map the same so that it doesn't need to be resized
	// often.
	for key := range r.rateLimitedAccounts {
		delete(r.rateLimitedAccounts, key)
	}
	r.checkStateShortTermRateLimiter.PruneRateLimits(ctx)
}
