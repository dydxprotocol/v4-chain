package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestValidateLockedShares(t *testing.T) {
	tests := map[string]struct {
		// LockedShares to validate.
		lockedShares types.LockedShares
		// Expected error.
		expectedErr string
	}{
		"Success": {
			lockedShares: types.LockedShares{
				OwnerAddress:      constants.Alice_Num0.Owner,
				TotalLockedShares: types.BigIntToNumShares(big.NewInt(77)),
				UnlockDetails: []types.UnlockDetail{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(77)),
						UnlockBlockHeight: 1,
					},
				},
			},
		},
		"Success - Multiple unlocks": {
			lockedShares: types.LockedShares{
				OwnerAddress:      constants.Alice_Num0.Owner,
				TotalLockedShares: types.BigIntToNumShares(big.NewInt(100)),
				UnlockDetails: []types.UnlockDetail{
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
			lockedShares: types.LockedShares{
				OwnerAddress:      "",
				TotalLockedShares: types.BigIntToNumShares(big.NewInt(77)),
				UnlockDetails: []types.UnlockDetail{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(77)),
						UnlockBlockHeight: 1,
					},
				},
			},
			expectedErr: "empty owner address",
		},
		"Failure - mismatched total locked shares and total shares to unlock": {
			lockedShares: types.LockedShares{
				OwnerAddress:      constants.AliceAccAddress.String(),
				TotalLockedShares: types.BigIntToNumShares(big.NewInt(100)),
				UnlockDetails: []types.UnlockDetail{
					{
						Shares:            types.BigIntToNumShares(big.NewInt(54)),
						UnlockBlockHeight: 1,
					},
					{
						Shares:            types.BigIntToNumShares(big.NewInt(45)),
						UnlockBlockHeight: 1234,
					},
				},
			},
			expectedErr: "total shares locked (100) not equal to total shares to unlock (99)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.lockedShares.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
