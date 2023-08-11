package lib

import (
	"encoding/binary"
	"strconv"
)

func Uint32ToBytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, i)
	return bytes
}

func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Uint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Uint32ToBytesForState converts the uint32 to 4 bytes + '/' to be used as state key prefixes
func Uint32ToBytesForState(id uint32) []byte {
	var key = make([]byte, 5)

	binary.LittleEndian.PutUint32(key, id)
	key[4] = '/'
	return key
}

func StringToUint32(s string) (uint32, error) {
	result, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint32(0), err
	}

	return uint32(result), nil
}

func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func BytesSliceToBytes32(b []byte) [32]byte {
	var byte32 [32]byte
	copy(byte32[:], b)
	return byte32
}
