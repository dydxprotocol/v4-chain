package util

import (
	"math/big"
	"time"

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
	limiterCapacityList []types.LimiterCapacity,
	timeSinceLastBlock time.Duration,
) (
	newCapacityList []dtypes.SerializableInt,
) {
	// Declare new capacity list to be populated.
	newCapacityList = make([]dtypes.SerializableInt, len(limiterCapacityList))

	for i, limiterCapacity := range limiterCapacityList {
		limiter := limiterCapacity.Limiter
		capacity := limiterCapacity.Capacity.BigInt()
		baseline := GetBaseline(bigTvl, limiter)
		capacityMinusBaseline := new(big.Int).Sub(capacity, baseline)

		// Calculate the absolute value of the capacity delta.
		capacityDiff := lib.BigMax(baseline, capacityMinusBaseline)
		capacityDiff.Mul(capacityDiff, lib.BigI(timeSinceLastBlock.Milliseconds()))
		capacityDiff.Div(capacityDiff, lib.BigI(limiter.Period.Milliseconds()))

		// Move the capacity towards the baseline by capacityDiff. Do not cross the baseline.
		// Capacity is modified in-place if necessary.
		if new(big.Int).Abs(capacityMinusBaseline).Cmp(capacityDiff) <= 0 {
			capacity = baseline
		} else if capacityMinusBaseline.Sign() < 0 {
			capacity.Add(capacity, capacityDiff)
		} else {
			capacity.Sub(capacity, capacityDiff)
		}
		newCapacityList[i] = dtypes.NewIntFromBigInt(capacity)
	}

	return newCapacityList
}
