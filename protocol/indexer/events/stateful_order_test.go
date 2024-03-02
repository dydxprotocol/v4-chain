package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

var (
	order          = constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	indexerOrder   = v1.OrderToIndexerOrder(order)
	orderId        = constants.OrderId_Alice_Num0_ClientId0_Clob0
	indexerOrderId = v1.OrderIdToIndexerOrderId(orderId)
	reason         = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_REPLACED
)

func TestLongTermOrderPlacementEvent_Success(t *testing.T) {
	longTermOrderPlacementEvent := events.NewLongTermOrderPlacementEvent(order)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_LongTermOrderPlacement{
			LongTermOrderPlacement: &events.StatefulOrderEventV1_LongTermOrderPlacementV1{
				Order: &indexerOrder,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, longTermOrderPlacementEvent)
}

func TestStatefulOrderRemovalEvent_Success(t *testing.T) {
	statefulOrderRemovalEvent := events.NewStatefulOrderRemovalEvent(orderId, reason)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_OrderRemoval{
			OrderRemoval: &events.StatefulOrderEventV1_StatefulOrderRemovalV1{
				RemovedOrderId: &indexerOrderId,
				Reason:         reason,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, statefulOrderRemovalEvent)
}

func TestConditionalOrderPlacementEvent_Success(t *testing.T) {
	conditionalOrderPlacementEvent := events.NewConditionalOrderPlacementEvent(order)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_ConditionalOrderPlacement{
			ConditionalOrderPlacement: &events.StatefulOrderEventV1_ConditionalOrderPlacementV1{
				Order: &indexerOrder,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, conditionalOrderPlacementEvent)
}

func TestConditionalOrderTriggeredEvent_Success(t *testing.T) {
	conditionalOrderTriggeredEvent := events.NewConditionalOrderTriggeredEvent(orderId)
	expectedStatefulOrderEventProto := &events.StatefulOrderEventV1{
		Event: &events.StatefulOrderEventV1_ConditionalOrderTriggered{
			ConditionalOrderTriggered: &events.StatefulOrderEventV1_ConditionalOrderTriggeredV1{
				TriggeredOrderId: &indexerOrderId,
			},
		},
	}
	require.Equal(t, expectedStatefulOrderEventProto, conditionalOrderTriggeredEvent)
}
