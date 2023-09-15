package time

import (
	"time"
)

// MustParseDuration turns a string into a duration.
// Panics if the string cannot be parsed.
func MustParseDuration(s string) time.Duration {
	v, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return v
}
