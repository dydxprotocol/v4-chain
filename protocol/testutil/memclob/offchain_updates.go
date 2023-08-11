package memclob

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func ExpectOnlyMessagesInOffchainUpdates(
	t *testing.T,
	ctx sdk.Context,
	offchainUpdates *types.OffchainUpdates,
	orderPlaces []off_chain_updates.OrderPlace,
	orderUpdates []off_chain_updates.OrderUpdate,
	orderRemoves []off_chain_updates.OrderRemove,
) {
	expectedMessages := make([]msgsender.Message, 0, len(orderUpdates))
	for _, orderPlace := range orderPlaces {
		message := off_chain_updates.MustCreateOrderPlaceMessage(
			ctx.Logger(),
			*orderPlace.Order,
		)
		require.Equal(t, message, offchainUpdates.PlaceMessages[orderPlace.Order.OrderId])
		expectedMessages = append(expectedMessages, message)
	}
	for _, orderUpdate := range orderUpdates {
		message := off_chain_updates.MustCreateOrderUpdateMessage(
			ctx.Logger(),
			*orderUpdate.OrderId,
			satypes.BaseQuantums(orderUpdate.TotalFilledQuantums),
		)
		require.Equal(t, message, offchainUpdates.UpdateMessages[*orderUpdate.OrderId])
		expectedMessages = append(expectedMessages, message)
	}
	for _, orderRemove := range orderRemoves {
		message := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
			ctx.Logger(),
			*orderRemove.RemovedOrderId,
			orderRemove.Reason,
			orderRemove.RemovalStatus,
		)
		require.Equal(t, message, offchainUpdates.RemoveMessages[*orderRemove.RemovedOrderId])
		expectedMessages = append(expectedMessages, message)
	}

	require.ElementsMatch(t, expectedMessages, offchainUpdates.GetMessages())
}

func RequireCancelOrderOffchainUpdate(
	t *testing.T,
	ctx sdk.Context,
	offchainUpdate *types.OffchainUpdates,
	orderId types.OrderId,
) {
	messages := offchainUpdate.GetMessages()
	require.Len(t, messages, 1)
	RequireCancelOrderMessage(t, ctx, &messages[0], orderId)
}

func RequireCancelOrderMessage(t *testing.T, ctx sdk.Context, message *msgsender.Message, orderId types.OrderId) {
	expectedMessage, _ := off_chain_updates.CreateOrderRemoveMessageWithReason(
		ctx.Logger(),
		orderId,
		off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_USER_CANCELED,
		off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
	)
	require.Equal(t, *message, expectedMessage)
}
