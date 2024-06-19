package memclob

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/typ.v4/slices"
)

func RequireCancelOrderOffchainUpdate(
	t testing.TB,
	ctx sdk.Context,
	offchainUpdate *types.OffchainUpdates,
	orderId types.OrderId,
) {
	messages := offchainUpdate.GetMessages()
	require.Len(t, messages, 1)
	RequireCancelOrderMessage(t, ctx, &messages[0], orderId)
}

func RequireCancelOrderMessage(t testing.TB, ctx sdk.Context, message *msgsender.Message, orderId types.OrderId) {
	expectedMessage, _ := off_chain_updates.CreateOrderRemoveMessageWithReason(
		ctx,
		orderId,
		indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
		ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
	)
	require.Equal(t, *message, expectedMessage)
}

// HasMessage checks if OffchainUpdates contains an message for the given OrderId & OffchainUpdateMessageType.
// Useful for when we want to verify an update message was added for a specific scenario.
func HasMessage(
	offchainUpdates *types.OffchainUpdates,
	orderId types.OrderId,
	messageType types.OffchainUpdateMessageType,
) bool {
	return slices.ContainsFunc(
		offchainUpdates.Messages,
		types.OffchainUpdateMessage{
			Message: msgsender.Message{},
			Type:    messageType,
			OrderId: orderId,
		},
		func(a, b types.OffchainUpdateMessage) bool {
			return a.OrderId == b.OrderId && a.Type == b.Type
		},
	)
}

// MessageCountOfType counts the number of messages in OffchainUpdates of a given type (Place/Remove/Update).
// Useful for when we want to verify we've added the correct number of messages.
func MessageCountOfType(
	offchainUpdates *types.OffchainUpdates,
	messageType types.OffchainUpdateMessageType,
) int {
	var count int = 0
	for _, msg := range offchainUpdates.Messages {
		if msg.Type == messageType {
			count += 1
		}
	}
	return count
}
