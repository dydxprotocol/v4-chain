package off_chain_updates_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	ocu "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	sharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	orderIdHash                    = constants.OrderIdHash_Alice_Number0_Id0
	order                          = constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	indexerOrder                   = v1.OrderToIndexerOrder(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15)
	totalFilledAmount              = satypes.BaseQuantums(5)
	orderStatus                    = clobtypes.Undercollateralized
	orderError               error = nil
	reason                         = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED
	status                         = ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED
	defaultRemovalReason           = sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR
	offchainUpdateOrderPlace       = ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderPlace{
			OrderPlace: &ocutypes.OrderPlaceV1{
				Order:           &indexerOrder,
				PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
			},
		},
	}
	offchainUpdateOrderReplace = ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderReplace{
			OrderReplace: &ocutypes.OrderReplaceV1{
				Order:           &indexerOrder,
				PlacementStatus: ocutypes.OrderPlaceV1_ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
			},
		},
	}
	offchainUpdateOrderUpdate = ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderUpdate{
			OrderUpdate: &ocutypes.OrderUpdateV1{
				OrderId:             &indexerOrder.OrderId,
				TotalFilledQuantums: totalFilledAmount.ToUint64(),
			},
		},
	}
	offchainUpdateOrderRemove = ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderRemove{
			OrderRemove: &ocutypes.OrderRemoveV1{
				RemovedOrderId: &indexerOrder.OrderId,
				Reason:         reason,
				RemovalStatus:  status,
			},
		},
	}
	offchainUpdateOrderRemoveWithDefaultRemovalReason = ocutypes.OffChainUpdateV1{
		UpdateMessage: &ocutypes.OffChainUpdateV1_OrderRemove{
			OrderRemove: &ocutypes.OrderRemoveV1{
				RemovedOrderId: &indexerOrder.OrderId,
				Reason:         defaultRemovalReason,
				RemovalStatus:  status,
			},
		},
	}
)

func TestCreateOrderPlaceMessage(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()
	actualMessage, success := ocu.CreateOrderPlaceMessage(
		ctx,
		order,
	)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderPlace)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestCreateOrderReplaceMessage(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()
	actualMessage, success := ocu.CreateOrderReplaceMessage(
		ctx,
		order,
	)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderReplace)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestCreateOrderUpdateMessage(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()

	actualMessage, success := ocu.CreateOrderUpdateMessage(ctx, order.OrderId, totalFilledAmount)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderUpdate)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestCreateOrderRemoveWithReason(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()

	actualMessage, success := ocu.CreateOrderRemoveMessage(
		ctx,
		order.OrderId,
		orderStatus,
		orderError,
		status,
	)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderRemove)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestCreateOrderRemoveMessageWithDefaultReason_HappyPath(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()

	require.NotEqual(
		t,
		offchainUpdateOrderRemove.GetOrderRemove().Reason,
		defaultRemovalReason,
		"defaultRemovalReason must be different than expectedMessage's removal reason for test to "+
			"be valid & useful.")
	actualMessage, success := ocu.CreateOrderRemoveMessageWithDefaultReason(
		ctx,
		order.OrderId,
		orderStatus,
		orderError,
		status,
		defaultRemovalReason,
	)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderRemove)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestCreateOrderRemoveMessageWithDefaultReason_DefaultReasonReturned(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()
	actualMessage, success := ocu.CreateOrderRemoveMessageWithDefaultReason(
		ctx,
		order.OrderId,
		clobtypes.Success,
		orderError,
		status,
		defaultRemovalReason,
	)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderRemoveWithDefaultRemovalReason)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestCreateOrderRemoveMessageWithDefaultReason_InvalidDefault(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()

	require.PanicsWithError(
		t,
		"Invalid parameter: defaultRemovalReason cannot be OrderRemove_ORDER_REMOVAL_REASON_UNSPECIFIED",
		func() {
			ocu.CreateOrderRemoveMessageWithDefaultReason(
				ctx,
				order.OrderId,
				clobtypes.Success,
				orderError,
				status,
				sharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNSPECIFIED,
			)
		},
	)
}

func TestCreateOrderRemoveWithReasonMessage(t *testing.T) {
	ctx, _, _ := sdk.NewSdkContextWithMultistore()

	actualMessage, success := ocu.CreateOrderRemoveMessageWithReason(
		ctx,
		order.OrderId,
		reason,
		status,
	)
	require.True(t, success)

	updateBytes, err := proto.Marshal(&offchainUpdateOrderRemove)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Key:   orderIdHash,
		Value: updateBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}

func TestNewOrderPlaceMessage(t *testing.T) {
	actualUpdateBytes, err := ocu.NewOrderPlaceMessage(
		order,
	)
	require.NoError(
		t,
		err,
		"Encoding OffchainUpdateV1 proto into bytes should not result in an error.",
	)
	actualUpdate := &ocutypes.OffChainUpdateV1{}
	err = proto.Unmarshal(actualUpdateBytes, actualUpdate)
	require.NoError(
		t,
		err,
		"Decoding OffchainUpdateV1 proto bytes should not result in an error.",
	)
	require.Equal(
		t,
		offchainUpdateOrderPlace,
		*actualUpdate,
		"Decoded OffchainUpdateV1 value should be equal to the expected OffchainUpdate proto message",
	)
}

func TestNewOrderUpdateMessage(t *testing.T) {
	actualUpdateBytes, err := ocu.NewOrderUpdateMessage(order.OrderId, totalFilledAmount)
	require.NoError(
		t,
		err,
		"Encoding OffchainUpdateV1 proto into bytes should not result in an error.",
	)
	actualUpdate := &ocutypes.OffChainUpdateV1{}
	err = proto.Unmarshal(actualUpdateBytes, actualUpdate)
	require.NoError(
		t,
		err,
		"Decoding OffchainUpdateV1 proto bytes should not result in an error.",
	)
	require.Equal(
		t,
		offchainUpdateOrderUpdate,
		*actualUpdate,
		"Decoded OffchainUpdateV1 value should be equal to the expected OffchainUpdate proto message",
	)
}

func TestNewOrderRemoveMessage(t *testing.T) {
	actualUpdateBytes, err := ocu.NewOrderRemoveMessage(order.OrderId, reason, status)
	require.NoError(
		t,
		err,
		"Encoding OffchainUpdateV1 proto into bytes should not result in an error.",
	)
	actualUpdate := &ocutypes.OffChainUpdateV1{}
	err = proto.Unmarshal(actualUpdateBytes, actualUpdate)
	require.NoError(
		t,
		err,
		"Decoding OffchainUpdateV1 proto bytes should not result in an error.",
	)
	require.Equal(
		t,
		offchainUpdateOrderRemove,
		*actualUpdate,
		"Decoded OffchainUpdateV1 value should be equal to the expected OffchainUpdate proto message",
	)
}

func TestGetOrderIdHash(t *testing.T) {
	tests := map[string]struct {
		orderId      clobtypes.OrderId
		expectedHash []byte
	}{
		"Can take SHA256 hash of an empty order id": {
			orderId:      clobtypes.OrderId{},
			expectedHash: constants.OrderIdHash_Empty,
		},
		"Can take SHA256 hash of a regular order id": {
			orderId:      order.OrderId,
			expectedHash: constants.OrderIdHash_Alice_Number0_Id0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hash, err := ocu.GetOrderIdHash(tc.orderId)
			require.NoError(t, err)
			require.Equal(t, tc.expectedHash, hash)
		})
	}
}

func TestShouldSendOrderRemovalOnReplay(t *testing.T) {
	tests := map[string]struct {
		// Input
		orderError error

		// Expectations
		expected bool
	}{
		"Returns false for ErrOrderReprocessed": {
			orderError: clobtypes.ErrOrderReprocessed,
			expected:   false,
		},
		"Returns false for ErrInvalidReplacement": {
			orderError: clobtypes.ErrInvalidReplacement,
			expected:   false,
		},
		"Returns false for ErrOrderFullyFilled": {
			orderError: clobtypes.ErrOrderFullyFilled,
			expected:   false,
		},
		"Returns false for ErrOrderIsCanceled": {
			orderError: clobtypes.ErrOrderIsCanceled,
			expected:   false,
		},
		"Returns false for ErrStatefulOrderAlreadyExists": {
			orderError: clobtypes.ErrStatefulOrderAlreadyExists,
			expected:   false,
		},
		"Returns false for ErrHeightExceedsGoodTilBlock": {
			orderError: clobtypes.ErrHeightExceedsGoodTilBlock,
			expected:   false,
		},
		"Returns false for ErrTimeExceedsGoodTilBlockTime": {
			orderError: clobtypes.ErrTimeExceedsGoodTilBlockTime,
			expected:   false,
		},
		"Returns false for wrapped ErrOrderReprocessed": {
			orderError: errorsmod.Wrapf(clobtypes.ErrOrderReprocessed, "wrapped error"),
			expected:   false,
		},
		"Returns false for wrapped ErrInvalidReplacement": {
			orderError: errorsmod.Wrapf(clobtypes.ErrInvalidReplacement, "wrapped error"),
			expected:   false,
		},
		"Returns false for wrapped ErrOrderFullyFilled": {
			orderError: errorsmod.Wrapf(clobtypes.ErrOrderFullyFilled, "wrapped error"),
			expected:   false,
		},
		"Returns false for wrapped ErrOrderIsCanceled": {
			orderError: errorsmod.Wrapf(clobtypes.ErrOrderIsCanceled, "wrapped error"),
			expected:   false,
		},
		"Returns false for wrapped ErrStatefulOrderAlreadyExists": {
			orderError: errorsmod.Wrapf(clobtypes.ErrStatefulOrderAlreadyExists, "wrapped error"),
			expected:   false,
		},
		"Returns false for wrapped ErrHeightExceedsGoodTilBlock": {
			orderError: errorsmod.Wrapf(clobtypes.ErrHeightExceedsGoodTilBlock, "wrapped error"),
			expected:   false,
		},
		"Returns false for wrapped ErrTimeExceedsGoodTilBlockTime": {
			orderError: errorsmod.Wrapf(clobtypes.ErrTimeExceedsGoodTilBlockTime, "wrapped error"),
			expected:   false,
		},
		"Returns false for ErrImmediateExecutionOrderAlreadyFilled": {
			orderError: clobtypes.ErrImmediateExecutionOrderAlreadyFilled,
			expected:   false,
		},
		"Returns true for other error": {
			orderError: clobtypes.ErrFokOrderCouldNotBeFullyFilled,
			expected:   true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			shouldSend := ocu.ShouldSendOrderRemovalOnReplay(tc.orderError)
			require.Equal(t, tc.expected, shouldSend)
		})
	}
}
