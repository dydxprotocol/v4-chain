package types

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
)

type OffchainUpdateMessageType int

// Enum used to track the types of messages, should correspond to the types of messages defined in
// https://github.com/dydxprotocol/v4-proto/blob/main/dydxprotocol/indexer/off_chain_updates/off_chain_updates.proto
const (
	PlaceMessageType OffchainUpdateMessageType = iota
	RemoveMessageType
	UpdateMessageType
	ReplaceMessageType
)

// Represents a single message added to the OffchainUpdates.
// It contains additional metadata needed for message manipulation for specific scenarios, such as
// Replay.
type OffchainUpdateMessage struct {
	Type    OffchainUpdateMessageType
	OrderId OrderId
	Message msgsender.Message
}

// Handles the collection of messages meant for Indexer's off-chain update ingestion.
type OffchainUpdates struct {
	Messages []OffchainUpdateMessage
}

// NewOffchainupdates creates a new OffchainUpdates struct and returns a pointer to it.
func NewOffchainUpdates() *OffchainUpdates {
	return &OffchainUpdates{
		Messages: make([]OffchainUpdateMessage, 0),
	}
}

// AddPlaceMessage adds an off-chain message for the placing of an order to the OffchainUpdates.
func (om *OffchainUpdates) AddPlaceMessage(orderId OrderId, message msgsender.Message) {
	om.Messages = append(om.Messages, OffchainUpdateMessage{PlaceMessageType, orderId, message})
}

// AddUpdateMessage adds an off-chain message for the update of an order to the OffchainUpdates.
func (om *OffchainUpdates) AddUpdateMessage(orderId OrderId, message msgsender.Message) {
	om.Messages = append(om.Messages, OffchainUpdateMessage{UpdateMessageType, orderId, message})
}

// AddRemoveMessage adds an off-chain message for the removal of an order to the OffchainUpdates.
func (om *OffchainUpdates) AddRemoveMessage(orderId OrderId, message msgsender.Message) {
	om.Messages = append(om.Messages, OffchainUpdateMessage{RemoveMessageType, orderId, message})
}

// AddReplaceMessage adds an off-chain message for the replacement of an order to the OffchainUpdates.
func (om *OffchainUpdates) AddReplaceMessage(orderId OrderId, message msgsender.Message) {
	om.Messages = append(om.Messages, OffchainUpdateMessage{ReplaceMessageType, orderId, message})
}

// CondenseMessageForReplay removes all but the last off-chain message for each OrderId from the
// slice of all off-chain messages tracked by the OffchainUpdates struct with the exception of
// OrderPlace messages.
// Intended for use after off-chain messages are generated when replaying multiple operations.
func (om *OffchainUpdates) CondenseMessagesForReplay() {
	seenOrderIds := mapset.NewSet[OrderId]()
	newMessages := make([]OffchainUpdateMessage, 0)
	for i := len(om.Messages) - 1; i >= 0; i-- {
		msg := om.Messages[i]
		if seenOrderIds.Contains(msg.OrderId) {
			// If we've seen an OrderId already, then we won't need to use the message found earlier
			// in the om.Messages slice.
			continue
		}
		seenOrderIds.Add(msg.OrderId)

		// Since we don't need to keep Place message types as Indexer will already have ingested
		// them, then...
		// 1. If the Place message was the finalmost message, we don't need to keep any messages for
		//    that OrderId, as we have all the information we need from it.
		// 2. If the final message was a Remove, we don't need to know about Place or Update messages
		//    because it'll just be removed anyway.
		// 3. Since Update messages only have an "amount filled" parameter, we only need the latest
		//    message.
		if msg.Type == PlaceMessageType {
			continue
		}

		// Note, messages are added in reverse order as we iterated over the initial set of messages
		// in reverse order. We will have to reverse this slice after.
		newMessages = append(newMessages, msg)
	}
	numNewMessages := len(newMessages)
	// Reverse the order of the messages, by swaping messages from the first half of the slice with
	// messages from the second half of the slice.
	for firstHalfIndex := 0; firstHalfIndex < numNewMessages/2; firstHalfIndex++ {
		secondHalfIndex := numNewMessages - firstHalfIndex - 1
		newMessages[firstHalfIndex], newMessages[secondHalfIndex] = newMessages[secondHalfIndex], newMessages[firstHalfIndex]
	}
	om.Messages = newMessages
}

// Append adds all of of the messages from another OffchainUpdates struct in order.
func (om *OffchainUpdates) Append(newMessages *OffchainUpdates) {
	om.Messages = append(om.Messages, newMessages.Messages...)
}

// GetMessages returns all the off-chaim messages in the OffchainUpdates struct.
func (om *OffchainUpdates) GetMessages() []msgsender.Message {
	messages := make([]msgsender.Message, 0, len(om.Messages))
	for _, message := range om.Messages {
		messages = append(messages, message.Message)
	}

	return messages
}
