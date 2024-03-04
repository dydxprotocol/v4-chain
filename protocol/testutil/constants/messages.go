package constants

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

func init() {
	// This package does not contain the `app/config` package in its import chain, and therefore needs to call
	// SetAddressPrefixes() explicitly in order to set the `dydx` address prefixes.
	config.SetAddressPrefixes()

	_ = TestTxBuilder.SetMsgs(Msg_PlaceOrder)
	Msg_PlaceOrder_TxBtyes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(Msg_CancelOrder)
	Msg_CancelOrder_TxBtyes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(Msg_BatchCancel)
	Msg_BatchCancel_TxBtyes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(Msg_Send)
	Msg_Send_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(Msg_Transfer_Invalid_SameSenderAndRecipient)
	Msg_Transfer_Invalid_SameSenderAndRecipient_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(Msg_Send, Msg_Transfer)
	Msg_SendAndTransfer_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

var (
	Msg_CancelOrder = &clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			ClientId:     0,
			SubaccountId: Alice_Num0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 10},
	}
	Msg_CancelOrder_TxBtyes  []byte
	Msg_CancelOrder_LongTerm = &clobtypes.MsgCancelOrder{
		OrderId:      LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 20},
	}
	Msg_CancelOrder_Conditional = &clobtypes.MsgCancelOrder{
		OrderId:      ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 20},
	}

	Msg_PlaceOrder = &clobtypes.MsgPlaceOrder{
		Order: Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
	}
	Msg_PlaceOrder_TxBtyes []byte

	Msg_BatchCancel = &clobtypes.MsgBatchCancel{
		SubaccountId: Alice_Num0,
		ShortTermCancels: []clobtypes.OrderBatch{
			{
				ClobPairId: 0,
				ClientIds:  []uint32{0},
			},
		},
		GoodTilBlock: 5,
	}
	Msg_BatchCancel_TxBtyes []byte

	Msg_PlaceOrder_LongTerm = &clobtypes.MsgPlaceOrder{
		Order: LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
	}
	Msg_PlaceOrder_Conditional = &clobtypes.MsgPlaceOrder{
		Order: ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
	}
	Msg_Transfer = &sendingtypes.MsgCreateTransfer{
		Transfer: &sendingtypes.Transfer{
			Sender:    Carl_Num0,
			Recipient: Dave_Num0,
			AssetId:   assettypes.AssetUsdc.Id,
			Amount:    500_000_000, // $500
		},
	}
	Msg_Transfer_Invalid_SameSenderAndRecipient = &sendingtypes.MsgCreateTransfer{
		Transfer: &sendingtypes.Transfer{
			Sender:    Alice_Num0,
			Recipient: Alice_Num0,
			AssetId:   assettypes.AssetUsdc.Id,
			Amount:    500_000_000, // $500
		},
	}
	Msg_Transfer_Invalid_SameSenderAndRecipient_TxBytes []byte

	Msg_Send = &banktypes.MsgSend{
		FromAddress: AliceAccAddress.String(),
		ToAddress:   BobAccAddress.String(),
		Amount: sdk.Coins{sdk.Coin{
			Denom:  "foo",
			Amount: sdkmath.OneInt(),
		}},
	}
	Msg_Send_TxBytes []byte

	Msg_SendAndTransfer_TxBytes []byte
)
