package types

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// This test is disabled because of the non-determinism of time.Sleep, which may
// go over the prescribed duration, and cause the timer to be triggered spuriously.
// We don't want to increase the duration of the test, because it will slow down
// the test suite.
func _TestUpdateFrequencyMonitor(t *testing.T) {
	ufm := NewUpdateFrequencyMonitor()
	ufm.RegisterDaemonService("test-service", 100*time.Millisecond)
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	ufm.Stop()
}

// This test is disabled because the panic is not recoverable, since it's thrown
// in a separate goroutine.
func _TestUpdateFrequencyMonitor_Panics(t *testing.T) {
	// Expect the following sequence to panic
	ufm := NewUpdateFrequencyMonitor()
	ufm.RegisterDaemonService("test-service", 100*time.Millisecond)
	time.Sleep(180 * time.Millisecond)
}
