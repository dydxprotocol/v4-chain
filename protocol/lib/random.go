package lib

import (
	"errors"
	"fmt"
	"math/rand"
)

// RandomBytesBetween returns a random byte slice that is in the range [start, end] when compared lexicographically.
// The slice will have a length in the range [len(start), len(end)].
// In the current implementation, all possible permutations are not equally likely.
// Nil slices for start and end will be treated as empty byte slices. Will panic if:
//   - start compares lexicographically greater than end
//   - nil rand is provided
func RandomBytesBetween(start []byte, end []byte, rand *rand.Rand) []byte {
	if rand == nil {
		panic(errors.New("rand expected to be non-nil."))
	}

	minLen := len(start)
	maxLen := len(end)
	if minLen > maxLen {
		minLen, maxLen = maxLen, minLen
	}

	bytes := make([]byte, maxLen)
	i := 0

	// Copy the common bytes between the two keys.
	for ; i < minLen; i++ {
		// Lexographically compare the byte.
		// If equal, copy the byte.
		// If not equal, then either panic or stop copying (depending on which byte is greater).
		if start[i] == end[i] {
			bytes[i] = start[i]
		} else if start[i] > end[i] {
			panic(fmt.Errorf("start %x compares lexicographically greater than end %x at position %d.", start, end, i))
		} else {
			break
		}
	}

	// If start == end then we are done and can return bytes.
	if i == maxLen {
		return bytes
	}

	// Remember the floor and ceiling starting at the first byte that differs between the two keys.
	// Note that if floor is -1, then len(start) <= len(bytes)
	isPrefixOfStart, isPrefixOfEnd := true, true
	floor := int32(0)
	if i < len(start) {
		floor = int32(start[i])
	}
	ceiling := int32(end[i])

	// Compute a random byte length that gives each byte string an equal probability.
	// Note that [0, 255] represents the possible values and 256 represents the "unset" byte.
	targetLength := maxLen
	for j := minLen; j < maxLen && rand.Int31n(257) == 256; j++ {
		targetLength--
	}

	// Generate the remainder of the random bytes producing a value that compares lexicographically
	// between start and end.
	for ; i < targetLength; i++ {
		current := floor + rand.Int31n(ceiling-floor+1)
		bytes[i] = byte(current)

		// Ensure that if bytes is a prefix of start that the next byte in start is the new floor.
		if isPrefixOfStart && current == floor && i+1 < len(start) {
			floor = int32(start[i+1])
		} else {
			floor = 0
			isPrefixOfStart = false
		}

		// Ensure that if bytes is a prefix of end that the next byte in end is the new ceiling.
		if isPrefixOfEnd && current == ceiling {
			// If bytes == end then we must return now as we can't generate any more bytes as
			// the result would be greater than end.
			if i+1 < len(end) {
				ceiling = int32(end[i+1])
			} else {
				return bytes[:i+1]
			}
		} else {
			ceiling = 255
			isPrefixOfEnd = false
		}
	}

	return bytes[:targetLength]
}
