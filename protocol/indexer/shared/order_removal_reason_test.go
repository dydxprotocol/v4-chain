package shared_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetOrderRemovalReason_Success(t *testing.T) {
	tests := map[string]struct {
		// Input
		orderStatus clobtypes.OrderStatus
		orderError  error

		// Expectations
		expectedReason sharedtypes.OrderRemovalReason
		expectedErr    error
	}{
		"Gets order removal reason for order status Undercollateralized": {
			orderStatus:    clobtypes.Undercollateralized,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED,
			expectedErr:    nil,
		},
		"Gets order removal reason for order status InternalError": {
			orderStatus:    clobtypes.InternalError,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
			expectedErr:    nil,
		},
		"Gets order removal reason for order status ImmediateOrCancelWouldRestOnBook": {
			orderStatus:    clobtypes.ImmediateOrCancelWouldRestOnBook,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK,
			expectedErr:    nil,
		},
		"Gets order removal reason for order status ViolatesIsolatedSubaccountConstraints": {
			orderStatus:    clobtypes.ViolatesIsolatedSubaccountConstraints,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS,
			expectedErr:    nil,
		},
		"Gets order removal reason for order error ErrFokOrderCouldNotBeFullyFilled": {
			orderError:     clobtypes.ErrFokOrderCouldNotBeFullyFilled,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED,
			expectedErr:    nil,
		},
		"Gets order removal reason for order error ErrPostOnlyWouldCrossMakerOrder": {
			orderError:     clobtypes.ErrPostOnlyWouldCrossMakerOrder,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
			expectedErr:    nil,
		},
		"Gets order removal reason for order error ErrReduceOnlyWouldIncreasePositionSize": {
			orderError:     clobtypes.ErrReduceOnlyWouldIncreasePositionSize,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE,
			expectedErr:    nil,
		},
		"Gets order removal reason for order error ErrWouldViolateIsolatedSubaccountConstraints": {
			orderError:     clobtypes.ErrWouldViolateIsolatedSubaccountConstraints,
			expectedReason: sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS,
			expectedErr:    nil,
		},
		"Returns error for order status Success": {
			orderStatus:    clobtypes.Success,
			orderError:     clobtypes.ErrNotImplemented,
			expectedReason: 0,
			expectedErr: fmt.Errorf(
				"unrecognized order status %d and error \"%w\"",
				clobtypes.Success,
				clobtypes.ErrNotImplemented,
			),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			reason, err := shared.GetOrderRemovalReason(tc.orderStatus, tc.orderError)
			require.Equal(t, tc.expectedReason, reason)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}
