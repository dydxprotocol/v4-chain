package maps

import (
	"fmt"
)

// MergeAllMapsMustHaveDistinctKeys merges all the maps into a single map.
// Panics if there are duplicate keys.
func MergeAllMapsMustHaveDistinctKeys[K comparable, V any](maps ...map[K]V) map[K]V {
	combinedMap := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			if _, exists := combinedMap[k]; exists {
				panic(fmt.Sprintf("duplicate key: %v", k))
			}
			combinedMap[k] = v
		}
	}
	return combinedMap
}

// ShallowCopy returns a shallow copy of originalMap.
func ShallowCopy[K comparable, V any](originalMap map[K]V) map[K]V {
	if originalMap == nil {
		return nil
	}
	newMap := make(map[K]V)
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}
