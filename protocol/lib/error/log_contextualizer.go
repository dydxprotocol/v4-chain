package error

import "cosmossdk.io/log"

type LogContextualizer interface {
	AddLoggingContext(logger log.Logger) log.Logger
	Unwrap() error
}
