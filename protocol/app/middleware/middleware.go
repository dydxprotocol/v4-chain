package middleware

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/module"
	"os"
	"runtime/debug"
	"strings"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	kitlog "github.com/go-kit/log"
)

var (
	Logger = log.NewLogger(kitlog.NewSyncWriter(os.Stdout))
)

func NewRunTxPanicLoggingMiddleware(moduleBasics module.BasicManager) baseapp.RecoveryHandler {
	return func(recoveryObj interface{}) error {
		stack := string(debug.Stack())

		var keyvals []interface{}

		for _, module := range moduleBasics {
			fullModuleName := "/x/" + module.Name()
			if strings.Contains(stack, fullModuleName) {
				keyvals = append(keyvals, fullModuleName, "true")
			}
		}

		keyvals = append(keyvals, "stack trace", stack)

		Logger.Error(
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
