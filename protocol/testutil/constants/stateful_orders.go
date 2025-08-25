package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price5_GTBT5 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
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
	LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   1,
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
	LongTermOrder_Carl_Num0_Id0_Clob0_WithOrderRouterAddress = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:               clobtypes.Order_SIDE_BUY,
		Quantums:           5,
		Subticks:           10,
		GoodTilOneof:       &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		OrderRouterAddress: AliceAccAddress.String(),
	}

	LongTermOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTBT5 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
	}
	LongTermOrder_Alice_Num1_Id1_Clob0_Buy02BTC_Price10_GTB15 = clobtypes.Order{
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
	LongTermOrder_Alice_Num1_Id2_Clob0_Sell02BTC_Price10_GTB15 = clobtypes.Order{
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
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15 = clobtypes.Order{
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
	LongTermOrder_Alice_Num1_Id4_Clob0_Buy10_Price45_GTBT20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     4,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     45,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
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
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_Alice_Num0_Id1_Clob0_Buy1BTC_Price50000_GTBT15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
	}
	LongTermOrder_Bob_Num0_Id0_Clob0_Sell2_Price5_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     2,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Bob_Num0_Id0_Clob0_Sell5_Price5_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Bob_Num0_Id1_Clob0_Sell5_Price10_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15 = clobtypes.Order{
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
	LongTermOrder_Bob_Num0_Id0_Clob1_Buy25_Price30_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   1,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Bob_Num0_Id0_Clob0_Buy35_Price30_GTBT11 = clobtypes.Order{
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
	LongTermOrder_Bob_Num1_Id3_Clob0_Buy10_Price40_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num1,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     49_500_000_000,
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
	LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}
	LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_001_000_000,
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
	LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
	}

	// Conditional orders.
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        20,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num1_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        20,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 10,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 25,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        15,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        15,
		Subticks:                        25,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 25,
	}
	ConditionalOrder_Alice_Num1_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        15,
		Subticks:                        25,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 25,
	}
	ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        15,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        15,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 5,
	}
	ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        20,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        20,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        20,
		Subticks:                        20,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num1_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        20,
		Subticks:                        20,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        20,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 10,
	}
	ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        25,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        25,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 25,
	}
	ConditionalOrder_Alice_Num1_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        25,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 25,
	}
	ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        25,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        25,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 10,
	}
	ConditionalOrder_Alice_Num0_Id3_Clob1_Buy25_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     3,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        25,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id1_Clob1_Buy5_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit30 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 30,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob1_Sell5_Price10_GTBT15_StopLoss20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob1_Sell5_Price10_GTBT15_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   1,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price50_GTBT10_StopLoss51_IOC = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        50,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_IOC,
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 51,
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
	ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 15,
	}
	ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTBT30_TakeProfit10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        50,
		Subticks:                        5,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 30},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 10,
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
	ConditionalOrder_Alice_Num1_Id1_Clob0_Sell50_Price5_GTB30_TakeProfit20 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        50,
		Subticks:                        5,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 20,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49700 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_700_000_000,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_995_000_000,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_999_000_000,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_IOC = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_999_000_000,
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_PO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_999_000_000,
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 50_001_000_000,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 50_005_000_000,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50300 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 50_300_000_000,
	}
	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 50_001_000_000,
	}
	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 50_005_000_000,
	}
	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50300 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 50_300_000_000,
	}
	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49700 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 49_700_000_000,
	}
	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 49_995_000_000,
	}
	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 49_999_000_000,
	}
	ConditionalOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
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

	ConditionalOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO_SL_15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        10,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 15_000_000,
	}
	ConditionalOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10_SL_15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        20,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 15_000_000,
	}
	ConditionalOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15_SL_15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        50,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 15_000_000,
	}
	ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_SL_15 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 15_000_000,
	}
	ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_50003 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        100_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 50_003_000_000,
	}
	ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003 = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        50_000_000, // 0.5 BTC
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 50_003_000_000,
	}
	ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        50_000_000, // 0.5 BTC
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_IOC,
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 50_003_000_000,
	}
	ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
			ClobPairId:   0,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        50_000_000, // 0.5 BTC
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 49_999_000_000,
	}

	// Conditional IOC RO orders.
	ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Alice_Num1,
			ClientId:     1,
			ClobPairId:   0,
			OrderFlags:   clobtypes.OrderIdFlags_Conditional,
		},
		Side:                            clobtypes.Order_SIDE_SELL,
		Quantums:                        50_000_000, // 0.5 BTC
		Subticks:                        500_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
		TimeInForce:                     clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:                      true,
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: 50_001_000_000,
	}

	// Long-Term post-only orders.
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_PO = clobtypes.Order{
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
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15_PO = clobtypes.Order{
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
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
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
	LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_PO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10_PO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Dave_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}

	// Long-Term reduce-only orders.
	LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     2,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}

	// Long-Term Immediate Or Cancel Orders.
	LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_IOC = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}

	// TWAP orders.
	TwapOrder_Bob_Num0_Id1_Clob0_Buy10_Price35_GTB20_RO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       300,
			Interval:       30,
			PriceTolerance: 0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ReduceOnly:   true,
	}

	TwapOrder_Bob_Num0_Id1_Clob0_Buy1000_Price35_GTB20_RO = clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: Bob_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_Twap,
			ClobPairId:   0,
		},
		TwapParameters: &clobtypes.TwapParameters{
			Duration:       300,
			Interval:       30,
			PriceTolerance: 0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     1000,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		ReduceOnly:   true,
	}
)
