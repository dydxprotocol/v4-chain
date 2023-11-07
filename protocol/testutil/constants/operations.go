package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(ValidEmptyMsgProposedOperations)
	ValidEmptyMsgProposedOperationsTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidProposedOperationsUnspecifiedOrderRemovalReason)
	InvalidProposedOperationsUnspecifiedOrderRemovalReasonTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(
		TestTxBuilder.GetTx())
}

var (
	ValidEmptyMsgProposedOperations        = &types.MsgProposedOperations{}
	ValidEmptyMsgProposedOperationsTxBytes []byte
	// InvalidProposedOperationsUnspecifiedOrderRemovalReason is invalid because the order removal reason is
	// unspecified.
	InvalidProposedOperationsUnspecifiedOrderRemovalReason = &types.MsgProposedOperations{
		OperationsQueue: []types.OperationRaw{
			{
				Operation: &types.OperationRaw_OrderRemoval{
					OrderRemoval: &types.OrderRemoval{
						OrderId:       LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
						RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNSPECIFIED,
					},
				},
			},
		},
	}
	InvalidProposedOperationsUnspecifiedOrderRemovalReasonTxBytes []byte
)
