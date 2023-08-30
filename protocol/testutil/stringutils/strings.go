package stringutils

import (
	"strings"
	"unicode"
)

// StripSpaces removes all whitespaces from a slice of bytes, including actual spaces, tabs,
// newlines, etc. An example usage would be to load a formatted, multiline json file exposed
// as human-readable test data and convert it into a compact form that can easily be compared with
// the stringified bytes of an unmarshalled json data structure for testing.
func StripSpaces(bytes []byte) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, string(bytes))
}
