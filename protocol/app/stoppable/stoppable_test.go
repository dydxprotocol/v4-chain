package stoppable

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type MockStoppable struct {
	StopCalled bool
}

func (m *MockStoppable) Stop() {
	m.StopCalled = true
}

func TestStopServices(t *testing.T) {
	mockStoppable := &MockStoppable{}
	mockStoppable2 := &MockStoppable{}
	mockStoppableSeparateTest := &MockStoppable{}
	RegisterServiceForTestCleanup("test", mockStoppable)
	RegisterServiceForTestCleanup("test", mockStoppable2)
	RegisterServiceForTestCleanup("test2", mockStoppableSeparateTest)

	// Stop test services, verify.
	StopServices(t, "test")

	// Verify test services stopped, test2 services unaffected.
	require.True(t, mockStoppable.StopCalled)
	require.True(t, mockStoppable2.StopCalled)
	require.False(t, mockStoppableSeparateTest.StopCalled)

	// Stop test2 services, verify.
	StopServices(t, "test2")
	require.True(t, mockStoppableSeparateTest.StopCalled)
}
