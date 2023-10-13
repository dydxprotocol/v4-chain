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

// Uint64LinearInterpolate interpolates value v0 towards v1 by a small constant value c, typically expected to
// be between 0 and 1. Here, the input value of c is represented in ppm. In order to avoid overflows, if
// 0 <= cPpm <= 1_000_000 then an error is returned.
func Uint64LinearInterpolate(v0 uint64, v1 uint64, cPpm uint32) (uint64, error) {
	if cPpm > OneMillion {
		return 0, fmt.Errorf("uint64 interpolation requires 0 <= cPpm <= 1_000_000, but received cPpm value of %v", cPpm)
	}
	// Note: Uint64MulPpm panics if the multiplication overflows an int64, but we've already prevented that from
	// happening by checking that cPpm <= 1_000_000.
	absDelta := Uint64MulPpm(AbsDiffUint64(v0, v1), cPpm)
	if v0 > v1 {
		return v0 - absDelta, nil
	} else {
		return v0 + absDelta, nil
	}
}

// AddUint32 returns the sum of a and b. If the sum underflows or overflows, this method returns an error.
func AddUint32(a int64, b uint32) (int64, error) {
	sum := a + int64(b)

	// This check should catch a + b overflows.
	if sum < a {
		return 0, fmt.Errorf("int64 overflow: %d + %d", a, b)
	}
	return sum, nil
}

// MustDivideUint32RoundUp returns the result of x/y, rounded up.
// Note: this method will panic if y == 0.
func MustDivideUint32RoundUp(x, y uint32) uint32 {
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

// Int64MulPpm multiplies an int64 by a scaling factor represented in parts per million. If the integer overflows,
// this method panics. This method rounds towards negative infinity.
func Int64MulPpm(x int64, ppm uint32) int64 {
	xMulPpm := BigIntMulPpm(big.NewInt(x), ppm)

	if !xMulPpm.IsInt64() {
		panic(fmt.Errorf("IntMulPpm (int = %d, ppm = %d) results in integer overflow", x, ppm))
	}

	return xMulPpm.Int64()
}

// Uint64MulPpm multiplies a uint64 value by a scaling factor represented in parts per million. If the integer
// overflows, this method panics.
func Uint64MulPpm(x uint64, ppm uint32) uint64 {
	xMulPpm := BigIntMulPpm(new(big.Int).SetUint64(x), ppm)

	if !xMulPpm.IsUint64() {
		panic(fmt.Errorf("UintMulPpm (uint = %d, ppm = %d) results in integer overflow", x, ppm))
	}

	return xMulPpm.Uint64()
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

// AvgInt32 returns average of the input int32 array. Result is rounded towards zero. Note: this method panics if
// the input array length exceeds AvgInt32MaxArrayLength, or if the result causes an int32 overflow.
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

// MustGetMedian is a wrapper around `Median` that panics if input length is zero.
func MustGetMedian[V uint64 | uint32 | int64 | int32](input []V) V {
	ret, err := Median(input)

	if err != nil {
		panic(err)
	}

	return ret
}

// Median is a generic median calculator.
// If the input has an even number of elements, then the average of the two middle numbers is rounded away from zero.
func Median[V uint64 | uint32 | int64 | int32](input []V) (V, error) {
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
