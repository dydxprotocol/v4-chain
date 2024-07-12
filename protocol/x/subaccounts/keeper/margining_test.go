package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perp_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetMarginedUpdates(t *testing.T) {
	perpInfos := perptypes.PerpInfos{
		0: perp_testutil.CreatePerpInfo(
			0,
			constants.BtcUsd_100PercentMarginRequirement.Params.AtomicResolution,
			constants.FiveBillion,
			constants.BtcUsdExponent,
		),
		1: perp_testutil.CreatePerpInfo(
			1,
			constants.EthUsd_100PercentMarginRequirement.Params.AtomicResolution,
			constants.ThreeBillion,
			constants.EthUsdExponent,
		),
	}

	tests := map[string]struct {
		subaccount       types.Subaccount
		assetUpdates     []types.AssetUpdate
		perpetualUpdates []types.PerpetualUpdate

		expectedAssetUpdates     []types.AssetUpdate
		expectedPerpetualUpdates []types.PerpetualUpdate
	}{
		"perpetual position is collateralized - no collateral updates": {
			subaccount: types.Subaccount{
				Id: &constants.Alice_Num0,
				AssetPositions: []*types.AssetPosition{
					&constants.Usdc_Asset_10_000,
				},
				PerpetualPositions: []*types.PerpetualPosition{
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(100_000_000), // 1 BTC
						big.NewInt(0),
						// $10,000, enough to cover 2 BTC total.
						big.NewInt(5_000_000_000-50_000_000_000),
					),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:          0,
					BigQuantumsDelta:     big.NewInt(100_000_000),
					BigQuoteBalanceDelta: big.NewInt(-50_000_000_000),
				},
			},
			expectedPerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:          0,
					BigQuantumsDelta:     big.NewInt(100_000_000),
					BigQuoteBalanceDelta: big.NewInt(-50_000_000_000),
				},
			},
		},
		`perpetual position is fully closed - remaining collateral moved back
		to main quote balance`: {
			subaccount: types.Subaccount{
				Id: &constants.Alice_Num0,
				AssetPositions: []*types.AssetPosition{
					&constants.Usdc_Asset_10_000,
				},
				PerpetualPositions: []*types.PerpetualPosition{
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(100_000_000), // 1 BTC
						big.NewInt(0),
						big.NewInt(10_000_000_000-50_000_000_000),
					),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      0,
					BigQuantumsDelta: big.NewInt(-100_000_000),
					// Position is closed at $50,000.
					BigQuoteBalanceDelta: big.NewInt(50_000_000_000),
				},
			},
			expectedAssetUpdates: []types.AssetUpdate{
				{
					AssetId: 0,
					// net $10,000 is moved back to the main quote balance.
					BigQuantumsDelta: big.NewInt(10_000_000_000),
				},
			},
			expectedPerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      0,
					BigQuantumsDelta: big.NewInt(-100_000_000),
					// net $10,000 is moved back to the main quote balance.
					BigQuoteBalanceDelta: big.NewInt(40_000_000_000),
				},
			},
		},
		`new perpetual position is under-collateralized - main quote balance has enough collateral
		for the new position`: {
			subaccount: types.Subaccount{
				Id: &constants.Alice_Num0,
				AssetPositions: []*types.AssetPosition{
					&constants.Usdc_Asset_10_000,
				},
				PerpetualPositions: []*types.PerpetualPosition{
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(100_000_000), // 1 BTC
						big.NewInt(0),
						big.NewInt(2_500_000_000-50_000_000_000),
					),
					testutil.CreateSinglePerpetualPosition(
						1,
						big.NewInt(1_000_000_000), // 1 ETH
						big.NewInt(0),
						// Extra collateral in ETH position.
						big.NewInt(10_000_000_000),
					),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:          0,
					BigQuantumsDelta:     big.NewInt(100_000_000),
					BigQuoteBalanceDelta: big.NewInt(-50_000_000_000),
				},
			},
			expectedAssetUpdates: []types.AssetUpdate{
				{
					AssetId: 0,
					// $2,500 is moved to the new position.
					BigQuantumsDelta: big.NewInt(-2_500_000_000),
				},
			},
			expectedPerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      0,
					BigQuantumsDelta: big.NewInt(100_000_000),
					// $2,500 is moved to the new position.
					BigQuoteBalanceDelta: big.NewInt(-50_000_000_000 + 2_500_000_000),
				},
			},
		},
		`new perpetual position is under-collateralized - main quote balance does not have enough collateral
		for the new position and rebalancing is needed`: {
			subaccount: types.Subaccount{
				Id: &constants.Alice_Num0,
				PerpetualPositions: []*types.PerpetualPosition{
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(100_000_000), // 1 BTC
						big.NewInt(0),
						big.NewInt(2_500_000_000-50_000_000_000),
					),
					testutil.CreateSinglePerpetualPosition(
						1,
						big.NewInt(1_000_000_000), // 1 ETH
						big.NewInt(0),
						// Extra collateral in ETH position.
						// NC = $10,000 + $3,000 = $13,000
						// MMR = $150
						// Free collateral = $13,000 - $150 = $12,850
						big.NewInt(10_000_000_000),
					),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:          0,
					BigQuantumsDelta:     big.NewInt(100_000_000),
					BigQuoteBalanceDelta: big.NewInt(-50_000_000_000),
				},
			},
			expectedAssetUpdates: []types.AssetUpdate{
				{
					AssetId:          0,
					BigQuantumsDelta: big.NewInt(12_850_000_000 - 2_500_000_000),
				},
			},
			expectedPerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      0,
					BigQuantumsDelta: big.NewInt(100_000_000),
					// $2,500 is moved to the new position.
					BigQuoteBalanceDelta: big.NewInt(-50_000_000_000 + 2_500_000_000),
				},
				{
					PerpetualId:      1,
					BigQuantumsDelta: big.NewInt(0),
					// $12,850 is moved out of the position.
					BigQuoteBalanceDelta: big.NewInt(-12_850_000_000),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualUpdates := keeper.GetMarginedUpdates(
				[]types.SettledUpdate{
					{
						SettledSubaccount: tc.subaccount,
						AssetUpdates:      tc.assetUpdates,
						PerpetualUpdates:  tc.perpetualUpdates,
					},
				}, perpInfos)
			require.Equal(t, tc.subaccount, actualUpdates[0].SettledSubaccount)
			require.Equal(t, tc.expectedAssetUpdates, actualUpdates[0].AssetUpdates)

			actualPerpetualUpdates := actualUpdates[0].PerpetualUpdates
			require.Equal(t, len(tc.expectedPerpetualUpdates), len(actualPerpetualUpdates))
			for i, expectedUpdate := range tc.expectedPerpetualUpdates {
				require.Equal(t, expectedUpdate.PerpetualId, actualPerpetualUpdates[i].PerpetualId)
				require.Equal(t, expectedUpdate.GetBigQuantums(), actualPerpetualUpdates[i].GetBigQuantums())
				require.Equal(t, expectedUpdate.GetBigQuoteBalance(), actualPerpetualUpdates[i].GetBigQuoteBalance())
			}
		})
	}
}
