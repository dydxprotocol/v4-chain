package log

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// // This structure describes a logger object that carries context with it.
// type DydxLogger struct {
// 	// original logger
// 	originalLogger *log.Logger

// 	// Store tag values that should be auto-added to any logs emitted.
// 	// Any tag values added here should be protocol-specific.
// 	// New values will clobber old values.
// 	keyValues map[string]interface{}
// }

func InfoLog(ctx sdk.Context, msg string, keyvals ...interface{}) {
	ctx.Logger().Info(msg, keyvals...)
}

func DebugLog(ctx sdk.Context, msg string, keyvals ...interface{}) {
	ctx.Logger().Debug(msg, keyvals...)
}

func ErrorLog(ctx sdk.Context, msg string, err error, keyvals ...interface{}) {
	ctx.Logger().Error(msg, append(keyvals, Error, err))
}

func AddPersistentTagsToLogger(ctx sdk.Context, keyvals ...interface{}) {
	ctx.WithLogger(ctx.Logger().With(keyvals))
}
