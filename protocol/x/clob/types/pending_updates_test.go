package types_test

import (
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

type perpetualFill struct {
	subaccountId         satypes.SubaccountId
	perpetualId          uint32
	isBuy                bool
	bigFillBaseQuantums  *big.Int
	bigFillQuoteQuantums *big.Int
	feePpm               int32
}

func TestPendingUpdates(t *testing.T) {
	tests := []struct {
		name            string
		perpetualFills  []perpetualFill
		expectedUpdates []satypes.Update
	}{
		{
			name:            "empty",
			perpetualFills:  []perpetualFill{},
			expectedUpdates: []satypes.Update{},
		},
		{
			name: "multiple fill amounts (no fees)",
			perpetualFills: []perpetualFill{
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                true,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(100),
				},
				{
					subaccountId:         constants.Alice_Num1,
					perpetualId:          uint32(1),
					isBuy:                false,
					bigFillBaseQuantums:  big.NewInt(200),
					bigFillQuoteQuantums: big.NewInt(200),
				},
			},
			expectedUpdates: []satypes.Update{
				{
					SubaccountId: constants.Alice_Num0,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100),
						},
					},
				},
				{
					SubaccountId: constants.Alice_Num1,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(200)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-200),
						},
					},
				},
			},
		},
		{
			name: "multiple fill amounts (with fees)",
			perpetualFills: []perpetualFill{
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                true,
					feePpm:               500,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(10_000), // fee = 5_000_000 / 1_000_000
				},
				{
					subaccountId:         constants.Alice_Num1,
					perpetualId:          uint32(1),
					isBuy:                false,
					feePpm:               200,
					bigFillBaseQuantums:  big.NewInt(200),
					bigFillQuoteQuantums: big.NewInt(20_000), // fee = 4_000_000 / 1_000_000
				},
				{
					subaccountId:         constants.Bob_Num0,
					perpetualId:          uint32(0),
					isBuy:                true,
					feePpm:               500,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(10_050), // 5_025_000 / 1_000_000) (round to 6)
				},
			},
			expectedUpdates: []satypes.Update{
				{
					SubaccountId: constants.Bob_Num0,
					// - 10_050 - (fee) 6
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-10_056)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100),
						},
					},
				},
				{
					SubaccountId: constants.Alice_Num0,
					// - 10_000 - (fee) 5
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-10005)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100),
						},
					},
				},
				{
					SubaccountId: constants.Alice_Num1,
					// 20_000 - (fee) 4
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(19_996)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-200),
						},
					},
				},
			},
		},
		{
			name: "multiple fill amounts for same account (no fees)",
			perpetualFills: []perpetualFill{
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                true,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(100),
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(1),
					isBuy:                true,
					bigFillBaseQuantums:  big.NewInt(200),
					bigFillQuoteQuantums: big.NewInt(200),
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                false,
					bigFillBaseQuantums:  big.NewInt(50),
					bigFillQuoteQuantums: big.NewInt(50),
				},
			},
			expectedUpdates: []satypes.Update{
				{
					SubaccountId: constants.Alice_Num0,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-250)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50),
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(200),
						},
					},
				},
			},
		},
		{
			name: "multiple fill amounts for same account (with fees)",
			perpetualFills: []perpetualFill{
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                true,
					feePpm:               200,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(15_000), // fee = 3_000_000 / 1_000_000
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(1),
					isBuy:                true,
					feePpm:               200,
					bigFillBaseQuantums:  big.NewInt(200),
					bigFillQuoteQuantums: big.NewInt(1_500), // fee = 300_000 / 1_000_000 (rounds to 1)
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                false,
					feePpm:               500,
					bigFillBaseQuantums:  big.NewInt(50),
					bigFillQuoteQuantums: big.NewInt(1_600), // fee = 800_000 / 1_000_000 (rounds to 1)
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(1),
					isBuy:                false,
					feePpm:               500,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(2_100), // fee = 1_050_000 / 1_000_000 (rounds to 2)
				},
			},
			expectedUpdates: []satypes.Update{
				{
					SubaccountId: constants.Alice_Num0,
					// - 15_000 - 1_500 + 1_600 + 2_100 - (fee) 7
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-12_807)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50),
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(100),
						},
					},
				},
			},
		},
		{
			name: "multiple fill amounts for same account (with negative fees)",
			perpetualFills: []perpetualFill{
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                true,
					feePpm:               -200,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(15_000), // fee = -3_000_000 / 1_000_000 (rounds to 0)
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(1),
					isBuy:                true,
					feePpm:               -200,
					bigFillBaseQuantums:  big.NewInt(200),
					bigFillQuoteQuantums: big.NewInt(1_500), // fee = -300_000 / 1_000_000 (rounds to 0)
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(0),
					isBuy:                false,
					feePpm:               500,
					bigFillBaseQuantums:  big.NewInt(50),
					bigFillQuoteQuantums: big.NewInt(1_600), // fee = 800_000 / 1_000_000 (rounds to 1)
				},
				{
					subaccountId:         constants.Alice_Num0,
					perpetualId:          uint32(1),
					isBuy:                false,
					feePpm:               500,
					bigFillBaseQuantums:  big.NewInt(100),
					bigFillQuoteQuantums: big.NewInt(2_100), // fee = 1_050_000 / 1_000_000 (rounds to 2)
				},
			},
			expectedUpdates: []satypes.Update{
				{
					SubaccountId: constants.Alice_Num0,
					// - 15_000 - 1_500 + 1_600 + 2_100 + (fee) 3
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-12_800)),
					PerpetualUpdates: []satypes.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50),
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(100),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run many times for determinism.
			for i := 0; i < 100; i++ {
				pendingUpdates := types.NewPendingUpdates()

				for _, perpetualFill := range tt.perpetualFills {
					pendingUpdates.AddPerpetualFill(
						perpetualFill.subaccountId,
						perpetualFill.perpetualId,
						perpetualFill.isBuy,
						perpetualFill.feePpm,
						perpetualFill.bigFillBaseQuantums,
						perpetualFill.bigFillQuoteQuantums,
					)
				}

				updates := pendingUpdates.ConvertToUpdates()

				require.Equal(t, tt.expectedUpdates, updates)
			}
		})
	}
}
