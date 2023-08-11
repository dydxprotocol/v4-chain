package constants

import (
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
)

var (
	// Long-term orders.
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 5},
	}
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
	}
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User1_Num0_Id2_Clob0_Sell02BTC_Price20_GTB15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20_000_000,
		Subticks:     20,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User1_Num0_Id1_Clob0_Sell02BTC_Price10_GTB15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20_000_000,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User1_Num0_Id3_Clob0_Buy02BTC_Price10_GTB15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20_000_000,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User1_Num1_Id1_Clob0_Buy02BTC_Price10_GTB15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20_000_000,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User1_Num1_Id2_Clob0_Sell02BTC_Price10_GTB15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20_000_000,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
	}
	LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   1,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     65,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 25},
	}
	LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     65,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 25},
	}
	LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     15,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_User1_Num1_Id1_Clob0_Sell50_Price30_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     50,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     50,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_Alice_Num1_Id2_Clob0_Buy10_Price40_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_User2_Num0_Id0_Clob0_Buy35_Price30_GTBT11 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     35,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 11},
	}
	LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     45,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Bob_Num0_Id2_Clob0_Buy15_Price5_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     15,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}

	// Conditional orders.
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
	}
	ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     50,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 30},
	}
	ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     50,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
	}
	ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}

	// Long-Term post-only orders.
	LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     65,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 25},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	LongTermOrder_User1_Num1_Id3_Clob0_Buy10_Price40_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
)
