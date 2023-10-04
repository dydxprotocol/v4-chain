package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "delaymsg"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// BlockMessageIdsPrefix is the prefix to retrieve all BlockMessageIds for a given block height.
	BlockMessageIdsPrefix = "block_message_ids/"

	// DelayedMessageKeyPrefix is the prefix to retrieve all DelayedMessages.
	DelayedMessageKeyPrefix = "delayed_message/"

	// NumDelayedMessagesKey is the prefix to retrieve the number of DelayedMessages.
	NumDelayedMessagesKey = "num_delayed_messages"
)
