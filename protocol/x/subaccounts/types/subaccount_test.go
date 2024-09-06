package types_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/sample"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestBaseQuantumsToBigInt(t *testing.T) {
	num := uint64(5)
	bq := types.BaseQuantums(5)

	require.Zero(t, new(big.Int).SetUint64(num).Cmp(bq.ToBigInt()))
}

func TestBaseQuantumsToUInt64(t *testing.T) {
	num := uint64(5)
	bq := types.BaseQuantums(5)

	require.Equal(t, num, bq.ToUint64())
}

func TestSubaccountIdValidate(t *testing.T) {
	tests := map[string]struct {
		owner         string
		number        uint32
		expectedError error
	}{
		"validates successfully": {
			owner:  "dydx1x2hd82qerp7lc0kf5cs3yekftupkrl620te6u2",
			number: 0,
		},
		"validates successfully with non-zero subaccount": {
			owner:  sample.AccAddress(),
			number: 127,
		},
		"invalid address": {
			owner:         "this is not a valid bech32 address",
			number:        0,
			expectedError: types.ErrInvalidSubaccountIdOwner,
		},
		"empty address": {
			owner:         "",
			number:        0,
			expectedError: types.ErrInvalidSubaccountIdOwner,
		},
		"invalid number": {
			owner:         sample.AccAddress(),
			number:        128_001,
			expectedError: types.ErrInvalidSubaccountIdNumber,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sid := &types.SubaccountId{
				Owner:  tc.owner,
				Number: tc.number,
			}

			err := sid.Validate()
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSubaccountIdMustGetAccAccount(t *testing.T) {
	tests := map[string]struct {
		owner  string
		number uint32
		panics bool
	}{
		"MustGetAccAccount successfully": {
			owner:  sample.AccAddress(),
			number: 0,
		},
		"invalid address": {
			owner:  "this is not a valid bech32 address",
			number: 0,
			panics: true,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sid := &types.SubaccountId{
				Owner:  tc.owner,
				Number: tc.number,
			}

			if tc.panics {
				require.Panics(t, func() {
					sid.MustGetAccAddress()
				})
			} else {
				require.NotPanics(t, func() {
					sid.MustGetAccAddress()
				})
			}
		})
	}
}

func TestSubaccountGetPerpetualPositionForId(t *testing.T) {
	expectedPerpetualPositions := []*types.PerpetualPosition{
		{
			PerpetualId: 0,
			Quantums:    dtypes.NewInt(100),
		},
		{
			PerpetualId: 1,
			Quantums:    dtypes.NewInt(100),
		},
	}
	subaccount := types.Subaccount{
		PerpetualPositions: expectedPerpetualPositions,
	}

	position, exists := subaccount.GetPerpetualPositionForId(0)
	require.True(t, exists)
	require.Equal(t, expectedPerpetualPositions[0], position)

	position, exists = subaccount.GetPerpetualPositionForId(1)
	require.True(t, exists)
	require.Equal(t, expectedPerpetualPositions[1], position)

	position, exists = subaccount.GetPerpetualPositionForId(2)
	require.False(t, exists)
	require.Nil(t, position)
}

func TestGetSubaccountQuoteBalance(t *testing.T) {
	tests := map[string]struct {
		subaccount           *types.Subaccount
		expectedQuoteBalance *big.Int
		panics               bool
	}{
		"can get positive quote balance": {
			subaccount:           &constants.Carl_Num0_599USD,
			expectedQuoteBalance: big.NewInt(599_000_000),
		},
		"can get negative quote balance": {
			subaccount: &types.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: []*types.AssetPosition{
					{
						AssetId:  assettypes.AssetTDai.Id,
						Quantums: dtypes.NewInt(-599_000_000), // $599
					},
				},
			},
			expectedQuoteBalance: big.NewInt(-599_000_000),
		},
		"can get from nil subaccount": {
			subaccount:           nil,
			expectedQuoteBalance: big.NewInt(0),
		},
		"can get from subaccount with no asset positions": {
			subaccount: &types.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedQuoteBalance: big.NewInt(0),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panics {
				require.Panics(
					t,
					func() {
						tc.subaccount.GetTDaiPosition()
					},
				)
			} else {
				require.NotPanics(t, func() {
					require.Equal(
						t,
						tc.expectedQuoteBalance,
						tc.subaccount.GetTDaiPosition(),
					)
				})
			}
		})
	}
}

func TestSetSubaccountQuoteBalance(t *testing.T) {
	tests := map[string]struct {
		subaccount             *types.Subaccount
		newQuoteBalance        *big.Int
		expectedAssetPositions []*types.AssetPosition
	}{
		"can set nil subaccount": {
			subaccount:      nil,
			newQuoteBalance: big.NewInt(0),
		},
		"can set positive quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: big.NewInt(10_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				&constants.TDai_Asset_10_000,
			},
		},
		"can set negative quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: big.NewInt(-10_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewInt(-10_000_000_000), // $10,000
				},
			},
		},
		"can set max quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: new(big.Int).SetUint64(math.MaxUint64),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewIntFromUint64(math.MaxUint64),
				},
			},
		},
		"can set min quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: constants.BigNegMaxUint64(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewIntFromBigInt(constants.BigNegMaxUint64()),
				},
			},
		},
		"can set zero quote balance and removes existing position from slice": {
			subaccount:             &constants.Carl_Num0_599USD,
			newQuoteBalance:        big.NewInt(0),
			expectedAssetPositions: []*types.AssetPosition{},
		},
		"can add TDai position to slice if non existent": {
			subaccount: &types.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: []*types.AssetPosition{
					&constants.Long_Asset_1BTC,
				},
			},
			newQuoteBalance: big.NewInt(10_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				&constants.TDai_Asset_10_000,
				&constants.Long_Asset_1BTC,
			},
		},
		"succeed if new quote balance overflows uint64": {
			subaccount: &types.Subaccount{
				Id:             &constants.Carl_Num0,
				AssetPositions: []*types.AssetPosition{},
			},
			newQuoteBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				big.NewInt(1),
			),
			expectedAssetPositions: []*types.AssetPosition{
				{
					Quantums: dtypes.NewIntFromBigInt(
						new(big.Int).Add(
							new(big.Int).SetUint64(math.MaxUint64),
							big.NewInt(1),
						),
					),
				},
			},
		},
		"returns error if abs new quote balance overflows uint64": {
			subaccount: &types.Subaccount{
				Id:             &constants.Carl_Num0,
				AssetPositions: []*types.AssetPosition{},
			},
			newQuoteBalance: new(big.Int).Add(
				constants.BigNegMaxUint64(),
				big.NewInt(-1),
			),
			expectedAssetPositions: []*types.AssetPosition{
				{
					Quantums: dtypes.NewIntFromBigInt(
						new(big.Int).Add(
							constants.BigNegMaxUint64(),
							big.NewInt(-1),
						),
					),
				},
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.subaccount.SetTDaiAssetPosition(tc.newQuoteBalance)
			if tc.subaccount != nil {
				require.Equal(
					t,
					tc.expectedAssetPositions,
					tc.subaccount.AssetPositions,
				)
			}
		})
	}
}

// ... existing imports ...

func TestGetTDaiPosition(t *testing.T) {
	tests := map[string]struct {
		subaccount           *types.Subaccount
		expectedQuoteBalance *big.Int
	}{
		"returns zero for nil subaccount": {
			subaccount:           nil,
			expectedQuoteBalance: new(big.Int),
		},
		"returns zero for subaccount with no asset positions": {
			subaccount:           &types.Subaccount{},
			expectedQuoteBalance: new(big.Int),
		},
		"returns correct positive balance": {
			subaccount: &types.Subaccount{
				AssetPositions: []*types.AssetPosition{
					{
						AssetId:  assettypes.AssetTDai.Id,
						Quantums: dtypes.NewInt(599_000_000),
					},
				},
			},
			expectedQuoteBalance: big.NewInt(599_000_000),
		},
		"returns correct negative balance": {
			subaccount: &types.Subaccount{
				AssetPositions: []*types.AssetPosition{
					{
						AssetId:  assettypes.AssetTDai.Id,
						Quantums: dtypes.NewInt(-10_000_000_000),
					},
				},
			},
			expectedQuoteBalance: big.NewInt(-10_000_000_000),
		},
		"gets TDai asset when there are multiple assets": {
			subaccount: &types.Subaccount{
				AssetPositions: []*types.AssetPosition{
					{
						AssetId:  assettypes.AssetTDai.Id,
						Quantums: dtypes.NewInt(1_000_000_000),
					},
					&constants.Long_Asset_1BTC,
				},
			},
			expectedQuoteBalance: big.NewInt(1_000_000_000),
		},
		"returns zero when TDai is not the first asset": {
			subaccount: &types.Subaccount{
				AssetPositions: []*types.AssetPosition{
					&constants.Long_Asset_1BTC,
					{
						AssetId:  assettypes.AssetTDai.Id,
						Quantums: dtypes.NewInt(1_000_000_000),
					},
				},
			},
			expectedQuoteBalance: new(big.Int),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.subaccount.GetTDaiPosition()
			require.Equal(t, 0, tc.expectedQuoteBalance.Cmp(result),
				"Expected quote balance %v. Got %v", tc.expectedQuoteBalance, result)
		})
	}
}

func TestSetTDaiAssetPosition(t *testing.T) {
	tests := map[string]struct {
		subaccount             *types.Subaccount
		newQuoteBalance        *big.Int
		expectedAssetPositions []*types.AssetPosition
	}{
		"sets nil subaccount": {
			subaccount:      nil,
			newQuoteBalance: big.NewInt(1_000_000_000),
		},
		"sets positive balance": {
			subaccount:      &types.Subaccount{},
			newQuoteBalance: big.NewInt(1_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewInt(1_000_000_000),
				},
			},
		},
		"sets negative balance": {
			subaccount:      &types.Subaccount{},
			newQuoteBalance: big.NewInt(-1_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewInt(-1_000_000_000),
				},
			},
		},
		"updates existing TDai position": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: big.NewInt(2_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewInt(2_000_000_000),
				},
			},
		},
		"removes TDai position when set to zero": {
			subaccount:             &constants.Carl_Num0_599USD,
			newQuoteBalance:        big.NewInt(0),
			expectedAssetPositions: []*types.AssetPosition{},
		},
		"removes TDai position when set to nil": {
			subaccount:             &constants.Carl_Num0_599USD,
			newQuoteBalance:        nil,
			expectedAssetPositions: []*types.AssetPosition{},
		},
		"adds TDai position to existing assets": {
			subaccount: &types.Subaccount{
				AssetPositions: []*types.AssetPosition{
					&constants.Long_Asset_1BTC,
				},
			},
			newQuoteBalance: big.NewInt(1_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  assettypes.AssetTDai.Id,
					Quantums: dtypes.NewInt(1_000_000_000),
				},
				&constants.Long_Asset_1BTC,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.subaccount.SetTDaiAssetPosition(tc.newQuoteBalance)
			if tc.subaccount != nil {
				require.Equal(t, tc.expectedAssetPositions, tc.subaccount.AssetPositions)
			}
		})
	}
}
