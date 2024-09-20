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

func TestMsgWithdrawFromMegavault_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgWithdrawFromMegavault
		expectedErr string
	}{
		"Success": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(1),
				},
				MinQuoteQuantums: dtypes.NewInt(1),
			},
		},
		"Success: zero quote quantums": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(1),
				},
				MinQuoteQuantums: dtypes.NewInt(0),
			},
		},
		"Success: max uint64 quote quantums": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(1),
				},
				MinQuoteQuantums: dtypes.NewIntFromUint64(math.MaxUint64),
			},
		},
		"Failure: quote quantums greater than max uint64": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(1),
				},
				MinQuoteQuantums: dtypes.NewIntFromBigInt(
					new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						new(big.Int).SetUint64(1),
					),
				),
			},
			expectedErr: types.ErrInvalidQuoteQuantums.Error(),
		},
		"Failure: negative quote quantums": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(1),
				},
				MinQuoteQuantums: dtypes.NewInt(-1),
			},
			expectedErr: types.ErrInvalidQuoteQuantums.Error(),
		},
		"Failure: zero shares": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(0),
				},
				MinQuoteQuantums: dtypes.NewInt(0),
			},
			expectedErr: types.ErrNonPositiveShares.Error(),
		},
		"Failure: negative shares": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: constants.Alice_Num0,
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(-1),
				},
				MinQuoteQuantums: dtypes.NewInt(0),
			},
			expectedErr: types.ErrNonPositiveShares.Error(),
		},
		"Failure: invalid subaccount owner": {
			msg: types.MsgWithdrawFromMegavault{
				SubaccountId: satypes.SubaccountId{
					Owner:  "invalid-owner",
					Number: 0,
				},
				Shares: types.NumShares{
					NumShares: dtypes.NewInt(1),
				},
				MinQuoteQuantums: dtypes.NewInt(0),
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
