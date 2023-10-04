package error_test

import (
	"fmt"
	"testing"

	liberror "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
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

func TestLogErrorWithBlockHeight(t *testing.T) {
	logger := &mocks.Logger{}
	err := fmt.Errorf("test error")

	// Expect that the block height will be appended to the error message.
	logger.On("Error", "Block height: 123, Callback: foobar: test error").Return()

	liberror.LogErrorWithBlockHeight(logger, err, 123, "foobar")

	logger.AssertExpectations(t)
}

func TestLogErrorWithBlockHeight_NilError(t *testing.T) {
	logger := &mocks.Logger{}

	// Expect that the block height will be appended to the error message.
	logger.On("Error", "LogErrorWithBlockHeight called with nil error").Return()

	liberror.LogErrorWithBlockHeight(logger, nil, 123, "foobar")

	logger.AssertExpectations(t)
}
