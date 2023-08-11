package maps_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4/lib/maps"
	"github.com/stretchr/testify/require"
)

func TestMergeAllMapsWithDistinctKeys(t *testing.T) {
	tests := map[string]struct {
		inputMaps []map[string]string

		expectedMap map[string]string
		expectedErr bool
	}{
		"Success: nil input": {
			inputMaps:   nil,
			expectedMap: map[string]string{},
			expectedErr: false,
		},
		"Success: single map": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2",
			},
			expectedErr: false,
		},
		"Success: single map, empty": {
			inputMaps:   []map[string]string{},
			expectedMap: map[string]string{},
			expectedErr: false,
		},
		"Success: multiple maps, all empty or nil": {
			inputMaps: []map[string]string{
				{}, nil,
			},
			expectedMap: map[string]string{},
			expectedErr: false,
		},
		"Success: multiple maps, some empty": {
			inputMaps: []map[string]string{
				{}, nil, {"a": "1", "b": "2"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2",
			},
			expectedErr: false,
		},
		"Success: multiple maps, no empty": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
				{"c": "3", "d": "4"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2", "c": "3", "d": "4",
			},
			expectedErr: false,
		},
		"Error: duplicate keys": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
				{"c": "3", "d": "4"},
				{"a": "5"}, // duplicate key
			},
			expectedMap: nil,
			expectedErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedErr {
				require.PanicsWithValue(
					t,
					"duplicate key: a",
					func() { maps.MergeAllMapsMustHaveDistinctKeys(tc.inputMaps...) })
			} else {
				actualMap := maps.MergeAllMapsMustHaveDistinctKeys(tc.inputMaps...)
				require.Equal(t, tc.expectedMap, actualMap)
			}
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
			actualResult := maps.GetSortedKeys(tc.inputMap)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestInvertMustHaveDistinctValues(t *testing.T) {
	tests := map[string]struct {
		inputMap       map[string]string
		expectedResult map[string]string
		expectedErr    error
	}{
		"Nil input": {
			inputMap:       nil,
			expectedResult: map[string]string{},
		},
		"Empty map": {
			inputMap:       map[string]string{},
			expectedResult: map[string]string{},
		},
		"Non-empty map": {
			inputMap: map[string]string{
				"a": "1", "b": "2", "c": "3", "d": "4",
			},
			expectedResult: map[string]string{
				"1": "a", "2": "b", "3": "c", "4": "d",
			},
		},
		"Error: duplicate values": {
			inputMap: map[string]string{
				"a": "1", "b": "2", "c": "3", "d": "3",
			},
			expectedErr: errors.New("duplicate map value: 3"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedErr != nil {
				require.PanicsWithValue(
					t,
					tc.expectedErr.Error(),
					func() { maps.InvertMustHaveDistinctValues(tc.inputMap) })
			} else {
				inverted := maps.InvertMustHaveDistinctValues(tc.inputMap)
				require.Equal(t, tc.expectedResult, inverted)
			}
		})
	}
}
