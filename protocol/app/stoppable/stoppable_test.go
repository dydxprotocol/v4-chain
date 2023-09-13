package stoppable_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/app/stoppable"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestStoppable struct {
	stopCalled bool
}

func (m *TestStoppable) Stop() {
	if m.stopCalled {
		panic("Stop called twice")
	}
	m.stopCalled = true
}

func TestStopServices(t *testing.T) {
	mockStoppable := &TestStoppable{}
	mockStoppable2 := &TestStoppable{}
	mockStoppableSeparateTest := &TestStoppable{}
	stoppable.RegisterServiceForTestCleanup("test", mockStoppable)
	stoppable.RegisterServiceForTestCleanup("test", mockStoppable2)
	stoppable.RegisterServiceForTestCleanup("test2", mockStoppableSeparateTest)

	// Stop test services, verify.
	stoppable.StopServices(t, "test")

	// Verify test services stopped, test2 services unaffected.
	require.True(t, mockStoppable.stopCalled)
	require.True(t, mockStoppable2.stopCalled)
	require.False(t, mockStoppableSeparateTest.stopCalled)

	// Stop test services again. This should not cause any panics.
	stoppable.StopServices(t, "test")

	// Stop test2 services, verify.
	stoppable.StopServices(t, "test2")
	require.True(t, mockStoppableSeparateTest.stopCalled)
}
