package locking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixLockKeys(t *testing.T) {
	baseKeyA := []byte("asdf")
	baseKeyB := []byte("023213")
	keyA := []byte("\x00asdf")
	keyB := []byte("\x00023213")

	keys := PrefixLockKeys(GlobalPrefix, [][]byte{baseKeyA, baseKeyB})
	require.Equal(t, [][]byte{keyA, keyB}, keys)
}
