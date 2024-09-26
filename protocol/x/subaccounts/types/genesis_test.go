package types_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/sample"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
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
		"valid with zero AssetYieldIndex": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						AssetYieldIndex: "1",
					},
				},
			},
			expectedError: nil,
		},
		"valid with positive AssetYieldIndex": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						AssetYieldIndex: "10",
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
							{
								AssetId: 0,
							},
							{
								AssetId: 1,
							},
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
							{
								AssetId:  1, // asset id must be zero (0 = TDAI).
								Quantums: dtypes.NewInt(1_000),
							},
						},
					},
				},
			},
			expectedError: types.ErrAssetPositionNotSupported,
		},
		"invalid: AssetYieldIndex is not a rational number 1": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						AssetPositions: []*types.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(1_000),
							},
						},
						AssetYieldIndex: "12abc34",
					},
				},
			},
			expectedError: errors.New("could not convert string to big.Rat"),
		},
		"invalid: AssetYieldIndex is negative": {
			genState: &types.GenesisState{
				Subaccounts: []types.Subaccount{
					{
						Id: &types.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(127),
						},
						AssetPositions: []*types.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(1_000),
							},
						},
						AssetYieldIndex: "-1/1",
					},
				},
			},
			expectedError: types.ErrNegativeAssetYieldIndexNotSupported,
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
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(0), // quantum cannot be zero.
							},
						},
					},
				},
			},
			shouldPanic: true,
			expectedError: errorsmod.Wrapf(
				types.ErrAssetPositionZeroQuantum,
				"asset position (asset Id: 0) has zero quantum",
			),
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
							{
								PerpetualId: 2, // out of order.
								Quantums:    dtypes.NewInt(1_000),
							},
							{
								PerpetualId: 1,
								Quantums:    dtypes.NewInt(1_000),
							},
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
							{
								PerpetualId: 0,
								Quantums:    dtypes.ZeroInt(), // quantum cannot be zero.
							},
						},
					},
				},
			},
			shouldPanic: true,
			expectedError: errorsmod.Wrapf(
				types.ErrPerpPositionZeroQuantum,
				"perpetual position (perpetual Id: 0) has zero quantum",
			),
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
