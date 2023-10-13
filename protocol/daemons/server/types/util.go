package types

import "time"

// MaximumAcceptableUpdateDelay computes the maximum acceptable update delay for a daemon service as a
// multiple of the loop delay.
func MaximumAcceptableUpdateDelay(loopDelayMs uint32) time.Duration {
	return MaximumLoopDelayMultiple * time.Duration(loopDelayMs) * time.Millisecond
}
