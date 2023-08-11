package types

import (
	"github.com/dydxprotocol/v4/indexer/msgsender"
)

type OffchainUpdates struct {
	PlaceMessages  map[OrderId]msgsender.Message
	UpdateMessages map[OrderId]msgsender.Message
	RemoveMessages map[OrderId]msgsender.Message
}

func NewOffchainUpdates() *OffchainUpdates {
	return &OffchainUpdates{
		PlaceMessages:  make(map[OrderId]msgsender.Message),
		UpdateMessages: make(map[OrderId]msgsender.Message),
		RemoveMessages: make(map[OrderId]msgsender.Message),
	}
}

func (om *OffchainUpdates) AddPlaceMessage(orderId OrderId, message msgsender.Message) {
	om.PlaceMessages[orderId] = message
}

func (om *OffchainUpdates) AddUpdateMessage(orderId OrderId, message msgsender.Message) {
	om.UpdateMessages[orderId] = message
}
func (om *OffchainUpdates) AddRemoveMessage(orderId OrderId, message msgsender.Message) {
	om.RemoveMessages[orderId] = message
}

func (om *OffchainUpdates) ClearPlaceMessages() {
	om.PlaceMessages = make(map[OrderId]msgsender.Message)
}

func (om *OffchainUpdates) BulkUpdate(newMessages *OffchainUpdates) {
	for orderId, placeMessage := range newMessages.PlaceMessages {
		om.PlaceMessages[orderId] = placeMessage
	}
	for orderId, updateMessage := range newMessages.UpdateMessages {
		om.UpdateMessages[orderId] = updateMessage
	}
	for orderId, removeMessage := range newMessages.RemoveMessages {
		om.RemoveMessages[orderId] = removeMessage
	}
}

func (om *OffchainUpdates) GetMessages() []msgsender.Message {
	count := len(om.PlaceMessages) + len(om.UpdateMessages) + len(om.RemoveMessages)
	messages := make([]msgsender.Message, 0, count)

	for _, placeMessage := range om.PlaceMessages {
		messages = append(messages, placeMessage)
	}
	for _, updateMessage := range om.UpdateMessages {
		messages = append(messages, updateMessage)
	}
	for _, removeMessage := range om.RemoveMessages {
		messages = append(messages, removeMessage)
	}

	return messages
}
