package lib

import (
	"encoding/binary"
	"strconv"
)

// IntToString converts any int type to a base-10 string.
func IntToString[T int | int32 | int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

// UintToString converts any uint type to a base-10 string.
func UintToString[T uint | uint32 | uint64](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func Uint32ToBytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func Uint64ToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
