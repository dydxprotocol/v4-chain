package constants

import (
	"github.com/dydxprotocol/v4/dtypes"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

var (
	// Subaccounts.
	Alice_Num0_10_000USD = satypes.Subaccount{
		Id: &Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Alice_Num1_10_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	User1_Num1_100_000USD = satypes.Subaccount{
		Id: &Alice_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num0_1BTC_Short = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
	}
	Carl_Num1_1BTC_Short = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
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
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
			},
		},
	}
	Carl_Num0_599USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_599,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num0_660USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_660,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num0_10000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num0_50000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num0_100000USD = satypes.Subaccount{
		Id: &Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_100_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num0_0USD = satypes.Subaccount{
		Id:                 &Carl_Num0,
		AssetPositions:     []*satypes.AssetPosition{},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num1_500USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num1_Short_500USD = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_500,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Carl_Num1_01BTC_Long_4600USD_Short = satypes.Subaccount{
		Id: &Carl_Num1,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_4_600,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC
			},
		},
	}
	Dave_Num0_1BTC_Long = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_50_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
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
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
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
	}
	Dave_Num0_599USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_599,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Dave_Num0_10000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_10_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Dave_Num0_500000USD = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Usdc_Asset_500_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{},
	}
	Dave_Num0_1BTC_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_46_000,
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId: 0,
				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
			},
		},
	}
	Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short = satypes.Subaccount{
		Id: &Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			&Short_Usdc_Asset_46_000,
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
