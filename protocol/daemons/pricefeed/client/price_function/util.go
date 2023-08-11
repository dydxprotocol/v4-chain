package price_function

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"math/big"
	"strconv"

	"github.com/dydxprotocol/v4/lib"
)

var (
	apiResponseValidator *validator.Validate
)

// validatePositiveNumericString is a custom validation function that ensures a particular string field in
// a struct being validated can be parsed into a positive-valued float. We register this function in order
// to ensure that returned numeric string values in the Kraken response do not represent zero or negative numbers.
// To see where this is used, note the `validate:"positive-float-string"` struct tag in the KrakenTickerResult.
func validatePositiveNumericString(fl validator.FieldLevel) bool {
	val, err := strconv.ParseFloat(fl.Field().String(), 64)
	if err != nil {
		return false
	}
	return val > 0
}

// GetApiResponseValidator returns a validator with custom logic registered to validate fields returned by
// various exchange API responses.
func GetApiResponseValidator() (*validator.Validate, error) {
	if apiResponseValidator == nil {
		validate := validator.New()
		err := validate.RegisterValidation("positive-float-string", validatePositiveNumericString)
		if err != nil {
			return nil, fmt.Errorf("kraken API response validation internal error (%w)", err)
		}
		apiResponseValidator = validate
	}
	return apiResponseValidator, nil
}

// ExtractFirstStringFromSliceField takes a generic unmarshalled JSON object, interprets it
// as a slice of strings, and returns the first element.
func ExtractFirstStringFromSliceField(strToSliceMap map[string]interface{}, fieldName string) (string, error) {
	if strToSliceMap == nil {
		return "", fmt.Errorf("expected non-nil map")
	}
	slice, ok := strToSliceMap[fieldName].([]interface{})
	if !ok || len(slice) < 1 {
		return "", fmt.Errorf("expected non-empty list for fieldname '%v'", fieldName)
	}
	valString, ok := slice[0].(string)
	if !ok {
		return "", fmt.Errorf("expected nonempty string value for field %v[0], but found %v", fieldName, slice[0])
	}

	return valString, nil
}

// GetStringOrFloatValuesFromMap returns string or float64 values that correspond to the specified
// keys from a JSON map.
func GetStringOrFloatValuesFromMap[V string | float64](
	responseJson map[string]interface{},
	keys []string,
) ([]V, error) {
	jsonValues := make([]V, 0, len(keys))
	for _, key := range keys {
		value, ok := responseJson[key].(V)
		if !ok {
			return nil, fmt.Errorf(
				"Value was either not present or not valid for field: %v",
				key,
			)
		}

		jsonValues = append(jsonValues, value)
	}

	return jsonValues, nil
}

// GetInnerMapFromMap returns JSON map that correspond to the specific key from a JSON map.
func GetInnerMapFromMap(
	responseJson map[string]interface{},
	key string,
) (map[string]interface{}, error) {
	value, ok := responseJson[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"Value was either not present or not a valid JSON map for key: %v",
			key,
		)
	}

	return value, nil
}

// GetOnlyMarketSymbolAndExponent returns the only market symbol and exponent in the provided
// `marketPriceExponent` map. If the map contains more than one key, an error is returned.
func GetOnlyMarketSymbolAndExponent(
	marketPriceExponent map[string]int32,
	exchange string,
) (string, int32, error) {
	// Verify exactly one market is expected from the response.
	if len(marketPriceExponent) != 1 {
		return "", 0, fmt.Errorf(
			"Invalid market price exponent map for %v price function of length: %v, expected length 1",
			exchange,
			len(marketPriceExponent),
		)
	}

	// Get market symbol and value of exponent.
	var marketSymbol string
	var exponent int32
	// Set `marketSymbol` and `exponent` explicitly from the for loop.
	// marketPriceExponent has only one entry so the for loop only runs once.
	for marketSymbol, exponent = range marketPriceExponent {
	}

	return marketSymbol, exponent, nil
}

// GetUint64MedianFromReverseShiftedBigFloatValues shifts all values in a slice of floats by an
// exponent, converts the shifted slice values to uint64 and then returns the median of the slice.
// 1) Verify length of slice is > 0.
// 2) Reverse shift big float price values by the exponent for the market.
// 3) Convert big float values to uint64 values.
// 4) Get the median value from the uint64 price values and return.
func GetUint64MedianFromReverseShiftedBigFloatValues(
	// 1) Verify length of slice is > 0.
	bigFloatSlice []*big.Float,
	exponent int32,
	medianizer lib.Medianizer,

) (uint64, error) {
	if len(bigFloatSlice) == 0 {
		return 0, errors.New("Invalid input: big float slice must contain values to medianize")
	}

	// 2) Reverse shift big float price values by the exponent for the market.
	updatedBigFloatSlice := reverseShiftBigFloatSlice(bigFloatSlice, exponent)

	// 3) Convert big float values to uint64 values.
	uint64Slice, err := lib.ConvertBigFloatSliceToUint64Slice(updatedBigFloatSlice)
	if err != nil {
		return 0, err
	}

	// 4) Get the median value from the uint64 price values and return.
	return medianizer.MedianUint64(uint64Slice)
}

// ReverseShiftBigFloat shifts the given float by exponent in the reverse direction.
// If the exponent is 0, then do nothing (i.e. `123.456` -> `123.456`)
// If the exponent is positive, then shift to the right (i.e. exponent = 1, `123.456` -> `12.3456`)
// If the exponent is negative, then shift to the left (i.e. exponent = -1, `123.456` -> `1234.56`)
func ReverseShiftBigFloat(
	value *big.Float,
	exponent int32,
) *big.Float {
	unsignedExponent := lib.AbsInt32(exponent)

	pow10 := new(big.Float).SetInt(lib.BigPow10(uint64(unsignedExponent)))
	return reverseShiftFloatWithPow10(value, pow10, exponent)
}

// reverseShiftBigFloatSlice shifts the given floats by exponent in the reverse direction.
// If the exponent is 0, then do nothing (i.e. `123.456` -> `123.456`)
// If the exponent is positive, then shift to the right (i.e. exponent = 1, `123.456` -> `12.3456`)
// If the exponent is negative, then shift to the left (i.e. exponent = -1, `123.456` -> `1234.56`)
func reverseShiftBigFloatSlice(
	values []*big.Float,
	exponent int32,
) []*big.Float {
	unsignedExponent := lib.AbsInt32(exponent)

	pow10 := new(big.Float).SetInt(lib.BigPow10(uint64(unsignedExponent)))
	updatedValues := make([]*big.Float, 0, len(values))
	for _, value := range values {
		updatedValues = append(updatedValues, reverseShiftFloatWithPow10(value, pow10, exponent))
	}
	return updatedValues
}

func reverseShiftFloatWithPow10(value *big.Float, pow10 *big.Float, exponent int32) *big.Float {
	if exponent == 0 {
		return value
	} else if exponent > 0 {
		return new(big.Float).Quo(value, pow10)
	} else { // exponent < 0
		return new(big.Float).Mul(value, pow10)
	}
}
