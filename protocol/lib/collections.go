package lib

import (
	"fmt"
	"sort"
)

// ContainsDuplicates returns true if the slice contains duplicates, false if not.
func ContainsDuplicates[V comparable](values []V) bool {
	// Optimize edge case. If values is nil, len(values) returns 0.
	if len(values) <= 1 {
		return false
	}

	// Store each value as a key in the mapping.
	seenValues := make(map[V]struct{}, len(values))
	for i, val := range values {
		// Add the value to the mapping.
		seenValues[val] = struct{}{}

		// Return early if the size of the mapping did not grow.
		if len(seenValues) <= i {
			return true
		}
	}

	return false
}

// MapToSortedSlice returns a slice of values from a map, sorted by key.
func MapToSortedSlice[R interface {
	~[]K
	sort.Interface
}, K comparable, V any](m map[K]V) []V {
	keys := GetSortedKeys[R](m)
	values := make([]V, 0, len(m))
	for _, key := range keys {
		values = append(values, m[key])
	}
	return values
}

// DedupeSlice deduplicates a slice of comparable values.
func DedupeSlice[V comparable](values []V) []V {
	seenValues := make(map[V]struct{})
	deduped := make([]V, 0)
	for _, val := range values {
		if _, seen := seenValues[val]; !seen {
			deduped = append(deduped, val)
		}
		seenValues[val] = struct{}{}
	}
	return deduped
}

// GetSortedKeys returns the keys of the map in sorted order.
func GetSortedKeys[R interface {
	~[]K
	sort.Interface
}, K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(R(keys))
	return keys
}

// UniqueSliceToSet converts a slice of unique values to a set.
// The function will panic if there are duplicate values.
func UniqueSliceToSet[K comparable](values []K) map[K]struct{} {
	set := make(map[K]struct{}, len(values))
	for _, sliceVal := range values {
		if _, exists := set[sliceVal]; exists {
			panic(
				fmt.Sprintf(
					"UniqueSliceToSet: duplicate value: %+v",
					sliceVal,
				),
			)
		}
		set[sliceVal] = struct{}{}
	}
	return set
}

// UniqueSliceToMap converts a slice to a map using the provided keyFunc to generate the key.
func UniqueSliceToMap[K comparable, V any](slice []V, keyFunc func(V) K) map[K]V {
	m := make(map[K]V)
	for _, v := range slice {
		k := keyFunc(v)
		if _, exists := m[k]; exists {
			panic(
				fmt.Sprintf(
					"UniqueSliceToMap: duplicate value: %+v",
					v,
				),
			)
		}
		m[k] = v
	}
	return m
}

// MapSlice takes a function and executes that function on each element of a slice, returning the result.
// Note the function must return one result for each element of the slice.
func MapSlice[V any, E any](values []V, mapFunc func(V) E) []E {
	mappedValues := make([]E, 0, len(values))
	for _, value := range values {
		mappedValues = append(mappedValues, mapFunc(value))
	}

	return mappedValues
}

// FilterSlice takes a function that returns a boolean on whether to include the element in the final
// result, and returns a slice of elements where the function returned true when called with each element.
func FilterSlice[V any](values []V, filterFunc func(V) bool) []V {
	filteredValues := make([]V, 0, len(values))
	for _, value := range values {
		if filterFunc(value) {
			filteredValues = append(filteredValues, value)
		}
	}

	return filteredValues
}

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

// MergeMaps merges all the maps into a single map.
// Does not require maps to have distinct keys.
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	combinedMap := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			combinedMap[k] = v
		}
	}
	return combinedMap
}

// Check if slice contains a particular value
func SliceContains[T comparable](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
