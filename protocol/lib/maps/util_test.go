package maps_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib/maps"
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

func TestShallowCopy(t *testing.T) {
	tests := map[string]struct {
		inputMap       map[string]int
		expectedOutput map[string]int
	}{
		"Nil map": {
			inputMap:       nil,
			expectedOutput: nil,
		},
		"Empty map": {
			inputMap:       map[string]int{},
			expectedOutput: map[string]int{},
		},
		"Single-element map": {
			inputMap: map[string]int{
				"a": 1,
			},
			expectedOutput: map[string]int{
				"a": 1,
			},
		},
		"Multi-element map": {
			inputMap: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			expectedOutput: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			copiedMap := maps.ShallowCopy(tc.inputMap)
			require.Equal(t, tc.expectedOutput, copiedMap)

			// If the original map is modified, it should not affect the copied map
			if tc.inputMap != nil {
				tc.inputMap["new_key"] = 100
				require.NotEqual(t, tc.inputMap, copiedMap)
			}
		})
	}
}
