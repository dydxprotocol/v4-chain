package types

import fmt "fmt"

// Module name and store keys
const (
	// ModuleName defines the module name
	// Use `ratelimit` instead of `ratelimit` to prevent potential key space conflicts with the IBC module.
	ModuleName = "ratelimit"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// DenomCapacityKeyPrefix is the prefix for the key-value store for DenomCapacity
	DenomCapacityKeyPrefix = "DenomCapacity:"

	// LimitParamsKeyPrefix is the prefix for the key-value store for LimitParams
	LimitParamsKeyPrefix = "LimitParams:"

	// PendingSendPacketPrefix is the prefix for the key-value store for PendingSendPacket.
	PendingSendPacketPrefix = "PendingSendPacket:"
)

// State
const ()

func GetPendingSendPacketKey(channelId string, sequenceNumber uint64) []byte {
	return []byte(fmt.Sprintf("%s_%d", channelId, sequenceNumber))
}
