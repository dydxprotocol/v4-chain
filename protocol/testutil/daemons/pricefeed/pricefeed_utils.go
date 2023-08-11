package pricefeed

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// ErrorMapsEqual is a testing method that takes any two maps of keys to errors and asserts that they have the same
// sets of keys, and that each associated error value has the same rendered message.
func ErrorMapsEqual[K comparable](t *testing.T, expected map[K]error, actual map[K]error) {
	require.Equal(t, len(expected), len(actual))
	for key, expectedError := range expected {
		error, ok := actual[key]
		require.True(t, ok)
		require.EqualError(t, expectedError, error.Error())
	}
}

// ErrorsEqual is a testing method that takes any two slices of errors and asserts that each actual error has
// the same rendered message as the expected error.
func ErrorsEqual(t *testing.T, expected []error, actual []error) {
	require.Equal(t, len(expected), len(actual))
	for i, expectedError := range expected {
		require.EqualError(t, expectedError, actual[i].Error())
	}
}
