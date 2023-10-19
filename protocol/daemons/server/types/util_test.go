package types

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMaximumAcceptableUpdateDelay(t *testing.T) {
	loopDelayMs := uint32(1000)
	expected := time.Duration(MaximumLoopDelayMultiple*loopDelayMs) * time.Millisecond
	actual := MaximumAcceptableUpdateDelay(loopDelayMs)
	require.Equal(t, expected, actual)
}
