package util

import (
	"math/big"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
)

// CalculateNewCapacityList calculates the new capacity list for the given current `tvl` and `limitParams“.
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
	prevCapapcityList []dtypes.SerializableInt,
	timeSinceLastBlock time.Duration,
) (newCapacityList []dtypes.SerializableInt) {
	// Declare new capacity list to be populated.
	newCapacityList = make([]dtypes.SerializableInt, len(prevCapapcityList))

	for i, limiter := range limitParams.Limiters {
		// For each limiter, calculate the current baseline.
		baseline := GetBaseline(bigTvl, limiter)

		capacityMinusBaseline := new(big.Int).Sub(
			prevCapapcityList[i].BigInt(), // array access is safe because of the invariant check above
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
		periodMilli := new(big.Int).Mul(
			new(big.Int).SetUint64(uint64(limiter.PeriodSec)),
			big.NewInt(1000),
		)
		operandR := new(big.Rat).SetFrac(
			new(big.Int).SetInt64(timeSinceLastBlock.Milliseconds()),
			periodMilli,
		)

		// Calculate: `capacity_diff = max(baseline, capacity-baseline) * (time_since_last_block / period)`
		// Since both operands > 0, `capacity_diff` is positive or zero (due to rounding).
		capacityDiff := new(big.Rat).Mul(
			operandL,
			operandR,
		)

		bigRatcapacityMinusBaseline := new(big.Rat).SetInt(capacityMinusBaseline)

		if new(big.Rat).Abs(bigRatcapacityMinusBaseline).Cmp(capacityDiff) < 0 {
			// if `abs(capacity - baseline) < capacity_diff` then `capacity = baseline``
			newCapacityList[i] = dtypes.NewIntFromBigInt(baseline)
		} else if capacityMinusBaseline.Sign() < 0 {
			// else if `capacity < baseline` then `capacity += capacity_diff`
			newCapacityList[i] = dtypes.NewIntFromBigInt(
				new(big.Int).Add(
					prevCapapcityList[i].BigInt(),
					lib.BigRatRound(capacityDiff, false), // rounds down `capacity_diff`
				),
			)
		} else {
			// else `capacity -= capacity_diff`
			newCapacityList[i] = dtypes.NewIntFromBigInt(
				new(big.Int).Sub(
					prevCapapcityList[i].BigInt(),
					lib.BigRatRound(capacityDiff, false), // rounds down `capacity_diff`
				),
			)
		}
	}

	return newCapacityList
}
