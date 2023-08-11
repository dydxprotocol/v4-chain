package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	makerOrder            = constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	takerOrder            = constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15
	liquidationTakerOrder = constants.LiquidationOrder_Carl_Num0_Clob0_Buy3_Price50_BTC
	fillAmount            = satypes.BaseQuantums(5)
)

func TestNewOrderFillEvent_Success(t *testing.T) {
	orderFillEvent := events.NewOrderFillEvent(makerOrder, takerOrder, fillAmount)

	expectedOrderFillEventProto := &events.OrderFillEvent{
		MakerOrder: makerOrder,
		TakerOrder: &events.OrderFillEvent_Order{
			Order: &takerOrder,
		},
		FillAmount: fillAmount.ToUint64(),
	}
	require.Equal(t, expectedOrderFillEventProto, orderFillEvent)
}

func TestNewLiquidationOrderFillEvent_Success(t *testing.T) {
	var matchableTakerOrder clobtypes.MatchableOrder = &liquidationTakerOrder
	liquidationOrderFillEvent := events.NewLiquidationOrderFillEvent(
		makerOrder,
		matchableTakerOrder,
		fillAmount,
	)

	expectedLiquidationOrder := events.LiquidationOrder{
		Liquidated:  liquidationTakerOrder.GetSubaccountId(),
		ClobPairId:  liquidationTakerOrder.GetClobPairId().ToUint32(),
		PerpetualId: liquidationTakerOrder.MustGetLiquidatedPerpetualId(),
		TotalSize:   uint64(liquidationTakerOrder.GetBaseQuantums()),
		IsBuy:       liquidationTakerOrder.IsBuy(),
		Subticks:    uint64(liquidationTakerOrder.GetOrderSubticks()),
	}
	expectedOrderFillEventProto := &events.OrderFillEvent{
		MakerOrder: makerOrder,
		TakerOrder: &events.OrderFillEvent_LiquidationOrder{
			LiquidationOrder: &expectedLiquidationOrder,
		},
		FillAmount: fillAmount.ToUint64(),
	}
	require.Equal(t, expectedOrderFillEventProto, liquidationOrderFillEvent)
}
