package off_chain_updates

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4/indexer/common"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

const (
	hashErrMsg   = "Cannot hash order id."
	createErrMsg = "Cannot create message."
)

// MustCreateOrderPlaceMessage invokes CreateOrderPlaceMessage and panics if creation was unsuccessful.
func MustCreateOrderPlaceMessage(
	logger log.Logger,
	order clobtypes.Order,
) msgsender.Message {
	msg, ok := CreateOrderPlaceMessage(logger, order)
	if !ok {
		panic(fmt.Errorf("Unable to create place order message for order %+v", order))
	}
	return msg
}

// CreateOrderPlaceMessage creates an off-chain update message for an order.
func CreateOrderPlaceMessage(
	logger log.Logger,
	order clobtypes.Order,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for placing order."
	errDetails := fmt.Sprintf("Order: %+v", order)

	orderIdHash, err := GetOrderIdHash(order.OrderId)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, hashErrMsg, err, errDetails))
		return msgsender.Message{}, false
	}

	update, err := newOrderPlaceMessage(&order)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, createErrMsg, err, errDetails))
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderUpdateMessage invokes CreateOrderUpdateMessage and panics if creation was unsuccessful.
func MustCreateOrderUpdateMessage(
	logger log.Logger,
	orderId clobtypes.OrderId,
	totalFilled satypes.BaseQuantums,
) msgsender.Message {
	msg, ok := CreateOrderUpdateMessage(logger, orderId, totalFilled)
	if !ok {
		panic(fmt.Errorf("Unable to create place order message for order id %+v", orderId))
	}
	return msg
}

// CreateOrderUpdateMessage creates an off-chain update message for an order being updated.
func CreateOrderUpdateMessage(
	logger log.Logger,
	orderId clobtypes.OrderId,
	totalFilled satypes.BaseQuantums,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for updating order."
	errDetails := fmt.Sprintf("OrderId: %+v, TotalFilled %+v", orderId, totalFilled)

	orderIdHash, err := GetOrderIdHash(orderId)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, hashErrMsg, err, errDetails))
		return msgsender.Message{}, false
	}

	update, err := newOrderUpdateMessage(&orderId, totalFilled)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, createErrMsg, err, errDetails))
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderRemoveMessageWithReason invokes CreateOrderRemoveMessageWithReason and panics if creation was
// unsuccessful.
func MustCreateOrderRemoveMessageWithReason(
	logger log.Logger,
	orderId clobtypes.OrderId,
	reason OrderRemove_OrderRemovalReason,
	removalStatus OrderRemove_OrderRemovalStatus,
) msgsender.Message {
	msg, ok := CreateOrderRemoveMessageWithReason(logger, orderId, reason, removalStatus)
	if !ok {
		panic(fmt.Errorf("Unable to create remove order message with reason for order id %+v", orderId))
	}
	return msg
}

// CreateOrderRemoveMessageWithReason creates an off-chain update message for an order being removed
// with a specific reason for the removal and the resulting removal status of the removed order.
func CreateOrderRemoveMessageWithReason(
	logger log.Logger,
	orderId clobtypes.OrderId,
	reason OrderRemove_OrderRemovalReason,
	removalStatus OrderRemove_OrderRemovalStatus,
) (message msgsender.Message, success bool) {
	errMessage := "Error creating off-chain update message for removing order."
	errDetails := fmt.Sprintf(
		"OrderId: %+v, Reason %d, Removal status %d",
		orderId,
		reason,
		removalStatus,
	)

	orderIdHash, err := GetOrderIdHash(orderId)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, hashErrMsg, err, errDetails))
		return msgsender.Message{}, false
	}

	update, err := newOrderRemoveMessage(orderId, reason, removalStatus)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, createErrMsg, err, errDetails))
		return msgsender.Message{}, false
	}

	return msgsender.Message{Key: orderIdHash, Value: update}, true
}

// MustCreateOrderRemoveMessage invokes CreateOrderRemoveMessage and panics if creation was unsuccessful.
func MustCreateOrderRemoveMessage(logger log.Logger,
	orderId clobtypes.OrderId,
	orderStatus clobtypes.OrderStatus,
	orderError error,
	removalStatus OrderRemove_OrderRemovalStatus,
) msgsender.Message {
	msg, ok := CreateOrderRemoveMessage(logger, orderId, orderStatus, orderError, removalStatus)
	if !ok {
		panic(fmt.Errorf("Unable to create remove order message for order id %+v", orderId))
	}
	return msg
}

// CreateOrderRemoveMessage creates an off-chain update message for an order being removed, with the
// order's status and the resulting removal status of the removed order.
func CreateOrderRemoveMessage(
	logger log.Logger,
	orderId clobtypes.OrderId,
	orderStatus clobtypes.OrderStatus,
	orderError error,
	removalStatus OrderRemove_OrderRemovalStatus,
) (message msgsender.Message, success bool) {
	errDetails := fmt.Sprintf(
		"OrderId: %+v, Removal status %d",
		orderId,
		removalStatus,
	)

	reason, err := GetOrderRemovalReason(orderStatus, orderError)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error creating off-chain update message for removing order. Invalid order removal "+
					"reason. Error: %+v %s\n",
				err,
				errDetails,
			),
		)
		return msgsender.Message{}, false
	}

	return CreateOrderRemoveMessageWithReason(logger, orderId, reason, removalStatus)
}

// CreateOrderRemoveMessageWithDefaultReason creates an off-chain update message for an order being
// removed with the resulting removal status of the removed order. It attempts to look up the removal
// reason using the given orderStatus & orderError. If the reason cannot be found, it logs an error
// and falls back to the defaultRemovalReason. If defaultRemovalReason is ...UNSPECIFIED, it panics.
func CreateOrderRemoveMessageWithDefaultReason(
	logger log.Logger,
	orderId clobtypes.OrderId,
	orderStatus clobtypes.OrderStatus,
	orderError error,
	removalStatus OrderRemove_OrderRemovalStatus,
	defaultRemovalReason OrderRemove_OrderRemovalReason,
) (message msgsender.Message, success bool) {
	if defaultRemovalReason == OrderRemove_ORDER_REMOVAL_REASON_UNSPECIFIED {
		panic(
			fmt.Errorf(
				"Invalid parameter: " +
					"defaultRemovalReason cannot be OrderRemove_ORDER_REMOVAL_REASON_UNSPECIFIED",
			),
		)
	}
	errDetails := fmt.Sprintf(
		"OrderId: %+v, Removal status %d",
		orderId,
		removalStatus,
	)

	reason, err := GetOrderRemovalReason(orderStatus, orderError)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error creating off-chain update message for removing order. Invalid order removal "+
					"reason. Error: %+v %s\n",
				err,
				errDetails,
			),
		)
		reason = defaultRemovalReason
	}

	return CreateOrderRemoveMessageWithReason(logger, orderId, reason, removalStatus)
}

// newOrderPlaceMessage returns an `OffChainUpdate` struct populated with an `OrderPlace` struct
// as the `UpdateMessage` parameter, encoded as a byte slice.
func newOrderPlaceMessage(
	order *clobtypes.Order,
) ([]byte, error) {
	update := OffChainUpdate{
		UpdateMessage: &OffChainUpdate_OrderPlace{
			&OrderPlace{
				Order: order,
				// Protocol will always send best effort opened messages to indexer.
				PlacementStatus: OrderPlace_ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
			},
		},
	}
	return marshalOffchainUpdate(update, &common.MarshalerImpl{})
}

// newOrderPlaceMessage returns an `OffChainUpdate` struct populated with an `OrderRemove`
// struct as the `UpdateMessage` parameter, encoded as a byte slice.
// The `OrderRemove` struct is instantiated with the given orderId, reason and status parameters.
func newOrderRemoveMessage(
	orderId clobtypes.OrderId,
	reason OrderRemove_OrderRemovalReason,
	status OrderRemove_OrderRemovalStatus,
) ([]byte, error) {
	update := OffChainUpdate{
		UpdateMessage: &OffChainUpdate_OrderRemove{
			&OrderRemove{
				RemovedOrderId: &orderId,
				Reason:         reason,
				RemovalStatus:  status,
			},
		},
	}
	return marshalOffchainUpdate(update, &common.MarshalerImpl{})
}

// NewOrderUpdateMessage returns an `OffChainUpdate` struct populated with an `OrderUpdate`
// struct as the `UpdateMessage` parameter, encoded as a byte slice.
// The `OrderUpdate` struct is instantiated with the given orderId and totalFilled parameters.
func newOrderUpdateMessage(
	orderId *clobtypes.OrderId,
	totalFilled satypes.BaseQuantums,
) ([]byte, error) {
	update := OffChainUpdate{
		UpdateMessage: &OffChainUpdate_OrderUpdate{
			&OrderUpdate{
				OrderId:             orderId,
				TotalFilledQuantums: totalFilled.ToUint64(),
			},
		},
	}
	return marshalOffchainUpdate(update, &common.MarshalerImpl{})
}

func marshalOffchainUpdate(
	offChainUpdate OffChainUpdate,
	marshaler common.Marshaler,
) ([]byte, error) {
	updateBytes, err := marshaler.Marshal(&offChainUpdate)
	return updateBytes, err
}

// GetOrderIdHash gets the SHA256 hash of an `OrderId`.
func GetOrderIdHash(orderId clobtypes.OrderId) ([]byte, error) {
	orderIdBytes, err := (&orderId).Marshal()
	if err != nil {
		return []byte{}, err
	}
	byteArray := sha256.Sum256(orderIdBytes)
	return byteArray[:], nil
}

// GetOrderRemovalReason gets the matching `OrderRemove_OrderRemovalReason` given the status of an
// order.
func GetOrderRemovalReason(
	orderStatus clobtypes.OrderStatus,
	orderError error,
) (OrderRemove_OrderRemovalReason, error) {
	switch {
	case errors.Is(orderError, clobtypes.ErrPostOnlyWouldCrossMakerOrder):
		return OrderRemove_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER, nil
	case errors.Is(orderError, clobtypes.ErrFokOrderCouldNotBeFullyFilled):
		return OrderRemove_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED, nil
	}

	switch orderStatus {
	case clobtypes.Undercollateralized:
		return OrderRemove_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED, nil
	case clobtypes.InternalError:
		return OrderRemove_ORDER_REMOVAL_REASON_INTERNAL_ERROR, nil
	case clobtypes.ImmediateOrCancelWouldRestOnBook:
		return OrderRemove_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK, nil
	case clobtypes.ReduceOnlyResized:
		return OrderRemove_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE, nil
	default:
		return 0, fmt.Errorf("unrecognized order status %d and error \"%w\"", orderStatus, orderError)
	}
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
		return false
	default:
		return true
	}
}
