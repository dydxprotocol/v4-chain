package shared_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetOrderRemovalReason_Success(t *testing.T) {
	tests := map[string]struct {
		// Input
		orderStatus clobtypes.OrderStatus
		orderError  error

		// Expectations
		expectedReason shared.OrderRemovalReason
		expectedErr    error
	}{
		"Gets order removal reason for order status Undercollateralized": {
			orderStatus:    clobtypes.Undercollateralized,
			expectedReason: shared.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED,
			expectedErr:    nil,
		},
		"Gets order removal reason for order status InternalError": {
			orderStatus:    clobtypes.InternalError,
			expectedReason: shared.OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
			expectedErr:    nil,
		},
		"Gets order removal reason for order status ImmediateOrCancelWouldRestOnBook": {
			orderStatus:    clobtypes.ImmediateOrCancelWouldRestOnBook,
			expectedReason: shared.OrderRemovalReason_ORDER_REMOVAL_REASON_IMMEDIATE_OR_CANCEL_WOULD_REST_ON_BOOK,
			expectedErr:    nil,
		},
		"Gets order removal reason for order error ErrFokOrderCouldNotBeFullyFilled": {
			orderError:     clobtypes.ErrFokOrderCouldNotBeFullyFilled,
			expectedReason: shared.OrderRemovalReason_ORDER_REMOVAL_REASON_FOK_ORDER_COULD_NOT_BE_FULLY_FULLED,
			expectedErr:    nil,
		},
		"Gets order removal reason for order error ErrPostOnlyWouldCrossMakerOrder": {
			orderError:     clobtypes.ErrPostOnlyWouldCrossMakerOrder,
			expectedReason: shared.OrderRemovalReason_ORDER_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
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
