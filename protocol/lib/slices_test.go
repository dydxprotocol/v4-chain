package lib_test

import (
	"math/big"
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"github.com/stretchr/testify/require"
)

func TestContainsDuplicates(t *testing.T) {
	// Empty case.
	require.False(t, lib.ContainsDuplicates([]types.OrderId{}))

	// Unique uint32 case.
	allUniqueUint32s := []uint32{1, 2, 3, 4}
	require.False(t, lib.ContainsDuplicates(allUniqueUint32s))

	// Duplicate uint32 case.
	containsDuplicateUint32 := append(allUniqueUint32s, 3)
	require.True(t, lib.ContainsDuplicates(containsDuplicateUint32))

	// Unique string case.
	allUniqueStrings := []string{"hello", "world", "h", "w"}
	require.False(t, lib.ContainsDuplicates(allUniqueStrings))

	// Duplicate string case.
	containsDuplicateString := append(allUniqueStrings, "world")
	require.True(t, lib.ContainsDuplicates(containsDuplicateString))
}

func TestGetSortedKeys(t *testing.T) {
	tests := map[string]struct {
		inputMap       map[string]string
		expectedResult []string
	}{
		"Nil input": {
			inputMap:       nil,
			expectedResult: []string{},
		},
		"Empty map": {
			inputMap:       map[string]string{},
			expectedResult: []string{},
		},
		"Non-empty map": {
			inputMap: map[string]string{
				"d": "4", "b": "2", "a": "1", "c": "3",
			},
			expectedResult: []string{"a", "b", "c", "d"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualResult := lib.GetSortedKeys[sort.StringSlice](tc.inputMap)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestMustRemoveIndex(t *testing.T) {
	// Remove first index of uint32 slice case and doesn't modify original slice.
	uint32Slice := []uint32{1, 2, 3, 4}
	require.Equal(t, []uint32{2, 3, 4}, lib.MustRemoveIndex(uint32Slice, 0))
	require.Equal(t, []uint32{1, 2, 3, 4}, uint32Slice)

	// Remove last index of uint32 case and doesn't modify original slice.
	require.Equal(t, []uint32{1, 2, 3}, lib.MustRemoveIndex(uint32Slice, 3))
	require.Equal(t, []uint32{1, 2, 3, 4}, uint32Slice)

	// Remove middle element of uint32 case and doesn't modify original slice.
	require.Equal(t, []uint32{1, 2, 4}, lib.MustRemoveIndex(uint32Slice, 2))
	require.Equal(t, []uint32{1, 2, 3, 4}, uint32Slice)

	// Remove first index of string slice case and doesn't modify original slice.
	stringSlice := []string{"h", "e", "l", "l", "o"}
	require.Equal(t, []string{"e", "l", "l", "o"}, lib.MustRemoveIndex(stringSlice, 0))
	require.Equal(t, []string{"h", "e", "l", "l", "o"}, stringSlice)

	// Remove last index of string case and doesn't modify original slice.
	require.Equal(t, []string{"h", "e", "l", "l"}, lib.MustRemoveIndex(stringSlice, 4))
	require.Equal(t, []string{"h", "e", "l", "l", "o"}, stringSlice)

	// Remove middle element of string case and doesn't modify original slice.
	require.Equal(t, []string{"h", "e", "l", "o"}, lib.MustRemoveIndex(stringSlice, 2))
	require.Equal(t, []string{"h", "e", "l", "l", "o"}, stringSlice)

	// Panics if provided index greater than slice length.
	require.PanicsWithValue(
		t,
		"MustRemoveIndex: index 0 is greater than array length 0",
		func() {
			lib.MustRemoveIndex([]types.OrderId{}, 0)
		},
	)
}

func TestMapSlice(t *testing.T) {
	// Can increment all numbers in a slice by 1, and change type to `uint64`.
	require.Equal(
		t,
		[]uint64{2, 3, 4, 5},
		lib.MapSlice(
			[]uint32{1, 2, 3, 4},
			func(a uint32) uint64 {
				return uint64(a + 1)
			},
		),
	)

	// Can return the length of all strings in a slice.
	require.Equal(
		t,
		[]int{1, 2, 3, 5, 0},
		lib.MapSlice(
			[]string{"1", "22", "333", "hello", ""},
			func(a string) int {
				return len(a)
			},
		),
	)

	// Works properly on empty slice.
	require.Equal(
		t,
		[]int{},
		lib.MapSlice(
			[]string{},
			func(a string) int {
				return 1000
			},
		),
	)

	// Works properly on constant function.
	require.Equal(
		t,
		[]bool{true, true, true},
		lib.MapSlice(
			[]string{"hello", "world", "hello"},
			func(a string) bool {
				return true
			},
		),
	)
}

func TestFilterSlice(t *testing.T) {
	// Can filter out all numbers less than 3.
	require.Equal(
		t,
		[]uint32{1, 2},
		lib.FilterSlice(
			[]uint32{1, 2, 3, 4},
			func(a uint32) bool {
				return a < 3
			},
		),
	)

	// Can filter out all strings that have length greater than 3.
	require.Equal(
		t,
		[]string{"hello"},
		lib.FilterSlice(
			[]string{"1", "22", "333", "hello"},
			func(a string) bool {
				return len(a) > 3
			},
		),
	)

	// Works properly on empty slice.
	require.Equal(
		t,
		[]string{},
		lib.FilterSlice(
			[]string{},
			func(a string) bool {
				return true
			},
		),
	)

	// Works properly on constant function that always returns true.
	require.Equal(
		t,
		[]string{"hello", "world", "hello"},
		lib.FilterSlice(
			[]string{"hello", "world", "hello"},
			func(a string) bool {
				return true
			},
		),
	)

	// Works properly on constant function that always returns false.
	require.Equal(
		t,
		[]string{},
		lib.FilterSlice(
			[]string{"hello", "world", "hello"},
			func(a string) bool {
				return false
			},
		),
	)
}

func TestSliceToSet(t *testing.T) {
	slice := make([]int, 0)
	for i := 0; i < 3; i++ {
		slice = append(slice, i)
	}
	set := lib.SliceToSet(slice)
	require.Equal(
		t,
		map[int]struct{}{
			0: {},
			1: {},
			2: {},
		},
		set,
	)
	stringSlice := []string{
		"one",
		"two",
	}
	stringSet := lib.SliceToSet(stringSlice)
	require.Equal(
		t,
		map[string]struct{}{
			"one": {},
			"two": {},
		},
		stringSet,
	)

	emptySlice := []types.OrderId{}
	emptySet := lib.SliceToSet(emptySlice)
	require.Equal(
		t,
		map[types.OrderId]struct{}{},
		emptySet,
	)
}

func TestSliceToSet_PanicOnDuplicate(t *testing.T) {
	stringSlice := []string{
		"one",
		"two",
		"one",
	}
	require.PanicsWithValue(
		t,
		"SliceToSet: duplicate value: one",
		func() {
			lib.SliceToSet(stringSlice)
		},
	)
}
