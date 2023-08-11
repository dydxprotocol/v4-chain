package lib_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/clob/types"

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

func TestConvertMapToSliceOfKeys(t *testing.T) {
	intMap := map[int]bool{0: true, 1: false, 2: false}
	intSlice := lib.ConvertMapToSliceOfKeys(intMap)
	require.ElementsMatch(t, []int{0, 1, 2}, intSlice)

	stringMap := map[string]int{"0": 0, "1": 1, "2": 2}
	stringSlice := lib.ConvertMapToSliceOfKeys(stringMap)
	require.ElementsMatch(t, []string{"0", "1", "2"}, stringSlice)

	emptyMap := map[string]bool{}
	stringSlice = lib.ConvertMapToSliceOfKeys(emptyMap)
	require.ElementsMatch(t, []string{}, stringSlice)
}

func TestContainsValue(t *testing.T) {
	// Empty case.
	require.False(t, lib.ContainsValue([]types.OrderId{}, types.OrderId{}))

	// Contains uint32 case.
	uint32Slice := []uint32{1, 2, 3, 4}
	require.True(t, lib.ContainsValue(uint32Slice, 3))

	// Doesn't contain uint32 case.
	require.False(t, lib.ContainsValue(uint32Slice, 0))

	// Contains string case.
	stringSlice := []string{"hello", "world", "h", "w"}
	require.True(t, lib.ContainsValue(stringSlice, "hello"))

	// Doesn't contain string case.
	require.False(t, lib.ContainsValue(stringSlice, "hh"))
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

func TestMustGetValue(t *testing.T) {
	// Gets middle value successfully.
	require.Equal(
		t,
		uint32(2),
		lib.MustGetValue(
			[]uint32{1, 2, 3, 4},
			1,
		),
	)

	// Gets 0th value successfully.
	require.Equal(
		t,
		"hello",
		lib.MustGetValue(
			[]string{"hello", "world", "hello"},
			0,
		),
	)

	// Gets last value successfully.
	require.Equal(
		t,
		big.NewInt(4),
		lib.MustGetValue(
			[]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4)},
			3,
		),
	)

	// Panics if index is equal to array length.
	require.PanicsWithValue(
		t,
		"MustGetValue: index 2 is greater than or equal to array length 2",
		func() {
			lib.MustGetValue([]uint32{1, 2}, 2)
		},
	)

	// Panics if index is greater than array length.
	require.PanicsWithValue(
		t,
		"MustGetValue: index 3 is greater than or equal to array length 2",
		func() {
			lib.MustGetValue([]uint32{1, 2}, 3)
		},
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
