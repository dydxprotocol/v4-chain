package types

import (
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

// The following tests may still intermittently fail on an overloaded system as they rely
// on `time.Sleep`, which is not guaranteed to wake up after the specified amount of time.
func TestRegisterDaemonService_Success(t *testing.T) {
	ufm := NewUpdateFrequencyMonitor()
	success := ufm.RegisterDaemonService("test-service", 200*time.Millisecond)
	require.True(t, success)

	// As long as responses come in before the 200ms deadline, no panic should occur.
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)

	ufm.Stop()
}

func TestRegisterDaemonServiceWithCallback_Success(t *testing.T) {
	callbackCalled := atomic.Bool{}

	ufm := NewUpdateFrequencyMonitor()
	success := ufm.RegisterDaemonServiceWithCallback("test-service", 200*time.Millisecond, func() {
		callbackCalled.Store(true)
	})
	require.True(t, success)

	// As long as responses come in before the 200ms deadline, no panic should occur.
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)

	require.False(t, callbackCalled.Load())

	ufm.Stop()
}

func TestRegisterDaemonService_DoubleRegistrationFails(t *testing.T) {
	ufm := NewUpdateFrequencyMonitor()
	success := ufm.RegisterDaemonService("test-service", 200*time.Millisecond)
	require.True(t, success)

	// Register the same daemon service again. This should fail, and 50ms update frequency should be ignored.
	success = ufm.RegisterDaemonService("test-service", 50*time.Millisecond)
	require.False(t, success)

	// Confirm that the original 200ms update frequency is still in effect. 50ms would have triggered a panic.
	// Note there is a possibility that 200ms will still cause a panic due to the semantics of Sleep, which is
	// not guaranteed to sleep for exactly the specified duration.
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	ufm.Stop()
}

func TestRegisterDaemonServiceWithCallback_DoubleRegistrationFails(t *testing.T) {
	// lock synchronizes callback flags and was added to avoid race test failures.
	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}

	ufm := NewUpdateFrequencyMonitor()
	// First registration should succeed.
	success := ufm.RegisterDaemonServiceWithCallback("test-service", 200*time.Millisecond, func() {
		callback1Called.Store(true)
	})
	require.True(t, success)

	// Register the same daemon service again. This should fail, and 50ms update frequency should be ignored.
	success = ufm.RegisterDaemonServiceWithCallback("test-service", 50*time.Millisecond, func() {
		callback2Called.Store(true)
	})
	require.False(t, success)

	// Validate that the original callback is still in effect for the original 200ms update frequency.
	// The 50ms update frequency should have invoked a callback if it were applied.
	time.Sleep(80 * time.Millisecond)
	require.False(t, callback1Called.Load())
	require.False(t, callback2Called.Load())

	// Validate no issues with RegisterValidResponse after a double registration was attempted.
	require.NoError(t, ufm.RegisterValidResponse("test-service"))

	// Sleep until the callback should be called.
	time.Sleep(250 * time.Millisecond)
	require.True(t, callback1Called.Load())
	require.False(t, callback2Called.Load())

	ufm.Stop()
}

func TestRegisterDaemonService_RegistrationFailsAfterStop(t *testing.T) {
	ufm := NewUpdateFrequencyMonitor()
	ufm.Stop()
	success := ufm.RegisterDaemonService("test-service", 50*time.Millisecond)
	require.False(t, success)

	// Any accidentally scheduled functions with panics should fire before this timer expires.
	time.Sleep(100 * time.Millisecond)
}

func TestRegisterDaemonServiceWithCallback_RegistrationFailsAfterStop(t *testing.T) {
	ufm := NewUpdateFrequencyMonitor()
	ufm.Stop()

	callbackCalled := atomic.Bool{}

	// Registering a daemon service with a callback should fail after the monitor has been stopped.
	success := ufm.RegisterDaemonServiceWithCallback("test-service", 50*time.Millisecond, func() {
		callbackCalled.Store(true)
	})
	require.False(t, success)

	// Wait until after the callback duration has expired. The callback should not be called.
	time.Sleep(75 * time.Millisecond)

	// Validate that the callback was not called.
	require.False(t, callbackCalled.Load())
}

func TestPanicServiceNotResponding(t *testing.T) {
	panicFunc := PanicServiceNotResponding("test-service")
	require.Panics(t, panicFunc)
}

// This test is disabled because the panic is not recoverable, since it's thrown
// in a separate goroutine.
//func TestUpdateFrequencyMonitor_Panics(t *testing.T) {
//	// Expect the following sequence to panic
//	ufm := NewUpdateFrequencyMonitor()
//	ufm.RegisterDaemonService("test-service", 100*time.Millisecond)
//	time.Sleep(150 * time.Millisecond)
//}
