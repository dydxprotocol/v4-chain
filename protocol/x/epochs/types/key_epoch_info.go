package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// EpochInfoKeyPrefix is the prefix to retrieve all EpochInfo
	EpochInfoKeyPrefix = "EpochInfo/value/"
)

// EpochInfoKey returns the store key to retrieve a EpochInfo from the index fields
func EpochInfoKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
