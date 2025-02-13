package streaming_test

import (
	"testing"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	streaming "github.com/dydxprotocol/v4-chain/protocol/streaming"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func OpenOrder(
	order *v1types.IndexerOrder,
	timestamp *time.Time,
) ocutypes.OffChainUpdateV1 {
	return ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderPlace{
			OrderPlace: &ocutypes.OrderPlaceV1{
				Order:           order,
				PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_OPENED,
				TimeStamp:       timestamp,
			},
		},
	}
}

func CancelOrder(
	removedOrderId *v1types.IndexerOrderId,
	timestamp *time.Time,
) ocutypes.OffChainUpdateV1 {
	return ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderRemove{
			OrderRemove: &ocutypes.OrderRemoveV1{
				RemovedOrderId: removedOrderId,
				Reason:         sharedtypes.OrderRemovalReason(ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_CANCELED),
				RemovalStatus:  ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_CANCELED,
				TimeStamp:      timestamp,
			},
		},
	}
}

func ReplaceOrder(
	oldOrderId *v1types.IndexerOrderId,
	newOrder *v1types.IndexerOrder,
	timestamp *time.Time,
) ocutypes.OffChainUpdateV1 {
	return ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderReplace{
			OrderReplace: &ocutypes.OrderReplaceV1{
				OldOrderId:      oldOrderId,
				Order:           newOrder,
				PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_OPENED,
				TimeStamp:       timestamp,
			},
		},
	}
}

func UpdateOrder(orderId *v1types.IndexerOrderId, totalFilledQuantums uint64) ocutypes.OffChainUpdateV1 {
	return ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderUpdate{
			OrderUpdate: &ocutypes.OrderUpdateV1{
				OrderId:             orderId,
				TotalFilledQuantums: totalFilledQuantums,
			},
		},
	}
}

func toStreamUpdate(snapshot bool, offChainUpdates ...ocutypes.OffChainUpdateV1) clobtypes.StreamUpdate {
	return clobtypes.StreamUpdate{
		BlockHeight: uint32(0),
		ExecMode:    uint32(sdktypes.ExecModeFinalize),
		UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
			OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
				Updates:  offChainUpdates,
				Snapshot: snapshot,
			},
		},
	}
}

func NewStreamOrderbookFill(
	blockHeight uint32,
	execMode uint32,
) *clobtypes.StreamUpdate {
	return &clobtypes.StreamUpdate{
		BlockHeight: blockHeight,
		ExecMode:    execMode,
		UpdateMessage: &clobtypes.StreamUpdate_OrderFill{
			OrderFill: nil,
		},
	}
}

func NewStreamTakerOrder(
	blockHeight uint32,
	execMode uint32,
	order *clobtypes.Order,
	remainingQuantums uint64,
	optimisticallyFilledQuantums uint64,
) *clobtypes.StreamUpdate {
	return &clobtypes.StreamUpdate{
		BlockHeight: blockHeight,
		ExecMode:    execMode,
		UpdateMessage: &clobtypes.StreamUpdate_TakerOrder{
			TakerOrder: &clobtypes.StreamTakerOrder{
				TakerOrder: &clobtypes.StreamTakerOrder_Order{
					Order: order,
				},
				TakerOrderStatus: &clobtypes.StreamTakerOrderStatus{
					OrderStatus:                  uint32(clobtypes.Success),
					RemainingQuantums:            remainingQuantums,
					OptimisticallyFilledQuantums: optimisticallyFilledQuantums,
				},
			},
		},
	}
}

func NewSubaccountUpdate(
	blockHeight uint32,
	execMode uint32,
	subaccountId *satypes.SubaccountId,
) *clobtypes.StreamUpdate {
	return &clobtypes.StreamUpdate{
		BlockHeight: blockHeight,
		ExecMode:    execMode,
		UpdateMessage: &clobtypes.StreamUpdate_SubaccountUpdate{
			SubaccountUpdate: &satypes.StreamSubaccountUpdate{
				SubaccountId:              subaccountId,
				UpdatedPerpetualPositions: []*satypes.SubaccountPerpetualPosition{},
				UpdatedAssetPositions:     []*satypes.SubaccountAssetPosition{},
				Snapshot:                  false,
			},
		},
	}
}

func NewPriceUpdate(
	blockHeight uint32,
	execMode uint32,
) *clobtypes.StreamUpdate {
	return &clobtypes.StreamUpdate{
		BlockHeight: blockHeight,
		ExecMode:    execMode,
		UpdateMessage: &clobtypes.StreamUpdate_PriceUpdate{
			PriceUpdate: &pricestypes.StreamPriceUpdate{
				MarketId: 1,
				Price: pricestypes.MarketPrice{
					Id:       1,
					Exponent: 6,
					Price:    1,
				},
				Snapshot: false,
			},
		},
	}
}

func NewIndexerOrderId(owner string, id uint32) v1types.IndexerOrderId {
	return v1types.IndexerOrderId{
		SubaccountId: v1types.IndexerSubaccountId{
			Owner:  owner,
			Number: id,
		},
		ClientId:   0,
		OrderFlags: 0,
		ClobPairId: 0,
	}
}

func NewOrderId(owner string, id uint32) clobtypes.OrderId {
	return clobtypes.OrderId{
		SubaccountId: satypes.SubaccountId{
			Owner:  owner,
			Number: id,
		},
		ClientId:   0,
		OrderFlags: 0,
		ClobPairId: 0,
	}
}

func NewIndexerOrder(id v1types.IndexerOrderId) v1types.IndexerOrder {
	return v1types.IndexerOrder{
		OrderId:  id,
		Side:     v1types.IndexerOrder_SIDE_BUY,
		Quantums: uint64(1_000_000),
		Subticks: 1,
		GoodTilOneof: &v1types.IndexerOrder_GoodTilBlock{
			GoodTilBlock: 1_000_000_000,
		},
		TimeInForce:                     1_000_000_000,
		ReduceOnly:                      false,
		ClientMetadata:                  0,
		ConditionType:                   v1types.IndexerOrder_CONDITION_TYPE_UNSPECIFIED,
		ConditionalOrderTriggerSubticks: 0,
	}
}

func NewOrder(id clobtypes.OrderId) *clobtypes.Order {
	return &clobtypes.Order{
		OrderId:  id,
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: uint64(1_000_000),
		Subticks: 1,
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: 1_000_000_000,
		},
		TimeInForce:                     1_000_000_000,
		ReduceOnly:                      false,
		ClientMetadata:                  0,
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_UNSPECIFIED,
		ConditionalOrderTriggerSubticks: 0,
	}
}

func NewLogger() *mocks.Logger {
	logger := mocks.Logger{}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	return &logger
}

type TestCase struct {
	updates         []clobtypes.StreamUpdate
	subaccountIds   []satypes.SubaccountId
	filteredUpdates []clobtypes.StreamUpdate
}

func TestFilterStreamUpdates(t *testing.T) {
	logger := NewLogger()

	subaccountId := satypes.SubaccountId{Owner: "me", Number: 1337}
	orderId := NewIndexerOrderId(subaccountId.Owner, subaccountId.Number)
	order := NewIndexerOrder(orderId)

	otherSubaccountId := satypes.SubaccountId{Owner: "we", Number: 2600}
	otherOrderId := NewIndexerOrderId(otherSubaccountId.Owner, otherSubaccountId.Number)
	otherOrder := NewIndexerOrder(otherOrderId)

	neverInScopeSubaccountId := satypes.SubaccountId{Owner: "them", Number: 404}

	newOrderId := order.OrderId
	newOrderId.ClientId += 1
	newOrder := NewIndexerOrder(newOrderId)

	otherNewOrderId := otherOrder.OrderId
	otherNewOrderId.ClientId += 1
	otherNewOrder := NewIndexerOrder(otherNewOrderId)

	orderPlaceTime := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	openOrder := OpenOrder(&order, &orderPlaceTime)
	orderCancelTime := orderPlaceTime.Add(time.Second)
	cancelOrder := CancelOrder(&orderId, &orderCancelTime)
	orderReplaceTime := orderPlaceTime.Add(time.Minute)
	replaceOrder := ReplaceOrder(&orderId, &newOrder, &orderReplaceTime)
	updateOrder := UpdateOrder(&orderId, uint64(1988))

	otherOpenOrder := OpenOrder(&otherOrder, &orderPlaceTime)
	otherCancelOrder := CancelOrder(&otherOrderId, &orderCancelTime)
	otherReplaceOrder := ReplaceOrder(&otherOrderId, &otherNewOrder, &orderReplaceTime)
	otherUpdateOrder := UpdateOrder(&otherOrderId, uint64(1999))

	baseStreamUpdate := toStreamUpdate(false, openOrder, cancelOrder, replaceOrder, updateOrder)
	snapshotBaseStreamUpdate := toStreamUpdate(true, openOrder, cancelOrder, replaceOrder, updateOrder)
	otherStreamUpdate := toStreamUpdate(false, otherOpenOrder, otherCancelOrder, otherReplaceOrder, otherUpdateOrder)
	bothStreamUpdate := toStreamUpdate(
		false,
		openOrder,
		cancelOrder,
		replaceOrder,
		updateOrder,
		otherOpenOrder,
		otherCancelOrder,
		otherReplaceOrder,
		otherUpdateOrder,
	)

	orderBookFillUpdate := NewStreamOrderbookFill(0, 0)
	clobOrder := NewOrder(NewOrderId("foo", 23))
	takerOrderUpdate := NewStreamTakerOrder(0, 0, clobOrder, 0, 0)
	subaccountUpdate := NewSubaccountUpdate(
		0,
		0,
		(*satypes.SubaccountId)(&orderId.SubaccountId),
	)
	priceUpdate := NewPriceUpdate(0, 0)

	tests := map[string]TestCase{
		"snapshotNotInScope": {
			updates:         []clobtypes.StreamUpdate{snapshotBaseStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{neverInScopeSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{snapshotBaseStreamUpdate},
		},
		"snapshotNoScope": {
			updates:         []clobtypes.StreamUpdate{snapshotBaseStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{snapshotBaseStreamUpdate},
		},
		"baseInScope": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate},
		},
		"baseNotInScope": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{neverInScopeSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{},
		},
		"otherInScope": {
			updates:         []clobtypes.StreamUpdate{otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{otherSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{otherStreamUpdate},
		},
		"otherNotInScope": {
			updates:         []clobtypes.StreamUpdate{otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{},
		},
		"bothInScope": {
			updates:         []clobtypes.StreamUpdate{bothStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId, otherSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{bothStreamUpdate},
		},
		"bothOtherInScope": {
			updates:         []clobtypes.StreamUpdate{bothStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{otherSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{bothStreamUpdate},
		},
		"bothSequentiallyOtherInScope": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{otherSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{otherStreamUpdate},
		},
		"bothBaseInScope": {
			updates:         []clobtypes.StreamUpdate{bothStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{bothStreamUpdate},
		},
		"bothSequentiallyBaseInScope": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate},
		},
		"bothNoneInScopeWrongId": {
			updates:         []clobtypes.StreamUpdate{bothStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{neverInScopeSubaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{},
		},
		"bothNoneInScopeNoId": {
			updates:         []clobtypes.StreamUpdate{bothStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{},
		},
		"noUpdates": {
			updates:         []clobtypes.StreamUpdate{},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{},
		},
		"noUpdatesNoId": {
			updates:         []clobtypes.StreamUpdate{},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{},
		},
		"orderBookFillUpdates": {
			updates:         []clobtypes.StreamUpdate{*orderBookFillUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*orderBookFillUpdate},
		},
		"orderBookFillUpdatesDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *orderBookFillUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*orderBookFillUpdate},
		},
		"orderBookFillUpdatesFilterUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *orderBookFillUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *orderBookFillUpdate},
		},
		"orderBookFillUpdatesFilterAndDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *orderBookFillUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *orderBookFillUpdate},
		},
		"takerOrderUpdates": {
			updates:         []clobtypes.StreamUpdate{*takerOrderUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*takerOrderUpdate},
		},
		"takerOrderUpdatesDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *takerOrderUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*takerOrderUpdate},
		},
		"takerOrderUpdatesFilterUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *takerOrderUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *takerOrderUpdate},
		},
		"takerOrderUpdatesFilterAndDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *takerOrderUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *takerOrderUpdate},
		},
		"subaccountUpdates": {
			updates:         []clobtypes.StreamUpdate{*subaccountUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*subaccountUpdate},
		},
		"subaccountUpdatesDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *subaccountUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*subaccountUpdate},
		},
		"subaccountUpdatesFilterUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *subaccountUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *subaccountUpdate},
		},
		"subaccountUpdatesFilterAndDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *subaccountUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *subaccountUpdate},
		},
		"priceUpdates": {
			updates:         []clobtypes.StreamUpdate{*priceUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*priceUpdate},
		},
		"priceUpdatesDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *priceUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{},
			filteredUpdates: []clobtypes.StreamUpdate{*priceUpdate},
		},
		"priceUpdatesFilterUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *priceUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *priceUpdate},
		},
		"priceUpdatesFilterAndDropUpdate": {
			updates:         []clobtypes.StreamUpdate{baseStreamUpdate, *priceUpdate, otherStreamUpdate},
			subaccountIds:   []satypes.SubaccountId{subaccountId},
			filteredUpdates: []clobtypes.StreamUpdate{baseStreamUpdate, *priceUpdate},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			filteredUpdates := streaming.FilterStreamUpdateBySubaccount(testCase.updates, testCase.subaccountIds, logger)
			require.Equal(t, testCase.filteredUpdates, filteredUpdates)
		})
	}
}
