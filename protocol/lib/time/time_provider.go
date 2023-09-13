package time

import (
	"time"
)

// Ensure the `TimeProviderImpl` is implemented at compile time.
var _ TimeProvider = &TimeProviderImpl{}

// TimeProvider is an interface that provides time.
type TimeProvider interface {
	Now() time.Time
}

// TimeProviderImpl implements TimeProvider interface.
type TimeProviderImpl struct{}

// Now returns current time.
func (t *TimeProviderImpl) Now() time.Time {
	return time.Now()
}
