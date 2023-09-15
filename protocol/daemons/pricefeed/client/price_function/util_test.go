package price_function

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

const (
	ETHUSDC        = "ETHUSDC"
	BTCUSDC        = "BTCUSDC"
	deltaPrecision = 0.000000001
)

var (
	resultJson = map[string]interface{}{
		"lastPrice":      "100.002",
		"bidPrice":       "200.3",
		"askPrice":       "300.102",
		"nonStringValue": 1000,
	}
	positiveTagValidationError = errors.New(
		"Key: 'TestPositiveValidation.PositiveFloatString' Error:Field validation for 'PositiveFloatString' " +
			"failed on the 'positive-float-string' tag",
	)
)

func TestIsExchangeError_Mixed(t *testing.T) {
	tests := map[string]struct {
		err             error
		isExchangeError bool
	}{
		"Exchange Error - server sent GOAWAY": {
			err:             fmt.Errorf(`http2: server sent GOAWAY and closed the connection`),
			isExchangeError: true,
		},
		"Exchange Error - server sent GOAWAY with extra text": {
			err:             fmt.Errorf(`http2: server sent GOAWAY and closed the connection blah blah blah`),
			isExchangeError: true,
		},
		"Exchange Error - internal error": {
			err:             fmt.Errorf("internal error: something went wrong"),
			isExchangeError: true,
		},
		"Exchange Error - Internal error": {
			err:             fmt.Errorf("Internal error: something went wrong"),
			isExchangeError: true,
		},
		"Exchange Error - INTERNAL_ERROR": {
			err:             fmt.Errorf("INTERNAL_ERROR: something went wrong"),
			isExchangeError: true,
		},
		"Exchange Error - generic": {
			err:             fmt.Errorf("Unexpected response status code of: 5"),
			isExchangeError: true,
		},
		"Not exchange error": {
			err:             fmt.Errorf("some other error"),
			isExchangeError: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.isExchangeError, IsGenericExchangeError(tc.err))
		})
	}
}

func TestGetApiResponseValidator_validatePositiveNumericString_Mixed(t *testing.T) {
	tests := map[string]struct {
		testValue     string
		expectedError error
	}{
		"Success - canonical float": {
			testValue: "12345.6",
		},
		"Failure - negative float": {
			testValue:     "-12345.6",
			expectedError: positiveTagValidationError,
		},
		"Failure - empty string": {
			testValue:     "",
			expectedError: positiveTagValidationError,
		},
		"Failure - text": {
			testValue:     "cat",
			expectedError: positiveTagValidationError,
		},
	}

	type TestPositiveValidation struct {
		PositiveFloatString string `validate:"positive-float-string"`
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			validator, err := GetApiResponseValidator()
			require.Nil(t, err)
			err = validator.Struct(TestPositiveValidation{
				PositiveFloatString: tc.testValue,
			})
			if tc.expectedError == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
				require.EqualError(t, tc.expectedError, err.Error())
			}
		})
	}
}

func TestExtractFirstStringFromSliceField(t *testing.T) {
	tests := map[string]struct {
		inputJson      string
		nilInput       bool
		expectedResult string
		expectedError  error
	}{
		"Failure: nil input": {
			nilInput:       true,
			expectedResult: "",
			expectedError:  errors.New("expected non-nil map"),
		},
		"Failure: field does not exist": {
			inputJson:      "{}",
			expectedResult: "",
			expectedError:  errors.New("expected non-empty list for fieldname 'a'"),
		},
		"Failure: field is not a slice": {
			inputJson:      `{"a": 1}`,
			expectedResult: "",
			expectedError:  errors.New("expected non-empty list for fieldname 'a'"),
		},
		"Failure: field is an empty slice": {
			inputJson:      `{"a": []}`,
			expectedResult: "",
			expectedError:  errors.New("expected non-empty list for fieldname 'a'"),
		},
		"Failure: field is not a slice of strings": {
			inputJson:      `{"a": [12.3]}`,
			expectedResult: "",
			expectedError:  errors.New("expected nonempty string value for field a[0], but found 12.3"),
		},
		"Success: field is a slice containing 1 string": {
			inputJson:      `{"a": ["12.3"]}`,
			expectedResult: "12.3",
			expectedError:  nil,
		},
		"Success: field is a slice containing multiple strings": {
			inputJson:      `{"a": ["12.3", "34.5", "67.8"]}`,
			expectedResult: "12.3",
			expectedError:  nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var responseBody map[string]interface{}
			if !tc.nilInput {
				err := json.Unmarshal([]byte(tc.inputJson), &responseBody)
				require.NoError(t, err)
			}
			result, err := ExtractFirstStringFromSliceField(responseBody, "a")
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Zero(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetInnerMapFromMap(t *testing.T) {
	responseJson := map[string]interface{}{
		"result":    resultJson,
		"nonMapKey": true,
	}

	tests := map[string]struct {
		// parameters
		key string

		// expectations
		expectedResult map[string]interface{}
		expectedError  error
	}{
		"nonMapKey": {
			key:            "result",
			expectedResult: resultJson,
		},
		"Failure - missing field": {
			key: "unknown key",
			expectedError: fmt.Errorf(
				"Value was either not present or not a valid JSON map for key: %v",
				"unknown key",
			),
		},
		"Failure - value is not a JSON map": {
			key: "nonMapKey",
			expectedError: fmt.Errorf(
				"Value was either not present or not a valid JSON map for key: %v",
				"nonMapKey",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GetInnerMapFromMap(responseJson, tc.key)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestStringValuesFromMap(t *testing.T) {
	resultJsonWithStrings := map[string]interface{}{
		"lastPrice":  "100.002",
		"bidPrice":   "200.3",
		"askPrice":   "300.102",
		"floatValue": 1000.0001,
	}

	tests := map[string]struct {
		// parameters
		fields []string

		// expectations
		expectedResult []string
		expectedError  error
	}{
		"Success": {
			fields: []string{
				"lastPrice",
				"bidPrice",
				"askPrice",
			},
			expectedResult: []string{
				"100.002",
				"200.3",
				"300.102",
			},
		},
		"Failure - missing field": {
			fields: []string{
				"lastPrice",
				"bidPrice",
				"askPrice",
				"superSecretField",
			},
			expectedError: fmt.Errorf(
				"Value was either not present or not valid for field: %v",
				"superSecretField",
			),
		},
		"Failure - value is not a string": {
			fields: []string{
				"lastPrice",
				"bidPrice",
				"askPrice",
				"floatValue",
			},
			expectedError: fmt.Errorf(
				"Value was either not present or not valid for field: %v",
				"floatValue",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GetStringOrFloatValuesFromMap[string](resultJsonWithStrings, tc.fields)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetFloatValuesFromMap(t *testing.T) {
	resultJsonWithFloats := map[string]interface{}{
		"lastPrice":   100.002,
		"bidPrice":    200.3,
		"askPrice":    300.102,
		"stringValue": "1000.0001",
	}

	tests := map[string]struct {
		// parameters
		fields []string

		// expectations
		expectedResult []float64
		expectedError  error
	}{
		"Success": {
			fields: []string{
				"lastPrice",
				"bidPrice",
				"askPrice",
			},
			expectedResult: []float64{
				100.002,
				200.3,
				300.102,
			},
		},
		"Failure - missing field": {
			fields: []string{
				"lastPrice",
				"bidPrice",
				"askPrice",
				"superSecretField",
			},
			expectedError: fmt.Errorf(
				"Value was either not present or not valid for field: %v",
				"superSecretField",
			),
		},
		"Failure - value is not a float64": {
			fields: []string{
				"lastPrice",
				"bidPrice",
				"askPrice",
				"stringValue",
			},
			expectedError: fmt.Errorf(
				"Value was either not present or not valid for field: %v",
				"stringValue",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := GetStringOrFloatValuesFromMap[float64](resultJsonWithFloats, tc.fields)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetOnlyTickerAndExponent(t *testing.T) {
	tests := map[string]struct {
		//parameters
		tickerToExponent map[string]int32
		exchange         string

		// expectations
		expectedTicker   string
		expectedExponent int32
		expectedError    error
	}{
		"Success - isPositive = true and exchange = Binance": {
			tickerToExponent: map[string]int32{
				ETHUSDC: 6,
			},
			exchange:         exchange_common.EXCHANGE_NAME_BINANCE,
			expectedTicker:   ETHUSDC,
			expectedExponent: 6,
		},

		"Success - isNegative = false and exchange = Bitfinex": {
			tickerToExponent: map[string]int32{
				ETHUSDC: -6,
			},
			exchange:         exchange_common.EXCHANGE_NAME_BITFINEX,
			expectedTicker:   ETHUSDC,
			expectedExponent: -6,
		},
		"Failure - no exponents": {
			tickerToExponent: map[string]int32{},
			exchange:         exchange_common.EXCHANGE_NAME_BINANCE,
			expectedError: errors.New(
				"Invalid market price exponent map for Binance price function of length: 0, expected length 1",
			),
		},

		"Failure - too many exponents": {
			tickerToExponent: map[string]int32{
				ETHUSDC: -6,
				BTCUSDC: -8,
			},
			exchange: exchange_common.EXCHANGE_NAME_BITFINEX,
			expectedError: errors.New(
				"Invalid market price exponent map for Bitfinex price function of length: 2, expected length 1",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ticker,
				exponent,
				err := GetOnlyTickerAndExponent(
				tc.tickerToExponent,
				tc.exchange,
			)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())

				require.Equal(t, tc.expectedTicker, ticker)
				require.Equal(t, tc.expectedExponent, exponent)
			} else {
				require.NoError(t, err)

				require.Equal(t, tc.expectedTicker, ticker)
				require.Equal(t, tc.expectedExponent, exponent)
			}
		})
	}
}

func TestGetUint64MedianFromShiftedBigFloatValues(t *testing.T) {
	tests := map[string]struct {
		// parameters
		bigFloatSlice []*big.Float
		exponent      int32

		// expectations
		expectedMedianValue uint64
		expectedError       error
	}{
		"Success - isPositive = false": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100.001),
				new(big.Float).SetFloat64(300.001),
				new(big.Float).SetFloat64(200.022),
			},
			exponent:            -2,
			expectedMedianValue: uint64(200_02),
		},
		"Success - isPositive = true": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100.001),
				new(big.Float).SetFloat64(300.001),
				new(big.Float).SetFloat64(200.002),
			},
			exponent:            2,
			expectedMedianValue: uint64(2),
		},
		"Success - one value": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100.0002),
			},
			exponent:            0,
			expectedMedianValue: uint64(100),
		},
		"Failure - empty bigFloatSlice": {
			bigFloatSlice: []*big.Float{},
			exponent:      0,
			expectedError: errors.New(
				"Invalid input: big float slice must contain values to medianize",
			),
		},
		"Failure - underflow": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(-100),
			},
			exponent: 0,
			expectedError: errors.New(
				"value underflows uint64",
			),
		},
		"Failure - overflow": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100),
			},
			exponent: -1000,
			expectedError: errors.New(
				"value overflows uint64",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			median, err := GetUint64MedianFromReverseShiftedBigFloatValues(
				tc.bigFloatSlice,
				tc.exponent,
				lib.Median[uint64],
			)

			if tc.expectedError != nil {
				require.Equal(t, uint64(0), median)
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedMedianValue, median)
				require.NoError(t, err)
			}
		})
	}
}

func TestReverseShiftBigFloat(t *testing.T) {
	tests := map[string]struct {
		// parameters
		floatValue *big.Float
		exponent   int32

		// expectations
		expectedUpdatedFloatValue *big.Float
	}{
		"Success with negative exponent": {
			floatValue:                new(big.Float).SetPrec(64).SetFloat64(100.123),
			exponent:                  -3,
			expectedUpdatedFloatValue: new(big.Float).SetPrec(64).SetFloat64(100_123),
		},
		"Success with positive exponent": {
			floatValue:                new(big.Float).SetPrec(64).SetFloat64(100.1),
			exponent:                  1,
			expectedUpdatedFloatValue: new(big.Float).SetPrec(64).SetFloat64(10.01),
		},
		"Success with exponent of 0": {
			floatValue:                new(big.Float).SetPrec(64).SetFloat64(100),
			exponent:                  0,
			expectedUpdatedFloatValue: new(big.Float).SetPrec(64).SetFloat64(100),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			updatedFloatValue := ReverseShiftBigFloat(
				tc.floatValue,
				tc.exponent,
			)

			require.InDeltaSlice(
				t,
				bigSliceToFloatSlice([]*big.Float{tc.expectedUpdatedFloatValue}),
				bigSliceToFloatSlice([]*big.Float{updatedFloatValue}),
				deltaPrecision,
			)
		})
	}
}

func TestReverseShiftBigFloatWithPow10(t *testing.T) {
	tests := map[string]struct {
		// parameters
		floatValue *big.Float
		exponent   int32

		// expectations
		expectedUpdatedFloatValue *big.Float
	}{
		"Success with negative exponent": {
			floatValue:                new(big.Float).SetPrec(64).SetFloat64(100.123),
			exponent:                  -3,
			expectedUpdatedFloatValue: new(big.Float).SetPrec(64).SetFloat64(100_123),
		},
		"Success with positive exponent": {
			floatValue:                new(big.Float).SetPrec(64).SetFloat64(100.1),
			exponent:                  1,
			expectedUpdatedFloatValue: new(big.Float).SetPrec(64).SetFloat64(10.01),
		},
		"Success with exponent of 0": {
			floatValue:                new(big.Float).SetPrec(64).SetFloat64(100),
			exponent:                  0,
			expectedUpdatedFloatValue: new(big.Float).SetPrec(64).SetFloat64(100),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			unsignedExponent := lib.AbsInt32(tc.exponent)

			pow10 := new(big.Float).SetInt(lib.BigPow10(uint64(unsignedExponent)))

			updatedFloatValue := reverseShiftFloatWithPow10(
				tc.floatValue,
				pow10,
				tc.exponent,
			)

			require.InDeltaSlice(
				t,
				bigSliceToFloatSlice([]*big.Float{tc.expectedUpdatedFloatValue}),
				bigSliceToFloatSlice([]*big.Float{updatedFloatValue}),
				deltaPrecision,
			)
		})
	}
}

func TestReverseShiftBigFloatSlice(t *testing.T) {
	tests := map[string]struct {
		// parameters
		floatValues []*big.Float
		exponent    int32

		// expectations
		expectedUpdatedFloatValues []*big.Float
	}{
		"Success with empty floatValues": {
			floatValues:                []*big.Float{},
			exponent:                   -3,
			expectedUpdatedFloatValues: []*big.Float{},
		},
		"Success with negative exponent": {
			floatValues:                []*big.Float{new(big.Float).SetPrec(64).SetFloat64(100.123)},
			exponent:                   -3,
			expectedUpdatedFloatValues: []*big.Float{new(big.Float).SetPrec(64).SetFloat64(100_123)},
		},
		"Success with multiple values and a negative exponent": {
			floatValues: []*big.Float{new(big.Float).SetFloat64(100.122), new(big.Float).SetFloat64(2)},
			exponent:    -3,
			expectedUpdatedFloatValues: []*big.Float{
				new(big.Float).SetPrec(64).SetFloat64(100_122),
				new(big.Float).SetPrec(64).SetFloat64(2_000),
			},
		},
		"Success with positive exponent": {
			floatValues:                []*big.Float{new(big.Float).SetPrec(64).SetFloat64(100.1)},
			exponent:                   1,
			expectedUpdatedFloatValues: []*big.Float{new(big.Float).SetPrec(64).SetFloat64(10.01)},
		},
		"Success with multiple values and a positive exponent": {
			floatValues: []*big.Float{
				new(big.Float).SetPrec(64).SetFloat64(100),
				new(big.Float).SetPrec(64).SetFloat64(20),
			},
			exponent: 1,
			expectedUpdatedFloatValues: []*big.Float{
				new(big.Float).SetPrec(64).SetFloat64(10),
				new(big.Float).SetPrec(64).SetFloat64(2),
			},
		},
		"Success with exponent of 0": {
			floatValues:                []*big.Float{new(big.Float).SetPrec(64).SetFloat64(100)},
			exponent:                   0,
			expectedUpdatedFloatValues: []*big.Float{new(big.Float).SetPrec(64).SetFloat64(100)},
		},
		"Success with multiple values and an exponent of 0": {
			floatValues: []*big.Float{
				new(big.Float).SetPrec(64).SetFloat64(100.1),
				new(big.Float).SetPrec(64).SetFloat64(20.0000012),
			},
			exponent: 0,
			expectedUpdatedFloatValues: []*big.Float{
				new(big.Float).SetPrec(64).SetFloat64(100.1),
				new(big.Float).SetPrec(64).SetFloat64(20.0000012),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			updatedFloatValues := reverseShiftBigFloatSlice(
				tc.floatValues,
				tc.exponent,
			)

			require.InDeltaSlice(
				t,
				bigSliceToFloatSlice(tc.expectedUpdatedFloatValues),
				bigSliceToFloatSlice(updatedFloatValues),
				deltaPrecision,
			)
		})
	}
}

func TestConvertFloat64ToString(t *testing.T) {
	tests := map[string]struct {
		// parameters
		float64Value float64

		// expectations
		expectedFloat64String string
	}{
		"Success with low precision number": {
			float64Value:          float64(1.23),
			expectedFloat64String: "1.23",
		},
		"Success with a high precision number": {
			float64Value:          float64(0.12345678987654321),
			expectedFloat64String: "0.12345678987654321",
		},
		"Success with a large positive number": {
			float64Value:          float64(123456789.12345),
			expectedFloat64String: "123456789.12345",
		},
		"Success with a large negative number": {
			float64Value:          float64(-123456789.12345),
			expectedFloat64String: "-123456789.12345",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			float64String := ConvertFloat64ToString(tc.float64Value)

			require.Equal(t, tc.expectedFloat64String, float64String)
		})
	}
}

func bigSliceToFloatSlice(bigFloat []*big.Float) []float64 {
	floatSlice := make([]float64, 0, len(bigFloat))
	for _, val := range bigFloat {
		floatVal, _ := val.Float64()
		floatSlice = append(floatSlice, floatVal)
	}

	return floatSlice
}
