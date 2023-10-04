package lib

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBit32ToBytes_Uint32(t *testing.T) {
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
			expected: []byte{0x0f, 0, 0, 0},
		},
		"max uint": {
			// Max uint32 = 4294967295.
			value:    math.MaxUint32,
			expected: []byte{0xff, 0xff, 0xff, 0xff},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := Bit32ToBytes(tc.value)
			require.Equal(t, tc.expected, result)
			require.Equal(t, BytesToUint32(result), tc.value)
		})
	}
}

func TestBit32ToBytes_Int32(t *testing.T) {
	tests := map[string]struct {
		value    int32
		expected []byte
	}{
		"value of -1": {
			value:    -1,
			expected: []byte{0xff, 0xff, 0xff, 0xff},
		},
		"value of zero": {
			value:    0,
			expected: []byte{0, 0, 0, 0},
		},
		"value of 15": {
			value:    15,
			expected: []byte{0x0f, 0, 0, 0},
		},
		"max int": {
			// Max int32 = 2147483647.
			value:    math.MaxInt32,
			expected: []byte{0xff, 0xff, 0xff, 0x7f},
		},
		"min int": {
			// Max int32 = -2147483648.
			value:    math.MinInt32,
			expected: []byte{0x00, 0x00, 0x00, 0x80},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := Bit32ToBytes(tc.value)
			require.Equal(t, tc.expected, result)
			require.Equal(t, BytesToInt32(result), tc.value)
		})
	}
}

func TestBit64ToBytes(t *testing.T) {
	tests := map[string]struct {
		value    int64
		expected []byte
	}{
		"value of -1": {
			value:    -1,
			expected: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		"value of zero": {
			value:    0,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		"value of 15": {
			value:    15,
			expected: []byte{0x0f, 0, 0, 0, 0, 0, 0, 0},
		},
		"max int": {
			value:    math.MaxInt64,
			expected: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f},
		},
		"min int": {
			value:    math.MinInt64,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := Bit64ToBytes(tc.value)
			require.Equal(t, tc.expected, result)
			require.Equal(t, int64(binary.LittleEndian.Uint64(result)), tc.value)
		})
	}
}

func TestIntToString(t *testing.T) {
	require.Equal(t, "15", IntToString(int(15)))
	require.Equal(t, "15", IntToString(int32(15)))
	require.Equal(t, "15", IntToString(int64(15)))
	require.Equal(t, "-15", IntToString(int(-15)))
	require.Equal(t, "-15", IntToString(int32(-15)))
	require.Equal(t, "-15", IntToString(int64(-15)))
	require.Equal(t, "9223372036854775807", IntToString(math.MaxInt64))
	require.Equal(t, "-9223372036854775808", IntToString(math.MinInt64))
}

func TestUintToString(t *testing.T) {
	require.Equal(t, "15", UintToString(uint(15)))
	require.Equal(t, "15", UintToString(uint32(15)))
	require.Equal(t, "15", UintToString(uint64(15)))
	require.Equal(t, "18446744073709551615", UintToString(uint64(math.MaxUint64)))
}
