package constants

import (
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

var (
	// Subaccounts.
	Alice_Num0_1USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_1,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Alice_Num0_10_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Alice_Num0_100_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Alice_Num1_10_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Alice_Num1_100_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Alice_Num0_1ISO_LONG_10_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  3,
				Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Alice_Num1_1BTC_Short_100_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Alice_Num1_1BTC_Long_500_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // +1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Bob_Num0_10_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Bob_Num0_100_000USD = satypes.Subaccount{
		Id: &Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_100BTC_Short_10100USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_100,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-10_000_000_000), // -100 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_2BTC_Short = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-200_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num1_1BTC_Short = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short_49999USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(49_999_000_000), // $49,999
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short_50000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(50_000_000_000), // $50,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_short_50000_non_TDai = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  1,
				Quantums: dtypes.NewInt(50_000_000_000), // non tdai asset
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
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
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short_54999USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(54_999_000_000), // $54,999
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Long_54999USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-54_999_000_000), // -$54,999
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short_55000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(55_000_000_000), // $55,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short_100000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(100_000_000_000), // $100,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_1BTC_Short_1ETH_Long_47000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(47_000_000_000), // $47,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
			{
				PerpetualId: 1,
				Quantums:    dtypes.NewInt(1_000_000_000), // 1 ETH
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num0_599USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_599,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_660USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_660,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_10000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_50000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_100000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_500000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_0USD = satypes.Subaccount{
		Id:                 &Carl_Num0,
		AssetPositions:     []*satypes.AssetPosition{},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num1_500USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_500,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num1_100000USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num1_Short_500USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_500,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Carl_Num0_1ISO_Short_49USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(49_000_000), // $49
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  3,
				Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ISO
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num1_01BTC_Long_4600USD_Short = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_4_600,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Carl_Num1_1BTC_Short_50499USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(50_499_000_000), // $50,499
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_01BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_50001USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(50_001_000_000), // $50,001
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Short_100000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_45000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-45_000_000_000), // -$45,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_45001USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-45_001_000_000), // -$45,001
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_49501USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-49_501_000_000), // -$49,501
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_50000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-50_000_000_000), // -$50,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_49999USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-49_999_000_000), // -$50,000
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_50001USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-50_001_000_000), // -$50,001
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  0,
				Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_599USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_599,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Dave_Num0_10000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Dave_Num0_500000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Dave_Num0_100BTC_Short_10200USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_200,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-10_000_000_000), // -100 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_100BTC_Long_9900USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_9_900,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(10_000_000_000), // 100 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1ISO2_Short_499USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(499_000_000), // $499
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  4,
				Quantums:     dtypes.NewInt(-10_000_000), // -1 ISO2
				FundingIndex: dtypes.NewInt(0),
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_49500USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_49_500,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
			{
				PerpetualId: 1,
				Quantums:    dtypes.NewInt(1_000_000_000), // 1 ETH
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_2_900,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(1), // tiny BTC
			},
			{
				PerpetualId: 1,
				Quantums:    dtypes.NewInt(1_000_000_000), // 1 ETH
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}

	Dave_Num1_1BTC_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_1ETH_Long_2900USD_Short = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_2_900,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 1,
				Quantums:    dtypes.NewInt(1_000_000_000), // 1 ETH
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_10_000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Dave_Num1_500000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
		AssetYieldIndex:    big.NewRat(1, 1).String(),
	}
	Dave_Num1_025BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(25_000_000), // 0.25 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_05BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(50_000_000), // 0.5 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_1BTC_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_100BTC_Short_10100USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_10_100,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-10_000_000_000), // -100 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_1BTC_Long_49501USD_Short = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(-49_501_000_000), // -$49,501
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num1_1ETH_Long_50000USD = satypes.Subaccount{
		Id: &Dave_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&TDai_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 1,
				Quantums:    dtypes.NewInt(1_000_000_000), // 1 ETH
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}
	Dave_Num2_1BTC_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num2,
		AssetPositions: []*satypes.AssetPosition{
			&Short_TDai_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
		AssetYieldIndex: big.NewRat(1, 1).String(),
	}

	// Quote balances.
	QuoteBalance_OneDollar = int64(1_000_000) // $1.

	// Asset Yield Index
	AssetYieldIndex_Zero = big.NewRat(1, 1).String()

	InvalidSubaccountIdNumber = satypes.SubaccountId{
		Owner:  Alice_Num0.Owner,
		Number: satypes.MaxSubaccountIdNumber + 1, // one over the limit
	}
	InvalidSubaccountIdOwner = satypes.SubaccountId{
		Owner:  "This is not valid",
		Number: 0,
	}
)
