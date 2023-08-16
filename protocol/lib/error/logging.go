package error

import (
	"errors"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/logging"
)

// LogErrorWithOptionalContext logs an error, optionally adding context to the logger iff the error implements
// the LogContextualizer interface.
func LogErrorWithOptionalContext(
	logger log.Logger,
	msg string,
	err error,
) {
	var logContextualizer LogContextualizer
	if ok := errors.As(err, &logContextualizer); ok {
		logger = logContextualizer.AddLoggingContext(logger)
		// Log the original error.
		err = logContextualizer.Unwrap()
	}

	logger.Error(msg, "error", err)
}

// WrapErrorWithPricesSourceModuleContext wraps an error with a LogContextualizer that the spercified source module.
// This is useful for logging the error within the process proposal handler (or any other location that uses
// LogErrorWithOptionalContext) with metadata that can be used to identify the source of the error.
func WrapErrorWithSourceModuleContext(err error, module string) error {
	return NewErrorWithLogContext(err).
		WithLogKeyValue(logging.SourceModuleKey, fmt.Sprintf("x/%v", module))
}
