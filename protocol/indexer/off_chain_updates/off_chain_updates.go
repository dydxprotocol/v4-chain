package off_chain_updates

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// MustCreateOrderPlaceMessage invokes CreateOrderPlaceMessage and panics if creation was unsuccessful.
func MustCreateOrderPlaceMessage(
	ctx sdk.Context,
	order clobtypes.Order,
) msgsender.Message {
	msg, ok := CreateOrderPlaceMessage(ctx, order)
	if !ok {
		panic(fmt.Errorf("Unable to create place order message for order %+v", order))
	}
	return msg
}

// CreateOrderPlaceMessage creates an off-chain update message for an order.
func CreateOrderPlaceMessage(
	ctx sdk.Context,
	order clobtypes.Order,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for placing order."

	orderIdHash, err := GetOrderIdHash(order.OrderId)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.Order, order,
		)
		return msgsender.Message{}, false
	}

	update, err := NewOrderPlaceMessage(order)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.Order, order,
		)
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderReplaceMessage invokes CreateOrderReplaceMessage and panics if creation was unsuccessful.
func MustCreateOrderReplaceMessage(
	ctx sdk.Context,
	order clobtypes.Order,
) msgsender.Message {
	msg, ok := CreateOrderReplaceMessage(ctx, order)
	if !ok {
		panic(fmt.Errorf("Unable to create place order message for order %+v", order))
	}
	return msg
}

// CreateOrderReplaceMessage creates an off-chain update message for an order.
func CreateOrderReplaceMessage(
	ctx sdk.Context,
	order clobtypes.Order,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for replacing order."

	orderIdHash, err := GetOrderIdHash(order.OrderId)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.Order, order,
		)
		return msgsender.Message{}, false
	}

	update, err := NewOrderReplaceMessage(order)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.Order, order,
		)
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderUpdateMessage invokes CreateOrderUpdateMessage and panics if creation was unsuccessful.
func MustCreateOrderUpdateMessage(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	totalFilled satypes.BaseQuantums,
) msgsender.Message {
	msg, ok := CreateOrderUpdateMessage(ctx, orderId, totalFilled)
	if !ok {
		panic(fmt.Errorf("Unable to create place order message for order id %+v", orderId))
	}
	return msg
}

// CreateOrderUpdateMessage creates an off-chain update message for an order being updated.
func CreateOrderUpdateMessage(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	totalFilled satypes.BaseQuantums,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for updating order."

	orderIdHash, err := GetOrderIdHash(orderId)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.OrderId, orderId,
			log.TotalFilled, totalFilled,
		)
		return msgsender.Message{}, false
	}

	update, err := NewOrderUpdateMessage(orderId, totalFilled)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.OrderId, orderId,
			log.TotalFilled, totalFilled,
		)
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderRemoveMessageWithReason invokes CreateOrderRemoveMessageWithReason and panics if creation was
// unsuccessful.
func MustCreateOrderRemoveMessageWithReason(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	reason sharedtypes.OrderRemovalReason,
	removalStatus ocutypes.OrderRemoveV1_OrderRemovalStatus,
) msgsender.Message {
	msg, ok := CreateOrderRemoveMessageWithReason(ctx, orderId, reason, removalStatus)
	if !ok {
		panic(fmt.Errorf("Unable to create remove order message with reason for order id %+v", orderId))
	}
	return msg
}

// CreateOrderRemoveMessageWithReason creates an off-chain update message for an order being removed
// with a specific reason for the removal and the resulting removal status of the removed order.
func CreateOrderRemoveMessageWithReason(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	reason sharedtypes.OrderRemovalReason,
	removalStatus ocutypes.OrderRemoveV1_OrderRemovalStatus,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for removing order."

	orderIdHash, err := GetOrderIdHash(orderId)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.OrderId, orderId,
			log.Reason, reason,
			log.RemovalStatus, removalStatus,
		)
		return msgsender.Message{}, false
	}

	update, err := NewOrderRemoveMessage(orderId, reason, removalStatus)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			errMessage,
			err,
			log.OrderId, orderId,
			log.Reason, reason,
			log.RemovalStatus, removalStatus,
		)
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderRemoveMessage invokes CreateOrderRemoveMessage and panics if creation was unsuccessful.
func MustCreateOrderRemoveMessage(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	orderStatus clobtypes.OrderStatus,
	orderError error,
	removalStatus ocutypes.OrderRemoveV1_OrderRemovalStatus,
) msgsender.Message {
	msg, ok := CreateOrderRemoveMessage(ctx, orderId, orderStatus, orderError, removalStatus)
	if !ok {
		panic(fmt.Errorf("Unable to create remove order message for order id %+v", orderId))
	}
	return msg
}

// CreateOrderRemoveMessage creates an off-chain update message for an order being removed, with the
// order's status and the resulting removal status of the removed order.
func CreateOrderRemoveMessage(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	orderStatus clobtypes.OrderStatus,
	orderError error,
	removalStatus ocutypes.OrderRemoveV1_OrderRemovalStatus,
) (message msgsender.Message, success bool) {
	reason, err := shared.GetOrderRemovalReason(orderStatus, orderError)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"Error creating off-chain update message for removing order. Invalid order removal reason.",
			err,
			log.OrderId, orderId,
			log.RemovalStatus, removalStatus,
		)
		return msgsender.Message{}, false
	}

	return CreateOrderRemoveMessageWithReason(ctx, orderId, reason, removalStatus)
}

// CreateOrderRemoveMessageWithDefaultReason creates an off-chain update message for an order being
// removed with the resulting removal status of the removed order. It attempts to look up the removal
// reason using the given orderStatus & orderError. If the reason cannot be found, it logs an error
// and falls back to the defaultRemovalReason. If defaultRemovalReason is ...UNSPECIFIED, it panics.
// TODO(CLOB-1051) take in ctx, not logger
func CreateOrderRemoveMessageWithDefaultReason(
	ctx sdk.Context,
	orderId clobtypes.OrderId,
	orderStatus clobtypes.OrderStatus,
	orderError error,
	removalStatus ocutypes.OrderRemoveV1_OrderRemovalStatus,
	defaultRemovalReason sharedtypes.OrderRemovalReason,
) (message msgsender.Message, success bool) {
	if defaultRemovalReason == sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNSPECIFIED {
		panic(
			fmt.Errorf(
				"Invalid parameter: " +
					"defaultRemovalReason cannot be OrderRemove_ORDER_REMOVAL_REASON_UNSPECIFIED",
			),
		)
	}

	reason, err := shared.GetOrderRemovalReason(orderStatus, orderError)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"Error creating off-chain update message for removing order. Invalid order removal reason.",
			err,
			log.OrderId, orderId,
			log.RemovalStatus, removalStatus,
		)
		reason = defaultRemovalReason
	}

	return CreateOrderRemoveMessageWithReason(ctx, orderId, reason, removalStatus)
}

// NewOrderPlaceMessage returns an `OffChainUpdate` struct populated with an `OrderPlace` struct
// as the `UpdateMessage` parameter, encoded as a byte slice.
func NewOrderPlaceMessage(
	order clobtypes.Order,
) ([]byte, error) {
	indexerOrder := v1.OrderToIndexerOrder(order)
	update := ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderPlace{
			OrderPlace: &ocutypes.OrderPlaceV1{
				Order: &indexerOrder,
				// Protocol will always send best effort opened messages to indexer.
				PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
			},
		},
	}
	return proto.Marshal(&update)
}

// NewOrderRemoveMessage returns an `OffChainUpdate` struct populated with an `OrderRemove`
// struct as the `UpdateMessage` parameter, encoded as a byte slice.
// The `OrderRemove` struct is instantiated with the given orderId, reason and status parameters.
func NewOrderRemoveMessage(
	orderId clobtypes.OrderId,
	reason sharedtypes.OrderRemovalReason,
	status ocutypes.OrderRemoveV1_OrderRemovalStatus,
) ([]byte, error) {
	indexerOrderId := v1.OrderIdToIndexerOrderId(orderId)
	update := ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderRemove{
			OrderRemove: &ocutypes.OrderRemoveV1{
				RemovedOrderId: &indexerOrderId,
				Reason:         reason,
				RemovalStatus:  status,
			},
		},
	}
	return proto.Marshal(&update)
}

// NewOrderUpdateMessage returns an `OffChainUpdate` struct populated with an `OrderUpdate`
// struct as the `UpdateMessage` parameter, encoded as a byte slice.
// The `OrderUpdate` struct is instantiated with the given orderId and totalFilled parameters.
func NewOrderUpdateMessage(
	orderId clobtypes.OrderId,
	totalFilled satypes.BaseQuantums,
) ([]byte, error) {
	indexerOrderId := v1.OrderIdToIndexerOrderId(orderId)
	update := ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderUpdate{
			OrderUpdate: &ocutypes.OrderUpdateV1{
				OrderId:             &indexerOrderId,
				TotalFilledQuantums: totalFilled.ToUint64(),
			},
		},
	}
	return proto.Marshal(&update)
}

// NewOrderReplaceMessage returns an `OffChainUpdate` struct populated with an `OrderReplace` struct
// as the `UpdateMessage` parameter, encoded as a byte slice.
func NewOrderReplaceMessage(
	order clobtypes.Order,
) ([]byte, error) {
	indexerOrder := v1.OrderToIndexerOrder(order)
	update := ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderReplace{
			OrderReplace: &ocutypes.OrderReplaceV1{
				Order: &indexerOrder,
				// Protocol will always send best effort opened messages to indexer.
				PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
			},
		},
	}
	return proto.Marshal(&update)
}

// GetOrderIdHash gets the SHA256 hash of the `IndexerOrderId` mapped from an `OrderId`.
func GetOrderIdHash(orderId clobtypes.OrderId) ([]byte, error) {
	indexerOrderId := v1.OrderIdToIndexerOrderId(orderId)
	orderIdBytes, err := (&indexerOrderId).Marshal()
	if err != nil {
		return []byte{}, err
	}
	byteArray := sha256.Sum256(orderIdBytes)
	return byteArray[:], nil
}

// ShouldSendOrderRemovalOnReplay returns a true/false for whether an order removal message should
// be sent given the error encountered while replaying an order.
// TODO(CLOB-518): Re-visit enumerating all the errors where an order removal should be / not be
// sent vs using the existence of an order nonce to determine if an order removal message should be
// sent.
func ShouldSendOrderRemovalOnReplay(
	orderError error,
) bool {
	switch {
	// Order was reprocessed, so should still be on the book.
	case errors.Is(orderError, clobtypes.ErrOrderReprocessed):
		fallthrough
	// Order was not replaced, the order or a newer replacement of it, is still on the book.
	case errors.Is(orderError, clobtypes.ErrInvalidReplacement):
		fallthrough
	// Order was fully filled, no need to remove.
	case errors.Is(orderError, clobtypes.ErrOrderFullyFilled):
		fallthrough
	// Order cancelation was already processed, no need to send a remove.
	case errors.Is(orderError, clobtypes.ErrOrderIsCanceled):
		fallthrough
	// Order already exists on the book, order is still on the book.
	case errors.Is(orderError, clobtypes.ErrStatefulOrderAlreadyExists):
		fallthrough
	// Order should have already been fully-filled or expired as the current height > GoodTilBlock.
	case errors.Is(orderError, clobtypes.ErrHeightExceedsGoodTilBlock):
		fallthrough
	// Order is not resting on the book if already filled, no need to send a remove.
	case errors.Is(orderError, clobtypes.ErrImmediateExecutionOrderAlreadyFilled):
		fallthrough
	// TODO(IND-199): Resolve edge case where the stateful order which has this error was never included
	// in a block and then expired. We do want to send the `OrderRemove` message as the order will not
	// be in state, and thus a stateful order expiration message will not be sent for the order.
	// Order should have already been fully-filled or expired as the  previous block time >= GoodTilBlockTime.
	case errors.Is(orderError, clobtypes.ErrTimeExceedsGoodTilBlockTime):
		return false
	default:
		return true
	}
}
