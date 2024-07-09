package types

import (
	"bytes"
	"fmt"
	"strconv"

	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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

	/* For the sDAI middleware */

	PoolAccount = ibctransfertypes.ModuleName

	// This is the name of the sDAI token
	SDaiDenom = "sDAI" // TODO: change

	// This is the name of the trading DAI token
	TradingDAIDenom = "tradingDAI"

	// sDAIKeyPrefix is the prefix for the key-value store forthe sDAI price
	SDAIKeyPrefix = "SDAIPrice:"

	// storedDripRate is the prefix for the key-value store for historical drip rate
	StoredDripRatePrefix = "StoredDripRatePrefix:"

	// DaiYieldEpochPrefix is the prefix for the key-value store for DaiYieldEpoch
	// The key vakue store is implemented as an array of size 100
	DaiYieldEpochPrefix = "DaiYieldEpoch:"
)

// State
const (

	// The number of ethereum blocks we store the sDAI rate for
	ETH_BLOCKS_TO_STORE = 5

	// base 10
	BASE_10 = 10

	// Maker RAY value which stores decimal points
	SDAI_DECIMALS = 27

	// The number of DaiYieldEpochParams we store
	DAI_YIELD_ARRAY_SIZE = 1000

	// The minimum number of blocks until we accept a new epoch
	DAI_YIELD_MIN_EPOCH_BLOCKS = 100000
)

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
