package error

import "github.com/cometbft/cometbft/libs/log"

type LogContextualizer interface {
	AddLoggingContext(logger log.Logger) log.Logger
	Unwrap() error
}
