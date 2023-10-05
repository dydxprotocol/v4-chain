package lib

import (
	"encoding/binary"
)

// Uint32ToKey
func Uint32ToKey(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(i))
	return bytes
}
