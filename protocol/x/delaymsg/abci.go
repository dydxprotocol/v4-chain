package delaymsg

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO(CORE-437): execute messages for this block and delete them from
	// the store. See
	// https://github.com/cosmos/cosmos-sdk/blob/208219a4283bad7fd6c9a3d93f50c96e7efbb3ae/x/gov/abci.go#L134
	// for example.
}
