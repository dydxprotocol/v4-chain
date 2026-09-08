package types

import (
	"encoding/binary"
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var ModuleAddress = authtypes.NewModuleAddress(ModuleName)

const (
	ModuleName = "canonicalusdc"
	StoreKey   = ModuleName

	ControlsKey      = "Controls"
	LedgerKey        = "Ledger"
	PendingCountKey  = "PendingCount"
	CompletedNextKey = "CompletedNext"

	ParticipantKeyPrefix = "Participant:"
	PendingKeyPrefix     = "Pending:"
	CompletedKeyPrefix   = "Completed:"
	CompletedSlotPrefix  = "CompletedSlot:"

	MaxParticipants       = 1_024
	HardMaxPending        = 100_000
	HardMaxCompleted      = 100_000
	MaxMemoLength         = 32_768
	SupportedMemoVersion  = 1
	FundBackingSwapAction = "fund_backing_swap"
)

func ParticipantKey(controller string) []byte {
	return []byte(controller)
}

// PendingKey is collision-free because validated channel identifiers cannot
// contain NUL. Route is included so state is explicitly keyed by
// route/channel/sequence as required by the settlement contract.
func PendingKey(route Route, channel string, sequence uint64) []byte {
	key := make([]byte, 1+len(channel)+1+8)
	key[0] = byte(route)
	copy(key[1:], channel)
	binary.BigEndian.PutUint64(key[len(key)-8:], sequence)
	return key
}

func ParsePendingKey(key []byte) (Route, string, uint64, error) {
	if len(key) < 11 || key[len(key)-9] != 0 {
		return Route_ROUTE_UNSPECIFIED, "", 0, fmt.Errorf("invalid pending settlement key")
	}
	return Route(key[0]), string(key[1 : len(key)-9]), binary.BigEndian.Uint64(key[len(key)-8:]), nil
}

func CompletedKey(channel string, sequence uint64) []byte {
	key := make([]byte, len(channel)+1+8)
	copy(key, channel)
	binary.BigEndian.PutUint64(key[len(key)-8:], sequence)
	return key
}

func CompletedSlotKey(slot uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, slot)
	return key
}
