package log

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InfoLog reports msg as an info level log with key vals.
func InfoLog(ctx sdk.Context, msg string, keyvals ...interface{}) {
	ctx.Logger().Info(msg, keyvals...)
}

// DebugLog reports msg as a debug level log with key vals.
func DebugLog(ctx sdk.Context, msg string, keyvals ...interface{}) {
	ctx.Logger().Debug(msg, keyvals...)
}

// ErrorLogWithError reports msg as a error log, as well as attaching the error
// object to the log for datadog error tracking.
func ErrorLogWithError(ctx sdk.Context, msg string, err error, keyvals ...interface{}) {
	ctx.Logger().Error(msg, append(keyvals, Error, err))
}

// ErrorLog reports msg as a error log. It constructs error object on the fly with
// the given message object.
// Please try to use define a new message and use `ErrorLogWithError` instead.
func ErrorLog(ctx sdk.Context, msg string, keyvals ...interface{}) {
	err := errors.New(msg)
	ErrorLogWithError(ctx, msg, err, keyvals...)
}

// AddPersistentTagsToLogger returns a new sdk.Context with a logger that has new persistent
// tags that are added to all logs emitted.
func AddPersistentTagsToLogger(ctx sdk.Context, keyvals ...interface{}) sdk.Context {
	return ctx.WithLogger(ctx.Logger().With(keyvals...))
}
