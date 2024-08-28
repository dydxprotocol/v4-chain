package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestValidateOwnerShareUnlocks(t *testing.T) {
	tests := map[string]struct {
		// OwnerShareUnlocks to validate.
		ownerShareUnlocks types.OwnerShareUnlocks
		// Expected error.
		expectedErr string
	}{
		"Success": {
			ownerShareUnlocks: types.OwnerShareUnlocks{
				OwnerAddress: constants.Alice_Num0.Owner,
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(77)),
						UnlockBlockHeight: 1,
					},
				},
			},
		},
		"Success - Multiple unlocks": {
			ownerShareUnlocks: types.OwnerShareUnlocks{
				OwnerAddress: constants.Alice_Num0.Owner,
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(40)),
						UnlockBlockHeight: 1,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(35)),
						UnlockBlockHeight: 1,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(25)),
						UnlockBlockHeight: 2,
					},
				},
			},
		},
		"Failure - empty owner address": {
			ownerShareUnlocks: types.OwnerShareUnlocks{
				OwnerAddress: "",
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(77)),
						UnlockBlockHeight: 1,
					},
				},
			},
			expectedErr: "Empty owner address",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.ownerShareUnlocks.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}

func TestGetTotalLockedShares(t *testing.T) {
	tests := map[string]struct {
		ownerShareUnlocks         types.OwnerShareUnlocks
		expectedTotalLockedShares *big.Int
	}{
		"0 unlocks": {
			ownerShareUnlocks: types.OwnerShareUnlocks{
				OwnerAddress: constants.AliceAccAddress.String(),
				ShareUnlocks: []types.ShareUnlock{},
			},
			expectedTotalLockedShares: big.NewInt(0),
		},
		"1 unlock": {
			ownerShareUnlocks: types.OwnerShareUnlocks{
				OwnerAddress: constants.BobAccAddress.String(),
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(123_456)),
						UnlockBlockHeight: 789_987,
					},
				},
			},
			expectedTotalLockedShares: big.NewInt(123_456),
		},
		"4 unlocks": {
			ownerShareUnlocks: types.OwnerShareUnlocks{
				OwnerAddress: constants.Alice_Num0.Owner,
				ShareUnlocks: []types.ShareUnlock{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(40)),
						UnlockBlockHeight: 1,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(35)),
						UnlockBlockHeight: 1,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(25)),
						UnlockBlockHeight: 2,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(122)),
						UnlockBlockHeight: 987,
					},
				},
			},
			expectedTotalLockedShares: big.NewInt(222),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedTotalLockedShares,
				tc.ownerShareUnlocks.GetTotalLockedShares(),
			)
		})
	}
}
