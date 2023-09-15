package types

import (
	errorsmod "cosmossdk.io/errors"
)

const (
	MaxShortTermOrdersPerNBlocksNumBlocks             = 1_000
	MaxShortTermOrdersPerNBlocksLimit                 = 10_000_000
	MaxShortTermOrderCancellationsPerNBlocksNumBlocks = 1_000
	MaxShortTermOrderCancellationsPerNBlocksLimit     = 10_000_000
	MaxStatefulOrdersPerNBlocksNumBlocks              = 10_000
	MaxStatefulOrdersPerNBlocksLimit                  = 1_000_000
)

// Validate validates each individual MaxPerNBlocksRateLimit.
// It returns an error if any of the rate limits fail the following validations:
//   - `Limit == 0` || `Limit > MaxShortTermOrdersPerNBlocksLimit` for short term order rate limits.
//   - `NumBlocks == 0` || `NumBlocks > MaxShortTermOrdersPerNBlocksNumBlocks` for short term order rate
//     limits.
//   - `Limit == 0` || `Limit > MaxStatefulOrdersPerNBlocksLimit` for stateful order rate limits.
//   - `NumBlocks == 0` || `NumBlocks > MaxStatefulOrdersPerNBlocksNumBlocks` for stateful order rate limits.
//   - `Limit == 0` || `Limit > MaxShortTermOrderCancellationsPerNBlocksNumBlocks` for short term order
//     cancellation rate limits.
//   - `NumBlocks == 0` || `NumBlocks > MaxShortTermOrderCancellationsPerNBlocksLimit` for short term order
//     cancellation rate limits.
//   - There are multiple rate limits for the same `NumBlocks` in `MaxShortTermOrdersPerNBlocks`,
//     `MaxStatefulOrdersPerNBlocks`, or `MaxShortTermOrderCancellationsPerNBlocks`.
func (lc BlockRateLimitConfiguration) Validate() error {
	if err := (maxPerNBlocksRateLimits)(lc.MaxShortTermOrdersPerNBlocks).validate(
		"MaxShortTermOrdersPerNBlocks",
		MaxShortTermOrdersPerNBlocksNumBlocks,
		MaxShortTermOrdersPerNBlocksLimit,
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
	if err := (maxPerNBlocksRateLimits)(lc.MaxShortTermOrderCancellationsPerNBlocks).validate(
		"MaxShortTermOrderCancellationsPerNBlocks",
		MaxShortTermOrderCancellationsPerNBlocksNumBlocks,
		MaxShortTermOrderCancellationsPerNBlocksLimit,
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
			return errorsmod.Wrapf(
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
		return errorsmod.Wrapf(
			ErrInvalidBlockRateLimitConfig,
			"%d is not a valid Limit for %s rate limit %+v",
			rl.Limit,
			field,
			rl)
	}
	if rl.NumBlocks == 0 || rl.NumBlocks > maxBlocks {
		return errorsmod.Wrapf(
			ErrInvalidBlockRateLimitConfig,
			"%d is not a valid NumBlocks for %s rate limit %+v",
			rl.NumBlocks,
			field,
			rl)
	}
	return nil
}
