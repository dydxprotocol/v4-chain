package types

import (
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Module name and store keys.
const (
	// The Account module uses "acc" as its module name.
	// KVStore keys cannot have other keys as prefixes so we prepend "dydx" to "accountplus"
	ModuleName = "dydxaccountplus"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)

// Prefix for account state.
const (
	AccountStateKeyPrefix = "AS/"
)

// Below key prefixes are for smart account implementation.
const (
	// SmartAccountKeyPrefix is the prefix key for all smart account store state.
	SmartAccountKeyPrefix = "SA/"

	// ParamsKeyPrefix is the prefix key for smart account params.
	ParamsKeyPrefix = SmartAccountKeyPrefix + "P/"

	// AuthenticatorKeyPrefix is the prefix key for all authenticators.
	AuthenticatorKeyPrefix = SmartAccountKeyPrefix + "A/"

	// AuthenticatorIdKeyPrefix is the prefix key for the next authenticator id.
	AuthenticatorIdKeyPrefix = SmartAccountKeyPrefix + "ID/"
)

func KeyAccountId(account sdk.AccAddress, id uint64) []byte {
	return BuildKey(account.String(), id)
}

// BuildKey creates a key by concatenating the provided elements with the key separator.
func BuildKey(elements ...interface{}) []byte {
	strElements := make([]string, len(elements))
	for i, element := range elements {
		strElements[i] = fmt.Sprint(element)
	}
	return []byte(strings.Join(strElements, "/") + "/")
}
