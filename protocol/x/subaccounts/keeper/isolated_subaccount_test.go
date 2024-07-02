package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetIsolatedPerpetualStateTransition(t *testing.T) {
	tests := map[string]struct {
		// parameters
		settledUpdateWithUpdatedSubaccount types.SettledUpdate
		perpetuals                         []perptypes.Perpetual

		// expectation
		expectedStateTransition *types.IsolatedPerpetualPositionStateTransition
		expectedErr             error
	}{
		`If no perpetual updates, nil state transition is returned`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id:                 &constants.Alice_Num0,
					PerpetualPositions: nil,
					AssetPositions:     nil,
				},
				PerpetualUpdates: nil,
				AssetUpdates:     nil,
			},
			perpetuals:              nil,
			expectedStateTransition: nil,
		},
		`If single non-isolated perpetual updates, nil state transition is returned`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id:                 &constants.Alice_Num0,
					PerpetualPositions: nil,
					AssetPositions:     nil,
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(0),
						BigQuantumsDelta: big.NewInt(-100),
					},
				},
				AssetUpdates: nil,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			expectedStateTransition: nil,
		},
		`If multiple non-isolated perpetual updates, nil state transition is returned`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id:                 &constants.Alice_Num0,
					PerpetualPositions: nil,
					AssetPositions:     nil,
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(0),
						BigQuantumsDelta: big.NewInt(-100),
					},
					{
						PerpetualId:      uint32(1),
						BigQuantumsDelta: big.NewInt(-200),
					},
				},
				AssetUpdates: nil,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			expectedStateTransition: nil,
		},
		`If multiple non-isolated perpetual positions, nil state transition is returned`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*types.PerpetualPosition{
						&constants.PerpetualPosition_OneBTCLong,
						&constants.PerpetualPosition_OneTenthEthLong,
					},
					AssetPositions: nil,
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(0),
						BigQuantumsDelta: big.NewInt(-100),
					},
				},
				AssetUpdates: nil,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			expectedStateTransition: nil,
		},
		`If single isolated perpetual update, no perpetual position, state transition is returned for closed position`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id:                 &constants.Alice_Num0,
					PerpetualPositions: nil,
					AssetPositions: []*types.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(3),
						BigQuantumsDelta: big.NewInt(-1_000_000_000),
					},
				},
				AssetUpdates: []types.AssetUpdate{
					{
						AssetId:          assettypes.AssetUsdc.Id,
						BigQuantumsDelta: big.NewInt(100_000_000),
					},
				},
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			expectedStateTransition: &types.IsolatedPerpetualPositionStateTransition{
				SubaccountId:  &constants.Alice_Num0,
				PerpetualId:   uint32(3),
				QuoteQuantums: constants.Usdc_Asset_10_000.GetBigQuantums(),
				Transition:    types.Closed,
			},
		},
		`If single isolated perpetual update, existing perpetual position with same size, state transition is returned for
		opened position`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*types.PerpetualPosition{
						&constants.PerpetualPosition_OneISOLong,
					},
					AssetPositions: []*types.AssetPosition{
						testutil.CreateSingleAssetPosition(
							assettypes.AssetUsdc.Id,
							big.NewInt(-40_000_000), // -$40
						),
					},
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(3),
						BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
					},
				},
				AssetUpdates: []types.AssetUpdate{
					{
						AssetId:          assettypes.AssetUsdc.Id,
						BigQuantumsDelta: big.NewInt(-50_000_000), // -$50
					},
				},
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			expectedStateTransition: &types.IsolatedPerpetualPositionStateTransition{
				SubaccountId:  &constants.Alice_Num0,
				PerpetualId:   uint32(3),
				QuoteQuantums: big.NewInt(10_000_000), // $-40 - (-$50)
				Transition:    types.Opened,
			},
		},
		`If single isolated perpetual update, existing perpetual position with different size, nil state transition
		returned`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*types.PerpetualPosition{
						&constants.PerpetualPosition_OneISOLong,
					},
					AssetPositions: []*types.AssetPosition{
						testutil.CreateSingleAssetPosition(
							assettypes.AssetUsdc.Id,
							big.NewInt(-40_000_000), // -$65
						),
					},
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(3),
						BigQuantumsDelta: big.NewInt(500_000_000), // 0.5 ISO
					},
				},
				AssetUpdates: []types.AssetUpdate{
					{
						AssetId:          assettypes.AssetUsdc.Id,
						BigQuantumsDelta: big.NewInt(-25_000_000), // -$25
					},
				},
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			expectedStateTransition: nil,
		},
		`Returns error if perpetual position was opened with no asset updates`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*types.PerpetualPosition{
						&constants.PerpetualPosition_OneISOLong,
					},
					AssetPositions: []*types.AssetPosition{
						testutil.CreateSingleAssetPosition(
							assettypes.AssetUsdc.Id,
							big.NewInt(50_000_000), // $50
						),
					},
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(3),
						BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
					},
				},
				AssetUpdates: nil,
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			expectedStateTransition: nil,
			expectedErr:             types.ErrFailedToUpdateSubaccounts,
		},
		`Returns error if perpetual position was opened with multiple asset updates`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*types.PerpetualPosition{
						&constants.PerpetualPosition_OneISOLong,
					},
					AssetPositions: []*types.AssetPosition{
						testutil.CreateSingleAssetPosition(
							assettypes.AssetUsdc.Id,
							big.NewInt(-40_000_000), // -$40
						),
						testutil.CreateSingleAssetPosition(
							constants.BtcUsd.Id,
							big.NewInt(100_000_000), // 1 BTC
						),
					},
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(3),
						BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
					},
				},
				AssetUpdates: []types.AssetUpdate{
					{
						AssetId:          assettypes.AssetUsdc.Id,
						BigQuantumsDelta: big.NewInt(-50_000_000), // -$50
					},
					{
						AssetId:          constants.BtcUsd.Id,
						BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
					},
				},
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			expectedStateTransition: nil,
			expectedErr:             types.ErrFailedToUpdateSubaccounts,
		},
		`Returns error if perpetual position was opened with non-usdc asset update`: {
			settledUpdateWithUpdatedSubaccount: types.SettledUpdate{
				SettledSubaccount: types.Subaccount{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*types.PerpetualPosition{
						&constants.PerpetualPosition_OneISOLong,
					},
					AssetPositions: []*types.AssetPosition{
						testutil.CreateSingleAssetPosition(
							assettypes.AssetUsdc.Id,
							big.NewInt(50_000_000), // $50
						),
						testutil.CreateSingleAssetPosition(
							constants.BtcUsd.Id,
							big.NewInt(100_000_000), // 1 BTC
						),
					},
				},
				PerpetualUpdates: []types.PerpetualUpdate{
					{
						PerpetualId:      uint32(3),
						BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
					},
				},
				AssetUpdates: []types.AssetUpdate{
					{
						AssetId:          constants.BtcUsd.Id,
						BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
					},
				},
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			expectedStateTransition: nil,
			expectedErr:             types.ErrFailedToUpdateSubaccounts,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				perpInfos := make(perptypes.PerpInfos, len(tc.perpetuals))
				for _, perp := range tc.perpetuals {
					perpInfos[perp.Params.Id] = perptypes.PerpInfo{
						Perpetual: perp,
					}
				}

				stateTransition, err := keeper.GetIsolatedPerpetualStateTransition(
					tc.settledUpdateWithUpdatedSubaccount,
					perpInfos,
				)
				if tc.expectedErr != nil {
					require.Error(t, tc.expectedErr, err)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedStateTransition, stateTransition)
				}
			},
		)
	}
}
