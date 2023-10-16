package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

var (
	zeroDuration = 0 * time.Second
)

func createTestMonitor() (*types.UpdateMonitor, *mocks.Logger) {
	logger := &mocks.Logger{}
	return types.NewUpdateFrequencyMonitor(zeroDuration, logger), logger
}

// The following tests may still intermittently fail on an overloaded system as they rely
// on `time.Sleep`, which is not guaranteed to wake up after the specified amount of time.
func TestRegisterDaemonService_Success(t *testing.T) {
	ufm, logger := createTestMonitor()
	err := ufm.RegisterDaemonService("test-service", 200*time.Millisecond)
	require.NoError(t, err)

	// As long as responses come in before the 200ms deadline, no errors should be logged.
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)

	ufm.Stop()
	// Assert: no calls to the logger were made.
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterDaemonService_SuccessfullyLogsError(t *testing.T) {
	ufm, logger := createTestMonitor()
	logger.On("Error", "daemon not responding", "service", "test-service").Once().Return()
	err := ufm.RegisterDaemonService("test-service", 1*time.Millisecond)
	require.NoError(t, err)
	time.Sleep(2 * time.Millisecond)
	ufm.Stop()

	// Assert: the logger was called with the expected arguments.
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterDaemonServiceWithCallback_Success(t *testing.T) {
	callbackCalled := atomic.Bool{}

	ufm, _ := createTestMonitor()
	err := ufm.RegisterDaemonServiceWithCallback("test-service", 200*time.Millisecond, func() {
		callbackCalled.Store(true)
	})
	require.NoError(t, err)

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
	ufm, logger := createTestMonitor()

	err := ufm.RegisterDaemonService("test-service", 200*time.Millisecond)
	require.NoError(t, err)

	// Register the same daemon service again. This should fail, and 50ms update frequency should be ignored.
	err = ufm.RegisterDaemonService("test-service", 50*time.Millisecond)
	require.ErrorContains(t, err, "service already registered")

	// Confirm that the original 200ms update frequency is still in effect. 50ms would have triggered an error log.
	// Note there is a possibility that 200ms will still cause an error log due to the semantics of Sleep, which is
	// not guaranteed to sleep for exactly the specified duration.
	time.Sleep(80 * time.Millisecond)
	require.NoError(t, ufm.RegisterValidResponse("test-service"))
	time.Sleep(80 * time.Millisecond)
	ufm.Stop()

	// Assert no calls to the logger were made.
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterDaemonServiceWithCallback_DoubleRegistrationFails(t *testing.T) {
	// lock synchronizes callback flags and was added to avoid race test failures.
	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}

	ufm, _ := createTestMonitor()
	// First registration should succeed.
	err := ufm.RegisterDaemonServiceWithCallback("test-service", 200*time.Millisecond, func() {
		callback1Called.Store(true)
	})
	require.NoError(t, err)

	// Register the same daemon service again. This should fail, and 50ms update frequency should be ignored.
	err = ufm.RegisterDaemonServiceWithCallback("test-service", 50*time.Millisecond, func() {
		callback2Called.Store(true)
	})
	require.ErrorContains(t, err, "service already registered")

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
	ufm, logger := createTestMonitor()
	ufm.Stop()
	err := ufm.RegisterDaemonService("test-service", 50*time.Millisecond)
	require.ErrorContains(t, err, "monitor has been stopped")

	// Any scheduled functions with error logs that were not cleaned up should trigger before this sleep finishes.
	time.Sleep(100 * time.Millisecond)
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterDaemonServiceWithCallback_RegistrationFailsAfterStop(t *testing.T) {
	ufm, _ := createTestMonitor()
	ufm.Stop()

	callbackCalled := atomic.Bool{}

	// Registering a daemon service with a callback should fail after the monitor has been stopped.
	err := ufm.RegisterDaemonServiceWithCallback("test-service", 50*time.Millisecond, func() {
		callbackCalled.Store(true)
	})
	require.ErrorContains(t, err, "monitor has been stopped")

	// Wait until after the callback duration has expired. The callback should not be called.
	time.Sleep(75 * time.Millisecond)

	// Validate that the callback was not called.
	require.False(t, callbackCalled.Load())
}

func TestRegisterValidResponse_NegativeUpdateDelay(t *testing.T) {
	ufm, logger := createTestMonitor()
	err := ufm.RegisterDaemonService("test-service", -50*time.Millisecond)
	require.ErrorContains(t, err, "update delay -50ms must be positive")

	// Sanity check: no calls to the logger should have been made.
	mock.AssertExpectationsForObjects(t, logger)
}

func TestRegisterValidResponseWithCallback_NegativeUpdateDelay(t *testing.T) {
	ufm, _ := createTestMonitor()
	err := ufm.RegisterDaemonServiceWithCallback("test-service", -50*time.Millisecond, func() {})
	require.ErrorContains(t, err, "update delay -50ms must be positive")
}

func TestPanicServiceNotResponding(t *testing.T) {
	panicFunc := types.PanicServiceNotResponding("test-service")
	require.Panics(t, panicFunc)
}

func TestLogErrorServiceNotResponding(t *testing.T) {
	logger := &mocks.Logger{}
	logger.On("Error", "daemon not responding", "service", "test-service").Return()
	logFunc := types.LogErrorServiceNotResponding("test-service", logger)
	logFunc()

	// Assert: the logger was called with the expected arguments.
	mock.AssertExpectationsForObjects(t, logger)
}
