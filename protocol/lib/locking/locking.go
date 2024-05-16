package locking

const (
	GlobalPrefix byte = iota
	AuthSignerPrefix
)

var (
	GlobalKey [][]byte = [][]byte{{GlobalPrefix}}
)

// PrefixLockKeys prefixes each key in the provided slice with the given prefix byte.
func PrefixLockKeys(prefix byte, keys [][]byte) [][]byte {
	result := make([][]byte, len(keys))
	for i, key := range keys {
		result[i] = append([]byte{prefix}, key...)
	}
	return result
}
