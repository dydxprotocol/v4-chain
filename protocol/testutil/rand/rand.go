package rand

import (
	"math/rand"
	"time"
)

// NewRand returns a new Rand that generates random values, seeded from the current
// unix time in nanoseconds.
func NewRand() *rand.Rand {
	s := rand.NewSource(time.Now().UnixNano())
	return rand.New(s)
}
