package types

import (
	"bytes"
	"fmt"
	"strconv"

	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
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

	// Addresses of tDAI and sDAI pools in x/bank module

	// TDaiPool: Address, where tDAI yield is held before it is claimed by subaccount.
	TDaiPoolAccount = ibctransfertypes.ModuleName
	// SDaiPool: Address, where bridged sDAI is held until it is bridged out.
	SDaiPoolAccount = "sDAIPoolAccount"

	// Denom of sDAI in the x/bank module
	SDaiDenom               = "ibc/DEEFE2DEFDC8EA8879923C4CCA42BB888C3CD03FF7ECFEFB1C2FEC27A732ACC8"
	SDaiBaseDenom           = "gsdai"
	SDaiBaseDenomPathPrefix = "transfer/channel-0"
	SDaiBaseDenomFullPath   = SDaiBaseDenomPathPrefix + "/" + SDaiBaseDenom
	SDaiDenomExponent       = -18

	// Denom of tDAI in the x/bank module
	TDaiDenom = assettypes.TDaiDenom

	// sDAIKeyPrefix is the prefix for the key-value store forthe sDAI price
	SDaiKeyPrefix = "SDAIPrice:"

	// SDAILastBlockUpdate is the prefix for the key-value store for the last block that the sDAI price was updated
	SDAILastBlockUpdate = "SDAILastBlockUpdate:"

	// AssetYieldIndexPrefix is the prefix for the key value store that tracks
	// the cumulative yield index across all yield epochs.
	AssetYieldIndexPrefix = "AssetYieldIndex:"
)

// State
const (

	// base 10
	BASE_10 = 10

	// Maker RAY value which stores decimal points
	SDAI_DECIMALS = 27

	SDAI_UPDATE_BLOCK_DELAY = 5000
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
