package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgBatchCancel_ValidateBasic(t *testing.T) {
	oneOverMax := []uint32{}
	for i := uint32(0); i < types.MaxMsgBatchCancelBatchSize+1; i++ {
		oneOverMax = append(oneOverMax, i)
	}

	tests := map[string]struct {
		msg types.MsgBatchCancel
		err error
	}{
		"invalid subaccount": {
			msg: *types.NewMsgBatchCancel(
				constants.InvalidSubaccountIdNumber,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds: []uint32{
							0, 1, 2, 3,
						},
					},
				},
				10,
			),
			err: satypes.ErrInvalidSubaccountIdNumber,
		},
		"over 100 cancels in for one clob pair id": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds:  oneOverMax,
					},
				},
				10,
			),
			err: types.ErrInvalidBatchCancel,
		},
		"over 100 cancels split over two clob pair id": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds:  oneOverMax[:types.MaxMsgBatchCancelBatchSize/2+2],
					},
					{
						ClobPairId: 1,
						ClientIds:  oneOverMax[:types.MaxMsgBatchCancelBatchSize/2+2],
					},
				},
				10,
			),
			err: types.ErrInvalidBatchCancel,
		},
		"success: two clob pair id, 100 cancels": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds:  oneOverMax[:types.MaxMsgBatchCancelBatchSize/2],
					},
					{
						ClobPairId: 1,
						ClientIds:  oneOverMax[:types.MaxMsgBatchCancelBatchSize/2],
					},
				},
				10,
			),
			err: nil,
		},
		"success: one clob pair id, 100 cancels": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds:  oneOverMax[:types.MaxMsgBatchCancelBatchSize],
					},
				},
				10,
			),
			err: nil,
		},
		"duplicate clob pair ids": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds: []uint32{
							0, 1, 2, 3,
						},
					},
					{
						ClobPairId: 0,
						ClientIds: []uint32{
							2, 3, 4,
						},
					},
				},
				10,
			),
			err: types.ErrInvalidBatchCancel,
		},
		"duplicate client ids": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds: []uint32{
							0, 1, 2, 3, 1,
						},
					},
				},
				10,
			),
			err: types.ErrInvalidBatchCancel,
		},
		"zero batches in cancel batch": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{},
				10,
			),
			err: types.ErrInvalidBatchCancel,
		},
		"zero client ids in cancel batch": {
			msg: *types.NewMsgBatchCancel(
				constants.Alice_Num0,
				[]types.OrderBatch{
					{
						ClobPairId: 0,
						ClientIds:  []uint32{},
					},
				},
				10,
			),
			err: types.ErrInvalidBatchCancel,
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
