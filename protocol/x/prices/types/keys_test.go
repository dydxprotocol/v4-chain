package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyPrefix(t *testing.T) {
	b := []byte("a")
	require.Equal(t, uint8(0x61), b[0])
	require.Equal(t, 1, len(b))
}
