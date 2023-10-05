package lib

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint32ToKey(t *testing.T) {
	tests := map[string]struct {
		value    uint32
		expected []byte
	}{
		"value of zero": {
			value:    0,
			expected: []byte{0, 0, 0, 0},
		},
		"value of 15": {
			value:    15,
			expected: []byte{0, 0, 0, 0x0f},
		},
		"max uint": {
			// Max uint32 = 4294967295.
			value:    math.MaxUint32,
			expected: []byte{0xff, 0xff, 0xff, 0xff},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := Uint32ToKey(tc.value)
			require.Equal(t, tc.expected, result)
			require.Equal(t, binary.BigEndian.Uint32(result), tc.value)
		})
	}
}
