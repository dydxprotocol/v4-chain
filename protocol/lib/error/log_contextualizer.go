package error

import "github.com/cometbft/cometbft/libs/log"

// LogContextualizer describes an object that can add context - that is, descriptive key-value pairs, to a logger.
// This interface is implemented by ErrorWithLogContext, which wraps some errors that are returned from the protocol,
// and can be used to add context to log statements from the process proposal handler that helps with debugging
// erroneous unexpected and/or non-deterministic behaviors.
type LogContextualizer interface {
	AddLoggingContext(logger log.Logger) log.Logger
	Unwrap() error
}
