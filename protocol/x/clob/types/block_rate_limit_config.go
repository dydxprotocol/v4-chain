package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxShortTermOrdersPerMarketPerNBlocksNumBlocks             = 1_000
	MaxShortTermOrdersPerMarketPerNBlocksLimit                 = 10_000_000
	MaxShortTermOrderCancellationsPerMarketPerNBlocksNumBlocks = 1_000
	MaxShortTermOrderCancellationsPerMarketPerNBlocksLimit     = 10_000_000
	MaxStatefulOrdersPerNBlocksNumBlocks                       = 10_000
	MaxStatefulOrdersPerNBlocksLimit                           = 1_000_000
)

// Validate validates each individual MaxPerNBlocksRateLimit.
// It returns an error if any of the rate limits fail the following validations:
//   - `Limit == 0` || `Limit > MaxShortTermOrdersPerMarketPerNBlocksLimit` for short term order rate limits.
//   - `NumBlocks == 0` || `NumBlocks > MaxShortTermOrdersPerMarketPerNBlocksNumBlocks` for short term order rate
//     limits.
//   - `Limit == 0` || `Limit > MaxStatefulOrdersPerNBlocksLimit` for stateful order rate limits.
//   - `NumBlocks == 0` || `NumBlocks > MaxStatefulOrdersPerNBlocksNumBlocks` for stateful order rate limits.
//   - `Limit == 0` || `Limit > MaxShortTermOrderCancellationsPerMarketPerNBlocksNumBlocks` for short term order
//     cancellation rate limits.
//   - `NumBlocks == 0` || `NumBlocks > MaxShortTermOrderCancellationsPerMarketPerNBlocksLimit` for short term order
//     cancellation rate limits.
//   - There are multiple rate limits for the same `NumBlocks` in `MaxShortTermOrdersPerMarketPerNBlocks`,
//     `MaxStatefulOrdersPerNBlocks`, or `MaxShortTermOrderCancellationsPerMarketPerNBlocks`.
func (lc BlockRateLimitConfiguration) Validate() error {
	if err := (maxPerNBlocksRateLimits)(lc.MaxShortTermOrdersPerMarketPerNBlocks).validate(
		"MaxShortTermOrdersPerMarketPerNBlocks",
		MaxShortTermOrdersPerMarketPerNBlocksNumBlocks,
		MaxShortTermOrdersPerMarketPerNBlocksLimit,
	); err != nil {
		return err
	}
	if err := (maxPerNBlocksRateLimits)(lc.MaxStatefulOrdersPerNBlocks).validate(
		"MaxStatefulOrdersPerNBlocks",
		MaxStatefulOrdersPerNBlocksNumBlocks,
		MaxStatefulOrdersPerNBlocksLimit,
	); err != nil {
		return err
	}
	if err := (maxPerNBlocksRateLimits)(lc.MaxShortTermOrderCancellationsPerMarketPerNBlocks).validate(
		"MaxShortTermOrderCancellationsPerMarketPerNBlocks",
		MaxShortTermOrderCancellationsPerMarketPerNBlocksNumBlocks,
		MaxShortTermOrderCancellationsPerMarketPerNBlocksLimit,
	); err != nil {
		return err
	}
	return nil
}

type maxPerNBlocksRateLimits []MaxPerNBlocksRateLimit

func (rl maxPerNBlocksRateLimits) validate(field string, maxBlocks uint32, maxOrders uint32) error {
	duplicates := make(map[uint32]MaxPerNBlocksRateLimit, 0)
	for _, rateLimit := range rl {
		if err := rateLimit.validate(
			field,
			maxBlocks,
			maxOrders); err != nil {
			return err
		}
		if existing, found := duplicates[rateLimit.NumBlocks]; found {
			return sdkerrors.Wrapf(
				ErrInvalidBlockRateLimitConfig,
				"Multiple rate limits %+v and %+v for the same block height found for %s",
				existing,
				rateLimit,
				field)
		}
		duplicates[rateLimit.NumBlocks] = rateLimit
	}
	return nil
}

func (rl MaxPerNBlocksRateLimit) validate(field string, maxBlocks uint32, maxOrders uint32) error {
	if rl.Limit == 0 || rl.Limit > maxOrders {
		return sdkerrors.Wrapf(
			ErrInvalidBlockRateLimitConfig,
			"%d is not a valid Limit for %s rate limit %+v",
			rl.Limit,
			field,
			rl)
	}
	if rl.NumBlocks == 0 || rl.NumBlocks > maxBlocks {
		return sdkerrors.Wrapf(
			ErrInvalidBlockRateLimitConfig,
			"%d is not a valid NumBlocks for %s rate limit %+v",
			rl.NumBlocks,
			field,
			rl)
	}
	return nil
}
