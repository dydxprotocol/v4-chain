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

// TxMode returns a textual representation of the tx mode, one of `CheckTx`, `ReCheckTx`, or `DeliverTx`.
func TxMode(ctx sdk.Context) string {
	if ctx.IsReCheckTx() {
		return "ReCheckTx"
	} else if ctx.IsCheckTx() {
		return "CheckTx"
	} else {
		return "DeliverTx"
	}
}
