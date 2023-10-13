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
