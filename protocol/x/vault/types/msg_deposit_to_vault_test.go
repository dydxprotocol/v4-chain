package types_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDepositToVault_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgDepositToVault
		expectedErr string
	}{
		"Success": {
			msg: types.MsgDepositToVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewInt(1),
			},
		},
		"Success: max uint64 quote quantums": {
			msg: types.MsgDepositToVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewIntFromUint64(math.MaxUint64),
			},
		},
		"Failure: quote quantums greater than max uint64": {
			msg: types.MsgDepositToVault{
				VaultId:      &constants.Vault_Clob0,
				SubaccountId: &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewIntFromBigInt(
					new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						new(big.Int).SetUint64(1),
					),
				),
			},
			expectedErr: "Deposit amount is invalid",
		},
		"Failure: zero quote quantums": {
			msg: types.MsgDepositToVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewInt(0),
			},
			expectedErr: "Deposit amount is invalid",
		},
		"Failure: negative quote quantums": {
			msg: types.MsgDepositToVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewInt(-1),
			},
			expectedErr: "Deposit amount is invalid",
		},
		"Failure: invalid subaccount owner": {
			msg: types.MsgDepositToVault{
				VaultId: &constants.Vault_Clob0,
				SubaccountId: &satypes.SubaccountId{
					Owner:  "invalid-owner",
					Number: 0,
				},
				QuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: "subaccount id owner is an invalid address",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
