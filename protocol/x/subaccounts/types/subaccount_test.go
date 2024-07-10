package types_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestBaseQuantums_ToBigInt(t *testing.T) {
	num := uint64(5)
	bq := types.BaseQuantums(5)

	require.Zero(t, new(big.Int).SetUint64(num).Cmp(bq.ToBigInt()))
}

func TestBaseQuantums_ToUInt64(t *testing.T) {
	num := uint64(5)
	bq := types.BaseQuantums(5)

	require.Equal(t, num, bq.ToUint64())
}

func TestSubaccountId_Validate(t *testing.T) {
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

func TestSubaccountId_MustGetAccAccount(t *testing.T) {
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

func TestSubaccount_DeepCopy(t *testing.T) {
	subaccount := constants.Alice_Num1_1BTC_Long_500_000USD
	deepCopy := subaccount.DeepCopy()

	require.Equal(t, subaccount, deepCopy)
	require.NotSame(t, &subaccount, &deepCopy)
}

func TestSubaccount_GetPerpetualPositionForId(t *testing.T) {
	expectedPerpetualPositions := []*types.PerpetualPosition{
		testutil.CreateSinglePerpetualPosition(
			0,
			big.NewInt(100),
			big.NewInt(0),
			big.NewInt(0),
		),
		testutil.CreateSinglePerpetualPosition(
			1,
			big.NewInt(100),
			big.NewInt(0),
			big.NewInt(0),
		),
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

func TestSubaccount_GetUsdcPosition(t *testing.T) {
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
					testutil.CreateSingleAssetPosition(
						assettypes.AssetUsdc.Id,
						big.NewInt(-599_000_000), // $599
					),
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
						tc.subaccount.GetUsdcPosition()
					},
				)
			} else {
				require.NotPanics(t, func() {
					require.Equal(
						t,
						tc.expectedQuoteBalance,
						tc.subaccount.GetUsdcPosition(),
					)
				})
			}
		})
	}
}

func TestSubaccount_SetUsdcAssetPosition(t *testing.T) {
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
				&constants.Usdc_Asset_10_000,
			},
		},
		"can set negative quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: big.NewInt(-10_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				testutil.CreateSingleAssetPosition(
					assettypes.AssetUsdc.Id,
					big.NewInt(-10_000_000_000), // $10,000
				),
			},
		},
		"can set max quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: new(big.Int).SetUint64(math.MaxUint64),
			expectedAssetPositions: []*types.AssetPosition{
				testutil.CreateSingleAssetPosition(
					assettypes.AssetUsdc.Id,
					big.NewInt(0).SetUint64(math.MaxUint64),
				),
			},
		},
		"can set min quote balance": {
			subaccount:      &constants.Carl_Num0_599USD,
			newQuoteBalance: constants.BigNegMaxUint64(),
			expectedAssetPositions: []*types.AssetPosition{
				testutil.CreateSingleAssetPosition(
					assettypes.AssetUsdc.Id,
					constants.BigNegMaxUint64(),
				),
			},
		},
		"can set zero quote balance and removes existing position from slice": {
			subaccount:             &constants.Carl_Num0_599USD,
			newQuoteBalance:        big.NewInt(0),
			expectedAssetPositions: []*types.AssetPosition{},
		},
		"can add usdc position to slice if non existent": {
			subaccount: &types.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: []*types.AssetPosition{
					&constants.Long_Asset_1BTC,
				},
			},
			newQuoteBalance: big.NewInt(10_000_000_000),
			expectedAssetPositions: []*types.AssetPosition{
				&constants.Usdc_Asset_10_000,
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
				testutil.CreateSingleAssetPosition(
					0,
					new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						big.NewInt(1),
					),
				),
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
				testutil.CreateSingleAssetPosition(
					0,
					new(big.Int).Add(
						constants.BigNegMaxUint64(),
						big.NewInt(-1),
					),
				),
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.subaccount.SetUsdcAssetPosition(tc.newQuoteBalance)
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
