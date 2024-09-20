package types_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRetrieveFromVault_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgRetrieveFromVault
		expectedErr error
	}{
		"Success": {
			msg: types.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(1),
			},
		},
		"Success: max uint64 quote quantums": {
			msg: types.MsgRetrieveFromVault{
				Authority:     constants.GovAuthority,
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewIntFromUint64(math.MaxUint64),
			},
		},
		"Failure: quote quantums greater than max uint64": {
			msg: types.MsgRetrieveFromVault{
				Authority: constants.GovAuthority,
				VaultId:   constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewIntFromBigInt(
					new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						new(big.Int).SetUint64(1),
					),
				),
			},
			expectedErr: types.ErrInvalidQuoteQuantums,
		},
		"Failure: zero quote quantums": {
			msg: types.MsgRetrieveFromVault{
				Authority:     constants.GovAuthority,
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(0),
			},
			expectedErr: types.ErrInvalidQuoteQuantums,
		},
		"Failure: negative quote quantums": {
			msg: types.MsgRetrieveFromVault{
				Authority:     constants.GovAuthority,
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(-1),
			},
			expectedErr: types.ErrInvalidQuoteQuantums,
		},
		"Failure: empty authority": {
			msg: types.MsgRetrieveFromVault{
				Authority:     "",
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(0),
			},
			expectedErr: types.ErrInvalidAuthority,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
