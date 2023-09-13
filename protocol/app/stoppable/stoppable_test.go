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
	testStoppable := &TestStoppable{}
	testStoppable2 := &TestStoppable{}
	testStoppableSeparateTest := &TestStoppable{}
	stoppable.RegisterServiceForTestCleanup("test", testStoppable)
	stoppable.RegisterServiceForTestCleanup("test", testStoppable2)
	stoppable.RegisterServiceForTestCleanup("test2", testStoppableSeparateTest)

	// Stop test services, verify.
	stoppable.StopServices(t, "test")

	// Verify test services stopped, test2 services unaffected.
	require.True(t, testStoppable.stopCalled)
	require.True(t, testStoppable2.stopCalled)
	require.False(t, testStoppableSeparateTest.stopCalled)

	// Stop test services again. This should not cause any panics.
	stoppable.StopServices(t, "test")

	// Stop test2 services, verify.
	stoppable.StopServices(t, "test2")
	require.True(t, testStoppableSeparateTest.stopCalled)
}
