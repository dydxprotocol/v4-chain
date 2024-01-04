package error

import "cosmossdk.io/log"

// LogContextualizer describes an object that can add context - that is, descriptive key-value pairs, to a logger.
// This interface is implemented by ErrorWithLogContext, which wraps some errors that are returned from the protocol.
// This can be used, for example, to add context to log statements from the process proposal handler that helps with
// debugging erroneous unexpected and/or non-deterministic behaviors.
type LogContextualizer interface {
	AddLoggingContext(logger log.Logger) log.Logger
	Unwrap() error
}
