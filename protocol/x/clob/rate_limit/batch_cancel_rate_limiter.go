package rate_limit

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// A RateLimiter which rate limits types.MsgBatchCancel.
//
// The rate limiting keeps track of short term and stateful orders placed during
// CheckTx.
type batchCancelRateLimiter struct {
	checkStateBatchCancelRateLimiter RateLimiter[string]
	// The set of rate limited accounts is only stored for telemetry purposes.
	rateLimitedAccounts map[string]bool
}

var _ RateLimiter[*types.MsgBatchCancel] = (*batchCancelRateLimiter)(nil)

// NewbatchCancelRateLimiter returns a RateLimiter which rate limits types.MsgBatchCancel based upon the provided
// types.BlockRateLimitConfiguration. The rate limiter currently supports limiting based upon:
//   - how many batch cancel orders per account (by using string).
//
// The rate limiting must only be used during `CheckTx` because the rate limiting information is not recovered
// on application restart preventing it from being deterministic during `DeliverTx`.
//
// Depending upon the provided types.BlockRateLimitConfiguration, the returned RateLimiter may rely on:
//   - `ctx.BlockHeight()` in RateLimit to track which block the rate limit should apply to.
//   - `ctx.BlockHeight()` in PruneRateLimits and should be invoked during `EndBlocker`. If invoked
//     during `PrepareCheckState` one must supply a `ctx` with the previous block height via
//     `ctx.WithBlockHeight(ctx.BlockHeight()-1)`.
func NewBatchCancelRateLimiter(config types.BlockRateLimitConfiguration) RateLimiter[*types.MsgBatchCancel] {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	// Return the no-op rate limiter if the configuration is empty.
	if len(config.MaxBatchCancelsPerNBlocks)+len(config.MaxStatefulOrdersPerNBlocks) == 0 {
		return noOpRateLimiter[*types.MsgBatchCancel]{}
	}

	r := batchCancelRateLimiter{
		rateLimitedAccounts: make(map[string]bool, 0),
	}
	if len(config.MaxBatchCancelsPerNBlocks) == 0 {
		r.checkStateBatchCancelRateLimiter = NewNoOpRateLimiter[string]()
	} else if len(config.MaxBatchCancelsPerNBlocks) == 1 &&
		config.MaxBatchCancelsPerNBlocks[0].NumBlocks == 1 {
		r.checkStateBatchCancelRateLimiter = NewSingleBlockRateLimiter[string](
			"MaxBatchCancelsPerNBlocks",
			config.MaxBatchCancelsPerNBlocks[0],
		)
	} else {
		r.checkStateBatchCancelRateLimiter = NewMultiBlockRateLimiter[string](
			"MaxBatchCancelsPerNBlocks",
			config.MaxBatchCancelsPerNBlocks,
		)
	}

	return &r
}

func (r *batchCancelRateLimiter) RateLimit(ctx sdk.Context, msg *types.MsgBatchCancel) (err error) {
	lib.AssertCheckTxMode(ctx)
	err = r.checkStateBatchCancelRateLimiter.RateLimit(
		ctx,
		msg.SubaccountId.Owner,
	)

	if err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.RateLimit, metrics.PlaceOrder, metrics.Count},
			1,
			[]metrics.Label{},
		)
		r.rateLimitedAccounts[msg.SubaccountId.Owner] = true
	}
	return err
}

func (r *batchCancelRateLimiter) PruneRateLimits(ctx sdk.Context) {
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
	r.checkStateBatchCancelRateLimiter.PruneRateLimits(ctx)
}
