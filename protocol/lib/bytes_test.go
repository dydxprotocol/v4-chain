package lib

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestUint32ArrayToBytes(t *testing.T) {
	require.Equal(
		t,
		[]byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 4},
		Uint32ArrayToBytes([]uint32{1, 2, 3, 4}),
	)
	require.Equal(
		t,
		[]byte{0, 0, 0, 0, 255, 255, 255, 255},
		Uint32ArrayToBytes([]uint32{0, 4294967295}),
	)
	require.Equal(
		t,
		[]byte{7, 91, 205, 21, 58, 222, 104, 177},
		Uint32ArrayToBytes([]uint32{
			123456789, // 111_01011011_11001101_00010101
			987654321, // 111010_11011110_01101000_10110001
		}),
	)
	require.Equal(
		t,
		[]byte{0, 0, 3, 233, 0, 0, 7, 210, 0, 0, 11, 187, 0, 0, 15, 164},
		Uint32ArrayToBytes([]uint32{
			1001, // 11_11101001
			2002, // 111_11010010
			3003, // 1011_10111011
			4004, // 1111_10100100
		}),
	)
}

func TestBytesToUint32Array(t *testing.T) {
	require.Equal(
		t,
		[]uint32{1, 2, 3, 4},
		BytesToUint32Array([]byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 4}),
	)
	require.Equal(
		t,
		[]uint32{0, 4294967295},
		BytesToUint32Array([]byte{0, 0, 0, 0, 255, 255, 255, 255}),
	)
	require.Equal(
		t,
		[]uint32{
			123456789, // 111_01011011_11001101_00010101
			987654321, // 111010_11011110_01101000_10110001
		},
		BytesToUint32Array([]byte{7, 91, 205, 21, 58, 222, 104, 177}),
	)
	require.Equal(
		t,
		[]uint32{
			1001, // 11_11101001
			2002, // 111_11010010
			3003, // 1011_10111011
			4004, // 1111_10100100
		},
		BytesToUint32Array([]byte{0, 0, 3, 233, 0, 0, 7, 210, 0, 0, 11, 187, 0, 0, 15, 164}),
	)
}
