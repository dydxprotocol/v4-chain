package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/protocol/v1"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/require"
)

var (
	order          = constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	indexerOrder   = v1.OrderToIndexerOrder(order)
	orderId        = constants.OrderId_Alice_Num0_ClientId0_Clob0
	indexerOrderId = v1.OrderIdToIndexerOrderId(orderId)
)

func TestStatefulOrderPlacementEvent_Success(t *testing.T) {
	statefulOrderPlacementEvent := events.NewStatefulOrderPlacementEvent(order)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_OrderPlace{
			OrderPlace: &events.StatefulOrderEventV1_StatefulOrderPlacementV1{
				Order: &indexerOrder,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderPlacementEvent)
}

func TestStatefulOrderCancelationEvent_Success(t *testing.T) {
	statefulOrderCancelationEvent := events.NewStatefulOrderCancelationEvent(orderId)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_OrderCancel{
			OrderCancel: &events.StatefulOrderEventV1_StatefulOrderCancelationV1{
				CanceledOrderId: &indexerOrderId,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderCancelationEvent)
}

func TestStatefulOrderExpirationEvent_Success(t *testing.T) {
	statefulOrderExpirationEvent := events.NewStatefulOrderExpirationEvent(orderId)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_OrderExpiration{
			OrderExpiration: &events.StatefulOrderEventV1_StatefulOrderExpirationV1{
				ExpiredOrderId: &indexerOrderId,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderExpirationEvent)
}
