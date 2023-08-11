package lib

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"

	"golang.org/x/exp/constraints"
)

const (
	AvgInt32MaxArrayLength = 2 << 31
)

// DivisionUint32RoundUp returns the result of x/y, rounded up.
func DivisionUint32RoundUp(x, y uint32) uint32 {
	// Cast to uint64 so that equation below can't overflow.
	uint64X := uint64(x)
	uint64Y := uint64(y)
	result := (uint64X + uint64Y - 1) / uint64Y
	return uint32(result)
}

func Max[T constraints.Ordered](x, y T) T {
	if x < y {
		return y
	}
	return x
}

func Min[T constraints.Ordered](x, y T) T {
	if x > y {
		return y
	}
	return x
}

func MaxUint32(x, y uint32) uint32 {
	if x < y {
		return y
	}

	return x
}

func Int64MulPpm(x int64, ppm uint32) int64 {
	xMulPpm := BigIntMulPpm(big.NewInt(x), ppm)

	if !xMulPpm.IsInt64() {
		panic(fmt.Errorf("IntMulPpm (int = %d, ppm = %d) results in integer overflow", x, ppm))
	}

	return xMulPpm.Int64()
}

func AbsInt32(i int32) uint32 {
	if i < 0 {
		return uint32(0 - i)
	}

	return uint32(i)
}

func AbsInt64(i int64) uint64 {
	if i < 0 {
		return uint64(0 - i)
	}

	return uint64(i)
}

func AbsDiffUint64(x uint64, y uint64) uint64 {
	if x > y {
		return x - y
	}
	return y - x
}

// AvgInt32 returns average of the input int32 array. Result is rounded towards zero.
func AvgInt32(nums []int32) int32 {
	sum := int64(0)

	if len(nums) > AvgInt32MaxArrayLength {
		panic(fmt.Errorf(
			"input array to AvgInt32() exceeded maximum acceptable length (%d), got length = %d",
			AvgInt32MaxArrayLength,
			len(nums),
		))
	}

	for _, num := range nums {
		// For this sum to cause an int64 overflow, assuming each num is MaxInt32,
		// the length of nums would need to be ~(2^32).
		sum += int64(num)
	}

	avg := sum / int64(len(nums))

	if (avg > math.MaxInt32) || (avg < math.MinInt32) {
		panic(fmt.Errorf("result from AvgInt32 (%d) causes an int32 overflow", avg))
	}
	return int32(avg)
}

// ChangeRateUint64 returns the rate of change between the original and the new values.
// result = (new - original) / original
// Note: the return value is truncated to fit float32 precision.
func ChangeRateUint64(originalV uint64, newV uint64) (float32, error) {
	if originalV == 0 {
		return 0.0, errors.New("original value cannot be zero since we cannot divide by zero")
	}

	bigOriginalV := new(big.Float).SetUint64(originalV)
	bigNewV := new(big.Float).SetUint64(newV)

	diff := new(big.Float).Sub(bigNewV, bigOriginalV)
	diffRate := new(big.Float).Quo(
		diff,
		bigOriginalV,
	)

	result, _ := diffRate.Float32()
	return result, nil
}

// MedianUint64 returns the median value of the input slice. Note that if the
// input has an even number of elements, then the returned median is rounded up
// the nearest uint64. For example, 6.5 is rounded up to 7.
func MedianUint64(input []uint64) (uint64, error) {
	return medianIntGeneric(input)
}

// MedianInt32 returns the median value of the input slice. Note that if the
// input has an even number of elements, then the returned median is rounded
// towards positive/negative infinity, to the nearest int32. For example,
// 6.5 is rounded to 7 and -4.5 is rounded to -5.
func MedianInt32(input []int32) (int32, error) {
	return medianIntGeneric(input)
}

// MustGetMedianInt32 is a wrapper around `MedianInt32` that panics if
// input length is zero.
func MustGetMedianInt32(input []int32) int32 {
	ret, err := MedianInt32(input)

	if err != nil {
		panic(err)
	}

	return ret
}

// medianIntGeneric is a generic median calculator.
// It currently supports `uint64`, `int32` and more types can be added.
func medianIntGeneric[V uint64 | int32](input []V) (V, error) {
	l := len(input)
	if l == 0 {
		return 0, errors.New("input cannot be empty")
	}

	inputCopy := make([]V, l)
	copy(inputCopy, input)
	sort.Slice(inputCopy, func(i, j int) bool { return inputCopy[i] < inputCopy[j] })

	midIdx := l / 2

	if l%2 == 1 {
		return inputCopy[midIdx], nil
	}

	// The median is an average of the two middle numbers. It's rounded away from zero
	// to the nearest integer.
	// Note x <= y since `inputCopy` is sorted.
	x := inputCopy[midIdx-1]
	y := inputCopy[midIdx]

	if x <= 0 && y >= 0 {
		// x and y have different signs, so x+y cannot overflow.
		sum := x + y
		return sum/2 + sum%2, nil
	}

	if y > 0 {
		// x and y are both positive.
		return y - (y-x)/2, nil
	}

	// x and y are both negative.
	return x + (y-x)/2, nil
}
