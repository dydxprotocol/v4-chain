package lib_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/stretchr/testify/require"
)

func TestLen(t *testing.T) {
	tests := map[string]struct {
		input    []int
		expected int
	}{
		"Nil": {
			input:    nil,
			expected: 0,
		},
		"Empty": {
			input:    []int{},
			expected: 0,
		},
		"Positive": {
			input:    []int{-1, 5, 0, 2},
			expected: 4,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, lib.Sortable[int](tc.input).Len())
		})
	}
}

func TestSwap(t *testing.T) {
	tests := map[string]struct {
		input    []int
		i        int
		j        int
		expected []int
	}{
		"Success": {
			input:    []int{1, 2, 3, 4},
			i:        1,
			j:        2,
			expected: []int{1, 3, 2, 4},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lib.Sortable[int](tc.input).Swap(tc.i, tc.j)
			require.Equal(t, tc.expected, tc.input)
		})
	}
}

func TestLess(t *testing.T) {
	tests := map[string]struct {
		input    []int
		i        int
		j        int
		expected bool
	}{
		"Less-Than": {
			input:    []int{1, 2, 3, 4},
			i:        1,
			j:        2,
			expected: true,
		},
		"Greater-Than": {
			input:    []int{1, 3, 2, 4},
			i:        1,
			j:        2,
			expected: false,
		},
		"Equal": {
			input:    []int{1, 3, 3, 1},
			i:        1,
			j:        2,
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, lib.Sortable[int](tc.input).Less(tc.i, tc.j))
		})
	}
}

func TestIntegration(t *testing.T) {
	tests := map[string]struct {
		input    []int
		expected []int
	}{
		"Empty": {
			input:    []int{},
			expected: []int{},
		},
		"Already Sorted": {
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		"Duplicates": {
			input:    []int{4, 4, 2, 3, 3, 1},
			expected: []int{1, 2, 3, 3, 4, 4},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sort.Sort(lib.Sortable[int](tc.input))
			require.Equal(t, tc.expected, tc.input)
		})
	}
}
