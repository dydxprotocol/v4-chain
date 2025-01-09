package streaming_test

import (
	"testing"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	pv1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	streaming "github.com/dydxprotocol/v4-chain/protocol/streaming"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

const (
	maxSubscriptionChannelSize = 2 ^ 10
	owner                      = "foo"
	noMessagesMaxSleep         = 10 * time.Millisecond
)

func OpenOrder(
	order *pv1types.IndexerOrder,
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
	removedOrderId *pv1types.IndexerOrderId,
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
	oldOrderId *pv1types.IndexerOrderId,
	newOrder *pv1types.IndexerOrder,
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

func UpdateOrder(orderId *pv1types.IndexerOrderId, totalFilledQuantums uint64) ocutypes.OffChainUpdateV1 {
	return ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderUpdate{
			OrderUpdate: &ocutypes.OrderUpdateV1{
				OrderId:             orderId,
				TotalFilledQuantums: totalFilledQuantums,
			},
		},
	}
}

func toStreamUpdate(offChainUpdates []ocutypes.OffChainUpdateV1, blockHeight uint32) clobtypes.StreamUpdate {
	return clobtypes.StreamUpdate{
		BlockHeight: blockHeight,
		ExecMode:    uint32(sdktypes.ExecModeFinalize),
		UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
			OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
				Updates:  offChainUpdates,
				Snapshot: true,
			},
		},
	}
}

type MockMessageSender struct{}

func (mms *MockMessageSender) Send(*clobtypes.StreamOrderbookUpdatesResponse) error {
	return nil
}

func NewOrderbookSubscription(
	ids []uint32,
	updatesChannel chan []clobtypes.StreamUpdate,
) *streaming.OrderbookSubscription {
	sIds := []satypes.SubaccountId{}
	for _, id := range ids {
		sIds = append(sIds, satypes.SubaccountId{Owner: owner, Number: id})
	}
	return streaming.NewOrderbookSubscription(
		0,
		[]uint32{},
		sIds,
		[]uint32{},
		&MockMessageSender{},
		updatesChannel,
	)
}

type TestCase struct {
	updates         *[]clobtypes.StreamUpdate
	subaccountIds   []uint32
	filteredUpdates *[]clobtypes.StreamUpdate
}

func TestFilterStreamUpdates(t *testing.T) {
	logger := &mocks.Logger{}

	subaccountIdNumber := uint32(1337)
	subaccountId := pv1types.IndexerSubaccountId{
		Owner:  "foo",
		Number: subaccountIdNumber,
	}
	orderId := pv1types.IndexerOrderId{
		SubaccountId: subaccountId,
		ClientId:     0,
		OrderFlags:   0,
		ClobPairId:   0,
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

	newOrderId := order.OrderId
	newOrderId.ClientId += 1

	newOrder := order
	newOrder.OrderId = newOrderId
	newOrder.Quantums += 10 ^ 6

	totalFilledQuantums := uint64(1988)

	tests := make(map[string]TestCase)

	orderPlaceTime := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	openOrder := OpenOrder(&order, &orderPlaceTime)

	orderCancelTime := orderPlaceTime.Add(time.Second)
	cancelOrder := CancelOrder(&orderId, &orderCancelTime)

	orderReplaceTime := orderPlaceTime.Add(time.Minute)
	replaceOrder := ReplaceOrder(&orderId, &newOrder, &orderReplaceTime)

	updateOrder := UpdateOrder(&orderId, totalFilledQuantums)

	baseOffChainUpdates := []ocutypes.OffChainUpdateV1{openOrder, cancelOrder, replaceOrder, updateOrder}
	baseStreamUpdates := []clobtypes.StreamUpdate{toStreamUpdate(baseOffChainUpdates, 0)}
	tests["baseInScope"] = TestCase{
		updates:         &baseStreamUpdates,
		subaccountIds:   []uint32{orderId.SubaccountId.Number},
		filteredUpdates: &baseStreamUpdates,
	}
	tests["baseNotInScope"] = TestCase{
		updates:         &baseStreamUpdates,
		subaccountIds:   []uint32{0},
		filteredUpdates: nil,
	}

	otherOrderId := orderId
	otherSubaccountIdNumber := subaccountIdNumber + uint32(1)
	otherOrderId.SubaccountId = pv1types.IndexerSubaccountId{
		Owner:  "bar",
		Number: otherSubaccountIdNumber,
	}
	otherOrder := order
	otherOrder.OrderId = otherOrderId

	otherNewOrderId := otherOrder.OrderId
	otherNewOrderId.ClientId += 1

	otherNewOrder := otherOrder
	otherNewOrder.OrderId = otherNewOrderId
	otherNewOrder.Quantums += 10 ^ 6

	otherOpenOrder := OpenOrder(&otherOrder, &orderPlaceTime)
	otherCancelOrder := CancelOrder(&otherOrderId, &orderCancelTime)
	otherReplaceOrder := ReplaceOrder(&otherOrderId, &otherNewOrder, &orderReplaceTime)
	otherUpdateOrder := UpdateOrder(&otherOrderId, totalFilledQuantums)

	otherOffChainUpdates := []ocutypes.OffChainUpdateV1{
		otherOpenOrder, otherCancelOrder, otherReplaceOrder, otherUpdateOrder,
	}
	otherStreamUpdates := []clobtypes.StreamUpdate{toStreamUpdate(otherOffChainUpdates, 0)}
	tests["otherInScope"] = TestCase{
		updates:         &otherStreamUpdates,
		subaccountIds:   []uint32{otherSubaccountIdNumber},
		filteredUpdates: &otherStreamUpdates,
	}
	tests["otherNotInScope"] = TestCase{
		updates:         &otherStreamUpdates,
		subaccountIds:   []uint32{subaccountIdNumber},
		filteredUpdates: nil,
	}

	bothUpdates := []clobtypes.StreamUpdate{
		toStreamUpdate(append(baseOffChainUpdates, otherOffChainUpdates...), 0),
	}
	tests["bothInScope"] = TestCase{
		updates:         &bothUpdates,
		subaccountIds:   []uint32{subaccountIdNumber, otherSubaccountIdNumber},
		filteredUpdates: &bothUpdates,
	}
	tests["bothOtherInScope"] = TestCase{
		updates:         &bothUpdates,
		subaccountIds:   []uint32{otherSubaccountIdNumber},
		filteredUpdates: &otherStreamUpdates,
	}
	tests["bothBaseInScope"] = TestCase{
		updates:         &bothUpdates,
		subaccountIds:   []uint32{subaccountIdNumber},
		filteredUpdates: &baseStreamUpdates,
	}
	tests["bothNoneInScopeWrongId"] = TestCase{
		updates:         &bothUpdates,
		subaccountIds:   []uint32{404},
		filteredUpdates: nil,
	}
	tests["bothNoneInScopeNoId"] = TestCase{
		updates:         &bothUpdates,
		subaccountIds:   []uint32{},
		filteredUpdates: nil,
	}

	tests["noUpdates"] = TestCase{
		updates:         &[]clobtypes.StreamUpdate{},
		subaccountIds:   []uint32{subaccountIdNumber},
		filteredUpdates: nil,
	}
	tests["noUpdatesNoId"] = TestCase{
		updates:         &[]clobtypes.StreamUpdate{},
		subaccountIds:   []uint32{},
		filteredUpdates: nil,
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			func() {
				filteredUpdatesChannel := make(chan []clobtypes.StreamUpdate, maxSubscriptionChannelSize)
				defer close(filteredUpdatesChannel)
				updatesChannel := make(chan []clobtypes.StreamUpdate, maxSubscriptionChannelSize)
				defer close(updatesChannel)

				subscription := NewOrderbookSubscription(testCase.subaccountIds, updatesChannel)
				go subscription.FilterSubaccountStreamUpdates(filteredUpdatesChannel, logger)
				updatesChannel <- *testCase.updates

				if testCase.filteredUpdates != nil {
					require.Equal(t, <-filteredUpdatesChannel, *testCase.filteredUpdates)
				} else {
					time.Sleep(noMessagesMaxSleep)
					require.Equal(t, len(filteredUpdatesChannel), 0)
				}
			}()
		})
	}
}
