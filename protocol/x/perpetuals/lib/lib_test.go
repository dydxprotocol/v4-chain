package lib_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
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
	}{
		"zero indexDelta": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(4_000_000_000),
			},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(4_000_000_000),
			expectedNetSettlementPpm: big.NewInt(0),
		},
		"positive indexDelta, positive quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(123_456_789_123),
			},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(-23_456_789_123_000_000),
		},
		"positive indexDelta, negative quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(123_456_789_123),
			},
			quantums:                 big.NewInt(-1_000_000),
			index:                    big.NewInt(100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(23_456_789_123_000_000),
		},
		"negative indexDelta, positive quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(-123_456_789_123),
			},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(-100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(23_456_789_123_000_000),
		},
		"negative indexDelta, negative quantums": {
			perpetual: types.Perpetual{
				FundingIndex: dtypes.NewInt(-123_456_789_123),
			},
			quantums:                 big.NewInt(-1_000_000),
			index:                    big.NewInt(-100_000_000_000),
			expectedNetSettlementPpm: big.NewInt(-23_456_789_123_000_000),
		},
		"is long, index went from negative to positive": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(100)},
			quantums:                 big.NewInt(30_000),
			index:                    big.NewInt(-100),
			expectedNetSettlementPpm: big.NewInt(-6_000_000),
		},
		"is long, index unchanged": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(100)},
			quantums:                 big.NewInt(1_000_000),
			index:                    big.NewInt(100),
			expectedNetSettlementPpm: big.NewInt(0),
		},
		"is long, index went from positive to zero": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(0)},
			quantums:                 big.NewInt(10_000_000),
			index:                    big.NewInt(100),
			expectedNetSettlementPpm: big.NewInt(1_000_000_000),
		},
		"is long, index went from positive to negative": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(-200)},
			quantums:                 big.NewInt(10_000_000),
			index:                    big.NewInt(100),
			expectedNetSettlementPpm: big.NewInt(3_000_000_000),
		},
		"is short, index went from negative to positive": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(100)},
			quantums:                 big.NewInt(-30_000),
			index:                    big.NewInt(-100),
			expectedNetSettlementPpm: big.NewInt(6_000_000),
		},
		"is short, index unchanged": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(100)},
			quantums:                 big.NewInt(-1_000),
			index:                    big.NewInt(100),
			expectedNetSettlementPpm: big.NewInt(0),
		},
		"is short, index went from positive to zero": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(0)},
			quantums:                 big.NewInt(-5_000_000),
			index:                    big.NewInt(100),
			expectedNetSettlementPpm: big.NewInt(-500_000_000),
		},
		"is short, index went from positive to negative": {
			perpetual:                types.Perpetual{FundingIndex: dtypes.NewInt(-50)},
			quantums:                 big.NewInt(-5_000_000),
			index:                    big.NewInt(100),
			expectedNetSettlementPpm: big.NewInt(-750_000_000),
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
			require.Equal(t, test.perpetual.FundingIndex.BigInt(), newFundingIndex)
		})
	}
}

func TestGetNetCollateralAndMarginRequirements(t *testing.T) {
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
	testLiquidityTier := types.LiquidityTier{
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 500_000,
	}
	tests := map[string]struct {
		perpetual     types.Perpetual
		marketPrice   pricestypes.MarketPrice
		liquidityTier types.LiquidityTier
		quantums      *big.Int
		quoteBalance  *big.Int
	}{
		"zero quantums": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(0),
			quoteBalance:  big.NewInt(0),
		},
		"positive quantums": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(1_000_000_000_000),
			quoteBalance:  big.NewInt(0),
		},
		"negative quantums": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(-1_000_000_000_000),
			quoteBalance:  big.NewInt(0),
		},
		"positive quote balance": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(-1_000_000_000_000),
			quoteBalance:  big.NewInt(1_000_000_000_000),
		},
		"negative quote balance": {
			perpetual:     testPerpetual,
			marketPrice:   testMarketPrice,
			liquidityTier: testLiquidityTier,
			quantums:      big.NewInt(1_000_000_000_000),
			quoteBalance:  big.NewInt(-1_000_000_000_000),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			enc := lib.GetNetNotionalInQuoteQuantums(
				test.perpetual,
				test.marketPrice,
				test.quantums,
			)
			eimr, emmr := lib.GetMarginRequirementsInQuoteQuantums(
				test.perpetual,
				test.marketPrice,
				test.liquidityTier,
				test.quantums,
				0,
			)
			risk := lib.GetNetCollateralAndMarginRequirements(
				test.perpetual,
				test.marketPrice,
				test.liquidityTier,
				test.quantums,
				test.quoteBalance,
				0,
			)
			require.Equal(t, 0, new(big.Int).Add(enc, test.quoteBalance).Cmp(risk.NC))
			require.Equal(t, eimr, risk.IMR)
			require.Equal(t, emmr, risk.MMR)
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
			0,
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
				0,
			)
			require.Equal(t, test.expectedImr, imr)
			require.Equal(t, test.expectedMmr, mmr)
		})
	}
}

func TestGetMarginRequirementsInQuoteQuantums_2(t *testing.T) {
	oneBip := math.Pow10(2)
	tests := map[string]struct {
		price                        uint64
		exponent                     int32
		baseCurrencyAtomicResolution int32
		bigBaseQuantums              *big.Int
		initialMarginPpm             uint32
		maintenanceFractionPpm       uint32
		openInterest                 *big.Int
		openInterestLowerCap         uint64
		openInterestUpperCap         uint64
		bigExpectedInitialMargin     *big.Int
		bigExpectedMaintenanceMargin *big.Int
	}{
		// TODO: Add back tests for positive and zero exponent once x/marketmap supports them
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, atomic resolution 4": {
			price:                        55_550,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -4,
			bigBaseQuantums:              big.NewInt(7_000),
			initialMarginPpm:             uint32(oneBip * 100),
			maintenanceFractionPpm:       uint32(500_000), // 50% of IM
			bigExpectedInitialMargin:     big.NewInt(38_885_000),
			bigExpectedMaintenanceMargin: big.NewInt(19_442_500),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, negative exponent, atomic resolution 6": {
			price:                        42_000_000,
			exponent:                     -2,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(-5_000),
			initialMarginPpm:             uint32(oneBip * 100),
			maintenanceFractionPpm:       uint32(500_000), // 50% of IM
			bigExpectedInitialMargin:     big.NewInt(21_000_000),
			bigExpectedMaintenanceMargin: big.NewInt(10_500_000),
		},
		"InitialMargin 10_000 BIPs (max), MaintenanceMargin 10_000 BIPs (max), atomic resolution 6": {
			price:                        55_550,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(7_000),
			initialMarginPpm:             uint32(oneBip * 10_000),
			maintenanceFractionPpm:       uint32(1_000_000), // 100% of IM
			bigExpectedInitialMargin:     big.NewInt(38_885_000),
			bigExpectedMaintenanceMargin: big.NewInt(38_885_000),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 100 BIPs, atomic resolution 6": {
			price:                        55_550,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(7_000),
			initialMarginPpm:             uint32(oneBip * 100),
			maintenanceFractionPpm:       uint32(1_000_000), // 100% of IM
			bigExpectedInitialMargin:     big.NewInt(388_850),
			bigExpectedMaintenanceMargin: big.NewInt(388_850),
		},
		"InitialMargin 0 BIPs (min), MaintenanceMargin 0 BIPs (min), atomic resolution 6": {
			price:                        55_550,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(7_000),
			initialMarginPpm:             uint32(oneBip * 0),
			maintenanceFractionPpm:       uint32(1_000_000), // 100% of IM,
			bigExpectedInitialMargin:     big.NewInt(0),
			bigExpectedMaintenanceMargin: big.NewInt(0),
		},
		"Price is zero, atomic resolution 6": {
			price:                        0,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(-7_000),
			initialMarginPpm:             uint32(oneBip * 1),
			maintenanceFractionPpm:       uint32(1_000_000), // 100% of IM,
			bigExpectedInitialMargin:     big.NewInt(0),
			bigExpectedMaintenanceMargin: big.NewInt(0),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              new(big.Int).SetUint64(math.MaxUint64),
			initialMarginPpm:             uint32(oneBip * 1),
			maintenanceFractionPpm:       uint32(1_000_000), // 100% of IM,
			bigExpectedInitialMargin: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843492", 10),
			),
			bigExpectedMaintenanceMargin: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843492", 10),
			),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, atomic resolution 6": {
			price:                        55_550,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(7_000),
			initialMarginPpm:             uint32(oneBip * 100),
			maintenanceFractionPpm:       uint32(500_000), // 50% of IM
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums / 1_000_000
			// = 10_000 * 38_885_000 / 1_000_000 ~= 388_850.
			bigExpectedInitialMargin:     big.NewInt(388_850),
			bigExpectedMaintenanceMargin: big.NewInt(388_850 / 2),
		},
		"InitialMargin 20%, MaintenanceMargin 10%, atomic resolution 6": {
			price:                        367_500,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(12_000),
			initialMarginPpm:             uint32(200_000),
			maintenanceFractionPpm:       uint32(500_000), // 50% of IM
			// quoteQuantums = 36_750 * 12_000 = 441_000_000
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums / 1_000_000
			// = 200_000 * 441_000_000 / 1_000_000 ~= 88_200_000
			bigExpectedInitialMargin:     big.NewInt(88_200_000),
			bigExpectedMaintenanceMargin: big.NewInt(88_200_000 / 2),
		},
		"InitialMargin 5%, MaintenanceMargin 3%, atomic resolution 6": {
			price:                        1_234_560,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(74_523),
			initialMarginPpm:             uint32(50_000),
			maintenanceFractionPpm:       uint32(600_000), // 60% of IM
			// quoteQuantums = 123_456 * 74_523 = 9_200_311_488
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums / 1_000_000
			// = 50_000 * 9_200_311_488 / 1_000_000 ~= 460_015_575
			bigExpectedInitialMargin:     big.NewInt(460_015_575),
			bigExpectedMaintenanceMargin: big.NewInt(276_009_345),
		},
		"InitialMargin 25%, MaintenanceMargin 15%, atomic resolution 6": {
			price:                        1_234_560,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(74_523),
			initialMarginPpm:             uint32(250_000),
			maintenanceFractionPpm:       uint32(600_000), // 60% of IM
			// quoteQuantums = 123_456 * 74_523 = 9_200_311_488
			bigExpectedInitialMargin:     big.NewInt(2_300_077_872),
			bigExpectedMaintenanceMargin: big.NewInt(1_380_046_724), // Rounded up
		},
		"OIMF: IM 20%, scaled to 60%, MaintenanceMargin 10%, atomic resolution 6": {
			price:                        367_500,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(12_000),
			initialMarginPpm:             uint32(200_000),
			maintenanceFractionPpm:       uint32(500_000),         // 50% of IM
			openInterest:                 big.NewInt(408_163_265), // 408.163265
			openInterestLowerCap:         10_000_000_000_000,
			openInterestUpperCap:         20_000_000_000_000,
			// openInterestNotional = 408_163_265 * 36_750 = 14_999_999_988_750
			// percentageOfCap = (openInterestNotional - lowerCap) / (upperCap - lowerCap) = 0.499999998875
			// adjustedIMF = (0.499999998875) * 0.8 + 0.2 = 0.5999999991 (rounded is 599_999 ppm)
			// bigExpectedInitialMargin = bigBaseQuantums * price * adjustedIMF = 264_599_559
			bigExpectedInitialMargin:     big.NewInt(264_599_559),
			bigExpectedMaintenanceMargin: big.NewInt(88_200_000 / 2),
		},
		"OIMF: IM 20%, scaled to 100%, MaintenanceMargin 10%, atomic resolution 6": {
			price:                        367_500,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(12_000),
			initialMarginPpm:             uint32(200_000),
			maintenanceFractionPpm:       uint32(500_000),           // 50% of IM
			openInterest:                 big.NewInt(1_000_000_000), // 1000 or ~$36mm notional
			openInterestLowerCap:         10_000_000_000_000,
			openInterestUpperCap:         20_000_000_000_000,
			// quoteQuantums = 36_750 * 12_000 = 441_000_000
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums / 1_000_000
			// = 200_000 * 441_000_000 / 1_000_000 ~= 88_200_000
			bigExpectedInitialMargin:     big.NewInt(441_000_000),
			bigExpectedMaintenanceMargin: big.NewInt(88_200_000 / 2),
		},
		"OIMF: IM 20%, lower_cap < realistic open interest < upper_cap, MaintenanceMargin 10%, atomic resolution 6": {
			price:                        367_500,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(12_000),
			initialMarginPpm:             uint32(200_000),
			maintenanceFractionPpm:       uint32(500_000),           // 50% of IM
			openInterest:                 big.NewInt(1_123_456_789), // 1123.456 or ~$41mm notional
			openInterestLowerCap:         25_000_000_000_000,
			openInterestUpperCap:         50_000_000_000_000,
			// openInterestNotional = 1_123_456_789 * 36_750 = 41_287_036_995_750
			// percentageOfCap = (openInterestNotional - lowerCap) / (upperCap - lowerCap) = 0.65148147983
			// adjustedIMF = (0.65148147983) * 0.8 + 0.2 = 0.721185183864 (rounded is 721_185 ppm)
			// bigExpectedInitialMargin = bigBaseQuantums * price * adjustedIMF = 318_042_585
			bigExpectedInitialMargin:     big.NewInt(318_042_585),
			bigExpectedMaintenanceMargin: big.NewInt(88_200_000 / 2),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Individual test setup.
			pc := keepertest.PerpetualsKeepers(t)
			// Create a new market param and price.
			marketId := keepertest.GetNumMarkets(t, pc.Ctx, pc.PricesKeeper)
			_, err := keepertest.CreateTestMarket(
				t,
				pc.Ctx,
				pc.PricesKeeper,
				pricestypes.MarketParam{
					Id:                 marketId,
					Pair:               "base-quote",
					Exponent:           tc.exponent,
					MinExchanges:       uint32(1),
					MinPriceChangePpm:  uint32(50),
					ExchangeConfigJson: "{}",
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    1_000, // leave this as a placeholder b/c we cannot set the price to 0
				},
			)
			require.NoError(t, err)

			// Update `Market.price`. By updating prices this way, we can simulate conditions where the oracle
			// price may become 0.
			err = pc.PricesKeeper.UpdateMarketPrices(
				pc.Ctx,
				[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{pricestypes.NewMarketPriceUpdate(
					marketId,
					tc.price,
				)},
			)
			require.NoError(t, err)

			// Create `LiquidityTier` struct.
			_, err = pc.PerpetualsKeeper.SetLiquidityTier(
				pc.Ctx,
				0,
				"name",
				tc.initialMarginPpm,
				tc.maintenanceFractionPpm,
				1, // dummy impact notional value
				tc.openInterestLowerCap,
				tc.openInterestUpperCap,
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := pc.PerpetualsKeeper.CreatePerpetual(
				pc.Ctx,
				0,                               // PerpetualId
				"getMarginRequirementsTicker",   // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
				types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS, // MarketType
			)
			require.NoError(t, err)

			// If test case contains non-nil open interest, set it up.
			if tc.openInterest != nil {
				require.NoError(t, pc.PerpetualsKeeper.ModifyOpenInterest(
					pc.Ctx,
					perpetual.Params.Id,
					tc.openInterest, // initialized as zero, so passing `openInterest` as delta amount.
				))
			}

			// Verify initial and maintenance margin requirements are calculated correctly.
			perpetual, marketPrice, liquidityTier, err := pc.PerpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(
				pc.Ctx,
				perpetual.Params.Id,
			)
			require.NoError(t, err)
			imr, mmr := lib.GetMarginRequirementsInQuoteQuantums(
				perpetual,
				marketPrice,
				liquidityTier,
				tc.bigBaseQuantums,
				0,
			)

			require.Equal(t, tc.bigExpectedInitialMargin, imr, "Initial margin mismatch")
			require.Equal(t, tc.bigExpectedMaintenanceMargin, mmr, "Maintenance margin mismatch")
		})
	}
}
