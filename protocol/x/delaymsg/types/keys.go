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
	BlockMessageIdsPrefix = "BlockMsgIds:"

	// DelayedMessageKeyPrefix is the prefix to retrieve all DelayedMessages.
	DelayedMessageKeyPrefix = "Msg:"

	// NextDelayedMessageIdKey is the key to retrieve next delayed message id.
	NextDelayedMessageIdKey = "NextDelayedMessageId"
)
