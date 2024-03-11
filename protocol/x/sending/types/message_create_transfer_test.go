package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateTransfer_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgCreateTransfer
		err  error
	}{
		{
			name: "Invalid sender owner",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender: satypes.SubaccountId{
						Owner:  "invalid_owner",
						Number: uint32(0),
					},
					Recipient: constants.Carl_Num0,
				},
			},
			err: satypes.ErrInvalidSubaccountIdOwner,
		},
		{
			name: "Invalid recipient owner",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender: constants.Carl_Num0,
					Recipient: satypes.SubaccountId{
						Owner:  "invalid_owner",
						Number: uint32(0),
					},
				},
			},
			err: satypes.ErrInvalidSubaccountIdOwner,
		},
		{
			name: "Invalid sender number",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender: satypes.SubaccountId{
						Owner:  sample.AccAddress(),
						Number: uint32(999_999),
					},
					Recipient: constants.Carl_Num0,
				},
			},
			err: satypes.ErrInvalidSubaccountIdNumber,
		},
		{
			name: "Invalid recipient number",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender: constants.Carl_Num0,
					Recipient: satypes.SubaccountId{
						Owner:  sample.AccAddress(),
						Number: uint32(999_999),
					},
				},
			},
			err: satypes.ErrInvalidSubaccountIdNumber,
		},
		{
			name: "Valid address",
			msg: types.MsgCreateTransfer{
				Transfer: &constants.Transfer_Carl_Num0_Dave_Num0_Quote_500,
			},
		},
		{
			name: "Same sender and recipient",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender:    constants.Carl_Num0,
					Recipient: constants.Carl_Num0,
					AssetId:   assettypes.AssetUsdc.Id,
					Amount:    uint64(500_000_000),
				},
			},
			err: types.ErrSenderSameAsRecipient,
		},
		{
			name: "Non-USDC asset transfer not supported",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender:    constants.Carl_Num0,
					Recipient: constants.Dave_Num0,
					AssetId:   uint32(1),
					Amount:    uint64(100),
				},
			},
			err: types.ErrNonUsdcAssetTransferNotImplemented,
		},
		{
			name: "Invalid amount",
			msg: types.MsgCreateTransfer{
				Transfer: &types.Transfer{
					Sender:    constants.Carl_Num0,
					Recipient: constants.Dave_Num0,
					AssetId:   assettypes.AssetUsdc.Id,
					Amount:    uint64(0),
				},
			},
			err: types.ErrInvalidTransferAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
