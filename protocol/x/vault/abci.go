package vault

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
)

func BeginBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	keeper.DecommissionNonPositiveEquityVaults(ctx)
}

func EndBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	keeper.RefreshAllVaultOrders(ctx)
}
