package keeper

import (
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Buckets hold millionths of an order so refill rates below one order per block still accrue.
const microOrdersPerOrder = 1_000_000

// Buckets are pruned, and those with rejections log their placements, every this many blocks.
const shortTermTokenBucketWindowBlocks = 100

// A bucket's refill and capacity are recomputed from current volume and collateral at most this often.
const shortTermTokenBucketRateBlocks = 100

// At most this many of an address's subaccounts count toward its collateral, which bounds the cost of a rate.
const shortTermTokenBucketMaxSubaccounts = 128

// A bucket with no placement for this many blocks is dropped even if it isn't full, so buckets whose refill is
// near zero don't stay in memory.
const shortTermTokenBucketIdleBlocks = 10_000

// shortTermTokenBucketRate is an address's refill per block and capacity in micro-orders, and the inputs they
// came from.
type shortTermTokenBucketRate struct {
	refill        uint64
	capacity      uint64
	volume        *big.Int
	netCollateral *big.Int
}

// shortTermTokenBucket is an owner address's short-term order budget in micro-orders, rated at rateHeight.
type shortTermTokenBucket struct {
	shortTermTokenBucketRate
	tokens     uint64
	height     int64
	rateHeight int64

	windowStart       int64
	windowPlaced      uint32
	windowRejected    uint32
	windowMaxPerBlock uint32
	blockPlaced       uint32
	blockReduceOnly   bool
}

// ValidateShortTermOrderTokenBucket charges one order to the owner address's token bucket and rejects the
// order when the bucket is empty, or only counts it in shadow mode. A new bucket starts with min-capacity orders
// and grows to its capacity at the refill rate, so moving money to a fresh address buys no burst. An empty bucket
// still lets one reduce-only order a block through, so a position can always be closed.
func (k Keeper) ValidateShortTermOrderTokenBucket(ctx sdk.Context, order types.Order) error {
	cfg := k.Flags
	if !cfg.ShortTermTokenBucketEnabled {
		return nil
	}
	owner := order.GetSubaccountId().Owner
	height := ctx.BlockHeight()
	bucket, found := k.shortTermTokenBuckets[owner]
	if !found {
		bucket = &shortTermTokenBucket{
			tokens:      uint64(cfg.ShortTermTokenBucketMinCapacity) * microOrdersPerOrder,
			height:      height,
			windowStart: height,
		}
		k.shortTermTokenBuckets[owner] = bucket
	}
	if bucket.height != height {
		bucket.tokens = refillTokens(bucket.tokens, bucket.capacity, bucket.refill, height-bucket.height)
		bucket.height, bucket.blockPlaced, bucket.blockReduceOnly = height, 0, false
	}
	if !found || height-bucket.rateHeight >= shortTermTokenBucketRateBlocks {
		bucket.shortTermTokenBucketRate = k.shortTermTokenBucketRate(ctx, owner)
		bucket.tokens = min(bucket.tokens, bucket.capacity)
		bucket.rateHeight = height
	}
	bucket.windowPlaced++
	bucket.blockPlaced++
	bucket.windowMaxPerBlock = max(bucket.windowMaxPerBlock, bucket.blockPlaced)

	if bucket.tokens >= microOrdersPerOrder {
		bucket.tokens -= microOrdersPerOrder
		return nil
	}
	if order.IsReduceOnly() && !bucket.blockReduceOnly {
		bucket.blockReduceOnly = true
		return nil
	}
	bucket.windowRejected++
	if cfg.ShortTermTokenBucketShadow {
		telemetry.IncrCounter(1, types.ModuleName, metrics.PlaceOrder, metrics.TokenBucket, metrics.ShadowRejected)
		return nil
	}
	return errorsmod.Wrapf(
		types.ErrShortTermOrderRateExceeded,
		"address %s refills %.4f orders per block up to %d",
		owner,
		float64(bucket.refill)/microOrdersPerOrder,
		bucket.capacity/microOrdersPerOrder,
	)
}

// PruneShortTermTokenBuckets closes a window every shortTermTokenBucketWindowBlocks blocks. Each bucket that
// rejected orders logs its window, then idle buckets and full buckets no larger than a new one are dropped.
// Larger full buckets are kept so a trader keeps the burst it has saved.
func (k Keeper) PruneShortTermTokenBuckets(ctx sdk.Context) {
	height := ctx.BlockHeight()
	if height%shortTermTokenBucketWindowBlocks != 0 {
		return
	}
	minCapacity := uint64(k.Flags.ShortTermTokenBucketMinCapacity) * microOrdersPerOrder
	for owner, bucket := range k.shortTermTokenBuckets {
		if bucket.windowRejected > 0 {
			k.logShortTermTokenBucketWindow(ctx, owner, bucket, height)
		}
		full := refillTokens(bucket.tokens, bucket.capacity, bucket.refill, height-bucket.height) == bucket.capacity
		if height-bucket.height >= shortTermTokenBucketIdleBlocks || (full && bucket.capacity <= minCapacity) {
			delete(k.shortTermTokenBuckets, owner)
		} else {
			bucket.windowStart, bucket.windowPlaced, bucket.windowRejected, bucket.windowMaxPerBlock = height, 0, 0, 0
		}
	}
	telemetry.SetGauge(float32(len(k.shortTermTokenBuckets)), types.ModuleName, metrics.TokenBucket, metrics.Count)
}

// shortTermTokenBucketRate rates an address from its 30-day volume and the net collateral of up to
// shortTermTokenBucketMaxSubaccounts of its subaccounts, each counted at no less than zero. A subaccount whose
// collateral can't be computed counts as zero; that takes invalid state, and the equity tier check rejects the
// placing subaccount first.
func (k Keeper) shortTermTokenBucketRate(ctx sdk.Context, owner string) shortTermTokenBucketRate {
	netCollateral := new(big.Int)
	for _, id := range k.subaccountsKeeper.GetSubaccountIdsByOwner(ctx, owner, shortTermTokenBucketMaxSubaccounts) {
		risk, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(ctx, satypes.Update{SubaccountId: id})
		if err == nil && risk.NC.Sign() > 0 {
			netCollateral.Add(netCollateral, risk.NC)
		}
	}
	stats := k.statsKeeper.GetUserStats(ctx, owner)
	volume := new(big.Int).SetUint64(stats.TakerNotional)
	volume.Add(volume, new(big.Int).SetUint64(stats.MakerNotional))

	cfg := k.Flags
	refill := microOrdersPerBlock(volume, cfg.ShortTermTokenBucketVolumePerOrder)
	refill.Add(refill, microOrdersPerBlock(netCollateral, cfg.ShortTermTokenBucketCollateralPerOrder))
	burst := new(big.Int).Mul(refill, new(big.Int).SetUint64(uint64(cfg.ShortTermTokenBucketBurstBlocks)))

	return shortTermTokenBucketRate{
		refill:        saturatingUint64(refill),
		capacity:      max(saturatingUint64(burst), uint64(cfg.ShortTermTokenBucketMinCapacity)*microOrdersPerOrder),
		volume:        volume,
		netCollateral: netCollateral,
	}
}

func (k Keeper) logShortTermTokenBucketWindow(
	ctx sdk.Context,
	owner string,
	bucket *shortTermTokenBucket,
	height int64,
) {
	k.Logger(ctx).Warn(
		"Short-term order token bucket window",
		"owner", owner,
		"window_blocks", height-bucket.windowStart,
		"placed", bucket.windowPlaced,
		"rejected", bucket.windowRejected,
		"max_per_block", bucket.windowMaxPerBlock,
		"volume_30d", bucket.volume,
		"net_collateral", bucket.netCollateral,
		"refill_per_block", float64(bucket.refill)/microOrdersPerOrder,
		"capacity", bucket.capacity/microOrdersPerOrder,
		"shadow", k.Flags.ShortTermTokenBucketShadow,
	)
}

// microOrdersPerBlock is amount/perOrder orders per block in micro-orders. A zero perOrder disables the term.
func microOrdersPerBlock(amount *big.Int, perOrder uint64) *big.Int {
	if perOrder == 0 || amount.Sign() <= 0 {
		return new(big.Int)
	}
	r := new(big.Int).Mul(amount, big.NewInt(microOrdersPerOrder))
	return r.Quo(r, new(big.Int).SetUint64(perOrder))
}

// refillTokens adds blocks × refill to tokens, capped at capacity.
func refillTokens(tokens, capacity, refill uint64, blocks int64) uint64 {
	tokens = min(tokens, capacity)
	if blocks <= 0 || refill == 0 {
		return tokens
	}
	if uint64(blocks) > (capacity-tokens)/refill {
		return capacity
	}
	return tokens + refill*uint64(blocks)
}

func saturatingUint64(x *big.Int) uint64 {
	if x.IsUint64() {
		return x.Uint64()
	}
	return math.MaxUint64
}
