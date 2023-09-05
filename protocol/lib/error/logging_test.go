package error

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib/logging"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockLogger struct {
	loggingContext map[string]interface{}

	// addedContext tracks if context was added before logging.
	addedContext bool
	// errorMsg and keyVals store the parameters of the last log.Error call.
	errorMsg string
	keyVals  []interface{}
}

func (ml *MockLogger) Debug(msg string, keyvals ...interface{}) {
	if !ml.addedContext {
		panic("attempted to log without adding context")
	}
}

func (ml *MockLogger) Info(msg string, keyvals ...interface{}) {
	if !ml.addedContext {
		panic("attempted to log without adding context")
	}
}

func (ml *MockLogger) Error(msg string, keyvals ...interface{}) {
	if !ml.addedContext {
		panic("attempted to log without adding context")
	}
	ml.errorMsg = msg
	ml.keyVals = keyvals
}

func (ml *MockLogger) With(keyvals ...interface{}) log.Logger {
	for i := 0; i < len(keyvals); i += 2 {
		ml.loggingContext[keyvals[i].(string)] = keyvals[i+1]
	}
	ml.addedContext = true
	return ml
}

func TestWrapErrorWithSourceModuleContext(t *testing.T) {
	err := fmt.Errorf("test error")
	wrappedErr := WrapErrorWithSourceModuleContext(err, "test-module")
	logger := &MockLogger{loggingContext: map[string]interface{}{}}

	LogErrorWithOptionalContext(logger, "test message", wrappedErr)

	// Assert that logging context was added to the logger.
	require.Len(t, logger.loggingContext, 1)
	source_module, ok := logger.loggingContext[logging.SourceModuleKey]
	require.True(t, ok)
	require.Equal(t, "x/test-module", source_module)

	// Assert that the error was logged with expected keyVals.
	require.Equal(t, "test message", logger.errorMsg)
	require.Len(t, logger.keyVals, 2)
	require.Equal(t, "error", logger.keyVals[0])
	require.Equal(t, err, logger.keyVals[1])
}
