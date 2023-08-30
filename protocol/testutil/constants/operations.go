package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(ValidEmptyMsgProposedOperations)
	ValidEmptyMsgProposedOperationsTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidProposedOperations)
	InvalidProposedOperationsTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

var (
	ValidEmptyMsgProposedOperations        = &types.MsgProposedOperations{}
	ValidEmptyMsgProposedOperationsTxBytes []byte
	// InvalidProposedOperations is invalid because the maker order for the match operation
	// does not have a corresponding order placement operation before it in the operations queue.
	InvalidProposedOperations = &types.MsgProposedOperations{
		OperationsQueue: []types.OperationRaw{
			{
				Operation: &types.OperationRaw_Match{
					Match: &types.ClobMatch{
						Match: &types.ClobMatch_MatchOrders{
							MatchOrders: &types.MatchOrders{
								TakerOrderId: Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								Fills: []types.MakerFill{
									{
										MakerOrderId: OrderId_Alice_Num0_ClientId0_Clob0,
										FillAmount:   100_000_000,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	InvalidProposedOperationsTxBytes []byte
)
