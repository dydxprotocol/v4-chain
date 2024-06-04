package constants

import (
	"math"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

var (
	// Perpetual Positions.
	Long_Perp_1BTC_PositiveFunding = satypes.PerpetualPosition{
		PerpetualId:  0,
		Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
		FundingIndex: dtypes.NewInt(0),
	}
	Short_Perp_1ETH_NegativeFunding = satypes.PerpetualPosition{
		PerpetualId:  1,
		Quantums:     dtypes.NewInt(-100_000_000), // 1 ETH
		FundingIndex: dtypes.NewInt(-1),
	}
	PerpetualPosition_OneBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(100_000_000), // 1 BTC, $50,000 notional.
	}
	PerpetualPosition_OneBTCShort = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(-100_000_000), // 1 BTC, -$50,000 notional.
	}
	PerpetualPosition_OneTenthBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
	}
	PerpetualPosition_OneTenthBTCShort = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(-10_000_000), // 0.1 BTC, -$5,000 notional.
	}
	PerpetualPosition_FourThousandthsBTCLong = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(400_000), // 0.004 BTC, $200 notional.
	}
	PerpetualPosition_FourThousandthsBTCShort = satypes.PerpetualPosition{
		PerpetualId: 0,
		Quantums:    dtypes.NewInt(-400_000), // 0.004 BTC, -$200 notional.
	}
	PerpetualPosition_OneTenthEthLong = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewInt(100_000_000), // 0.1 ETH, $300 notional.
	}
	PerpetualPosition_OneTenthEthShort = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewInt(-100_000_000), // 0.1 ETH, -$300 notional.
	}
	PerpetualPosition_MaxUint64EthLong = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewIntFromUint64(math.MaxUint64), // 18,446,744,070 ETH, $55,340,232,210,000 notional.
	}
	PerpetualPosition_MaxUint64EthShort = satypes.PerpetualPosition{
		PerpetualId: 1,
		Quantums:    dtypes.NewIntFromBigInt(BigNegMaxUint64()), // 18,446,744,070 ETH, -$55,340,232,210,000 notional.
	}
	// Asset Positions
	Usdc_Asset_0 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(0), // $0
	}
	Usdc_Asset_500 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(500_000_000), // $500
	}
	Short_Usdc_Asset_500 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-500_000_000), // -$500
	}
	Usdc_Asset_599 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(599_000_000), // $599
	}
	Usdc_Asset_660 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(660_000_000), // $660
	}
	Short_Usdc_Asset_4_600 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-4_600_000_000), // -$4,600
	}
	Short_Usdc_Asset_46_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-46_000_000_000), // -$46,000
	}
	Short_Usdc_Asset_9_900 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(-9_900_000_000), // $-9,900
	}
	Usdc_Asset_1 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(1_000_000), // $1
	}
	Usdc_Asset_10_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(10_000_000_000), // $10,000
	}
	Usdc_Asset_10_100 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(10_100_000_000), // $10,100
	}
	Usdc_Asset_10_200 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(10_200_000_000), // $10,200
	}
	Usdc_Asset_50_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(50_000_000_000), // $50,000
	}
	Usdc_Asset_99_999 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(99_999_000_000), // $99,999
	}
	Usdc_Asset_100_000 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(100_000_000_000), // $100,000
	}
	Usdc_Asset_100_499 = satypes.AssetPosition{
		AssetId:  0,
		Quantums: dtypes.NewInt(100_499_000_000), // $100,499
	}
	Usdc_Asset_500_000 = satypes.AssetPosition{
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
