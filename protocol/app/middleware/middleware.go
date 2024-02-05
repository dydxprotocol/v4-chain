package middleware

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"

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
		if err, isError := recoveryObj.(error); isError {
			keyvals = append(keyvals, "error", err)
		}

		keyvals = append(keyvals, "stack_trace", stack)

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
