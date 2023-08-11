package types

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/msgsender"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
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

	require.Equal(t, message0, offchainUpdates.PlaceMessages[orderId0])
	require.Equal(t, []msgsender.Message{message0}, offchainUpdates.GetMessages())
}

func TestAddUpdateMessage(t *testing.T) {
	offchainUpdates := NewOffchainUpdates()
	offchainUpdates.AddUpdateMessage(orderId0, message0)

	require.Equal(t, message0, offchainUpdates.UpdateMessages[orderId0])
	require.Equal(t, []msgsender.Message{message0}, offchainUpdates.GetMessages())
}

func TestAddRemoveMessage(t *testing.T) {
	offchainUpdates := NewOffchainUpdates()
	offchainUpdates.AddRemoveMessage(orderId0, message0)

	require.Equal(t, message0, offchainUpdates.RemoveMessages[orderId0])
	require.Equal(t, []msgsender.Message{message0}, offchainUpdates.GetMessages())
}

func TestBulkUpdate(t *testing.T) {
	tests := map[string]struct {
		// Inputs
		updates    *OffchainUpdates
		newUpdates *OffchainUpdates

		// Expectations
		expectedUpdates *OffchainUpdates
	}{
		"Adds new messages from input updates": {
			updates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId0: message0,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId1: message1,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId2: message2,
				},
			},
			newUpdates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId1: message1,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId2: message2,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId0: message0,
				},
			},
			expectedUpdates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId0: message0,
					orderId1: message1,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId1: message1,
					orderId2: message2,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId2: message2,
					orderId0: message0,
				},
			},
		},
		"Replaces messages from input updates with matching orderId": {
			updates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId0: message0,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId1: message1,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId2: message2,
				},
			},
			newUpdates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId0: message1,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId1: message2,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId2: message0,
				},
			},
			expectedUpdates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId0: message1,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId1: message2,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId2: message0,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.updates.BulkUpdate(tc.newUpdates)
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
		"Updates are ordered by place, update, remove": {
			updates: &OffchainUpdates{
				PlaceMessages: map[OrderId]msgsender.Message{
					orderId0: message0,
				},
				UpdateMessages: map[OrderId]msgsender.Message{
					orderId1: message1,
				},
				RemoveMessages: map[OrderId]msgsender.Message{
					orderId2: message2,
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
