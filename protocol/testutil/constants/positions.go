package constants

import (
	"math"
	"math/big"

	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
)

var (
	// Perpetual Positions.
	Long_Perp_1BTC_PositiveFunding = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(100_000_000), // 1 BTC
		big.NewInt(0),
		big.NewInt(0),
	)
	Short_Perp_1ETH_NegativeFunding = *testutil.CreateSinglePerpetualPosition(
		1,
		big.NewInt(-100_000_000), // 1 ETH
		big.NewInt(-1),
		big.NewInt(0),
	)
	PerpetualPosition_OneBTCLong = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(100_000_000), // 1 BTC, $50,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneBTCShort = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(-100_000_000), // 1 BTC, -$50,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneTenthBTCLong = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneTenthBTCShort = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(-10_000_000), // 0.1 BTC, -$5,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneHundredthBTCLong = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(1_000_000), // 0.01 BTC, $500 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneHundredthBTCShort = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(-1_000_000), // 0.01 BTC, -$500 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_FourThousandthsBTCLong = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(400_000), // 0.004 BTC, $200 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_FourThousandthsBTCShort = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(-400_000), // 0.004 BTC, -$200 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneAndHalfBTCLong = *testutil.CreateSinglePerpetualPosition(
		0,
		big.NewInt(150_000_000), // 1.5 BTC, $75,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneTenthEthLong = *testutil.CreateSinglePerpetualPosition(
		1,
		big.NewInt(100_000_000), // 0.1 ETH, $300 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneTenthEthShort = *testutil.CreateSinglePerpetualPosition(
		1,
		big.NewInt(-100_000_000), // 0.1 ETH, -$300 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_MaxUint64EthLong = *testutil.CreateSinglePerpetualPosition(
		1,
		big.NewInt(0).SetUint64(math.MaxUint64), // 18,446,744,070 ETH, $55,340,232,210,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_MaxUint64EthShort = *testutil.CreateSinglePerpetualPosition(
		1,
		BigNegMaxUint64(), // 18,446,744,070 ETH, -$55,340,232,210,000 notional.
		big.NewInt(0),
		big.NewInt(0),
	)
	// SOL positions
	PerpetualPosition_OneSolLong = *testutil.CreateSinglePerpetualPosition(
		2,
		big.NewInt(100_000_000_000), // 1 SOL
		big.NewInt(0),
		big.NewInt(0),
	)
	// Long position for arbitrary isolated market
	PerpetualPosition_OneISOLong = *testutil.CreateSinglePerpetualPosition(
		3,
		big.NewInt(1_000_000_000),
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneISO2Long = *testutil.CreateSinglePerpetualPosition(
		4,
		big.NewInt(10_000_000),
		big.NewInt(0),
		big.NewInt(0),
	)
	// Short position for arbitrary isolated market
	PerpetualPosition_OneISOShort = *testutil.CreateSinglePerpetualPosition(
		3,
		big.NewInt(-100_000_000),
		big.NewInt(0),
		big.NewInt(0),
	)
	PerpetualPosition_OneISO2Short = *testutil.CreateSinglePerpetualPosition(
		4,
		big.NewInt(-10_000_000),
		big.NewInt(0),
		big.NewInt(0),
	)

	// Asset Positions
	Usdc_Asset_0 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(0), // $0
	)
	Usdc_Asset_1 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(1_000_000), // $1
	)
	Usdc_Asset_500 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(500_000_000), // $500
	)
	Short_Usdc_Asset_500 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(-500_000_000), // -$500
	)
	Usdc_Asset_599 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(599_000_000), // $599
	)
	Usdc_Asset_660 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(660_000_000), // $660
	)
	Short_Usdc_Asset_4_600 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(-4_600_000_000), // -$4,600
	)
	Short_Usdc_Asset_46_000 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(-46_000_000_000), // -$46,000
	)
	Short_Usdc_Asset_9_900 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(-9_900_000_000), // $-9,900
	)
	Usdc_Asset_10_000 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(10_000_000_000), // $10,000
	)
	Usdc_Asset_10_100 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(10_100_000_000), // $10,100
	)
	Usdc_Asset_10_200 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(10_200_000_000), // $10,200
	)
	Usdc_Asset_50_000 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(50_000_000_000), // $50,000
	)
	Usdc_Asset_99_999 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(99_999_000_000), // $99,999
	)
	Usdc_Asset_100_000 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(100_000_000_000), // $100,000
	)
	Usdc_Asset_100_499 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(100_499_000_000), // $100,499
	)
	Usdc_Asset_500_000 = *testutil.CreateSingleAssetPosition(
		0,
		big.NewInt(500_000_000_000), // $500,000
	)
	Long_Asset_1BTC = *testutil.CreateSingleAssetPosition(
		1,
		big.NewInt(100_000_000), // 1 BTC
	)
	Short_Asset_1BTC = *testutil.CreateSingleAssetPosition(
		1,
		big.NewInt(-100_000_000), // 1 BTC
	)
	Long_Asset_1ETH = *testutil.CreateSingleAssetPosition(
		2,
		big.NewInt(1_000_000_000), // 1 ETH
	)
	Short_Asset_1ETH = *testutil.CreateSingleAssetPosition(
		2,
		big.NewInt(-1_000_000_000), // 1 ETH
	)
)
