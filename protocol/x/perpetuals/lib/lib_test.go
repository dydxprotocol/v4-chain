package lib_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func BenchmarkGetSettlementPpmWithPerpetual(b *testing.B) {
	perpetual := types.Perpetual{
		FundingIndex: dtypes.NewInt(64_123_456_789),
	}
	quantums := big.NewInt(45_454_545_454)
	index := big.NewInt(87_654_321_99)
	for i := 0; i < b.N; i++ {
		lib.GetSettlementPpmWithPerpetual(
			perpetual,
			quantums,
			index,
		)
	}
}

func TestGetSettlementPpmWithPerpetual(t *testing.T) {
	tests := map[string]struct {
		perpetual                types.Perpetual
		quantums                 *big.Int
		index                    *big.Int
		expectedNetSettlementPpm *big.Int
		expectedNewFundingIndex  *big.Int
	}{
		"zero indexDelta": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(4_000_000_000),
			},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(4_000_000_000),
			expectedNetSettlementPpm: big.NewInt(0),
			expectedNewFundingIndex:  big.NewInt(4_000_000_000),
		},
		"positive indexDelta, positive quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(123_456_789_123),
			},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(-23_456_789_123_000_000),
			expectedNewFundingIndex:  big.NewInt(123_456_789_123),
		},
		"positive indexDelta, negative quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(123_456_789_123),
			},
			quantums:                 big.NewInt(-1_000_000),
			index:                    big.NewInt(100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(23_456_789_123_000_000),
			expectedNewFundingIndex:  big.NewInt(123_456_789_123),
		},
		"negative indexDelta, positive quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(-123_456_789_123),
			},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(-100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(23_456_789_123_000_000),
			expectedNewFundingIndex:  big.NewInt(-123_456_789_123),
		},
		"negative indexDelta, negative quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(-123_456_789_123),
			},
			quantums:                 big.NewInt(-1_000_000),
			index:                    big.NewInt(-100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(-23_456_789_123_000_000),
			expectedNewFundingIndex:  big.NewInt(-123_456_789_123),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			netSettlementPpm, newFundingIndex := lib.GetSettlementPpmWithPerpetual(
				test.perpetual,
				test.quantums,
				test.index,
			)
			require.Equal(t, test.expectedNetSettlementPpm, netSettlementPpm)
			require.Equal(t, test.expectedNewFundingIndex, newFundingIndex)
		})
	}
}

func BenchmarkGetNetNotionalInQuoteQuantums(b *testing.B) {
	perpetual := types.Perpetual{
		Params: types.PerpetualParams{
			AtomicResolution: 12,
		},
	}
	marketPrice := pricestypes.MarketPrice{
		Price:    123_456_789_123,
		Exponent: -8,
	}
	quantums := big.NewInt(987_654_321_321)
	for i := 0; i < b.N; i++ {
		lib.GetNetNotionalInQuoteQuantums(
			perpetual,
			marketPrice,
			quantums,
		)
	}
}

func TestGetNetNotionalInQuoteQuantums(t *testing.T) {
	testPerpetual := types.Perpetual{
		Params: types.PerpetualParams{
			AtomicResolution: -16,
		},
		OpenInterest: dtypes.NewInt(1_000_000_000_000),
	}
	testMarketPrice := pricestypes.MarketPrice{
		Price:    123_456_789_123,
		Exponent: -5,
	}
	tests := map[string]struct {
		perpetual           types.Perpetual
		marketPrice         pricestypes.MarketPrice
		quantums            *big.Int
		expectedNetNotional *big.Int
	}{
		"zero quantums": {
			perpetual:           testPerpetual,
			marketPrice:         testMarketPrice,
			quantums:            big.NewInt(0),
			expectedNetNotional: big.NewInt(1).SetUint64(0), // non-nil natural
		},
		"positive quantums": {
			perpetual:           testPerpetual,
			marketPrice:         testMarketPrice,
			quantums:            big.NewInt(1_000_000_000_000),
			expectedNetNotional: big.NewInt(123_456_789),
		},
		"negative quantums": {
			perpetual:           testPerpetual,
			marketPrice:         testMarketPrice,
			quantums:            big.NewInt(-1_000_000_000_000),
			expectedNetNotional: big.NewInt(-123_456_789),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			netNotional := lib.GetNetNotionalInQuoteQuantums(
				test.perpetual,
				test.marketPrice,
				test.quantums,
			)
			require.Equal(t, test.expectedNetNotional, netNotional)
		})
	}
}

func BenchmarkGetMarginRequirementsInQuoteQuantums(b *testing.B) {
	perpetual := types.Perpetual{
		Params: types.PerpetualParams{
			AtomicResolution: 8,
		},
		OpenInterest: dtypes.NewInt(1_000_000_000_000),
	}
	marketPrice := pricestypes.MarketPrice{
		Price:    123_456_789_123,
		Exponent: -5,
	}
	liquidityTier := types.LiquidityTier{
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 500_000,
	}
	quantums := big.NewInt(1_000_000_000_000)
	for i := 0; i < b.N; i++ {
		lib.GetMarginRequirementsInQuoteQuantums(
			perpetual,
			marketPrice,
			liquidityTier,
			quantums,
		)
	}
}

func TestGetMarginRequirementsInQuoteQuantums(t *testing.T) {
	testPerpetual := types.Perpetual{
		Params: types.PerpetualParams{
			AtomicResolution: 4,
		},
		OpenInterest: dtypes.NewInt(1_000_000_000_000),
	}
	testMarketPrice := pricestypes.MarketPrice{
		Price:    123_456_789_123,
		Exponent: -5,
	}
	testLiquidityTier := types.LiquidityTier{
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 500_000,
	}
	tests := map[string]struct {
		perpetual     types.Perpetual
		marketPrice   pricestypes.MarketPrice
		liquidityTier types.LiquidityTier
		quantums      *big.Int
		expectedImr   *big.Int
		expectedMmr   *big.Int
	}{
		"zero quantums": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(0),
			expectedImr:   big.NewInt(0),
			expectedMmr:   big.NewInt(0),
		},
		"positive quantums": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(1),
			expectedImr:   big_testutil.MustFirst(new(big.Int).SetString("2469135782460000", 10)),
			expectedMmr:   big_testutil.MustFirst(new(big.Int).SetString("1234567891230000", 10)),
		},
		"positive quantums, open interest above cap": {
			perpetual:   testPerpetual,
			marketPrice: testMarketPrice,
			liquidityTier: types.LiquidityTier{
				InitialMarginPpm:       200_000,
				MaintenanceFractionPpm: 500_000,
				OpenInterestLowerCap:   1_000_000_000_000,
				OpenInterestUpperCap:   2_000_000_000_000,
			},
			quantums:    big.NewInt(1),
			expectedImr: big_testutil.MustFirst(new(big.Int).SetString("12345678912300000", 10)),
			expectedMmr: big_testutil.MustFirst(new(big.Int).SetString("1234567891230000", 10)),
		},
		"negative quantums": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(-1),
			expectedImr:   big_testutil.MustFirst(new(big.Int).SetString("2469135782460000", 10)),
			expectedMmr:   big_testutil.MustFirst(new(big.Int).SetString("1234567891230000", 10)),
		},
		"negative quantums, open interest above cap": {
			perpetual:   testPerpetual,
			marketPrice: testMarketPrice,
			liquidityTier: types.LiquidityTier{
				InitialMarginPpm:       200_000,
				MaintenanceFractionPpm: 500_000,
				OpenInterestLowerCap:   1_000_000_000_000,
				OpenInterestUpperCap:   2_000_000_000_000,
			},
			quantums:    big.NewInt(-1),
			expectedImr: big_testutil.MustFirst(new(big.Int).SetString("12345678912300000", 10)),
			expectedMmr: big_testutil.MustFirst(new(big.Int).SetString("1234567891230000", 10)),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			imr, mmr := lib.GetMarginRequirementsInQuoteQuantums(
				test.perpetual,
				test.marketPrice,
				test.liquidityTier,
				test.quantums,
			)
			require.Equal(t, test.expectedImr, imr)
			require.Equal(t, test.expectedMmr, mmr)
		})
	}
}
