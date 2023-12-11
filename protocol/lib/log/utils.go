package log

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
