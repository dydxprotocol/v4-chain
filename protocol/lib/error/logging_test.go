package error_test

import (
	"fmt"
	"testing"

	liberror "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
)

func TestWrapErrorWithSourceModuleContext_ErrorWithLogContext(t *testing.T) {
	err := fmt.Errorf("test error")
	wrappedErr := liberror.WrapErrorWithSourceModuleContext(err, "test-module")
	logger := &mocks.Logger{}
	ctx, _, _ := sdk.NewSdkContextWithMultistore()
	ctx = ctx.WithLogger(logger)

	// Expect that source module context will be added to the logger,
	// and then the original error will be logged.
	call := logger.On("With", liberror.SourceModuleKey, "x/test-module").Return(logger)
	logger.On("Error", "test message", "error", err).Return().NotBefore(call)

	liberror.LogErrorWithOptionalContext(ctx, "test message", wrappedErr)

	logger.AssertExpectations(t)
}

func TestLogErrorWithOptionalContext_PlainError(t *testing.T) {
	logger := &mocks.Logger{}
	err := fmt.Errorf("test error")
	ctx, _, _ := sdk.NewSdkContextWithMultistore()
	ctx = ctx.WithLogger(logger)

	// Plain error messages will be logged without any additional context.
	logger.On("Error", "test message", "error", err).Return()

	liberror.LogErrorWithOptionalContext(ctx, "test message", err)

	logger.AssertExpectations(t)
}
