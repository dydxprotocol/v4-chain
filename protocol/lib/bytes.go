package lib

import (
	"encoding/binary"
	"strconv"
)

// Bit32ToBytes converts 32-bit types into a 4-byte slice in little-endian format.
func Bit32ToBytes[T uint32 | int32](i T) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

// Bit64ToBytes converts 64-bit types into an 8-byte slice in little-endian format.
func Bit64ToBytes[T uint64 | int64](i T) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	return bytes
}

// BytesToUint32 converts a byte slice (of length at least 4) into a uint32.
// The first 4 bytes of the slice are interpreted as little-endian format.
func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

// BytesToInt32 converts a byte slice (of length at least 4) into an int32.
// The first 4 bytes of the slice are interpreted as little-endian format.
func BytesToInt32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(b))
}

// IntToString converts any int type to a base-10 string.
func IntToString[T int | int32 | int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

// UintToString converts any uint type to a base-10 string.
func UintToString[T uint | uint32 | uint64](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}
