package constants

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var (
	// Subaccounts.
	Alice_Num0_1USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_1,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Alice_Num0_10_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: nil,
	}
	Alice_Num0_100_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: nil,
	}
	Alice_Num0_1BTC_LONG_10_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				1,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Alice_Num0_1ISO_LONG_10_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				3,
				big.NewInt(1_000_000_000), // 1 ISO
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Alice_Num1_10_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: nil,
	}
	Alice_Num1_100_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: nil,
	}
	Alice_Num1_1BTC_Short_100_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Alice_Num1_1BTC_Long_500_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // +1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Bob_Num0_1USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_1,
		},
		PerpetualPositions: nil,
	}
	Bob_Num0_10_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: nil,
	}
	Bob_Num0_50_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: nil,
	}
	Bob_Num0_100_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: nil,
	}
	Bob_Num0_1ISO_LONG_10_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				3,
				big.NewInt(1_000_000_000), // 1 ISO
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Bob_Num0_1ISO2_LONG_10_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				4,
				big.NewInt(10_000_000), // 1 ISO2
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_100BTC_Short_10100USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_100,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-10_000_000_000), // -100 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num1_1BTC_Short = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_49999USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(49_999_000_000)), // $49,999
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_50000USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(50_000_000_000)), // $50,000
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_50499USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(50_499_000_000), // $50,499
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_54999USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(54_999_000_000)), // $54,999
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Long_54999USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-54_999_000_000)), // -$54,999
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_55000USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(55_000_000_000)), // $55,000
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_100000USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(100_000_000_000)), // $100,000
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1BTC_Short_1ETH_Long_47000USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(47_000_000_000)), // $47,000
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
			testutil.CreateSinglePerpetualPosition(
				1,
				big.NewInt(1_000_000_000), // 1 ETH
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_1ISO_Short_49USD = satypes.Subaccount{
		Id:             &Carl_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(49_000_000)), // $49
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				3,
				big.NewInt(-1_000_000_000), // -1 ISO
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num0_599USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_599,
		},
		PerpetualPositions: nil,
	}
	Carl_Num0_660USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_660,
		},
		PerpetualPositions: nil,
	}
	Carl_Num0_10000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: nil,
	}
	Carl_Num0_50000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: nil,
	}
	Carl_Num0_100000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: nil,
	}
	Carl_Num0_500000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500_000,
		},
		PerpetualPositions: nil,
	}
	Carl_Num0_0USD = satypes.Subaccount{
		Id:                 &Carl_Num0,
		AssetPositions:     []*satypes.AssetPosition{},
		PerpetualPositions: nil,
	}
	Carl_Num1_500USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500,
		},
		PerpetualPositions: nil,
	}
	Carl_Num1_100000USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: nil,
	}
	Carl_Num1_Short_500USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_500,
		},
		PerpetualPositions: nil,
	}
	Carl_Num1_01BTC_Long_4600USD_Short = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_4_600,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(10_000_000), // 0.1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Carl_Num1_1BTC_Short_50499USD = satypes.Subaccount{
		Id:             &Carl_Num1,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(50_499_000_000)), // $50,499
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_01BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(10_000_000), // 0.1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_50001USD = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(50_001_000_000)), // $50,001
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Short_100000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000), // -1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_45000USD_Short = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-45_000_000_000)), // -$45,000
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_45001USD_Short = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-45_001_000_000)), // -$45,001
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_49501USD_Short = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-49_501_000_000)), // -$49,501
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_50000USD_Short = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-50_000_000_000)), // -$50,000
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_50001USD_Short = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-50_001_000_000)), // -$50,001
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1ISO_Long_50USD_Short = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-50_000_000)), // -$50
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				3,
				big.NewInt(1_000_000_000), // 1 ISO
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1ISO2_Short_499USD = satypes.Subaccount{
		Id:             &Dave_Num0,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(499_000_000)), // $499
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				4,
				big.NewInt(-10_000_000), // -1 ISO2
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_599USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_599,
		},
		PerpetualPositions: nil,
	}
	Dave_Num0_10000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: nil,
	}
	Dave_Num0_500000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500_000,
		},
		PerpetualPositions: nil,
	}
	Dave_Num0_100BTC_Short_10200USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_200,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-10_000_000_000), // -100 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_100BTC_Long_9900USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_9_900,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(10_000_000_000), // 100 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
			testutil.CreateSinglePerpetualPosition(
				1,
				big.NewInt(1_000_000_000), // 1 ETH
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num1_10_000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: nil,
	}
	Dave_Num1_500000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500_000,
		},
		PerpetualPositions: nil,
	}
	Dave_Num1_025BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(25_000_000), // 0.25 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num1_05BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(50_000_000), // 0.5 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num1_1BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num1_100BTC_Short_10100USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_100,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-10_000_000_000), // -100 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num1_1BTC_Long_49501USD_Short = satypes.Subaccount{
		Id:             &Dave_Num1,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(-49_501_000_000)), // -$49,501
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}
	Dave_Num1_1ETH_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				1,
				big.NewInt(1_000_000_000), // 1 ETH
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}

	// Quote balances.
	QuoteBalance_OneDollar = int64(1_000_000) // $1.

	InvalidSubaccountIdNumber = satypes.SubaccountId{
		Owner:  Alice_Num0.Owner,
		Number: satypes.MaxSubaccountIdNumber + 1, // one over the limit
	}
	InvalidSubaccountIdOwner = satypes.SubaccountId{
		Owner:  "This is not valid",
		Number: 0,
	}
)
