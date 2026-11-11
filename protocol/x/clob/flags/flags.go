package flags

import (
	"fmt"
	"strings"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type ClobFlags struct {
	MaxLiquidationAttemptsPerBlock      uint32
	MaxDeleveragingAttemptsPerBlock     uint32
	MaxDeleveragingSubaccountsToIterate uint32

	MevTelemetryEnabled    bool
	MevTelemetryHosts      []string
	MevTelemetryIdentifier string

	ShortTermTokenBucketEnabled            bool
	ShortTermTokenBucketShadow             bool
	ShortTermTokenBucketMinCapacity        uint32
	ShortTermTokenBucketBurstBlocks        uint32
	ShortTermTokenBucketVolumePerOrder     uint64
	ShortTermTokenBucketCollateralPerOrder uint64
}

// List of CLI flags.
const (
	// Liquidations and deleveraging.
	MaxLiquidationAttemptsPerBlock      = "max-liquidation-attempts-per-block"
	MaxDeleveragingAttemptsPerBlock     = "max-deleveraging-attempts-per-block"
	MaxDeleveragingSubaccountsToIterate = "max-deleveraging-subaccounts-to-iterate"

	// Mev.
	MevTelemetryEnabled    = "mev-telemetry-enabled"
	MevTelemetryHosts      = "mev-telemetry-hosts"
	MevTelemetryIdentifier = "mev-telemetry-identifier"

	// Short-term order token bucket.
	ShortTermTokenBucketEnabled            = "short-term-token-bucket-enabled"
	ShortTermTokenBucketShadow             = "short-term-token-bucket-shadow"
	ShortTermTokenBucketMinCapacity        = "short-term-token-bucket-min-capacity"
	ShortTermTokenBucketBurstBlocks        = "short-term-token-bucket-burst-blocks"
	ShortTermTokenBucketVolumePerOrder     = "short-term-token-bucket-volume-per-order-quote-quantums"
	ShortTermTokenBucketCollateralPerOrder = "short-term-token-bucket-collateral-per-order-quote-quantums"
)

// Default values.

const (
	DefaultMaxLiquidationAttemptsPerBlock      = 50
	DefaultMaxDeleveragingAttemptsPerBlock     = 10
	DefaultMaxDeleveragingSubaccountsToIterate = 500

	DefaultMevTelemetryEnabled    = false
	DefaultMevTelemetryHostsFlag  = ""
	DefaultMevTelemetryIdentifier = ""

	DefaultShortTermTokenBucketEnabled            = true
	DefaultShortTermTokenBucketShadow             = false
	DefaultShortTermTokenBucketMinCapacity        = 100
	DefaultShortTermTokenBucketBurstBlocks        = 1000
	DefaultShortTermTokenBucketVolumePerOrder     = 50_000_000_000 // $50k of 30-day volume
	DefaultShortTermTokenBucketCollateralPerOrder = 500_000_000    // $500 of net collateral
)

var DefaultMevTelemetryHosts = []string{}

// AddFlagsToCmd adds flags to app initialization.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --non-validating-full-node true`.
func AddClobFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Uint32(
		MaxLiquidationAttemptsPerBlock,
		DefaultMaxLiquidationAttemptsPerBlock,
		fmt.Sprintf(
			"Sets the maximum number of liquidation orders to process per block. Default = %d",
			DefaultMaxLiquidationAttemptsPerBlock,
		),
	)
	cmd.Flags().Uint32(
		MaxDeleveragingAttemptsPerBlock,
		DefaultMaxDeleveragingAttemptsPerBlock,
		fmt.Sprintf(
			"Sets the maximum number of attempted deleveraging events per block. Default = %d",
			DefaultMaxDeleveragingAttemptsPerBlock,
		),
	)
	cmd.Flags().Uint32(
		MaxDeleveragingSubaccountsToIterate,
		DefaultMaxDeleveragingSubaccountsToIterate,
		fmt.Sprintf(
			"Sets the maximum number of subaccounts iterated for each deleveraging event. Default = %d",
			DefaultMaxDeleveragingSubaccountsToIterate,
		),
	)
	cmd.Flags().Bool(
		MevTelemetryEnabled,
		DefaultMevTelemetryEnabled,
		"Runs the MEV Telemetry collection agent if true.",
	)
	cmd.Flags().String(
		MevTelemetryHosts,
		DefaultMevTelemetryHostsFlag,
		"Sets the addresses (comma-delimited) to connect to the MEV Telemetry collection agents.",
	)
	cmd.Flags().String(
		MevTelemetryIdentifier,
		DefaultMevTelemetryIdentifier,
		"Sets the identifier to use for MEV Telemetry collection agents.",
	)
	cmd.Flags().Bool(
		ShortTermTokenBucketEnabled,
		DefaultShortTermTokenBucketEnabled,
		"Rejects short-term orders in CheckTx once the owner address's token bucket is empty. The bucket "+
			"refills each block from the address's 30-day volume and the net collateral across its subaccounts.",
	)
	cmd.Flags().Bool(
		ShortTermTokenBucketShadow,
		DefaultShortTermTokenBucketShadow,
		"Counts orders the token bucket would reject without rejecting them.",
	)
	cmd.Flags().Uint32(
		ShortTermTokenBucketMinCapacity,
		DefaultShortTermTokenBucketMinCapacity,
		fmt.Sprintf(
			"Short-term orders any address can place in a burst. Default = %d",
			DefaultShortTermTokenBucketMinCapacity,
		),
	)
	cmd.Flags().Uint32(
		ShortTermTokenBucketBurstBlocks,
		DefaultShortTermTokenBucketBurstBlocks,
		fmt.Sprintf(
			"Blocks of refill a bucket holds when that exceeds the minimum capacity. Default = %d",
			DefaultShortTermTokenBucketBurstBlocks,
		),
	)
	cmd.Flags().Uint64(
		ShortTermTokenBucketVolumePerOrder,
		DefaultShortTermTokenBucketVolumePerOrder,
		fmt.Sprintf(
			"30-day volume in quote quantums that refills one short-term order per block. 0 disables the "+
				"volume term. Default = %d",
			DefaultShortTermTokenBucketVolumePerOrder,
		),
	)
	cmd.Flags().Uint64(
		ShortTermTokenBucketCollateralPerOrder,
		DefaultShortTermTokenBucketCollateralPerOrder,
		fmt.Sprintf(
			"Net collateral in quote quantums that refills one short-term order per block. 0 disables the "+
				"collateral term. Default = %d",
			DefaultShortTermTokenBucketCollateralPerOrder,
		),
	)
}

func GetDefaultClobFlags() ClobFlags {
	return ClobFlags{
		MaxLiquidationAttemptsPerBlock:         DefaultMaxLiquidationAttemptsPerBlock,
		MaxDeleveragingAttemptsPerBlock:        DefaultMaxDeleveragingAttemptsPerBlock,
		MaxDeleveragingSubaccountsToIterate:    DefaultMaxDeleveragingSubaccountsToIterate,
		MevTelemetryEnabled:                    DefaultMevTelemetryEnabled,
		MevTelemetryHosts:                      DefaultMevTelemetryHosts,
		MevTelemetryIdentifier:                 DefaultMevTelemetryIdentifier,
		ShortTermTokenBucketEnabled:            DefaultShortTermTokenBucketEnabled,
		ShortTermTokenBucketShadow:             DefaultShortTermTokenBucketShadow,
		ShortTermTokenBucketMinCapacity:        DefaultShortTermTokenBucketMinCapacity,
		ShortTermTokenBucketBurstBlocks:        DefaultShortTermTokenBucketBurstBlocks,
		ShortTermTokenBucketVolumePerOrder:     DefaultShortTermTokenBucketVolumePerOrder,
		ShortTermTokenBucketCollateralPerOrder: DefaultShortTermTokenBucketCollateralPerOrder,
	}
}

// GetFlagValuesFromOptions gets values from the `AppOptions` struct which contains values
// from the command-line flags.
func GetClobFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) ClobFlags {
	// Create default result.
	result := GetDefaultClobFlags()

	// Populate the flags if they exist.
	if option := appOpts.Get(MevTelemetryEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.MevTelemetryEnabled = v
		}
	}

	if option := appOpts.Get(MevTelemetryHosts); option != nil {
		if v, err := cast.ToStringE(option); err == nil && len(v) > 0 {
			result.MevTelemetryHosts = strings.Split(v, ",")
		}
	}

	if option := appOpts.Get(MevTelemetryIdentifier); option != nil {
		if v, err := cast.ToStringE(option); err == nil && len(v) > 0 {
			result.MevTelemetryIdentifier = v
		}
	}

	if option := appOpts.Get(MaxLiquidationAttemptsPerBlock); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.MaxLiquidationAttemptsPerBlock = v
		}
	}

	if option := appOpts.Get(MaxDeleveragingAttemptsPerBlock); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.MaxDeleveragingAttemptsPerBlock = v
		}
	}

	if option := appOpts.Get(MaxDeleveragingSubaccountsToIterate); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.MaxDeleveragingSubaccountsToIterate = v
		}
	}

	if option := appOpts.Get(ShortTermTokenBucketEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.ShortTermTokenBucketEnabled = v
		}
	}

	if option := appOpts.Get(ShortTermTokenBucketShadow); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.ShortTermTokenBucketShadow = v
		}
	}

	if option := appOpts.Get(ShortTermTokenBucketMinCapacity); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.ShortTermTokenBucketMinCapacity = v
		}
	}

	if option := appOpts.Get(ShortTermTokenBucketBurstBlocks); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.ShortTermTokenBucketBurstBlocks = v
		}
	}

	if option := appOpts.Get(ShortTermTokenBucketVolumePerOrder); option != nil {
		if v, err := cast.ToUint64E(option); err == nil {
			result.ShortTermTokenBucketVolumePerOrder = v
		}
	}

	if option := appOpts.Get(ShortTermTokenBucketCollateralPerOrder); option != nil {
		if v, err := cast.ToUint64E(option); err == nil {
			result.ShortTermTokenBucketCollateralPerOrder = v
		}
	}

	return result
}
