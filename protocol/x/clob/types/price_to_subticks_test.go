package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestPriceToSubticks(t *testing.T) {
	tests := map[string]struct {
		marketPrice           pricestypes.MarketPrice
		clobPair              types.ClobPair
		baseAtomicResolution  int32
		quoteAtomicResolution int32
		bigExpectedSubticks   *big.Rat
	}{
		"typical BTC configuration, at $10_000": {
			marketPrice: pricestypes.MarketPrice{
				Price:    1_000_000_000, // $10_000
				Exponent: -5,
			},
			clobPair: types.ClobPair{
				QuantumConversionExponent: -8,
			},
			baseAtomicResolution:  -10,
			quoteAtomicResolution: -6,
			bigExpectedSubticks:   big.NewRat(100_000_000, 1),
		},
		"typical ETH configuration, at $1_200": {
			marketPrice: pricestypes.MarketPrice{
				Price:    1_200_000_000, // $1_200
				Exponent: -6,
			},
			clobPair: types.ClobPair{
				QuantumConversionExponent: -9,
			},
			baseAtomicResolution:  -9,
			quoteAtomicResolution: -6,
			bigExpectedSubticks:   big.NewRat(1_200_000_000, 1),
		},
		"retains digits if not divisible": {
			marketPrice: pricestypes.MarketPrice{
				Price:    1_200_000_000, // $1_200
				Exponent: -6,
			},
			clobPair: types.ClobPair{
				QuantumConversionExponent: -9,
			},
			baseAtomicResolution:  -18,
			quoteAtomicResolution: -6,
			bigExpectedSubticks:   big.NewRat(12, 10),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			subticks := types.PriceToSubticks(
				tc.marketPrice,
				tc.clobPair,
				tc.baseAtomicResolution,
				tc.quoteAtomicResolution,
			)
			if tc.bigExpectedSubticks.Cmp(subticks) != 0 {
				t.Fatalf(
					"%s: bigExpectedSubticks: %s, subticks: %s",
					name,
					tc.bigExpectedSubticks.String(),
					subticks.String(),
				)
			}
		})
	}
}

func TestSubticksToPrice(t *testing.T) {
	tests := map[string]struct {
		subticks              types.Subticks
		clobPair              types.ClobPair
		baseAtomicResolution  int32
		quoteAtomicResolution int32
		expectedMarketPrice   pricestypes.MarketPrice
	}{
		"typical BTC configuration, at $10_000": {
			subticks: 100_000_000,
			clobPair: types.ClobPair{
				QuantumConversionExponent: -8,
			},
			baseAtomicResolution:  -10,
			quoteAtomicResolution: -6,
			expectedMarketPrice: pricestypes.MarketPrice{
				Price:    1_000_000_000, // $10_000
				Exponent: -5,
			},
		},
		"typical ETH configuration, at $1_200": {
			subticks: 1_200_000_000,
			clobPair: types.ClobPair{
				QuantumConversionExponent: -9,
			},
			baseAtomicResolution:  -9,
			quoteAtomicResolution: -6,
			expectedMarketPrice: pricestypes.MarketPrice{
				Price:    1_200_000_000, // $1_200
				Exponent: -6,
			},
		},
		"high base atomic resolution, at $1000": {
			subticks: 1,
			clobPair: types.ClobPair{
				QuantumConversionExponent: -9,
			},
			baseAtomicResolution:  -18,
			quoteAtomicResolution: -6,
			expectedMarketPrice: pricestypes.MarketPrice{
				Price:    1_000_000_000, // $1_000
				Exponent: -6,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			price := types.SubticksToPrice(
				tc.subticks,
				tc.expectedMarketPrice.Exponent,
				tc.clobPair,
				tc.baseAtomicResolution,
				tc.quoteAtomicResolution,
			)
			if tc.expectedMarketPrice.Price != price {
				t.Fatalf(
					"%s: expected market price: %+v, price: %+v",
					name,
					tc.expectedMarketPrice.Price,
					price,
				)
			}
		})
	}
}
