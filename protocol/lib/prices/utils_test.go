package prices_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/prices"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvert(t *testing.T) {
	tests := map[string]struct {
		price       uint64
		exponent    int32
		expected    uint64
		expectedErr error
	}{
		"Invert 1 = 1 (expected for USD-USDT)": {
			price:    10_000_000_000,
			exponent: -10,
			expected: 10_000_000_000,
		},
		"Invert 1 = 1 (small exponent)": {
			price:    1,
			exponent: 0,
			expected: 1,
		},
		"Invert 1 = 1 (large exponent)": {
			price:    10_000_000_000_000_000_000,
			exponent: -19,
			expected: 10_000_000_000_000_000_000,
		},
		"Invert 1.5 = 0.666666666": {
			price:    1_500_000_000,
			exponent: -9,
			expected: 666_666_666,
		},
		"Invert .015 = 66.666666": {
			price:    15_000_000_000,
			exponent: -12,
			expected: 66666666666666,
		},
		"Invert .5 = 2": {
			price:    500_000_000,
			exponent: -9,
			expected: 2_000_000_000,
		},
		"Error: Invert value < .001": {
			price:       1,
			exponent:    -3,
			expectedErr: fmt.Errorf("price 0.001 is outside of invertible range"),
		},
		"Error: Invert value > 100": {
			price:       1000,
			exponent:    0,
			expectedErr: fmt.Errorf("price 1000 is outside of invertible range"),
		},
		"Error: Invert overflows uint64 (due to >> exponent)": {
			price:       10000000000000000000,
			exponent:    -21,
			expectedErr: fmt.Errorf("inverted price overflows uint64"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := prices.Invert(test.price, test.exponent)
			if test.expectedErr != nil {
				require.Equal(t, uint64(0), actual)
				require.EqualError(t, err, test.expectedErr.Error())
			} else {
				require.Equal(t, test.expected, actual)
				require.NoError(t, err)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := map[string]struct {
		price            uint64
		exponent         types.Exponent
		adjustByPrice    uint64
		adjustByExponent types.Exponent
		expectedPrice    uint64
	}{
		"Large currency example: BTC-USD = BTC-USDT * USDT-USD": {
			// 1 BTC = 29,203.10 USDT
			price:    2_920_310_000,
			exponent: -5,
			// 1 USDT = 0.999765 USD
			adjustByPrice:    999_765_000,
			adjustByExponent: -9,
			// 1 BTC = $29,196.24
			expectedPrice: 2_919_623_727,
		},
		"Swap prices, same result": {
			// 1 USDT = 0.999765 USD
			price:    999_765_000,
			exponent: -9,
			// 1 BTC = 29,203.10 USDT
			adjustByPrice:    2_920_310_000,
			adjustByExponent: -5,
			// 1 BTC = $29,196.24, with -9 as exponent.
			expectedPrice: 29_196_237_271_500,
		},
		"Medium large currency example: ETH-USD = ETH-USDT * USDT-USD": {
			// 1 ETH = 1,862.41 USDT
			price:    1_862_410_000,
			exponent: -6,
			// 1 USDT = 0.999765 USD
			adjustByPrice:    999_765_000,
			adjustByExponent: -9,
			// 1 ETH = $1,861.972333
			expectedPrice: 1_861_972_334,
		},
		"Small currency example: 1INCH-USDT = 1INCH-USD * USD-USDT (two large exponents)": {
			// 1 1INCH = .3138 USDT
			price:    3_138_000_000,
			exponent: -10,
			// 1 USDT = 0.999765 USD
			adjustByPrice:    999_765_000,
			adjustByExponent: -9,
			// 1 1INCH = $0.313726257
			expectedPrice: 3_137_262_570,
		},
		"Micro currency example: XLM-USD = XLM-USDT * USDT-USD": {
			// 1 XLM = 0.1596 USDT
			price:    15_960_000_000,
			exponent: -11,
			// 1 USDT = 0.999765 USD
			adjustByPrice:    999_765_000,
			adjustByExponent: -9,
			// 1 XLM = $0.159562494
			expectedPrice: 15_956_249_400,
		},
		"Multiply rounds to nearest (round up)": {
			price:            1_000_000_000,
			exponent:         -9,
			adjustByPrice:    10_000_000_009, // this last digit will cause a round up
			adjustByExponent: -10,
			expectedPrice:    1_000_000_001,
		},
		"Multiply rounds to nearest (round down)": {
			price:            1_000_000_000,
			exponent:         -9,
			adjustByPrice:    10_000_000_001, // this last digit will cause a round down
			adjustByExponent: -10,
			expectedPrice:    1_000_000_000,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := prices.Multiply(tc.price, tc.exponent, tc.adjustByPrice, tc.adjustByExponent)
			require.Equal(t, tc.expectedPrice, actual)
		})
	}
}

func TestDivide(t *testing.T) {
	tests := map[string]struct {
		adjustByPrice    uint64
		adjustByExponent types.Exponent
		price            uint64
		exponent         types.Exponent
		expectedPrice    uint64
		expectedErr      error
	}{
		"Real example: USDT-USD = BTC-USD / BTC-USDT": {
			// 1 BTC-USD = $29,172.85
			adjustByPrice:    2_917_285_000,
			adjustByExponent: -5,
			// In practice, BTC-USDT would be represented with USDT market's exponent.
			// 1 BTC-USDT = 29,203.10 USDT
			price:    29_203_100_000_000,
			exponent: -9,
			// 1 USDT = .998964151 USD
			expectedPrice: 998_964_151,
		},
		"Real example: USDT-USD = ETH-USD / ETH-USDT": {
			// 1 ETH = $1,862.41
			adjustByPrice:    1_853_410_000,
			adjustByExponent: -6,
			// 1 ETH = 1,853.39 USDT, using USDT's exponent.
			price:         1_854_390_000_000,
			exponent:      -9,
			expectedPrice: 999_471_524,
		},
		"Edge case example w/smaller currency: USDT-USD = 1INCH-USD / 1INCH-USDT": {
			// 1 1INCH = .310997 USD
			adjustByPrice:    3_109_970_000,
			adjustByExponent: -10,
			// 1 1INCH = 0.3123 USDT, using USDT's exponent.
			price:         312_300_000,
			exponent:      -9,
			expectedPrice: 995_827_729,
		},
		"Failure: price < .01": {
			adjustByPrice:    2_917_285_000, // 1 BTC-USD = $29,172.85
			adjustByExponent: -5,
			price:            290_031_000_000, // 1 BTC-USDT = 290.03 USDT -- too small
			exponent:         -9,
			expectedErr: fmt.Errorf(
				"prices 290.031 and 29172.85 are too many orders of magnitude apart for accurate division",
			),
		},
		"Failure: price > 100": {
			price:            2_917_285_000, // 1 BTC-USD = $29,172.85
			exponent:         -5,
			adjustByPrice:    290_031_000_000, // 1 BTC-USDT = 290.03 USDT -- too small
			adjustByExponent: -9,
			expectedErr: fmt.Errorf(
				"prices 29172.85 and 290.031 are too many orders of magnitude apart for accurate division",
			),
		},
		"Failure: Divide by 0": {
			adjustByPrice:    1,
			adjustByExponent: 0,
			price:            0,
			exponent:         0,
			expectedErr:      fmt.Errorf("cannot divide by 0"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := prices.Divide(tc.adjustByPrice, tc.adjustByExponent, tc.price, tc.exponent)
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
				require.Equal(t, uint64(0), actual)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedPrice, actual)
			}
		})
	}
}

func TestPriceToFloat32ForLogging(t *testing.T) {
	tests := map[string]struct {
		price    uint64
		exponent types.Exponent
		expected float32
	}{
		"larger negative exp: BTC": {
			price:    2307961000,
			exponent: -5,
			expected: 23079.61,
		},
		"large negative exp: ETH": {
			price:    1853410000,
			exponent: -6,
			expected: 1853.41,
		},
		"median negative exp: LINK": {
			price:    751380000,
			exponent: -8,
			expected: 7.5138,
		},
		"small negative exp: 1INCH": {
			price:    3109970000,
			exponent: -10,
			expected: 0.310997,
		},
		"smaller negative exp: XLM": {
			price:    17263500000,
			exponent: -11,
			expected: 0.172635,
		},
		"positive exponent": {
			price:    48576,
			exponent: 5,
			expected: 4857600000,
		},
		"larger positive exponent": {
			price:    23,
			exponent: 10,
			expected: 230000000000,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := prices.PriceToFloat32ForLogging(tc.price, tc.exponent)
			require.Equal(t, tc.expected, actual)
		})
	}
}
