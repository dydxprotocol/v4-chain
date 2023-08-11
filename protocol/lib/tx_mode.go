package lib

import sdk "github.com/cosmos/cosmos-sdk/types"

func AssertDeliverTxMode(ctx sdk.Context) {
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		panic("assert deliverTx mode failed")
	}
}

func IsDeliverTxMode(ctx sdk.Context) bool {
	return !ctx.IsCheckTx() && !ctx.IsReCheckTx()
}

func AssertCheckTxMode(ctx sdk.Context) {
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		panic("assert checkTx mode failed")
	}
}
