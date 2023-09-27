package error_test

import (
	"fmt"
	liberror "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"testing"
)

func TestWrapErrorWithSourceModuleContext_ErrorWithLogContext(t *testing.T) {
	err := fmt.Errorf("test error")
	wrappedErr := liberror.WrapErrorWithSourceModuleContext(err, "test-module")
	logger := &mocks.Logger{}

	// Expect that source module context will be added to the logger,
	// and then the original error will be logged.
	call := logger.On("With", liberror.SourceModuleKey, "x/test-module").Return(logger)
	logger.On("Error", "test message", "error", err).Return().NotBefore(call)

	liberror.LogErrorWithOptionalContext(logger, "test message", wrappedErr)

	logger.AssertExpectations(t)
}

func TestLogErrorWithOptionalContext_PlainError(t *testing.T) {
	logger := &mocks.Logger{}
	err := fmt.Errorf("test error")

	// Plain error messages will be logged without any additional context.
	logger.On("Error", "test message", "error", err).Return()

	liberror.LogErrorWithOptionalContext(logger, "test message", err)

	logger.AssertExpectations(t)
}
