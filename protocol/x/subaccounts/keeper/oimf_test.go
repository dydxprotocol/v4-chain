package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	keeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	aliceSubaccountId = &types.SubaccountId{
		Owner: "Alice",
	}
	bobSubaccountId = &types.SubaccountId{
		Owner: "Bob",
	}
)

func TestGetDeltaOpenInterestFromPerpMatchUpdates(t *testing.T) {
	tests := map[string]struct {
		settledUpdates        []keeper.SettledUpdate
		expectedDelta         *big.Int
		expectedUpdatedPerpId uint32
		panicErr              string
	}{
		"Invalid: 1 update": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      0,
							BigQuantumsDelta: big.NewInt(1_000),
						},
					},
				},
			},
			panicErr: types.ErrMatchUpdatesMustHaveTwoUpdates,
		},
		"Invalid: one of the updates contains no perp update": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      0,
							BigQuantumsDelta: big.NewInt(1_000),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
					},
				},
			},
			panicErr: types.ErrMatchUpdatesMustUpdateOnePerp,
		},
		"Invalid: updates are on different perpetuals": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      0,
							BigQuantumsDelta: big.NewInt(1_000),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(1_000),
						},
					},
				},
			},
			panicErr: types.ErrMatchUpdatesMustBeSamePerpId,
		},
		"Invalid: updates don't have opposite signs": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(500),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(500),
						},
					},
				},
			},
			panicErr: types.ErrMatchUpdatesInvalidSize,
		},
		"Invalid: updates don't have equal absolute base quantums": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(500),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(-499),
						},
					},
				},
			},
			panicErr: types.ErrMatchUpdatesInvalidSize,
		},
		"Valid: 0 -> -500, 0 -> 500, delta = 500": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(500),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1,
							BigQuantumsDelta: big.NewInt(-500),
						},
					},
				},
			},
			expectedUpdatedPerpId: 1,
			expectedDelta:         big.NewInt(500),
		},
		"Valid: 500 -> 0, 0 -> 500, delta = 0": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id:                 aliceSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(500),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{
							{
								PerpetualId: 1000,
								Quantums:    dtypes.NewInt(500),
							},
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(-500),
						},
					},
				},
			},
			expectedUpdatedPerpId: 1000,
			expectedDelta:         big.NewInt(0),
		},
		"Valid: 500 -> 350, 0 -> 150, delta = 0": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{
							{
								PerpetualId: 1000,
								Quantums:    dtypes.NewInt(500),
							},
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(-150),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id:                 bobSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(150),
						},
					},
				},
			},
			expectedUpdatedPerpId: 1000,
			expectedDelta:         big.NewInt(0),
		},
		"Valid: -100 -> 200, 250 -> -50, delta = -50": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{
							{
								PerpetualId: 1000,
								Quantums:    dtypes.NewInt(-100),
							},
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(300),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{
							{
								PerpetualId: 1000,
								Quantums:    dtypes.NewInt(250),
							},
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(-300),
						},
					},
				},
			},
			expectedUpdatedPerpId: 1000,
			expectedDelta:         big.NewInt(-50),
		},
		"Valid: -3100 -> -5000, 1000 -> 2900, delta = 1900": {
			settledUpdates: []keeper.SettledUpdate{
				{
					SettledSubaccount: types.Subaccount{
						Id: aliceSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{
							{
								PerpetualId: 1000,
								Quantums:    dtypes.NewInt(-3100),
							},
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(-1900),
						},
					},
				},
				{
					SettledSubaccount: types.Subaccount{
						Id: bobSubaccountId,
						PerpetualPositions: []*types.PerpetualPosition{
							{
								PerpetualId: 1000,
								Quantums:    dtypes.NewInt(1000),
							},
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      1000,
							BigQuantumsDelta: big.NewInt(+1900),
						},
					},
				},
			},
			expectedUpdatedPerpId: 1000,
			expectedDelta:         big.NewInt(1900),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panicErr != "" {
				require.PanicsWithValue(t,
					fmt.Sprintf(
						tc.panicErr,
						tc.settledUpdates,
					), func() {
						keeper.GetDeltaOpenInterestFromPerpMatchUpdates(tc.settledUpdates)
					},
				)
				return
			}

			updatedPerpId, deltaOpenInterest := keeper.GetDeltaOpenInterestFromPerpMatchUpdates(tc.settledUpdates)
			require.Equal(
				t,
				tc.expectedUpdatedPerpId,
				updatedPerpId,
			)
			require.Zerof(
				t,
				tc.expectedDelta.Cmp(deltaOpenInterest),
				"deltaOpenInterest: %v, tc.expectedDelta: %v",
				deltaOpenInterest,
				tc.expectedDelta,
			)
		})
	}
}
