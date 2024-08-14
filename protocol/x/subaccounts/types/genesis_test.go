package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	duplicateOwnerId := sample.AccAddress()
	tests := map[string]struct {
		genState      *types.GenesisState
		shouldPanic   bool
		expectedError error
	}{
		"valid: default": {
			genState:      types.DefaultGenesis(),
			expectedError: nil,
		},
		"valid": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
				},
			},
			expectedError: nil,
		},
		"valid: duplicate owner ids with separate numbers": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  duplicateOwnerId,
							Number: uint32(9),
						},
					},
					{
						Id: &types.SubaccountId{
							Owner:  duplicateOwnerId,
							Number: uint32(11),
						},
					},
				},
			},
			expectedError: nil,
		},
		"invalid: id owner is empty": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  "",
							Number: uint32(0),
						},
					},
				},
			},
			expectedError: types.ErrInvalidSubaccountIdOwner,
		},
		"invalid: id owner is invalid": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  "this is not a valid bech32 address",
							Number: uint32(0),
						},
					},
				},
			},
			expectedError: types.ErrInvalidSubaccountIdOwner,
		},
		"invalid: id number is greater than 128_000": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(128_001),
						},
					},
				},
			},
			expectedError: types.ErrInvalidSubaccountIdNumber,
		},
		"invalid: duplicate subaccount ids": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  duplicateOwnerId,
							Number: uint32(42),
						},
					},
					{
						Id: &types.SubaccountId{
							Owner:  duplicateOwnerId,
							Number: uint32(42),
						},
					},
				},
			},
			expectedError: types.ErrDuplicateSubaccountIds,
		},
		"invalid: multiple asset positions": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						AssetPositions: []*types.AssetPosition{ // multiple asset positions.
							testutil.CreateSingleAssetPosition(0, new(big.Int)),
							testutil.CreateSingleAssetPosition(1, new(big.Int)),
						},
					},
				},
			},
			expectedError: types.ErrMultAssetPositionsNotSupported,
		},
		// TODO(DEC-582): once we support different assets, add a test case for the asset ordering.
		"invalid: asset position id != 0": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						AssetPositions: []*types.AssetPosition{
							testutil.CreateSingleAssetPosition(
								1, // asset id must be zero (0 = USDC).
								big.NewInt(1_000),
							),
						},
					},
				},
			},
			expectedError: types.ErrAssetPositionNotSupported,
		},
		"invalid: asset position quantum == 0": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						AssetPositions: []*types.AssetPosition{
							testutil.CreateSingleAssetPosition(
								0,
								big.NewInt(0), // quantum cannot be zero.
							),
						},
					},
				},
			},
			shouldPanic:   false,
			expectedError: types.ErrAssetPositionZeroQuantum,
		},
		"invalid: perpetual positions out of order": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						PerpetualPositions: []*types.PerpetualPosition{
							testutil.CreateSinglePerpetualPosition(
								2, // out of order.
								big.NewInt(1_000),
								big.NewInt(0),
								big.NewInt(0),
							),
							testutil.CreateSinglePerpetualPosition(
								1,
								big.NewInt(1_000),
								big.NewInt(0),
								big.NewInt(0),
							),
						},
					},
				},
			},
			expectedError: types.ErrPerpPositionsOutOfOrder,
		},
		"invalid: perpetual position quantum == 0": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						PerpetualPositions: []*types.PerpetualPosition{
							testutil.CreateSinglePerpetualPosition(
								0,
								big.NewInt(0), // quantum cannot be zero.
								big.NewInt(0),
								big.NewInt(0),
							),
						},
					},
				},
			},
			shouldPanic:   false,
			expectedError: types.ErrPerpPositionZeroQuantum,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				require.PanicsWithError(t, tc.expectedError.Error(), func() {
					// nolint:errcheck
					tc.genState.Validate()
				})
				return
			}

			err := tc.genState.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}
