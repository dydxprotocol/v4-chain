package app

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis file.
//
// Deprecated: This is a legacy feature of cosmos that is known to be unstable, so we
// explicitly do not support its usage.
func (app *App) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string, modulesToExport []string,
) (servertypes.ExportedApp, error) {
	panic("ExportAppStateAndValidators not supported")
}
