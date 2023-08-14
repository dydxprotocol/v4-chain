package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	OrderId_Alice_Num0_ClientId0_Clob0         = clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 0, ClobPairId: 0}
	OrderId_Alice_Num0_ClientId1_Clob0         = clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 1, ClobPairId: 0}
	OrderId_Alice_Num0_ClientId2_Clob0         = clobtypes.OrderId{SubaccountId: Alice_Num0, ClientId: 2, ClobPairId: 0}
	OrderId_Bob_Num0_ClientId0_Clob0           = clobtypes.OrderId{SubaccountId: Bob_Num0, ClientId: 0, ClobPairId: 0}
	LongTermOrderId_Alice_Num0_ClientId0_Clob0 = clobtypes.OrderId{
		SubaccountId: Alice_Num0,
		ClientId:     0,
		ClobPairId:   0,
		OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
	}
	LongTermOrderId_Alice_Num1_ClientId3_Clob1 = clobtypes.OrderId{
		SubaccountId: Alice_Num1,
		ClientId:     3,
		ClobPairId:   1,
		OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
	}

	// Invalid clobPairId
	InvalidClobPairId_Long_Term_Order = clobtypes.OrderId{
		SubaccountId: Alice_Num0,
		ClientId:     0,
		OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
		ClobPairId:   99999,
	}

	// Invalid order ids
	InvalidSubaccountIdOwner_OrderId = clobtypes.OrderId{
		SubaccountId: InvalidSubaccountIdOwner,
		ClientId:     0,
		ClobPairId:   0,
	}
	InvalidSubaccountIdNumber_OrderId = clobtypes.OrderId{
		SubaccountId: InvalidSubaccountIdNumber,
		ClientId:     0,
		ClobPairId:   0,
	}
)
