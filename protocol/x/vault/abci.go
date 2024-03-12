package vault

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
)

// func PrepareCheckState(
// 	ctx sdk.Context,
// 	keeper *keeper.Keeper,
// ) {
// 	keeper.ProvideLiquidity(ctx)
// }

func EndBlocker(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	keeper.ProvideLiquidity(ctx)
}
