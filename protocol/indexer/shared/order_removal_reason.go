package shared

import (
	"errors"
	"fmt"

	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// ConvertOrderRemovalReasonToIndexerOrderRemovalReason converts a `OrderRemoval_RemovalReason` to indexer's
// `OrderRemoveV1_OrderRemovalReason`. This is helpful in the memclob logic where we handle
// a bulk of order removals and generate offchain updates for each order removal.
func ConvertOrderRemovalReasonToIndexerOrderRemovalReason(
	removalReason clobtypes.OrderRemoval_RemovalReason,
) sharedtypes.OrderRemovalReason {
	var reason sharedtypes.OrderRemovalReason
	switch removalReason {
	case clobtypes.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED
	case clobtypes.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE
	case clobtypes.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER
	case clobtypes.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_SELF_TRADE_ERROR
	case clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED
	case clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK
	case clobtypes.OrderRemoval_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS:
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS
	case clobtypes.OrderRemoval_REMOVAL_REASON_PERMISSIONED_KEY_EXPIRED:
		// This is a special case where the order is no longer valid because the permissioned key used to placed
		// the order has expired.
		reason = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_EXPIRED
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
) (sharedtypes.OrderRemovalReason, error) {
	switch {
	case errors.Is(orderError, clobtypes.ErrReduceOnlyWouldIncreasePositionSize):
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE, nil
	case errors.Is(orderError, clobtypes.ErrPostOnlyWouldCrossMakerOrder):
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER, nil
	case errors.Is(orderError, clobtypes.ErrFokOrderCouldNotBeFullyFilled):
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED, nil
	case errors.Is(orderError, clobtypes.ErrOrderWouldExceedMaxOpenOrdersEquityTierLimit):
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_EQUITY_TIER, nil
	case errors.Is(orderError, clobtypes.ErrWouldViolateIsolatedSubaccountConstraints):
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS, nil
	}

	switch orderStatus {
	case clobtypes.Undercollateralized:
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED, nil
	case clobtypes.InternalError:
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR, nil
	case clobtypes.ImmediateOrCancelWouldRestOnBook:
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK, nil
	case clobtypes.ReduceOnlyResized:
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE, nil
	case clobtypes.ViolatesIsolatedSubaccountConstraints:
		return sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS, nil
	default:
		return 0, fmt.Errorf("unrecognized order status %d and error \"%w\"", orderStatus, orderError)
	}
}
