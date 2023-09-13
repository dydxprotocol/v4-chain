package types

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateFrequencyMonitor_Success(t *testing.T) {
	ufm := NewUpdateFrequencyMonitor()
	ufm.RegisterDaemonService("test-service", 200*time.Millisecond)
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
//func TestUpdateFrequencyMonitor_Panics(t *testing.T) {
//	// Expect the following sequence to panic
//	ufm := NewUpdateFrequencyMonitor()
//	ufm.RegisterDaemonService("test-service", 100*time.Millisecond)
//	time.Sleep(180 * time.Millisecond)
//}
