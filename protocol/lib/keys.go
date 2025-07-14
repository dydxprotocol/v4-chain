package lib

import (
	"encoding/binary"
)

// Uint32ToKey converts a uint32 to a 4-byte slice in big-endian format.
// The slices can be ordered lexicographically
func Uint32ToKey(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(i))
	return bytes
}

// BytesToUint32 converts a 4-byte slice in big-endian format to a uint32.
func BytesToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}
