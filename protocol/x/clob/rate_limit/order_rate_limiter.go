package rate_limit

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	BATCH_CANCEL_RATE_LIMIT_WEIGHT = uint32(2)
)

// A RateLimiter which rate limits types.MsgPlaceOrder, types.MsgCancelOrder, and
// types.MsgBatchCancel.
//
// The rate limiting keeps track of short term and stateful orders placed during
// CheckTx.
type placeAndCancelOrderRateLimiter struct {
	checkStateShortTermOrderPlaceCancelRateLimiter RateLimiter[string]
	checkStateStatefulOrderRateLimiter             RateLimiter[string]
	// The set of rate limited accounts is only stored for telemetry purposes.
	rateLimitedAccounts map[string]bool
}

var _ RateLimiter[sdk.Msg] = (*placeAndCancelOrderRateLimiter)(nil)

// NewPlaceCancelOrderRateLimiter returns a RateLimiter which rate limits types.MsgPlaceOrder, types.MsgCancelOrder,
// types.MsgBatchCancel based upon the provided types.BlockRateLimitConfiguration. The rate limiter currently
// supports limiting based upon:
//   - how many short term place/cancel orders per account (by using string).
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
func NewPlaceCancelOrderRateLimiter(config types.BlockRateLimitConfiguration) RateLimiter[sdk.Msg] {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	// Return the no-op rate limiter if the configuration is empty.
	if len(config.MaxShortTermOrdersAndCancelsPerNBlocks)+len(config.MaxStatefulOrdersPerNBlocks) == 0 {
		return noOpRateLimiter[sdk.Msg]{}
	}

	r := placeAndCancelOrderRateLimiter{
		rateLimitedAccounts: make(map[string]bool, 0),
	}
	if len(config.MaxShortTermOrdersAndCancelsPerNBlocks) == 0 {
		r.checkStateShortTermOrderPlaceCancelRateLimiter = NewNoOpRateLimiter[string]()
	} else if len(config.MaxShortTermOrdersAndCancelsPerNBlocks) == 1 &&
		config.MaxShortTermOrdersAndCancelsPerNBlocks[0].NumBlocks == 1 {
		r.checkStateShortTermOrderPlaceCancelRateLimiter = NewSingleBlockRateLimiter[string](
			"MaxShortTermOrdersAndCancelsPerNBlocks",
			config.MaxShortTermOrdersAndCancelsPerNBlocks[0],
		)
	} else {
		r.checkStateShortTermOrderPlaceCancelRateLimiter = NewMultiBlockRateLimiter[string](
			"MaxShortTermOrdersAndCancelsPerNBlocks",
			config.MaxShortTermOrdersAndCancelsPerNBlocks,
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

func (r *placeAndCancelOrderRateLimiter) RateLimit(ctx sdk.Context, msg sdk.Msg) (err error) {
	lib.AssertCheckTxMode(ctx)
	switch castedMsg := (msg).(type) {
	case *types.MsgCancelOrder:
		err = r.RateLimitCancelOrder(ctx, *castedMsg)
	case *types.MsgPlaceOrder:
		err = r.RateLimitPlaceOrder(ctx, *castedMsg)
	case *types.MsgBatchCancel:
		err = r.RateLimitBatchCancelOrder(ctx, *castedMsg)
	}
	return err
}

func (r *placeAndCancelOrderRateLimiter) RateLimitIncrBy(ctx sdk.Context, msg sdk.Msg, incrBy uint32) (err error) {
	panic("PlaceAndCancelOrderRateLimiter is a top-level rate limiter. It should not use IncrBy.")
}

func (r *placeAndCancelOrderRateLimiter) RateLimitPlaceOrder(ctx sdk.Context, msg types.MsgPlaceOrder) (err error) {
	lib.AssertCheckTxMode(ctx)
	if msg.Order.IsShortTermOrder() {
		err = r.checkStateShortTermOrderPlaceCancelRateLimiter.RateLimit(
			ctx,
			msg.Order.OrderId.SubaccountId.Owner,
		)
	} else {
		msg.Order.MustBeStatefulOrder()
		err = r.checkStateStatefulOrderRateLimiter.RateLimit(ctx, msg.Order.OrderId.SubaccountId.Owner)
	}

	if err != nil {
		metrics.IncrCounterWithLabels(
			metrics.ClobRateLimitPlaceOrderCount,
			1,
			msg.Order.GetOrderLabels()...,
		)
		r.rateLimitedAccounts[msg.Order.OrderId.SubaccountId.Owner] = true
	}
	return err
}

func (r *placeAndCancelOrderRateLimiter) RateLimitCancelOrder(
	ctx sdk.Context,
	msg types.MsgCancelOrder,
) (err error) {
	lib.AssertCheckTxMode(ctx)

	if msg.OrderId.IsShortTermOrder() {
		err = r.checkStateShortTermOrderPlaceCancelRateLimiter.RateLimit(
			ctx,
			msg.OrderId.SubaccountId.Owner,
		)
	}
	if err != nil {
		metrics.IncrCounterWithLabels(
			metrics.ClobRateLimitCancelOrderCount,
			1,
			msg.OrderId.GetOrderIdLabels()...,
		)
		r.rateLimitedAccounts[msg.OrderId.SubaccountId.Owner] = true
	}
	return err
}

func (r *placeAndCancelOrderRateLimiter) RateLimitBatchCancelOrder(
	ctx sdk.Context,
	msg types.MsgBatchCancel,
) (err error) {
	lib.AssertCheckTxMode(ctx)

	// TODO(CT-688) Use a scaling function such as (1 + ceil(0.1 * #cancels)) to calculate batch
	// cancel rate limit weights.
	err = r.checkStateShortTermOrderPlaceCancelRateLimiter.RateLimitIncrBy(
		ctx,
		msg.SubaccountId.Owner,
		BATCH_CANCEL_RATE_LIMIT_WEIGHT,
	)
	if err != nil {
		metrics.IncrCounterWithLabels(
			metrics.ClobRateLimitBatchCancelCount,
			1,
		)
		r.rateLimitedAccounts[msg.SubaccountId.Owner] = true
	}
	return err
}

func (r *placeAndCancelOrderRateLimiter) PruneRateLimits(ctx sdk.Context) {
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
	r.checkStateShortTermOrderPlaceCancelRateLimiter.PruneRateLimits(ctx)
	r.checkStateStatefulOrderRateLimiter.PruneRateLimits(ctx)
}

// NewUpdateLeverageRateLimiter returns a RateLimiter which rate limits leverage updates
// based upon the provided types.BlockRateLimitConfiguration. The rate limiter currently
// supports limiting based upon how many leverage updates per account per N blocks.
//
// The rate limiting must only be used during `CheckTx` because the rate limiting information is not recovered
// on application restart preventing it from being deterministic during `DeliverTx`.
//
// Depending upon the provided types.BlockRateLimitConfiguration, the returned RateLimiter may rely on:
//   - `ctx.BlockHeight()` in RateLimit to track which block the rate limit should apply to.
//   - `ctx.BlockHeight()` in PruneRateLimits and should be invoked during `EndBlocker`.
func NewUpdateLeverageRateLimiter(config types.BlockRateLimitConfiguration) RateLimiter[string] {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	// Return the no-op rate limiter if the configuration is empty.
	if len(config.MaxLeverageUpdatesPerNBlocks) == 0 {
		return noOpRateLimiter[string]{}
	}

	// Create the appropriate rate limiter based on configuration
	if len(config.MaxLeverageUpdatesPerNBlocks) == 1 &&
		config.MaxLeverageUpdatesPerNBlocks[0].NumBlocks == 1 {
		return NewSingleBlockRateLimiter[string](
			"UpdateLeverage",
			config.MaxLeverageUpdatesPerNBlocks[0],
		)
	} else {
		return NewMultiBlockRateLimiter[string](
			"UpdateLeverage",
			config.MaxLeverageUpdatesPerNBlocks,
		)
	}
}
