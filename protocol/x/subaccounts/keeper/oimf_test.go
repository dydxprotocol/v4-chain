package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
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

func TestGetDeltaOpenInterestFromUpdates(t *testing.T) {
	tests := map[string]struct {
		settledUpdates []types.SettledUpdate
		updateType     types.UpdateType
		expectedVal    *perptypes.OpenInterestDelta
		panicErr       string
	}{
		"Invalid: 1 update": {
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			expectedVal: &perptypes.OpenInterestDelta{
				PerpetualId:  1,
				BaseQuantums: big.NewInt(500),
			},
		},
		"Valid: 500 -> 0, 0 -> 500, delta = 0": {
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			expectedVal: nil, // delta is 0
		},
		"Not Match update, return nil": {
			updateType: types.CollatCheck,
			settledUpdates: []types.SettledUpdate{
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
			},
			expectedVal: nil,
		},
		"Valid: 500 -> 350, 0 -> 150, delta = 0": {
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			expectedVal: nil, // delta is 0
		},
		"Valid: -100 -> 200, 250 -> -50, delta = -50": {
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			expectedVal: &perptypes.OpenInterestDelta{
				PerpetualId:  1000,
				BaseQuantums: big.NewInt(-50),
			},
		},
		"Valid: -3100 -> -5000, 1000 -> 2900, delta = 1900": {
			updateType: types.Match,
			settledUpdates: []types.SettledUpdate{
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
			expectedVal: &perptypes.OpenInterestDelta{
				PerpetualId:  1000,
				BaseQuantums: big.NewInt(1900),
			},
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
						keeper.GetDeltaOpenInterestFromUpdates(
							tc.settledUpdates,
							tc.updateType,
						)
					},
				)
				return
			}

			perpOpenInterestDelta := keeper.GetDeltaOpenInterestFromUpdates(
				tc.settledUpdates,
				tc.updateType,
			)
			require.Equal(
				t,
				tc.expectedVal,
				perpOpenInterestDelta,
			)
		})
	}
}
