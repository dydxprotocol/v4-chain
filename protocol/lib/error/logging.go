package error

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	SourceModuleKey = "source_module"
)

// LogErrorWithOptionalContext logs an error, optionally adding context to the logger iff the error implements
// the LogContextualizer interface. This method is appropriate for logging errors that may or may not be wrapped
// in an ErrorWithLogContext.
func LogErrorWithOptionalContext(
	ctx sdk.Context,
	msg string,
	err error,
) {
	logger := ctx.Logger()
	var logContextualizer LogContextualizer
	if ok := errors.As(err, &logContextualizer); ok {
		logger = logContextualizer.AddLoggingContext(logger)
		// Log the original error.
		err = logContextualizer.Unwrap()
	}

	logger.Error(msg, "error", err)
}

// WrapErrorWithPricesSourceModuleContext wraps an error with a LogContextualizer that contains a key-value pair for
// the specified source module. This can be used, for example, for wrapping validation failure errors that are logged
// within the process proposal handler (or from any other location that uses LogErrorWithOptionalContext) with metadata
// that helps to more easily identify the source of the error and correlate logs.
func WrapErrorWithSourceModuleContext(err error, module string) error {
	return NewErrorWithLogContext(err).
		WithLogKeyValue(SourceModuleKey, fmt.Sprintf("x/%v", module))
}
