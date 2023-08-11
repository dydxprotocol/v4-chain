package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

func NewRunTxPanicLoggingMiddleware(logger log.Logger) baseapp.RecoveryHandler {
	return func(recoveryObj interface{}) error {
		logger.Error(
			fmt.Sprintf(
				"runTx panic'ed with %v",
				recoveryObj,
			),
			"stack trace",
			string(debug.Stack()),
		)
		// Return nil to indicate that this error was not processed.
		// Pass it on to the next middleware in chain.
		return nil
	}
}
