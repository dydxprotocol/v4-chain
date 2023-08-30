package sim_helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// RandSliceShuffle returns a shuffled version of the input slice.
func RandSliceShuffle[T any](r *rand.Rand, s []T) []T {
	ret := make([]T, len(s))
	perm := r.Perm(len(s))
	for i, randIndex := range perm {
		ret[i] = s[randIndex]
	}
	return ret
}

// RandWithWeight takes in a map of options, each with a weight and returns
// each element with a probability (weight/totalWeight).
// Panics if weight is negative.
func RandWithWeight[T comparable](r *rand.Rand, optionsWithWeight map[T]int) T {
	sum := 0
	for _, weight := range optionsWithWeight {
		if weight < 0 {
			panic(fmt.Errorf("RandWithWeight weights cannot be negative"))
		}
		sum += weight
	}
	// generate random number from [0, sumOfWeights)
	randNum := int(r.Int63n(int64(sum)))
	for option, weight := range optionsWithWeight {
		randNum -= weight
		if randNum < 0 {
			return option
		}
	}
	panic(fmt.Errorf("error computing RandWithWeight"))
}

// RandBool randomly returns true or false where each return value has 50% chance.
func RandBool(r *rand.Rand) bool {
	// Note: the current impl of RandIntBetween is not inclusive, because it relies on golang's `rand.Intn`
	// which is half-open interval.
	return simtypes.RandIntBetween(r, 0, 2) == 0
}

// GetRandomBucketValue randomly selects a min/max based on the given bucket boundaries and returns a random
// value within the chosen min/max. Each bucket has an equal probability of getting selected.
func GetRandomBucketValue(r *rand.Rand, bucketBoundaries []int) int {
	numBuckets := len(bucketBoundaries) - 1
	if numBuckets < 2 {
		panic(errors.New("you must supply at least 3 or more bucket boundary values"))
	}
	for i := 0; i < numBuckets; i++ {
		if bucketBoundaries[i+1] <= bucketBoundaries[i] {
			panic(errors.New("bucket boundary values must be strictly increasing"))
		}
	}

	bucketLowerBoundIdx := simtypes.RandIntBetween(r, 0, numBuckets)
	return simtypes.RandIntBetween(
		r,
		bucketBoundaries[bucketLowerBoundIdx],
		bucketBoundaries[bucketLowerBoundIdx+1],
	)
}

// RandPositiveUint64 returns a randomized number between 1 and the maximum `uint64`, inclusive.
func RandPositiveUint64(r *rand.Rand) uint64 {
	n := uint64(0)
	// While `n` is 0, keep generating random `uint64` values.
	for n == 0 {
		n = r.Uint64()
	}
	return n
}

// MakeRange returns a slice of uint32 where the values are sequential.
func MakeRange(n uint32) []uint32 {
	result := make([]uint32, n)
	for i := uint32(0); i < n; i++ {
		result[i] = i
	}
	return result
}

// GetMin returns minimum of the two inputs.
func GetMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// ShouldGenerateReasonableGenesis determines whether or not the Genesis parameters should be
// reasonable based on the ReasonableGenesisWeight (0-100). Weight is deterministic on genisisTime,
// meaning all modules should generate the same result for a simulation run.
func ShouldGenerateReasonableGenesis(r *rand.Rand, genesisTime time.Time) bool {
	r.Seed(genesisTime.UnixNano())
	weight := simtypes.RandIntBetween(r, 0, 100)
	return ReasonableGenesisWeight > weight
}

// PickGenesisParameter picks the Reasonable genesis parameter if `shouldUseReasonable = true`,
// and the Valid parameter otherwise.
func PickGenesisParameter[T any](
	param GenesisParameters[T],
	shouldUseReasonable bool,
) T {
	if shouldUseReasonable {
		return param.Reasonable
	}
	return param.Valid
}
