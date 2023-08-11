package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/require"
)

var (
	order   = constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	orderId = constants.OrderId_Alice_Num0_ClientId0_Clob0
)

func TestStatefulOrderPlacementEvent_Success(t *testing.T) {
	statefulOrderPlacementEvent := events.NewStatefulOrderPlacementEvent(order)
	expectedStatefulOrderEventProto := &events.StatefulOrderEvent{
		Event: &events.StatefulOrderEvent_OrderPlace{
			OrderPlace: &events.StatefulOrderEvent_StatefulOrderPlacement{
				Order: &order,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderPlacementEvent)
}

func TestStatefulOrderCancelationEvent_Success(t *testing.T) {
	statefulOrderCancelationEvent := events.NewStatefulOrderCancelationEvent(orderId)
	expectedStatefulOrderEventProto := &events.StatefulOrderEvent{
		Event: &events.StatefulOrderEvent_OrderCancel{
			OrderCancel: &events.StatefulOrderEvent_StatefulOrderCancelation{
				CanceledOrderId: &orderId,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderCancelationEvent)
}

func TestStatefulOrderExpirationEvent_Success(t *testing.T) {
	statefulOrderExpirationEvent := events.NewStatefulOrderExpirationEvent(orderId)
	expectedStatefulOrderEventProto := &events.StatefulOrderEvent{
		Event: &events.StatefulOrderEvent_OrderExpiration{
			OrderExpiration: &events.StatefulOrderEvent_StatefulOrderExpiration{
				ExpiredOrderId: &orderId,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderExpirationEvent)
}
