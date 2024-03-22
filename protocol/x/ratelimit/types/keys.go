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

func SplitPendingSendPacketKey(key []byte) (string, uint64, error) {
	err := error(nil)
	parts := bytes.Split(key, []byte("_"))
	if len(parts) != 2 {
		err = fmt.Errorf("unexpected PendingSendPacket key format: %s", key)
		return "", 0, err
	}
	channelId := string(parts[0])

	// convert parts[1] to uint64 parts[1] is a byte array with numeric characters of variable length
	sequenceNumberInt, _ := strconv.Atoi(string(parts[1]))
	sequenceNumber := uint64(sequenceNumberInt)
	return channelId, sequenceNumber, err
}
