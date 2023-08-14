package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	CancelOrder_Alice_Num0_Id12_Clob0_GTB5 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     12,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 5},
	}
	CancelOrder_Alice_Num1_Id13_Clob0_GTB25 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     13,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 25},
	}
	CancelOrder_Alice_Num1_Id13_Clob0_GTB30 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     13,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 30},
	}
	CancelOrder_Alice_Num1_Id13_Clob0_GTB35 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     13,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 35},
	}
	CancelOrder_Bob_Num0_Id2_Clob1_GTB5 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   1,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 5},
	}
	CancelOrder_Alice_Num0_Id10_Clob0_GTB20 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     10,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 20},
	}
	CancelOrder_Alice_Num1_Id11_Clob1_GTB20 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     11,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   1,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 20},
	}
	CancelOrder_Bob_Num1_Id11_Clob1_GTB20 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num1,
			ClientId:     11,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   1,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 20},
	}

	// Long Term Order Cancellations
	CancelLongTermOrder_Alice_Num1_Id1_Clob0_GTBT_20 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 25},
	}
	CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 5},
	}
	CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	CancelLongTermOrder_Bob_Num0_Id0_Clob0_GTBT5 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 5},
	}
	// Conditional Order Cancellations
	CancelConditionalOrder_Alice_Num1_Id0_Clob0_GTBT15 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15 = clobtypes.MsgCancelOrder{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		GoodTilOneof: &clobtypes.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
)
