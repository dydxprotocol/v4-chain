package middleware

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/dydxprotocol/v4/app/basic_manager"
)

func NewRunTxPanicLoggingMiddleware(logger log.Logger) baseapp.RecoveryHandler {
	return func(recoveryObj interface{}) error {
		stack := string(debug.Stack())

		var keyvals []interface{}

		for _, module := range basic_manager.ModuleBasics {
			fullModuleName := "/x/" + module.Name()
			if strings.Contains(stack, fullModuleName) {
				keyvals = append(keyvals, fullModuleName, "true")
			}
		}

		keyvals = append(keyvals, "stack trace", stack)

		logger.Error(
			fmt.Sprintf(
				"runTx panic'ed with %v",
				recoveryObj,
			),
			keyvals...,
		)

		// Return nil to indicate that this error was not processed.
		// Pass it on to the next middleware in chain.
		return nil
	}
}
