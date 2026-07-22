package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// ConditionalOrderTriggerConfig is the consensus-level configuration that gates the bounded
// conditional-order trigger path in MaybeTriggerConditionalOrders.
//
// It is persisted in the clob KVStore (consensus state) under ConditionalOrderTriggerConfigKey.
// When the key is absent, the zero value applies: Enabled == false, which runs the LEGACY
// full-scan trigger path (byte-for-byte identical to the pre-fix implementation). This is what
// makes the fix rolling-deployable: every node ships the new binary but continues to run the
// legacy path — and produces identical app hashes — until governance flips Enabled at a
// coordinated height.
//
// When Enabled is true:
//   - MaxTriggersPerBlock bounds the total conditional-order triggers per block (per-block budget).
//     The index is keyed by crossing proximity (triggerSubticks), so the budget is drained
//     nearest-crossing-first; far-from-market orders that cannot cross are never visited, so they
//     do not contribute to per-block work.
//   - MaxRemovalsPerBlock bounds the number of expired stateful orders removed per block, keeping
//     the EndBlocker expiry-scan a bounded per-block workload.
//   - MaxUntriggeredConditionalOrdersGlobal and MaxUntriggeredConditionalOrdersPerSubaccount bound
//     the resting untriggered conditional order set (admission caps, governance-settable).
type ConditionalOrderTriggerConfig struct {
	Enabled                                      bool
	MaxTriggersPerBlock                          uint32
	MaxRemovalsPerBlock                          uint32
	MaxUntriggeredConditionalOrdersGlobal        uint32
	MaxUntriggeredConditionalOrdersPerSubaccount uint32
}

// conditionalOrderTriggerConfigVersion is the layout version byte for the on-disk encoding.
// Layout: Enabled(1) + MaxTriggersPerBlock(4) + MaxRemovalsPerBlock(4) +
// MaxUntriggeredConditionalOrdersGlobal(4) + MaxUntriggeredConditionalOrdersPerSubaccount(4).
const conditionalOrderTriggerConfigVersion byte = 0x01

// conditionalOrderTriggerConfigLen is the fixed blob length: version(1) + enabled(1) + 4×uint32.
const conditionalOrderTriggerConfigLen = 1 + 1 + 4 + 4 + 4 + 4

// MaxConditionalRemovalsPerBlock is the default maximum number of expired stateful orders that
// RemoveExpiredStatefulOrders will process per block when the mitigation is enabled.
// Chosen to match MaxConditionalTriggersPerBlock so both budgets share the same default scale.
const MaxConditionalRemovalsPerBlock = 1000

// DefaultConditionalOrderTriggerConfig returns the config that applies when no config has been
// set in state: disabled (legacy behavior). All numeric fields default to the package constants
// so that if governance enables the config without specifying individual limits, sane bounds apply.
func DefaultConditionalOrderTriggerConfig() ConditionalOrderTriggerConfig {
	return ConditionalOrderTriggerConfig{
		Enabled:                               false,
		MaxTriggersPerBlock:                   MaxConditionalTriggersPerBlock,
		MaxRemovalsPerBlock:                   MaxConditionalRemovalsPerBlock,
		MaxUntriggeredConditionalOrdersGlobal: MaxUntriggeredConditionalOrdersGlobal,
		MaxUntriggeredConditionalOrdersPerSubaccount: MaxUntriggeredConditionalOrdersPerSubaccount,
	}
}

// encodeConditionalOrderTriggerConfig serializes the config into a fixed-layout byte blob:
//
//	<version:1> <enabled:1> <maxTriggersPerBlock:4 BE> <maxRemovalsPerBlock:4 BE>
//	<maxUntriggeredGlobal:4 BE> <maxUntriggeredPerSubaccount:4 BE>
//
// A fixed layout (rather than a proto message) avoids a proto/codegen migration for this field
// while remaining fully deterministic across nodes.
func encodeConditionalOrderTriggerConfig(cfg ConditionalOrderTriggerConfig) []byte {
	b := make([]byte, conditionalOrderTriggerConfigLen)
	b[0] = conditionalOrderTriggerConfigVersion
	if cfg.Enabled {
		b[1] = 0x01
	}
	binary.BigEndian.PutUint32(b[2:6], cfg.MaxTriggersPerBlock)
	binary.BigEndian.PutUint32(b[6:10], cfg.MaxRemovalsPerBlock)
	binary.BigEndian.PutUint32(b[10:14], cfg.MaxUntriggeredConditionalOrdersGlobal)
	binary.BigEndian.PutUint32(b[14:18], cfg.MaxUntriggeredConditionalOrdersPerSubaccount)
	return b
}

// decodeConditionalOrderTriggerConfig deserializes a config blob. An unknown version or short blob
// returns the default (disabled) config so the chain fails safe to the legacy path rather than
// panicking in a consensus-critical EndBlocker.
func decodeConditionalOrderTriggerConfig(b []byte) ConditionalOrderTriggerConfig {
	if len(b) != conditionalOrderTriggerConfigLen || b[0] != conditionalOrderTriggerConfigVersion {
		return DefaultConditionalOrderTriggerConfig()
	}
	return ConditionalOrderTriggerConfig{
		Enabled:                               b[1] == 0x01,
		MaxTriggersPerBlock:                   binary.BigEndian.Uint32(b[2:6]),
		MaxRemovalsPerBlock:                   binary.BigEndian.Uint32(b[6:10]),
		MaxUntriggeredConditionalOrdersGlobal: binary.BigEndian.Uint32(b[10:14]),
		MaxUntriggeredConditionalOrdersPerSubaccount: binary.BigEndian.Uint32(b[14:18]),
	}
}

// GetConditionalOrderTriggerConfig reads the trigger config from consensus state. When no config
// has been set, it returns the default (disabled) config, so a chain that has never set the
// config runs the legacy trigger path.
func (k Keeper) GetConditionalOrderTriggerConfig(ctx sdk.Context) ConditionalOrderTriggerConfig {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ConditionalOrderTriggerConfigKey))
	if b == nil {
		return DefaultConditionalOrderTriggerConfig()
	}
	return decodeConditionalOrderTriggerConfig(b)
}

// SetConditionalOrderTriggerConfig writes the trigger config to consensus state.
// Intended to be driven by a governance message or an upgrade handler at a coordinated height.
// Normalizes the config so that zero values fall back to sane defaults.
func (k Keeper) SetConditionalOrderTriggerConfig(ctx sdk.Context, cfg ConditionalOrderTriggerConfig) {
	if cfg.MaxTriggersPerBlock == 0 {
		cfg.MaxTriggersPerBlock = MaxConditionalTriggersPerBlock
	}
	if cfg.MaxRemovalsPerBlock == 0 {
		cfg.MaxRemovalsPerBlock = MaxConditionalRemovalsPerBlock
	}
	if cfg.MaxUntriggeredConditionalOrdersGlobal == 0 {
		cfg.MaxUntriggeredConditionalOrdersGlobal = MaxUntriggeredConditionalOrdersGlobal
	}
	if cfg.MaxUntriggeredConditionalOrdersPerSubaccount == 0 {
		cfg.MaxUntriggeredConditionalOrdersPerSubaccount = MaxUntriggeredConditionalOrdersPerSubaccount
	}

	// Detect the disabled->enabled activation transition. While disabled, the trigger-price index
	// (consensus state) is NOT maintained by the per-placement hooks, so any conditional orders
	// resting at activation have no index entries. We must (re)build the index at the transition,
	// otherwise the enabled bounded trigger path would never evaluate those pre-activation orders.
	// This runs deterministically on every node at the same activation block (the flag is flipped
	// via a governance action / upgrade handler), so all nodes reconcile identically.
	wasEnabled := k.GetConditionalOrderTriggerConfig(ctx).Enabled

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.ConditionalOrderTriggerConfigKey), encodeConditionalOrderTriggerConfig(cfg))

	if cfg.Enabled && !wasEnabled {
		k.BackfillConditionalOrderTriggerPriceIndex(ctx)
	}
}

// SetConditionalOrderTriggerConfigParams is a scalar-parameter wrapper around
// SetConditionalOrderTriggerConfig, invoked by the governance message handler. The config struct
// lives in this package, so the ClobKeeper interface (package types) cannot reference it directly
// without an import cycle; the handler therefore passes scalars.
func (k Keeper) SetConditionalOrderTriggerConfigParams(
	ctx sdk.Context,
	enabled bool,
	maxTriggersPerBlock uint32,
	maxRemovalsPerBlock uint32,
	maxUntriggeredConditionalOrdersGlobal uint32,
	maxUntriggeredConditionalOrdersPerSubaccount uint32,
) {
	k.SetConditionalOrderTriggerConfig(ctx, ConditionalOrderTriggerConfig{
		Enabled:                               enabled,
		MaxTriggersPerBlock:                   maxTriggersPerBlock,
		MaxRemovalsPerBlock:                   maxRemovalsPerBlock,
		MaxUntriggeredConditionalOrdersGlobal: maxUntriggeredConditionalOrdersGlobal,
		MaxUntriggeredConditionalOrdersPerSubaccount: maxUntriggeredConditionalOrdersPerSubaccount,
	})
}
