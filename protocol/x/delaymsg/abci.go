package delaymsg

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k types.DelayMsgKeeper) {
	k.DispatchMessagesForBlock(ctx)
}
