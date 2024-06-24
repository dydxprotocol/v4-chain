package types_test

import (
	"math/big"
	"testing"

	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	satest "github.com/dydxprotocol/v4-chain/protocol/testutil/subaccounts"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		update types.Update
		err    error
	}{
		"empty update": {
			update: types.Update{},
			err:    nil,
		},
		"valid update": {
			update: types.Update{
				AssetUpdates:     testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
				PerpetualUpdates: satest.CreatePerpetualUpdate(1, big.NewInt(1)),
			},
			err: nil,
		},
		"duplicate asset update": {
			update: types.Update{
				AssetUpdates: append(
					testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
					testutil.CreateUsdcAssetUpdate(big.NewInt(1))...,
				),
				PerpetualUpdates: satest.CreatePerpetualUpdate(1, big.NewInt(1)),
			},
			err: types.ErrNonUniqueUpdatesPosition,
		},
		"duplicate perpetual update": {
			update: types.Update{
				AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
				PerpetualUpdates: append(
					satest.CreatePerpetualUpdate(1, big.NewInt(1)),
					satest.CreatePerpetualUpdate(1, big.NewInt(1))...,
				),
			},
			err: types.ErrNonUniqueUpdatesPosition,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.update.Validate()
			if err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			} else {
				require.NoError(t, tc.err)
			}
		})
	}
}

func TestGetErrorFromUpdateResults(t *testing.T) {
	tests := map[string]struct {
		success          bool
		successPerUpdate []types.UpdateResult
		updates          []types.Update
		expectedErr      error
		expectPanic      bool
	}{
		"success = true": {
			success:     true,
			expectedErr: nil,
		},
		"success = false": {
			success:          false,
			successPerUpdate: []types.UpdateResult{types.NewlyUndercollateralized},
			updates: []types.Update{{
				SubaccountId: types.SubaccountId{
					Owner: "owner",
				},
				AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
			}},
			expectedErr: types.ErrFailedToUpdateSubaccounts,
		},
		"success = false, but successPerUpdate contains no failure": {
			success:          false,
			successPerUpdate: []types.UpdateResult{types.Success},
			updates: []types.Update{{
				SubaccountId: types.SubaccountId{
					Owner: "owner",
				},
				AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
			}},
			expectPanic: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectPanic {
				require.Panics(
					t,
					func() {
						err := types.GetErrorFromUpdateResults(tc.success, tc.successPerUpdate, tc.updates)
						require.NoError(t, err)
					},
				)
				return
			}

			err := types.GetErrorFromUpdateResults(tc.success, tc.successPerUpdate, tc.updates)
			if tc.expectedErr == nil {
				require.Equal(t, nil, err)
				return
			}
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestUpdateResultString(t *testing.T) {
	tests := map[string]struct {
		value          types.UpdateResult
		expectedResult string
	}{
		"Success": {
			value:          types.Success,
			expectedResult: "Success",
		},
		"NewlyUndercollateralized": {
			value:          types.NewlyUndercollateralized,
			expectedResult: "NewlyUndercollateralized",
		},
		"StillUndercollateralized": {
			value:          types.StillUndercollateralized,
			expectedResult: "StillUndercollateralized",
		},
		"UpdateCausedError": {
			value:          types.UpdateCausedError,
			expectedResult: "UpdateCausedError",
		},
		"ViolatesIsolatedSubaccountConstraints": {
			value:          types.ViolatesIsolatedSubaccountConstraints,
			expectedResult: "ViolatesIsolatedSubaccountConstraints",
		},
		"UnexpectedError": {
			value:          types.UpdateResult(6),
			expectedResult: "UnexpectedError",
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.value.String()
			require.Equal(t, result, tc.expectedResult)
		})
	}
}

func TestUpdateResultIsSuccess(t *testing.T) {
	tests := map[string]struct {
		value          types.UpdateResult
		expectedResult bool
	}{
		"Success": {
			value:          types.Success,
			expectedResult: true,
		},
		"NewlyUndercollateralized": {
			value:          types.NewlyUndercollateralized,
			expectedResult: false,
		},
		"StillUndercollateralized": {
			value:          types.StillUndercollateralized,
			expectedResult: false,
		},
		"WithdrawalsAndTransfersBlocked": {
			value:          types.WithdrawalsAndTransfersBlocked,
			expectedResult: false,
		},
		"UpdateCausedError": {
			value:          types.UpdateCausedError,
			expectedResult: false,
		},
		"UnexpectedError": {
			value:          types.UpdateResult(5),
			expectedResult: false,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.value.IsSuccess()
			require.Equal(t, result, tc.expectedResult)
		})
	}
}

func TestUpdateTypeString(t *testing.T) {
	tests := map[string]struct {
		value          types.UpdateType
		expectedResult string
	}{
		"Withdrawal": {
			value:          types.Withdrawal,
			expectedResult: "Withdrawal",
		},
		"Transfer": {
			value:          types.Transfer,
			expectedResult: "Transfer",
		},
		"Deposit": {
			value:          types.Deposit,
			expectedResult: "Deposit",
		},
		"Match": {
			value:          types.Match,
			expectedResult: "Match",
		},
		"UnexpectedError": {
			value:          types.UpdateType(999),
			expectedResult: "UnexpectedUpdateTypeError",
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.value.String()
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
