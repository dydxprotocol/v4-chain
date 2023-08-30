package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeNewEpoch = "new_epoch"

	AttributeKeyEpochInfoName       = "epoch_info_name"
	AttributeKeyEpochNumber         = "epoch_number"
	AttributeKeyEpochStartTickTime  = "epoch_start_tick_time"
	AttributeKeyEpochStartBlockTime = "epoch_start_block_time"
	AttributeKeyEpochStartBlock     = "epoch_start_block"
)

// NewEpochEvent constructs a new_epoch sdk.Event.
func NewEpochEvent(ctx sdk.Context, epoch EpochInfo, currentTick uint32) sdk.Event {
	return sdk.NewEvent(
		EventTypeNewEpoch,
		sdk.NewAttribute(AttributeKeyEpochInfoName, epoch.Name),
		sdk.NewAttribute(AttributeKeyEpochNumber, fmt.Sprint(epoch.CurrentEpoch)),
		sdk.NewAttribute(AttributeKeyEpochStartTickTime, fmt.Sprint(currentTick)),
		sdk.NewAttribute(AttributeKeyEpochStartBlockTime, fmt.Sprint(ctx.BlockTime().Unix())),
		sdk.NewAttribute(AttributeKeyEpochStartBlock, fmt.Sprint(epoch.CurrentEpochStartBlock)),
	)
}
