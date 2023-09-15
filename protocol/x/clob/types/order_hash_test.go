package types_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"github.com/stretchr/testify/require"
)

var hash0 = types.OrderHash{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}
var hash1 = types.OrderHash{
	0, 1, 2, 3, 4, 5, 6, 7,
	8, 9, 10, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 20, 21, 22, 23,
	24, 25, 26, 27, 28, 29, 30, 31,
}
var hash2 = types.OrderHash{
	31, 30, 29, 28, 27, 26, 25, 24,
	23, 22, 21, 20, 19, 18, 17, 16,
	15, 14, 13, 12, 11, 10, 9, 8,
	7, 6, 5, 4, 3, 2, 1, 0,
}
var hash3 = types.OrderHash{
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
}

func TestToBytes(t *testing.T) {
	tests := map[string]struct {
		input    types.OrderHash
		expected []byte
	}{
		"Hash1": {
			input:    hash1,
			expected: hash1[:],
		},
		"Hash2": {
			input:    hash2,
			expected: hash2[:],
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.input.ToBytes())
		})
	}
}

func TestLen(t *testing.T) {
	tests := map[string]struct {
		input    []types.OrderHash
		expected int
	}{
		"Nil": {
			input:    nil,
			expected: 0,
		},
		"Empty": {
			input:    []types.OrderHash{},
			expected: 0,
		},
		"Positive": {
			input:    []types.OrderHash{hash1, hash1, hash2, hash1},
			expected: 4,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, types.SortedOrderHashes(tc.input).Len())
		})
	}
}

func TestSwap(t *testing.T) {
	tests := map[string]struct {
		input    []types.OrderHash
		i        int
		j        int
		expected []types.OrderHash
	}{
		"Success": {
			input:    []types.OrderHash{hash1, hash2, hash3},
			i:        0,
			j:        1,
			expected: []types.OrderHash{hash2, hash1, hash3},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			types.SortedOrderHashes(tc.input).Swap(tc.i, tc.j)
			require.Equal(t, tc.expected, tc.input)
		})
	}
}

func TestLess(t *testing.T) {
	tests := map[string]struct {
		input    []types.OrderHash
		i        int
		j        int
		expected bool
	}{
		"Less-Than": {
			input:    []types.OrderHash{hash0, hash1, hash2, hash3},
			i:        1,
			j:        2,
			expected: true,
		},
		"Greater-Than": {
			input:    []types.OrderHash{hash0, hash2, hash1, hash3},
			i:        1,
			j:        2,
			expected: false,
		},
		"Equal": {
			input:    []types.OrderHash{hash1, hash1, hash1, hash1},
			i:        1,
			j:        2,
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, types.SortedOrderHashes(tc.input).Less(tc.i, tc.j))
		})
	}
}

func TestIntegration(t *testing.T) {
	tests := map[string]struct {
		input    []types.OrderHash
		expected []types.OrderHash
	}{
		"Empty": {
			input:    []types.OrderHash{},
			expected: []types.OrderHash{},
		},
		"Already Sorted": {
			input:    []types.OrderHash{hash0, hash1, hash2, hash3},
			expected: []types.OrderHash{hash0, hash1, hash2, hash3},
		},
		"Duplicates": {
			input:    []types.OrderHash{hash3, hash3, hash0, hash2, hash1, hash2, hash3},
			expected: []types.OrderHash{hash0, hash1, hash2, hash2, hash3, hash3, hash3},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sort.Sort(types.SortedOrderHashes(tc.input))
			require.Equal(t, tc.expected, tc.input)
		})
	}
}
