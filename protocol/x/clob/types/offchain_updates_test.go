package types

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	subaccountId = satypes.SubaccountId{
		Owner:  "dydx1x2hd82qerp7lc0kf5cs3yekftupkrl620te6u2",
		Number: 0,
	}
	orderId0 = OrderId{
		SubaccountId: subaccountId,
		ClientId:     0,
	}
	orderId1 = OrderId{
		SubaccountId: subaccountId,
		ClientId:     1,
	}
	orderId2 = OrderId{
		SubaccountId: subaccountId,
		ClientId:     2,
	}
	message0 = msgsender.Message{
		Key:   []byte("key0"),
		Value: []byte("value0"),
	}
	message1 = msgsender.Message{
		Key:   []byte("key1"),
		Value: []byte("value0"),
	}
	message2 = msgsender.Message{
		Key:   []byte("key2"),
		Value: []byte("value0"),
	}
)

func TestAddPlaceMessage(t *testing.T) {
	offchainUpdates := NewOffchainUpdates()
	offchainUpdates.AddPlaceMessage(orderId0, message0)

	require.Equal(t, message0, offchainUpdates.Messages[0].Message)
	require.Equal(t, PlaceMessageType, offchainUpdates.Messages[0].Type)
	require.Equal(t, orderId0, offchainUpdates.Messages[0].OrderId)
	require.Equal(t, []msgsender.Message{message0}, offchainUpdates.GetMessages())
}

func TestAddUpdateMessage(t *testing.T) {
	offchainUpdates := NewOffchainUpdates()
	offchainUpdates.AddUpdateMessage(orderId0, message0)

	require.Equal(t, message0, offchainUpdates.Messages[0].Message)
	require.Equal(t, UpdateMessageType, offchainUpdates.Messages[0].Type)
	require.Equal(t, orderId0, offchainUpdates.Messages[0].OrderId)
	require.Equal(t, []msgsender.Message{message0}, offchainUpdates.GetMessages())
}

func TestAddRemoveMessage(t *testing.T) {
	offchainUpdates := NewOffchainUpdates()
	offchainUpdates.AddRemoveMessage(orderId0, message0)

	require.Equal(t, message0, offchainUpdates.Messages[0].Message)
	require.Equal(t, RemoveMessageType, offchainUpdates.Messages[0].Type)
	require.Equal(t, orderId0, offchainUpdates.Messages[0].OrderId)
	require.Equal(t, []msgsender.Message{message0}, offchainUpdates.GetMessages())
}

func TestAppend(t *testing.T) {
	tests := map[string]struct {
		// Inputs
		updates    *OffchainUpdates
		newUpdates *OffchainUpdates

		// Expectations
		expectedUpdates *OffchainUpdates
	}{
		"Adds new messages from input updates": {
			updates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{PlaceMessageType, orderId0, message0},
					{UpdateMessageType, orderId1, message1},
					{RemoveMessageType, orderId2, message2},
				},
			},
			newUpdates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{PlaceMessageType, orderId1, message1},
					{UpdateMessageType, orderId2, message2},
					{RemoveMessageType, orderId0, message0},
				},
			},
			expectedUpdates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{PlaceMessageType, orderId0, message0},
					{UpdateMessageType, orderId1, message1},
					{RemoveMessageType, orderId2, message2},
					{PlaceMessageType, orderId1, message1},
					{UpdateMessageType, orderId2, message2},
					{RemoveMessageType, orderId0, message0},
				},
			},
		},
		"Can have multiple messages of the same type with matching orderIds": {
			updates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{PlaceMessageType, orderId0, message0},
					{UpdateMessageType, orderId1, message1},
					{RemoveMessageType, orderId2, message2},
				},
			},
			newUpdates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{PlaceMessageType, orderId0, message1},
					{UpdateMessageType, orderId1, message2},
					{RemoveMessageType, orderId2, message0},
				},
			},
			expectedUpdates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{PlaceMessageType, orderId0, message0},
					{UpdateMessageType, orderId1, message1},
					{RemoveMessageType, orderId2, message2},
					{PlaceMessageType, orderId0, message1},
					{UpdateMessageType, orderId1, message2},
					{RemoveMessageType, orderId2, message0},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.updates.Append(tc.newUpdates)
			require.Equal(t, tc.updates, tc.expectedUpdates)
		})
	}
}

func TestGetMessages(t *testing.T) {
	tests := map[string]struct {
		// Inputs
		updates *OffchainUpdates

		// Expectations
		expectedMessages []msgsender.Message
	}{
		"Empty updates": {
			updates:          NewOffchainUpdates(),
			expectedMessages: []msgsender.Message{},
		},
		"Updates are ordered by how they were added": {
			updates: &OffchainUpdates{
				Messages: []OffchainUpdateMessage{
					{RemoveMessageType, orderId0, message0},
					{PlaceMessageType, orderId0, message1},
					{RemoveMessageType, orderId0, message2},
				},
			},
			expectedMessages: []msgsender.Message{message0, message1, message2},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedMessages, tc.updates.GetMessages())
		})
	}
}

func TestCondenseMessagesForReplay(t *testing.T) {
	messageFirstPlace := msgsender.Message{Key: []byte("PLACE#1"), Value: []byte("1st placement")}
	messageSecondPlace := msgsender.Message{Key: []byte("PLACE#2"), Value: []byte("2nd placement")}
	messageFirstRemove := msgsender.Message{Key: []byte("REMOVE#1"), Value: []byte("1st removal")}
	messageSecondRemove := msgsender.Message{Key: []byte("REMOVE#2"), Value: []byte("2nd removal")}
	messageFirstUpdate := msgsender.Message{Key: []byte("UPDATE#1"), Value: []byte("1st update")}
	messageSecondUpdate := msgsender.Message{Key: []byte("UPDATE#2"), Value: []byte("2nd update")}
	messageThirdUpdate := msgsender.Message{Key: []byte("UPDATE#3"), Value: []byte("3rd update")}
	messageFourthUpdate := msgsender.Message{Key: []byte("UPDATE#4"), Value: []byte("4th update")}

	orderId := OrderId{SubaccountId: subaccountId, ClientId: 1}

	tests := map[string]struct {
		// Inputs
		messages []struct {
			messageType OffchainUpdateMessageType
			message     msgsender.Message
		}

		// Expectations
		expectedMessages []OffchainUpdateMessage
	}{
		"No updates, placed only": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"No updates, removed and then placed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"No updates, placed then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: RemoveMessageType, message: messageFirstRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageFirstRemove,
				},
			},
		},
		"No updates, placed then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"No updates, placed then replaced then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Single update": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageFirstUpdate,
				},
			},
		},
		"Updated then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageFirstRemove,
				},
			},
		},
		"Updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Placed, then updated": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageFirstUpdate,
				},
			},
		},
		"Replaced, then updated": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageFirstUpdate,
				},
			},
		},
		"Placed, updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageFirstRemove,
				},
			},
		},
		"Placed, updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Replaced, updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Replaced, updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: PlaceMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Updated, then updated again": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageSecondUpdate,
				},
			},
		},
		"Updated, then updated and removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageFirstRemove,
				},
			},
		},
		"Updated, then updated and replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Placed, then updated, then updated again": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageSecondUpdate,
				},
			},
		},
		"Replaced, then updated, then updated again": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageSecondUpdate,
				},
			},
		},
		"Placed, then updated, then updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageFirstRemove,
				},
			},
		},
		"Placed, then updated, then updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Replaced, then updated, then updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Replaced, then updated, then updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Updated, replaced, then updated": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageSecondUpdate,
				},
			},
		},
		"Updated, replaced, then updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Updated, replaced, then updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Updated, replaced, then updated, then updated again": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageThirdUpdate,
				},
			},
		},
		"Updated, replaced, then updated, then updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Updated, replaced, then updated, then updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Updated, then updated, then replaced, then updated again": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageThirdUpdate,
				},
			},
		},
		"Updated, then updated, then replaced, then updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Updated, then updated, then replaced, then updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Updated, then updated, then replaced, then updated, then updated again": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: UpdateMessageType, message: messageFourthUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageFourthUpdate,
				},
			},
		},
		"Updated, then updated, then replaced, then updated, then updated, then removed": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: UpdateMessageType, message: messageFourthUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					RemoveMessageType,
					orderId,
					messageSecondRemove,
				},
			},
		},
		"Updated, then updated, then replaced, then updated, then updated, then replaced": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: UpdateMessageType, message: messageFourthUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
			},
			expectedMessages: []OffchainUpdateMessage{},
		},
		"Updated, then replaced, then updated, then replaced, then updated": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageThirdUpdate,
				},
			},
		},
		"Updated, then replaced, then updated, then replaced, then updated, then updated": {
			messages: []struct {
				messageType OffchainUpdateMessageType
				message     msgsender.Message
			}{
				{messageType: UpdateMessageType, message: messageFirstUpdate},
				{messageType: RemoveMessageType, message: messageFirstRemove},
				{messageType: PlaceMessageType, message: messageFirstPlace},
				{messageType: UpdateMessageType, message: messageSecondUpdate},
				{messageType: RemoveMessageType, message: messageSecondRemove},
				{messageType: PlaceMessageType, message: messageSecondPlace},
				{messageType: UpdateMessageType, message: messageThirdUpdate},
				{messageType: UpdateMessageType, message: messageFourthUpdate},
			},
			expectedMessages: []OffchainUpdateMessage{
				{
					UpdateMessageType,
					orderId,
					messageFourthUpdate,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			offchainUpdates := NewOffchainUpdates()
			for _, message := range tc.messages {
				if message.messageType == PlaceMessageType {
					offchainUpdates.AddPlaceMessage(orderId, message.message)
				} else if message.messageType == RemoveMessageType {
					offchainUpdates.AddRemoveMessage(orderId, message.message)
				} else if message.messageType == UpdateMessageType {
					offchainUpdates.AddUpdateMessage(orderId, message.message)
				} else {
					panic("unknown message type")
				}
			}
			offchainUpdates.CondenseMessagesForReplay()
			require.Equal(t, tc.expectedMessages, offchainUpdates.Messages)
		})
	}
}

func TestCondenseMessagesForReplay_MultipleOrderIds(t *testing.T) {
	messageFirstPlace := msgsender.Message{Key: []byte("PLACE#1"), Value: []byte("1st placement")}
	messageFirstRemove := msgsender.Message{Key: []byte("REMOVE#1"), Value: []byte("1st removal")}
	messageFirstUpdate := msgsender.Message{Key: []byte("UPDATE#1"), Value: []byte("1st update")}
	messageSecondUpdate := msgsender.Message{Key: []byte("UPDATE#2"), Value: []byte("2nd update")}
	messageThirdUpdate := msgsender.Message{Key: []byte("UPDATE#3"), Value: []byte("3rd update")}

	offchainUpdates := NewOffchainUpdates()

	orderId_PR := OrderId{SubaccountId: subaccountId, ClientId: 1}
	orderId_PUU := OrderId{SubaccountId: subaccountId, ClientId: 2}
	orderId_URPUU := OrderId{SubaccountId: subaccountId, ClientId: 3}

	// Place then remove.
	offchainUpdates.AddPlaceMessage(orderId_PR, messageFirstPlace)
	offchainUpdates.AddRemoveMessage(orderId_PR, messageFirstRemove)

	// Place then update twice.
	offchainUpdates.AddPlaceMessage(orderId_PUU, messageFirstPlace)
	offchainUpdates.AddUpdateMessage(orderId_PUU, messageFirstUpdate)
	offchainUpdates.AddUpdateMessage(orderId_PUU, messageSecondUpdate)

	// Update then replace then update twice.
	offchainUpdates.AddUpdateMessage(orderId_URPUU, messageFirstUpdate)
	offchainUpdates.AddPlaceMessage(orderId_URPUU, messageFirstPlace)
	offchainUpdates.AddRemoveMessage(orderId_URPUU, messageFirstRemove)
	offchainUpdates.AddUpdateMessage(orderId_URPUU, messageSecondUpdate)
	offchainUpdates.AddUpdateMessage(orderId_URPUU, messageThirdUpdate)

	offchainUpdates.CondenseMessagesForReplay()

	expectedMessages := []OffchainUpdateMessage{
		{
			RemoveMessageType,
			orderId_PR,
			messageFirstRemove,
		},
		{
			UpdateMessageType,
			orderId_PUU,
			messageSecondUpdate,
		},
		{
			UpdateMessageType,
			orderId_URPUU,
			messageThirdUpdate,
		},
	}
	require.Equal(t, expectedMessages, offchainUpdates.Messages)
}
