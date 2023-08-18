package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// SubaccountKeyPrefix is the prefix to retrieve all Subaccount
	SubaccountKeyPrefix = "Subaccount/value/"
	// SubaccountWithDecreasedNetCollateralKeyPrefix is the prefix to retrieve all subaccounts with decreased
	// net collateral.
	SubaccountWithDecreasedNetCollateralKeyPrefix = "SubaccountWithDecreasedNetCollateral/value/"
)

// SubaccountKey returns the store key to retrieve a Subaccount from the index fields
func SubaccountKey(
	id SubaccountId,
) []byte {
	var key []byte

	idBytes, err := id.Marshal()
	if err != nil {
		panic(err)
	}
	key = append(key, idBytes...)
	key = append(key, []byte("/")...)

	return key
}
