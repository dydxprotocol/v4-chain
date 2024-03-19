package types

import (
	"bytes"
	"fmt"
	"strconv"
)

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

func SplitPendingSendPacketKey(key []byte) (channelId string, sequenceNumber uint64) {
	parts := bytes.Split(key, []byte("_"))
	if len(parts) != 2 {
		panic(fmt.Sprintf("unexpected key format: %s", key))
	}
	channelId = string(parts[0])
	// convert parts[1] to uint64 parts[1] is is a byte array with numeric characters of variable length

	sequenceNumberInt, _ := strconv.Atoi(string(parts[1]))
	sequenceNumber = uint64(sequenceNumberInt)
	return
}
