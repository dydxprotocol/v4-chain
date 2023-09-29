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

func Int32ToBytes(i int32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

func Int64ToBytes(i int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	return bytes
}

func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func BytesToInt32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(b))
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Uint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

func StringToUint32(s string) (uint32, error) {
	result, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint32(0), err
	}

	return uint32(result), nil
}
