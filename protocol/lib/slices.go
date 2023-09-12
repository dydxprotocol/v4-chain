package lib

import (
	"fmt"
	"sort"
)

// ContainsDuplicates returns true if the slice contains duplicates, false if not.
func ContainsDuplicates[V comparable](values []V) bool {
	seenValues := make(map[V]bool)
	for _, val := range values {
		if _, exists := seenValues[val]; exists {
			return true
		}

		seenValues[val] = true
	}

	return false
}

// GetSortedKeys returns the keys of the map in sorted order.
func GetSortedKeys[R interface {
	~[]K
	sort.Interface
}, K comparable, V any](m map[K]V) R {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(R(keys))
	return keys
}

// ContainsValue returns true if the slice contains the provided value, false if not.
func ContainsValue[V comparable](values []V, value V) bool {
	for _, sliceVal := range values {
		if sliceVal == value {
			return true
		}
	}

	return false
}

// SliceToSet converts a slice to a set. Function will panic if there are duplicate values.
func SliceToSet[K comparable](values []K) map[K]struct{} {
	set := make(map[K]struct{}, len(values))
	for _, sliceVal := range values {
		if _, exists := set[sliceVal]; exists {
			panic(
				fmt.Sprintf(
					"SliceToSet: duplicate value: %+v",
					sliceVal,
				),
			)
		}
		set[sliceVal] = struct{}{}
	}
	return set
}

// MustRemoveIndex returns a copy of the provided slice with the value at `index` removed. This function
// will not change the ordering of other elements within the original slice.
// Note that function will panic if `index >= len(values)`.
func MustRemoveIndex[V any](values []V, index uint) []V {
	numValues := uint(len(values))
	if numValues <= index {
		panic(
			fmt.Sprintf(
				"MustRemoveIndex: index %d is greater than array length %d",
				index,
				numValues,
			),
		)
	}
	ret := make([]V, 0, numValues-1)
	ret = append(ret, values[:index]...)
	return append(ret, values[index+1:]...)
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

// MustGetValue returns the element at `index` position in a slice and panics if `index` is greater than or
// equal to slice length.
func MustGetValue[V any](values []V, index uint) V {
	if index >= uint(len(values)) {
		panic(
			fmt.Sprintf(
				"MustGetValue: index %d is greater than or equal to array length %d",
				index,
				len(values),
			),
		)
	}
	return values[index]
}
