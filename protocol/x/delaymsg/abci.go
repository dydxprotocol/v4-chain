package delaymsg

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k types.DelayMsgKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	keeper.DispatchMessagesForBlock(k, ctx)
}
