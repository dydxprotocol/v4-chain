package types

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCancelOrder_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg MsgCancelOrder
		err error
	}{
		"invalid subaccountId owner": {
			msg: MsgCancelOrder{
				OrderId: OrderId{
					SubaccountId: satypes.SubaccountId{
						Owner:  "invalid_owner",
						Number: uint32(0),
					},
				},
			},
			err: satypes.ErrInvalidSubaccountIdOwner,
		},
		"invalid subaccountId number": {
			msg: MsgCancelOrder{
				OrderId: OrderId{
					SubaccountId: satypes.SubaccountId{
						Owner:  sample.AccAddress(),
						Number: uint32(999_999),
					},
				},
			},
			err: satypes.ErrInvalidSubaccountIdNumber,
		},
		"invalid 0 valued GoodTilBlock, short term order": {
			msg: *NewMsgCancelOrderShortTerm(OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(0),
				},
				OrderFlags: OrderIdFlags_ShortTerm,
			}, 0),
			err: ErrInvalidOrderGoodTilBlock,
		},
		"short term order valid": {
			msg: *NewMsgCancelOrderShortTerm(OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(0),
				},
				OrderFlags: OrderIdFlags_ShortTerm,
			}, 100),
		},
		"long term order invalid 0 valued GoodTilBlockTime": {
			msg: *NewMsgCancelOrderStateful(OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(0),
				},
				OrderFlags: OrderIdFlags_LongTerm,
			}, 0),
			err: ErrInvalidStatefulOrderGoodTilBlockTime,
		},
		"long term order valid": {
			msg: *NewMsgCancelOrderStateful(OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(0),
				},
				OrderFlags: OrderIdFlags_LongTerm,
			}, 100),
		},
		"conditional term order invalid 0 valued GoodTilBlockTime": {
			msg: *NewMsgCancelOrderStateful(OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(0),
				},
				OrderFlags: OrderIdFlags_Conditional,
			}, 0),
			err: ErrInvalidStatefulOrderGoodTilBlockTime,
		},
		"conditional term order valid": {
			msg: *NewMsgCancelOrderStateful(OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(0),
				},
				OrderFlags: OrderIdFlags_Conditional,
			}, 100),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
