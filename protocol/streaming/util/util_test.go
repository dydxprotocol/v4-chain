package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	pv1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	stypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/util"
)

func _ToPtr[V any](v V) *V {
	return &v
}

func TestGetOffChainUpdateV1SubaccountIdNumber(t *testing.T) {
	subaccountIdNumber := uint32(1337)
	orderId := pv1types.IndexerOrderId{
		SubaccountId: pv1types.IndexerSubaccountId{
			Owner:  "foo",
			Number: uint32(subaccountIdNumber),
		},
		ClientId:   0,
		OrderFlags: 0,
		ClobPairId: 0,
	}
	order := pv1types.IndexerOrder{
		OrderId:  orderId,
		Side:     pv1types.IndexerOrder_SIDE_BUY,
		Quantums: uint64(10 ^ 6),
		Subticks: 1,
		GoodTilOneof: &pv1types.IndexerOrder_GoodTilBlock{
			GoodTilBlock: 10 ^ 9,
		},
		TimeInForce:                     10 ^ 9,
		ReduceOnly:                      false,
		ClientMetadata:                  0,
		ConditionType:                   pv1types.IndexerOrder_CONDITION_TYPE_UNSPECIFIED,
		ConditionalOrderTriggerSubticks: 0,
	}
	newOrder := order
	newOrder.Quantums += 10 ^ 6

	orderPlaceTime := time.Now()
	fillQuantums := uint64(1988)

	tests := map[string]struct {
		update ocutypes.OffChainUpdateV1
		id     uint32
		err    error
	}{
		"OrderPlace": {
			update: ocutypes.OffChainUpdateV1{
				UpdateMessage: &ocutypes.OffChainUpdateV1_OrderPlace{
					OrderPlace: &ocutypes.OrderPlaceV1{
						Order:           &order,
						PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
						TimeStamp:       _ToPtr(orderPlaceTime),
					},
				},
			},
			id:  subaccountIdNumber,
			err: nil,
		},
		"OrderRemove": {
			update: ocutypes.OffChainUpdateV1{
				UpdateMessage: &ocutypes.OffChainUpdateV1_OrderRemove{
					OrderRemove: &ocutypes.OrderRemoveV1{
						RemovedOrderId: &orderId,
						Reason:         stypes.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
						RemovalStatus:  ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_CANCELED,
						TimeStamp:      _ToPtr(orderPlaceTime.Add(1 * time.Second)),
					},
				},
			},
			id:  subaccountIdNumber,
			err: nil,
		},
		"OrderUpdate": {
			update: ocutypes.OffChainUpdateV1{
				UpdateMessage: &ocutypes.OffChainUpdateV1_OrderUpdate{
					OrderUpdate: &ocutypes.OrderUpdateV1{
						OrderId:             &orderId,
						TotalFilledQuantums: fillQuantums,
					},
				},
			},
			id:  subaccountIdNumber,
			err: nil,
		},
		"OrderReplace": {
			update: ocutypes.OffChainUpdateV1{
				UpdateMessage: &ocutypes.OffChainUpdateV1_OrderReplace{
					OrderReplace: &ocutypes.OrderReplaceV1{
						OldOrderId:      &orderId,
						Order:           &newOrder,
						PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_OPENED,
						TimeStamp:       _ToPtr(orderPlaceTime.Add(3 * time.Second)),
					},
				},
			},
			id:  subaccountIdNumber,
			err: nil,
		},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := util.GetOffChainUpdateV1SubaccountIdNumber(testCase.update)
			fmt.Println("expected", id)
			require.Equal(t, err, testCase.err)
			require.Equal(t, id, testCase.id)
		})
	}
}
