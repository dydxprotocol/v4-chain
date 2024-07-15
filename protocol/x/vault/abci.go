package vault

import (
	"runtime/debug"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
)

func BeginBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	// Panic is not expected in BeginBlocker, but we should recover instead of
	// halting the chain.
	defer func() {
		if r := recover(); r != nil {
			log.ErrorLog(
				ctx,
				"panic in vault BeginBlocker",
				"stack",
				string(debug.Stack()),
			)
		}
	}()

	keeper.DecommissionNonPositiveEquityVaults(ctx)
}

func EndBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	// Panic is not expected in EndBlocker, but we should recover instead of
	// halting the chain.
	defer func() {
		if r := recover(); r != nil {
			log.ErrorLog(
				ctx,
				"panic in vault EndBlocker",
				"stack",
				string(debug.Stack()),
			)
		}
	}()

	keeper.RefreshAllVaultOrders(ctx)
}
