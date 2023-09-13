package lib

// MedianizerImpl is the struct that implements the `Medianizer` interface.
type MedianizerImpl struct{}

// Ensure the `MedianizerImpl` struct is implemented at compile time.
var _ Medianizer = (*MedianizerImpl)(nil)

// Medianizer is an interface that encapsulates the lib.math function `MedianUint64`.
type Medianizer interface {
	MedianUint64(input []uint64) (uint64, error)
}

// MedianUint64 wraps `lib.Median` which gets the median of a uint64 slice.
func (r *MedianizerImpl) MedianUint64(input []uint64) (uint64, error) {
	return Median(input)
}
