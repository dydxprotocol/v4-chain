package lib_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"github.com/stretchr/testify/require"
)

func TestDedupeSlice(t *testing.T) {
	tests := map[string]struct {
		input  []types.OrderId
		output []types.OrderId
	}{
		"Empty": {
			input:  []types.OrderId{},
			output: []types.OrderId{},
		},
		"No dupes": {
			input: []types.OrderId{
				constants.CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15.OrderId,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
			},
			output: []types.OrderId{
				constants.CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15.OrderId,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
			},
		},
		"Dedupe one": {
			input: []types.OrderId{
				constants.CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15.OrderId,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
			},
			output: []types.OrderId{
				constants.CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15.OrderId,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
			},
		},
		"Dedupe multiple": {
			input: []types.OrderId{
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.OrderId,
				constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.OrderId,
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10.OrderId,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10.OrderId,
			},
			output: []types.OrderId{
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.OrderId,
				constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId,
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10.OrderId,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.output, lib.DedupeSlice(tc.input))
		})
	}
}

func BenchmarkContainsDuplicates_True(b *testing.B) {
	var result bool
	input := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 3, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for i := 0; i < b.N; i++ {
		result = lib.ContainsDuplicates(input)
	}
	require.True(b, result)
}

func BenchmarkContainsDuplicates_False(b *testing.B) {
	var result bool
	input := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for i := 0; i < b.N; i++ {
		result = lib.ContainsDuplicates(input)
	}
	require.False(b, result)
}

func TestContainsDuplicates(t *testing.T) {
	tests := map[string]struct {
		input    []uint32
		expected bool
	}{
		"Nil": {
			input:    nil,
			expected: false,
		},
		"Empty": {
			input:    []uint32{},
			expected: false,
		},
		"One Item": {
			input:    []uint32{10},
			expected: false,
		},
		"False": {
			input:    []uint32{1, 2, 3, 4},
			expected: false,
		},
		"True": {
			input:    []uint32{1, 2, 3, 4, 3},
			expected: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, lib.ContainsDuplicates(tc.input))
		})
	}
}

func BenchmarkMapToSortedSlice(b *testing.B) {
	input := map[string]string{
		"d": "4",
		"b": "2",
		"a": "1",
		"c": "3",
		"e": "5",
		"f": "6",
		"g": "7",
		"h": "8",
		"i": "9",
		"j": "10",
	}
	for i := 0; i < b.N; i++ {
		_ = lib.MapToSortedSlice[sort.StringSlice, string, string](input)
	}
}

func TestMapToSortedSlice(t *testing.T) {
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
		"Single item": {
			inputMap:       map[string]string{"a": "1"},
			expectedResult: []string{"1"},
		},
		"Multiple items": {
			inputMap: map[string]string{
				"d": "4",
				"b": "2",
				"a": "1",
				"c": "3",
			},
			expectedResult: []string{"1", "2", "3", "4"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualResult := lib.MapToSortedSlice[sort.StringSlice](tc.inputMap)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
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

func TestUniqueSliceToSet(t *testing.T) {
	tests := map[string]struct {
		input     []string
		expected  map[string]struct{}
		panicWith string
	}{
		"Empty": {
			input:    []string{},
			expected: map[string]struct{}{},
		},
		"Basic": {
			input: []string{"0", "1", "2"},
			expected: map[string]struct{}{
				"0": {},
				"1": {},
				"2": {},
			},
		},
		"Duplicate": {
			input:     []string{"one", "2", "two", "one", "4"},
			panicWith: "UniqueSliceToSet: duplicate value: one",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panicWith != "" {
				require.PanicsWithValue(
					t,
					tc.panicWith,
					func() { lib.UniqueSliceToSet[string](tc.input) },
				)
			} else {
				require.Equal(t, tc.expected, lib.UniqueSliceToSet[string](tc.input))
			}
		})
	}
}

func TestUniqueSliceToMap(t *testing.T) {
	type testStruct struct {
		Id uint32
	}

	tests := map[string]struct {
		input     []testStruct
		expected  map[uint32]testStruct
		panicWith string
	}{
		"Empty": {
			input:    []testStruct{},
			expected: map[uint32]testStruct{},
		},
		"Basic": {
			input: []testStruct{
				{Id: 0}, {Id: 1}, {Id: 2},
			},
			expected: map[uint32]testStruct{
				0: {Id: 0},
				1: {Id: 1},
				2: {Id: 2},
			},
		},
		"Duplicate": {
			input: []testStruct{
				{Id: 0}, {Id: 0},
			},
			panicWith: "UniqueSliceToMap: duplicate value: {Id:0}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panicWith != "" {
				require.PanicsWithValue(
					t,
					tc.panicWith,
					func() {
						lib.UniqueSliceToMap(tc.input, func(t testStruct) uint32 { return t.Id })
					},
				)
			} else {
				require.Equal(
					t,
					tc.expected,
					lib.UniqueSliceToMap(tc.input, func(t testStruct) uint32 { return t.Id }),
				)
			}
		})
	}
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

func TestMergeAllMapsWithDistinctKeys(t *testing.T) {
	tests := map[string]struct {
		inputMaps []map[string]string

		expectedMap map[string]string
		expectedErr bool
	}{
		"Success: nil input": {
			inputMaps:   nil,
			expectedMap: map[string]string{},
		},
		"Success: single map": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2",
			},
		},
		"Success: single map, empty": {
			inputMaps:   []map[string]string{},
			expectedMap: map[string]string{},
		},
		"Success: multiple maps, all empty or nil": {
			inputMaps: []map[string]string{
				{}, nil,
			},
			expectedMap: map[string]string{},
		},
		"Success: multiple maps, some empty": {
			inputMaps: []map[string]string{
				{}, nil, {"a": "1", "b": "2"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2",
			},
		},
		"Success: multiple maps, no empty": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
				{"c": "3", "d": "4"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2", "c": "3", "d": "4",
			},
		},
		"Error: duplicate keys": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
				{"c": "3", "d": "4"},
				{"a": "5"}, // duplicate key
			},
			expectedErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedErr {
				require.PanicsWithValue(
					t,
					"duplicate key: a",
					func() { lib.MergeAllMapsMustHaveDistinctKeys(tc.inputMaps...) })
			} else {
				actualMap := lib.MergeAllMapsMustHaveDistinctKeys(tc.inputMaps...)
				require.Equal(t, tc.expectedMap, actualMap)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	require.True(
		t,
		lib.SliceContains([]uint32{1, 2}, 1),
	)

	require.False(
		t,
		lib.SliceContains([]uint32{1, 2}, 3),
	)
}
