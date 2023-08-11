package events

import (
	"github.com/dydxprotocol/v4/indexer/protocol/v1"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
)

func NewStatefulOrderPlacementEvent(
	order clobtypes.Order,
) *StatefulOrderEventV1 {
	indexerOrder := v1.OrderToIndexerOrder(order)
	orderPlace := StatefulOrderEventV1_StatefulOrderPlacementV1{
		Order: &indexerOrder,
	}
	return &StatefulOrderEventV1{
		Event: &StatefulOrderEventV1_OrderPlace{
			OrderPlace: &orderPlace,
		},
	}
}

func NewStatefulOrderCancelationEvent(
	canceledOrderId clobtypes.OrderId,
) *StatefulOrderEventV1 {
	orderId := v1.OrderIdToIndexerOrderId(canceledOrderId)
	orderCancel := StatefulOrderEventV1_StatefulOrderCancelationV1{
		CanceledOrderId: &orderId,
	}
	return &StatefulOrderEventV1{
		Event: &StatefulOrderEventV1_OrderCancel{
			OrderCancel: &orderCancel,
		},
	}
}

func NewStatefulOrderExpirationEvent(
	expiredOrderId clobtypes.OrderId,
) *StatefulOrderEventV1 {
	orderId := v1.OrderIdToIndexerOrderId(expiredOrderId)
	orderExpire := StatefulOrderEventV1_StatefulOrderExpirationV1{
		ExpiredOrderId: &orderId,
	}
	return &StatefulOrderEventV1{
		Event: &StatefulOrderEventV1_OrderExpiration{
			OrderExpiration: &orderExpire,
		},
	}
}
