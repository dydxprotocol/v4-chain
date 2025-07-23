package constants

import (
	"math"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	// Short-term orders.
	Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 16},
	}
	Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 2},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Buy5_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     6,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Buy35_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     35,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20_BuilderCode = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		BuilderCodeParameters: &clobtypes.BuilderCodeParameters{
			BuilderAddress: Bob_Num0.Owner,
			FeePpm:         1000,
		},
	}
	Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 2, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 3, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id4_Clob2_Buy25_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 4, ClobPairId: 2},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 5, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 6, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 7, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 8, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     math.MaxUint64,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 9, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     15,
		Subticks:     45,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 19},
	}
	Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id10_Clob0_Sell35_Price15_GTB25 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     35,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 25},
	}
	Order_Alice_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20_000_000_000, // 200 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price102_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     102_000_000,    // $102
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price106_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     106_000_000,    // $106
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id0_Clob0_Sell10_Price10_GTB30 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
	}
	Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 2, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id2_Clob1_Buy10_Price10_GTB26 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 2, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 26},
	}
	Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 2, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     67,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id3_Clob1_Buy7_Price5 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 3, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     7,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 4, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     45,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id5_Clob1_Sell50_Price40_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 5, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     50,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 6, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     15,
		Subticks:     22,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
	}
	Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 7, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     35,
		Subticks:     math.MaxUint64,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
	}
	Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 8, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     15,
		Subticks:     25,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 31},
	}
	Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 9, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 31},
	}
	Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 31},
	}
	Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 32},
	}
	Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     6,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 32},
	}
	Order_Alice_Num1_Id10_Clob0_Buy7_Price30_GTB33 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     7,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 33},
	}
	Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 33},
	}
	Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     15,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 33},
	}
	Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 34},
	}
	Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 34},
	}
	Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 11, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     45,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 12, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 25},
	}
	Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 13, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     30,
		Subticks:     50,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 25},
	}
	Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 13, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     50,
		Subticks:     50,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
	}
	Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100,
		Subticks:     51_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id0_Clob0_Sell100_Price100000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100,
		Subticks:     1_000_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id3_Clob0_Sell100_Price100000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 3, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100,
		Subticks:     1_000_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Alice_Num1_Id5_Clob1_Buy10_Price15_GTB23 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 5, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 23},
	}
	Order_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000, // 1 BTC
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 0, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id0_Clob2_Sell10_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 0, ClobPairId: 2},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB18 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     11,
		Subticks:     16,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 18},
	}
	Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     11,
		Subticks:     16,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 2, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     12,
		Subticks:     13,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 3, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 5, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 6, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20,
		Subticks:     1000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 7, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20,
		Subticks:     10000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id8_Clob1_Sell5_Price10_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 8, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id8_Clob1_Sell20_Price10_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 8, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 8, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id9_Clob0_Sell20_Price1000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 9, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     1000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id10_Clob0_Sell20_Price10000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 10, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     10000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
	}
	Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 11, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 11, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id12_Clob0_Buy5_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 12, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id12_Clob0_Buy5_Price40_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 12, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 12, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 11, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     40,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 32},
	}
	Order_Bob_Num0_Id12_Clob0_Sell20_Price5_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 12, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id12_Clob0_Sell20_Price15_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 12, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id12_Clob0_Sell20_Price35_GTB32 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 12, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 32},
	}
	Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 13, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     35,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 30},
	}
	Order_Bob_Num0_Id14_Clob0_Sell10_Price10_GTB25 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 14, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 25},
	}
	Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     35,
		Subticks:     55,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 32},
	}
	Order_Bob_Num0_Id2_Clob0_Sell25_Price95_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     95,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price98_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     98_000_000,     // $98
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     99_000_000,     // $99
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20_000_000_000, // 200 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Bob_Num1_Id1_Clob1_Sell25_Price85_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25,
		Subticks:     85,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price5subticks_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price5000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     5_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25_000_000,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     70,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy110_Price500000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     110,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     110,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     80,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id0_Clob0_Buy10_Price50000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     110,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id2_Clob0_Sell5_Price10_GTB15 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 15},
	}
	Order_Carl_Num0_Id1_Clob0_Buy01BTC_Price49500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000,
		Subticks:     49_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     49_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49800_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     49_800_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     200_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id2_Clob0_Buy1BTC_Price50500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id1_Clob0_Buy1BTC_Price49999 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     49_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id1_Clob0_Buy1BTC_WithValidOrderRouter = clobtypes.Order{
		OrderId:            clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 1, ClobPairId: 0},
		Side:               clobtypes.Order_SIDE_BUY,
		Quantums:           100_000_000,
		Subticks:           49_000_000_000,
		GoodTilOneof:       &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
		OrderRouterAddress: AliceAccAddress.String(),
	}
	Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     50_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price49500 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 3, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25_000_000,
		Subticks:     49_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price49800 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 3, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25_000_000,
		Subticks:     49_800_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 3, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id4_Clob0_Buy05BTC_Price40000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 4, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     50_000_000,
		Subticks:     40_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id5_Clob0_Buy2BTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 5, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     200_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id2_Clob1_Buy10ETH_Price3000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 2, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000,
		Subticks:     3_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id3_Clob1_Buy1ETH_Price3000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 3, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     1_000_000_000,
		Subticks:     3_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id4_Clob1_Buy01ETH_Price3000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     3_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy10QtBTC_Price100000QuoteQt = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     100_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id0_Clob0_Buy10QtBTC_Price100001QuoteQt = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     100_001_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id0_Clob0_Sell1kQtBTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     1_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id1_Clob0_Sell1kQtBTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     1_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     99_000_000,     // $99
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price100_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     100_000_000,    // $100
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Carl_Num1_Id0_Clob0_Buy1kQtBTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     1_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num1_Id0_Clob0_Buy1kQtBTC_Price60000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     1_000_000,
		Subticks:     60_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     1_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     49_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_003_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id2_Clob0_Sell1BTC_Price49500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     49_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price49999_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     49_999_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_498_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id1_Clob0_Sell01BTC_Price50500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10_000_000,
		Subticks:     50_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price60000_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     60_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 11},
	}
	Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_498_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 11},
	}
	Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 11},
	}
	// Replacement for the above order
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB12 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 12},
	}
	Order_Dave_Num0_Id1_Clob3_Sell025ISO_Price50_GTB11 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 3},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     250_000_000,
		Subticks:     5_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 11},
	}
	Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 12},
	}
	Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 12},
	}
	Order_Dave_Num0_Id0_Clob0_Buy100BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Dave_Num0_Id0_Clob0_Buy100BTC_Price102_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     102_000_000,    // $102
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Dave_Num0_Id1_Clob0_Buy100BTC_Price104_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     104_000_000,    // $104
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 3, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     1_000_000_000,
		Subticks:     3_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3000 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     1_000_000_000,
		Subticks:     3_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3020 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     1_000_000_000,
		Subticks:     3_020_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3030 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     1_000_000_000,
		Subticks:     3_030_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     49_500_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     49_997_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     49_999_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     101_000_000,    // $101
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}
	Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price102_GTB20 = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10_000_000_000, // 100 BTC
		Subticks:     102_000_000,    // $102
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
	}

	// IOC orders.
	Order_Alice_Num0_Id1_Clob0_Buy5_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num0_Id1_Clob1_Buy5_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num0_Id1_Clob1_Sell5_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     5,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num0_Id1_Clob1_Buy10_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB21_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num0_Id0_Clob1_Buy10_Price15_GTB20_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000, // 1 BTC
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}
	Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     50_000_000, // 0.5 BTC
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
	}

	// IOC + RO orders.
	Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id1_Clob1_Buy10_Price15_GTB20_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id1_Clob0_Buy10_Price15_GTB20_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     15,
		Subticks:     500_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     110,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id0_Clob0_Buy110_Price50000_GTB21_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     110,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		ReduceOnly:   true,
	}

	// Reduce-only orders.
	Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20,
		Subticks:     30,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id3_Clob1_Buy30_Price35_GTB25_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 3, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     30,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 25},
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 4, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     15,
		Subticks:     20,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id5_Clob1_Sell10_Price15_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 5, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 6, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     10,
		Subticks:     5,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Bob_Num0_Id1_Clob0_Sell15_Price50_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     15,
		Subticks:     50,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 3, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     20,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		ReduceOnly:   true,
	}
	Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Carl_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
		ReduceOnly:   true,
	}
	Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 0, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 10},
		ReduceOnly:   true,
	}
	Order_Dave_Num0_Id2_Clob0_Sell25BTC_Price50000_GTB12_RO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Dave_Num0, ClientId: 2, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     25_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 12},
		ReduceOnly:   true,
	}

	// Post-only orders.
	Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     15,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 18},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 0},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     15,
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 18},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_PO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 1, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	Order_Alice_Num1_Id4_Clob1_Sell10_Price15_GTB20_PO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Alice_Num1, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     10,
		Subticks:     15,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
	Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO = clobtypes.Order{
		OrderId:      clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 4, ClobPairId: 1},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     20,
		Subticks:     35,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 22},
		TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
	}
)
