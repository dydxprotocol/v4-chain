package rate_limit

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// Used as the key for per market short term order limits.
type subaccountIdAndClobPairId struct {
	subaccountId satypes.SubaccountId
	clobPairId   types.ClobPairId
}

// A RateLimiter which rate limits types.MsgPlaceOrder.
//
// The rate limiting keeps track of short term and stateful orders placed during
// CheckTx separately from when they are placed during DeliverTx modes.
type placeOrderRateLimiter struct {
	checkStateShortTermOrderRateLimiter  RateLimiter[subaccountIdAndClobPairId]
	checkStateStatefulOrderRateLimiter   RateLimiter[satypes.SubaccountId]
	deliverStateStatefulOrderRateLimiter RateLimiter[satypes.SubaccountId]
	// The set of rate limited subaccounts is only stored for telemetry purposes.
	rateLimitedSubaccounts map[satypes.SubaccountId]bool
}

var _ RateLimiter[*types.MsgPlaceOrder] = (*placeOrderRateLimiter)(nil)

// NewPlaceOrderRateLimiter returns a RateLimiter which rate limits types.MsgPlaceOrder based upon the provided
// types.BlockRateLimitConfiguration. The rate limiter currently supports limiting based upon:
//   - how many short term orders per market and subaccount (by using the union type subaccountIdAndClobPairId).
//   - how many stateful order per subaccount (by using satypes.SubaccountId).
//
// The rate limiting keeps track of orders placed during CheckTx separately from when they
// are placed during DeliverTx modes.
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
	if len(config.MaxShortTermOrdersPerMarketPerNBlocks)+len(config.MaxStatefulOrdersPerNBlocks) == 0 {
		return noOpRateLimiter[*types.MsgPlaceOrder]{}
	}

	r := placeOrderRateLimiter{
		rateLimitedSubaccounts: make(map[satypes.SubaccountId]bool, 0),
	}
	if len(config.MaxShortTermOrdersPerMarketPerNBlocks) == 0 {
		r.checkStateShortTermOrderRateLimiter = NewNoOpRateLimiter[subaccountIdAndClobPairId]()
	} else if len(config.MaxShortTermOrdersPerMarketPerNBlocks) == 1 &&
		config.MaxShortTermOrdersPerMarketPerNBlocks[0].NumBlocks == 1 {
		r.checkStateShortTermOrderRateLimiter = NewSingleBlockRateLimiter[subaccountIdAndClobPairId](
			"MaxShortTermOrdersPerMarketPerNBlocks",
			config.MaxShortTermOrdersPerMarketPerNBlocks[0],
		)
	} else {
		r.checkStateShortTermOrderRateLimiter = NewMultiBlockRateLimiter[subaccountIdAndClobPairId](
			"MaxShortTermOrdersPerMarketPerNBlocks",
			config.MaxShortTermOrdersPerMarketPerNBlocks,
		)
	}
	if len(config.MaxStatefulOrdersPerNBlocks) == 0 {
		r.checkStateStatefulOrderRateLimiter = NewNoOpRateLimiter[satypes.SubaccountId]()
		r.deliverStateStatefulOrderRateLimiter = NewNoOpRateLimiter[satypes.SubaccountId]()
	} else if len(config.MaxStatefulOrdersPerNBlocks) == 1 &&
		config.MaxStatefulOrdersPerNBlocks[0].NumBlocks == 1 {
		r.checkStateStatefulOrderRateLimiter = NewSingleBlockRateLimiter[satypes.SubaccountId](
			"MaxStatefulOrdersPerNBlocks",
			config.MaxStatefulOrdersPerNBlocks[0],
		)
		r.deliverStateStatefulOrderRateLimiter = NewSingleBlockRateLimiter[satypes.SubaccountId](
			"MaxStatefulOrdersPerNBlocks",
			config.MaxStatefulOrdersPerNBlocks[0],
		)
	} else {
		r.checkStateStatefulOrderRateLimiter = NewMultiBlockRateLimiter[satypes.SubaccountId](
			"MaxStatefulOrdersPerNBlocks",
			config.MaxStatefulOrdersPerNBlocks,
		)
		r.deliverStateStatefulOrderRateLimiter = NewMultiBlockRateLimiter[satypes.SubaccountId](
			"MaxStatefulOrdersPerNBlocks",
			config.MaxStatefulOrdersPerNBlocks,
		)
	}

	return &r
}

func (r *placeOrderRateLimiter) RateLimit(ctx sdk.Context, msg *types.MsgPlaceOrder) (err error) {
	if lib.IsDeliverTxMode(ctx) {
		// Short-Term orders don't go through AnteHandler in DeliverTx since they're placed
		// as part of MsgProposedOperations and don't need to be rate limited since the user
		// will pay fees.
		if msg.Order.IsShortTermOrder() {
			return nil
		} else {
			msg.Order.MustBeStatefulOrder()
			err = r.deliverStateStatefulOrderRateLimiter.RateLimit(
				// We specifically pass in `height-1` since we want the deliverTx rate limiting to happen
				// as if the order was placed in the last block so that PruneRateLimits during EndBlocker
				// doesn't immediately clear it out.
				ctx.WithBlockHeight(ctx.BlockHeight()-1),
				msg.Order.GetSubaccountId(),
			)
		}
	} else {
		if msg.Order.IsShortTermOrder() {
			err = r.checkStateShortTermOrderRateLimiter.RateLimit(
				ctx,
				subaccountIdAndClobPairId{
					subaccountId: msg.Order.GetSubaccountId(),
					clobPairId:   msg.Order.GetClobPairId(),
				},
			)
		} else {
			msg.Order.MustBeStatefulOrder()
			err = r.checkStateStatefulOrderRateLimiter.RateLimit(ctx, msg.Order.GetSubaccountId())
		}
	}
	if err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.RateLimit, metrics.PlaceOrder, metrics.Count},
			1,
			msg.Order.GetOrderLabels(),
		)
		r.rateLimitedSubaccounts[msg.Order.GetSubaccountId()] = true
	}
	return err
}

func (r *placeOrderRateLimiter) PruneRateLimits(ctx sdk.Context) {
	telemetry.IncrCounter(
		float32(len(r.rateLimitedSubaccounts)),
		types.ModuleName,
		metrics.RateLimit,
		metrics.PlaceOrderSubaccounts,
		metrics.Count,
	)
	// Note that this method for clearing the map is optimized by the go compiler significantly
	// and will leave the relative size of the map the same so that it doesn't need to be resized
	// often.
	for key := range r.rateLimitedSubaccounts {
		delete(r.rateLimitedSubaccounts, key)
	}
	r.checkStateShortTermOrderRateLimiter.PruneRateLimits(ctx)
	r.checkStateStatefulOrderRateLimiter.PruneRateLimits(ctx)
	r.deliverStateStatefulOrderRateLimiter.PruneRateLimits(ctx)
}

// A RateLimiter which rate limits types.MsgCancelOrder.
//
// The rate limiting keeps track of short term order cancellations during CheckTx.
type cancelOrderRateLimiter struct {
	checkStateShortTermRateLimiter RateLimiter[subaccountIdAndClobPairId]
	// The set of rate limited subaccounts is only stored for telemetry purposes.
	rateLimitedSubaccounts map[satypes.SubaccountId]bool
}

var _ RateLimiter[*types.MsgCancelOrder] = (*cancelOrderRateLimiter)(nil)

// NewCancelOrderRateLimiter returns a RateLimiter which rate limits types.MsgCancelOrder based upon the provided
// types.BlockRateLimitConfiguration. The rate limiter currently supports limiting based upon:
//   - how many short term order cancellations per market and subaccount (by using the union type
//     subaccountIdAndClobPairId).
//
// The rate limiting keeps track of order cancellations placed during CheckTx.
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
	if len(config.MaxShortTermOrderCancellationsPerMarketPerNBlocks) == 0 {
		return noOpRateLimiter[*types.MsgCancelOrder]{}
	} else if len(config.MaxShortTermOrderCancellationsPerMarketPerNBlocks) == 1 &&
		config.MaxShortTermOrderCancellationsPerMarketPerNBlocks[0].NumBlocks == 1 {
		return &cancelOrderRateLimiter{
			rateLimitedSubaccounts: make(map[satypes.SubaccountId]bool, 0),
			checkStateShortTermRateLimiter: NewSingleBlockRateLimiter[subaccountIdAndClobPairId](
				"MaxShortTermOrdersPerMarketPerNBlocks",
				config.MaxShortTermOrderCancellationsPerMarketPerNBlocks[0],
			),
		}
	} else {
		return &cancelOrderRateLimiter{
			rateLimitedSubaccounts: make(map[satypes.SubaccountId]bool, 0),
			checkStateShortTermRateLimiter: NewMultiBlockRateLimiter[subaccountIdAndClobPairId](
				"MaxShortTermOrdersPerMarketPerNBlocks",
				config.MaxShortTermOrderCancellationsPerMarketPerNBlocks,
			),
		}
	}
}

func (r *cancelOrderRateLimiter) RateLimit(ctx sdk.Context, msg *types.MsgCancelOrder) (err error) {
	// Short-Term order cancellations don't go through AnteHandler in DeliverTx since they are removed
	// from the orderbook immediately which prevents them from being matched and we don't perform
	// any rate limiting on stateful order cancellation since the order must exist in state for it be
	// accepted and will be rejected otherwise so there is no need to rate limit either of them.
	if lib.IsDeliverTxMode(ctx) {
		return nil
	}

	if msg.OrderId.IsShortTermOrder() {
		err = r.checkStateShortTermRateLimiter.RateLimit(
			ctx,
			subaccountIdAndClobPairId{
				subaccountId: msg.OrderId.GetSubaccountId(),
				clobPairId:   types.ClobPairId(msg.OrderId.ClobPairId),
			},
		)
	}
	if err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.RateLimit, metrics.CancelOrder, metrics.Count},
			1,
			msg.OrderId.GetOrderIdLabels(),
		)
		r.rateLimitedSubaccounts[msg.OrderId.GetSubaccountId()] = true
	}
	return err
}

func (r *cancelOrderRateLimiter) PruneRateLimits(ctx sdk.Context) {
	telemetry.IncrCounter(
		float32(len(r.rateLimitedSubaccounts)),
		types.ModuleName,
		metrics.RateLimit,
		metrics.CancelOrderSubaccounts,
		metrics.Count,
	)
	// Note that this method for clearing the map is optimized by the go compiler significantly
	// and will leave the relative size of the map the same so that it doesn't need to be resized
	// often.
	for key := range r.rateLimitedSubaccounts {
		delete(r.rateLimitedSubaccounts, key)
	}
	r.checkStateShortTermRateLimiter.PruneRateLimits(ctx)
}
