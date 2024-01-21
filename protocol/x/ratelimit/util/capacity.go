package util

import (
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

// CalculateNewCapacityList calculates the new capacity list for the given current `tvl` and `limitParams“.
// Input invariant: `len(prevCapacityList) == len(limitParams.Limiters)`
// Detailed math for calculating the updated capacity:
//
//	`baseline = max(baseline_minimum, baseline_tvl_ppm * tvl)`
//	`capacity_diff = max(baseline, capacity-baseline) * (time_since_last_block / period)`
//
// This is basically saying that the capacity returns to the baseline over the course of the `period`.
// Usually in a linear way, but if the `capacity` is more than twice the `baseline`, then in an exponential way.
//
//	`capacity =`
//	    if `abs(capacity - baseline) < capacity_diff` then `capacity = baseline`
//	    else if `capacity < baseline` then `capacity += capacity_diff`
//	    else `capacity -= capacity_diff`
//
// On a high level, `capacity` trends towards `baseline` by `capacity_diff` but does not “cross” it.
func CalculateNewCapacityList(
	bigTvl *big.Int,
	limitParams types.LimitParams,
	prevCapacityList []dtypes.SerializableInt,
	timeSinceLastBlock time.Duration,
) (
	newCapacityList []dtypes.SerializableInt,
	err error,
) {
	// Declare new capacity list to be populated.
	newCapacityList = make([]dtypes.SerializableInt, len(prevCapacityList))

	if len(limitParams.Limiters) != len(prevCapacityList) {
		// This violates an invariant. Since this is in the `EndBlocker`, we return an error instead of panicking.
		return nil, errorsmod.Wrapf(
			types.ErrMismatchedCapacityLimitersLength,
			"denom = %v, len(limiters) = %v, len(prevCapacityList) = %v",
			limitParams.Denom,
			len(limitParams.Limiters),
			len(prevCapacityList),
		)
	}

	for i, limiter := range limitParams.Limiters {
		// For each limiter, calculate the current baseline.
		baseline := GetBaseline(bigTvl, limiter)

		capacityMinusBaseline := new(big.Int).Sub(
			prevCapacityList[i].BigInt(), // array access is safe because of input invariant
			baseline,
		)

		// Calculate left operand: `max(baseline, capacity-baseline)`. This equals `baseline` when `capacity <= 2 * baseline`
		operandL := new(big.Rat).SetInt(
			lib.BigMax(
				baseline,
				capacityMinusBaseline,
			),
		)

		// Calculate right operand: `time_since_last_block / period`
		operandR := new(big.Rat).SetFrac64(
			timeSinceLastBlock.Milliseconds(),
			limiter.Period.Milliseconds(),
		)

		// Calculate: `capacity_diff = max(baseline, capacity-baseline) * (time_since_last_block / period)`
		// Since both operands > 0, `capacity_diff` is positive or zero (due to rounding).
		capacityDiffRat := new(big.Rat).Mul(operandL, operandR)
		capacityDiff := lib.BigRatRound(capacityDiffRat, false) // rounds down `capacity_diff`

		if new(big.Int).Abs(capacityMinusBaseline).Cmp(capacityDiff) <= 0 {
			// if `abs(capacity - baseline) < capacity_diff` then `capacity = baseline``
			newCapacityList[i] = dtypes.NewIntFromBigInt(baseline)
		} else if capacityMinusBaseline.Sign() < 0 {
			// else if `capacity < baseline` then `capacity += capacity_diff`
			newCapacityList[i] = dtypes.NewIntFromBigInt(
				new(big.Int).Add(
					prevCapacityList[i].BigInt(),
					capacityDiff,
				),
			)
		} else {
			// else `capacity -= capacity_diff`
			newCapacityList[i] = dtypes.NewIntFromBigInt(
				new(big.Int).Sub(
					prevCapacityList[i].BigInt(),
					capacityDiff,
				),
			)
		}
	}

	return newCapacityList, nil
}
