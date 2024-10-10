package constants

import (
	"math"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

var (
	// Perpetual Positions.
	Long_Perp_1BTC_PositiveFunding = satypes.PerpetualPosition{
		PerpetualId:  0,
		Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
		FundingIndex: dtypes.NewInt(0),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	Short_Perp_1ETH_NegativeFunding = satypes.PerpetualPosition{
		PerpetualId:  1,
		Quantums:     dtypes.NewInt(-100_000_000), // 1 ETH
		FundingIndex: dtypes.NewInt(-1),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(100_000_000), // 1 BTC, $50,000 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneBTCShort = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(-100_000_000), // 1 BTC, -$50,000 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneTenthBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneTenthBTCShort = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(-10_000_000), // 0.1 BTC, -$5,000 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneHundredthBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(1_000_000), // 0.01 BTC, $500 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneAndHalfBTCLong = satypes.PerpetualPosition{
		PerpetualId:  0,
		Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC, $75,000 notional.
		FundingIndex: dtypes.NewInt(0),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	PerpetualPosition_FourThousandthsBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(400_000), // 0.004 BTC, $200 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_FourThousandthsBTCShort = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(-400_000), // 0.004 BTC, -$200 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneTenthEthLong = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewInt(100_000_000), // 0.1 ETH, $300 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneTenthEthShort = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewInt(-100_000_000), // 0.1 ETH, -$300 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_MaxUint64EthLong = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewIntFromUint64(math.MaxUint64), // 18,446,744,070 ETH, $55,340,232,210,000 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	PerpetualPosition_MaxUint64EthShort = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewIntFromBigInt(BigNegMaxUint64()), // 18,446,744,070 ETH, -$55,340,232,210,000 notional.
		YieldIndex:  big.NewRat(0, 1).String(),
	}
	// Long position for arbitrary isolated market
	PerpetualPosition_OneISOLong = satypes.PerpetualPosition{
		PerpetualId:  3,
		Quantums:     dtypes.NewInt(1_000_000_000),
		FundingIndex: dtypes.NewInt(0),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneISO2Long = satypes.PerpetualPosition{
		PerpetualId:  4,
		Quantums:     dtypes.NewInt(10_000_000),
		FundingIndex: dtypes.NewInt(0),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	// Short position for arbitrary isolated market
	PerpetualPosition_OneISOShort = satypes.PerpetualPosition{
		PerpetualId:  3,
		Quantums:     dtypes.NewInt(-100_000_000),
		FundingIndex: dtypes.NewInt(0),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	PerpetualPosition_OneISO2Short = satypes.PerpetualPosition{
		PerpetualId:  4,
		Quantums:     dtypes.NewInt(-10_000_000),
		FundingIndex: dtypes.NewInt(0),
		YieldIndex:   big.NewRat(0, 1).String(),
	}
	// Asset Positions
	TDai_Asset_0 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(0), // $0
	}
	TDai_Asset_500 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(500_000_000), // $500
	}
	Short_TDai_Asset_500 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-500_000_000), // -$500
	}
	TDai_Asset_599 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(599_000_000), // $599
	}
	TDai_Asset_660 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(660_000_000), // $660
	}
	Short_TDai_Asset_4_600 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-4_600_000_000), // -$4,600
	}
	Short_TDai_Asset_2_900 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-2_900_000_000), // -$2,900
	}
	Short_TDai_Asset_46_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-46_000_000_000), // -$46,000
	}
	Short_TDai_Asset_49_500 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-49_500_000_000), // -$49,500
	}
	Short_TDai_Asset_9_900 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-9_900_000_000), // $-9,900
	}
	TDai_Asset_1 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(1_000_000), // $1
	}
	TDai_Asset_10_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(10_000_000_000), // $10,000
	}
	TDai_Asset_10_100 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(10_100_000_000), // $10,100
	}
	TDai_Asset_10_200 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(10_200_000_000), // $10,200
	}
	TDai_Asset_50_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(50_000_000_000), // $50,000
	}
	TDai_Asset_99_999 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(99_999_000_000), // $99,999
	}
	TDai_Asset_100_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(100_000_000_000), // $100,000
	}
	TDai_Asset_100_499 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(100_499_000_000), // $100,499
	}
	TDai_Asset_500_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(500_000_000_000), // $500,000
	}
	Long_Asset_1BTC = satypes.AssetPosition{
		AssetId:  1,
		Quantums: dtypes.NewInt(100_000_000), // 1 BTC
	}
	Short_Asset_1BTC = satypes.AssetPosition{
		AssetId:  1,
		Quantums: dtypes.NewInt(-100_000_000), // 1 BTC
	}
	Long_Asset_1ETH = satypes.AssetPosition{
		AssetId:  2,
		Quantums: dtypes.NewInt(1_000_000_000), // 1 ETH
	}
	Short_Asset_1ETH = satypes.AssetPosition{
		AssetId:  2,
		Quantums: dtypes.NewInt(-1_000_000_000), // 1 ETH
	}
)
