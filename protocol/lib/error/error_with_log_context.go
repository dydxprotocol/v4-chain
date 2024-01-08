package error

import "cosmossdk.io/log"

// ErrorWithLogContext wraps an error with log context and implements the LogContextualizer interface.
// ErrorWithLogContext unwraps to the original error.
type ErrorWithLogContext struct {
	err error
	LogContextualizer
	// Store key-values as a slice of ordered pairs in order to maintain the flexibility of the logging library,
	// which allows for keys of type interface{}.
	keyValues []interface{}
}

// NewErrorWithLogContext returns a new ErrorWithLogContext wrapping the given error.
func NewErrorWithLogContext(err error) *ErrorWithLogContext {
	errWithLogContext := &ErrorWithLogContext{
		keyValues: []interface{}{},
		err:       err,
	}
	return errWithLogContext
}

// Unwrap returns the underlying error.
func (ewlc *ErrorWithLogContext) Unwrap() error {
	return ewlc.err
}

// WithLogKeyValue adds a key-value pair to the error's log context. The returned ErrorWithLogContext is the same
// error.
func (ewlc *ErrorWithLogContext) WithLogKeyValue(key interface{}, value interface{}) *ErrorWithLogContext {
	ewlc.keyValues = append(ewlc.keyValues, key, value)
	return ewlc
}

// AddLoggingContext returns a modified logger with the error's logging context added.
func (ewlc *ErrorWithLogContext) AddLoggingContext(logger log.Logger) log.Logger {
	return logger.With(ewlc.keyValues...)
}

// Error returns the underlying error's Error() string.
func (ewlc *ErrorWithLogContext) Error() string {
	return ewlc.err.Error()
}
