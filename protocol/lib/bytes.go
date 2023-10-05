package lib

import (
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
