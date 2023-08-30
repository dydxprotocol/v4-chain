package shared

import (
	"errors"
	"fmt"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// ConvertOrderRemovalReasonToIndexerOrderRemovalReason converts a `OrderRemoval_RemovalReason` to indexer's
// `OrderRemoveV1_OrderRemovalReason`. This is helpful in the memclob logic where we handle
// a bulk of order removals and generate offchain updates for each order removal.
func ConvertOrderRemovalReasonToIndexerOrderRemovalReason(
	removalReason clobtypes.OrderRemoval_RemovalReason,
) OrderRemovalReason {
	var reason OrderRemovalReason
	switch removalReason {
	case clobtypes.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED:
		reason = OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED
	case clobtypes.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY:
		reason = OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE
	case clobtypes.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER:
		reason = OrderRemovalReason_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER
	case clobtypes.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE:
		reason = OrderRemovalReason_ORDER_REMOVAL_REASON_SELF_TRADE_ERROR
	case clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED:
		reason = OrderRemovalReason_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED
	case clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK:
		reason = OrderRemovalReason_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK
	default:
		panic("ConvertOrderRemovalReasonToIndexerOrderRemovalReason: unspecified removal reason not allowed")
	}
	return reason
}

// GetOrderRemovalReason gets the matching `OrderRemove_OrderRemovalReason` given the status of an
// order.
func GetOrderRemovalReason(
	orderStatus clobtypes.OrderStatus,
	orderError error,
) (OrderRemovalReason, error) {
	switch {
	case errors.Is(orderError, clobtypes.ErrPostOnlyWouldCrossMakerOrder):
		return OrderRemovalReason_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER, nil
	case errors.Is(orderError, clobtypes.ErrFokOrderCouldNotBeFullyFilled):
		return OrderRemovalReason_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED, nil
	case errors.Is(orderError, clobtypes.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit):
		return OrderRemovalReason_ORDER_REMOVAL_REASON_EQUITY_TIER, nil
	}

	switch orderStatus {
	case clobtypes.Undercollateralized:
		return OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED, nil
	case clobtypes.InternalError:
		return OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR, nil
	case clobtypes.ImmediateOrCancelWouldRestOnBook:
		return OrderRemovalReason_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK, nil
	case clobtypes.ReduceOnlyResized:
		return OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE, nil
	default:
		return 0, fmt.Errorf("unrecognized order status %d and error \"%w\"", orderStatus, orderError)
	}
}
