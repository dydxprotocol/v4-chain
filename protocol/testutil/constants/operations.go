package constants

import (
	"github.com/dydxprotocol/v4/x/clob/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(ValidEmptyMsgProposedOperations)
	ValidEmptyMsgProposedOperationsTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidProposedOperations)
	InvalidProposedOperationsTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

var (
	ValidEmptyMsgProposedOperations        = &types.MsgProposedOperations{Proposer: "foobar"}
	ValidEmptyMsgProposedOperationsTxBytes []byte
	// InvalidProposedOperations is invalid because the maker order for the match operation
	// has not been previously placed or placed in the set of operations preceeding the match.
	InvalidProposedOperations = &types.MsgProposedOperations{
		Proposer: "foobar",
		OperationsQueue: []types.Operation{
			{
				Operation: &types.Operation_OrderPlacement{
					OrderPlacement: &types.MsgPlaceOrder{
						Order: Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					},
				},
			},
			{Operation: &types.Operation_Match{
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
							TakerOrderHash: Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderHash().ToBytes(),
						},
					},
				},
			}},
		},
	}
	InvalidProposedOperationsTxBytes []byte
)
