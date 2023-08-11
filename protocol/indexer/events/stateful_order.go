package events

import (
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
)

func NewStatefulOrderPlacementEvent(
	order clobtypes.Order,
) *StatefulOrderEvent {
	orderPlace := StatefulOrderEvent_StatefulOrderPlacement{
		Order: &order,
	}
	return &StatefulOrderEvent{
		Event: &StatefulOrderEvent_OrderPlace{
			OrderPlace: &orderPlace,
		},
	}
}

func NewStatefulOrderCancelationEvent(
	canceledOrderId clobtypes.OrderId,
) *StatefulOrderEvent {
	orderCancel := StatefulOrderEvent_StatefulOrderCancelation{
		CanceledOrderId: &canceledOrderId,
	}
	return &StatefulOrderEvent{
		Event: &StatefulOrderEvent_OrderCancel{
			OrderCancel: &orderCancel,
		},
	}
}

func NewStatefulOrderExpirationEvent(
	expiredOrderId clobtypes.OrderId,
) *StatefulOrderEvent {
	orderExpire := StatefulOrderEvent_StatefulOrderExpiration{
		ExpiredOrderId: &expiredOrderId,
	}
	return &StatefulOrderEvent{
		Event: &StatefulOrderEvent_OrderExpiration{
			OrderExpiration: &orderExpire,
		},
	}
}
