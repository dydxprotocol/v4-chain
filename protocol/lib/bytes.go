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

// Uint32ArrayToBytes converts a slice of uint32 to a byte slice.
func Uint32ArrayToBytes(arr []uint32) []byte {
	buf := make([]byte, len(arr)*4)
	for i, v := range arr {
		binary.BigEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}

// BytesToUint32Array converts a byte slice to a slice of uint32.
func BytesToUint32Array(b []byte) []uint32 {
	arr := make([]uint32, len(b)/4)
	for i := range arr {
		arr[i] = binary.BigEndian.Uint32(b[i*4:])
	}
	return arr
}
