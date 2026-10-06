package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// BlockLimitsConfig_Default uses 0 for MaxStatefulOrderRemovalsPerBlock,
	// which means "no cap" - process all expired orders.
	// This provides backward compatibility with existing chains.
	// Governance can set a specific cap later if needed.
	BlockLimitsConfig_Default = BlockLimitsConfig{
		MaxStatefulOrderRemovalsPerBlock: 0,
	}
)

// BlockLimitsConfigKeeper is an interface that encapsulates all reads and writes to the
// block limits configuration values written to state.
type BlockLimitsConfigKeeper interface {
	GetBlockLimitsConfig(
		ctx sdk.Context,
	) BlockLimitsConfig
	UpdateBlockLimitsConfig(
		ctx sdk.Context,
		config BlockLimitsConfig,
	) error
}
