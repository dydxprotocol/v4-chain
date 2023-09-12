package lib

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
)

// RandomBytesBetween returns a random byte slice that is in the range [lo, hi] when compared lexicographically.
// The slice will have a length at most max(len(lo), len(hi)).
// Nil slices for lo and hi will be treated as empty byte slices. Will panic if:
//   - lo compares lexicographically greater than hi
//   - nil rand is provided
func RandomBytesBetween(lo []byte, hi []byte, rand *rand.Rand) []byte {
	if rand == nil {
		panic(errors.New("rand expected to be non-nil."))
	}

	if bytes.Compare(lo, hi) > 0 {
		panic(fmt.Errorf("lo %x compares lexicographically greater than hi %x", lo, hi))
	}

	// Determine the maximum length.
	maxLen := Max(len(lo), len(hi))

	// Allocate the bytes.
	bytes := make([]byte, maxLen)

	// Track if written bytes is a prefix of lo or hi.
	isLoPrefix, isHiPrefix := true, true

	for i := 0; i < maxLen; i++ {
		// Get the minimum and maximum values.
		a, b := byte(0), byte(255)
		if isLoPrefix && i < len(lo) {
			a = lo[i]
		}
		if isHiPrefix && i < len(hi) {
			b = hi[i]
		}

		// If we are in the process of copying a common prefix, then continue to copy.
		if isLoPrefix && isHiPrefix && a == b {
			bytes[i] = a
			continue
		}

		// Number of possibilities.
		numPossibilities := int32(b) - int32(a) + 1

		// If we are not a prefix of lo, then we may return early.
		// The probability of returning early is equal to the probability of any unique byte.
		if !isLoPrefix && rand.Int31n(numPossibilities+1) == 0 {
			return bytes[:i]
		}

		// Determine the random byte
		cur := byte(int32(a) + rand.Int31n(numPossibilities))
		bytes[i] = cur

		// Check if we need to set either prefix variable to false.
		if isLoPrefix && i < len(lo) && cur == lo[i] {
			isLoPrefix = false
		}
		if isHiPrefix && i < len(hi) && cur == hi[i] {
			isHiPrefix = false
		}
	}

	return bytes
}
