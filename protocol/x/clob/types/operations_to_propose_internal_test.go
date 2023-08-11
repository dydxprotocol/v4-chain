package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustGetNonceFromOperation_PanicsOnNoNonce(t *testing.T) {
	otp := NewOperationsToPropose()
	orderPlacementOperation := NewOrderPlacementOperation(Order{})
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			orderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.insertOperationIntoOperationsToPropose(orderPlacementOperation)
		},
	)
}

func TestInsertIntoOperationsQueue_PanicsOnDuplicateNonce(t *testing.T) {
	otp := NewOperationsToPropose()
	cancel := MsgCancelOrder{}
	cancel2 := MsgCancelOrder{OrderId: OrderId{ClientId: 2}}
	cancelOperation := NewOrderCancellationOperation(&cancel)
	cancel2Operation := NewOrderCancellationOperation(&cancel2)

	otp.AddOrderCancellationToOperationsQueue(cancel)
	// Set the next available nonce to 0 so that it's incorrectly re-used.
	otp.NextAvailableNonce = 0
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"insertOperationIntoOperationsToPropose: an operation with nonce %d already exists "+
				"in the operations to propose. New operation: (%+v). Existing operation: (%+v).",
			0,
			cancel2Operation.GetOperationTextString(),
			cancelOperation.GetOperationTextString(),
		),
		func() {
			otp.AddOrderCancellationToOperationsQueue(cancel2)
		},
	)
}
