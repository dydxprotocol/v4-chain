package lib

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"golang.org/x/exp/constraints"
)

func MustConvertIntegerToUint32[T constraints.Integer](value T) uint32 {
	if value < 0 {
		panic(fmt.Errorf("value %d underflows uint32", value))
	}
	if uint64(value) > math.MaxUint32 {
		panic(fmt.Errorf("value %d overflows uint32", value))
	}
	return uint32(value)
}

// ConvertBigFloatToUint64 attempts to convert a big Float into a uint64 and returns an error
// if the conversion would fail. Note: during conversion, will always round down.
func ConvertBigFloatToUint64(value *big.Float) (uint64, error) {
	if value.Cmp(BigFloatMaxUint64()) == 1 {
		return 0, errors.New("value overflows uint64")
	}

	if value.Sign() < 0 {
		return 0, errors.New("value underflows uint64")
	}

	uint64Val, _ := value.Uint64()
	return uint64Val, nil
}

// ConvertStringSliceToBigFloatSlice attempts to convert all values in a slice
// from string to big Float and return an error if any conversions fail.
func ConvertStringSliceToBigFloatSlice(values []string) ([]*big.Float, error) {
	convertedValues := make([]*big.Float, 0, len(values))
	for _, value := range values {
		bigValue, success := new(big.Float).SetString(value)
		if !success {
			return nil, fmt.Errorf("invalid, value is not a number: %v", value)
		}

		convertedValues = append(convertedValues, bigValue)
	}

	return convertedValues, nil
}

// ConvertBigFloatSliceToUint64Slice attempts to convert all values in a slice
// from big Float to uint64 and return an error if any conversions fail. Note: during conversion,
// will always round down.
func ConvertBigFloatSliceToUint64Slice(values []*big.Float) ([]uint64, error) {
	convertedValues := make([]uint64, 0, len(values))
	for _, value := range values {
		uint64Value, err := ConvertBigFloatToUint64(value)
		if err != nil {
			return nil, err
		}

		convertedValues = append(convertedValues, uint64Value)
	}

	return convertedValues, nil
}
