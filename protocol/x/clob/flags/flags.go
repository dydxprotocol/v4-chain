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
)

// Default values.

const (
	DefaultMaxLiquidationAttemptsPerBlock      = 50
	DefaultMaxDeleveragingAttemptsPerBlock     = 10
	DefaultMaxDeleveragingSubaccountsToIterate = 500

	DefaultMevTelemetryEnabled    = false
	DefaultMevTelemetryHostsFlag  = ""
	DefaultMevTelemetryIdentifier = ""
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
}

func GetDefaultClobFlags() ClobFlags {
	return ClobFlags{
		MaxLiquidationAttemptsPerBlock:      DefaultMaxLiquidationAttemptsPerBlock,
		MaxDeleveragingAttemptsPerBlock:     DefaultMaxDeleveragingAttemptsPerBlock,
		MaxDeleveragingSubaccountsToIterate: DefaultMaxDeleveragingSubaccountsToIterate,
		MevTelemetryEnabled:                 DefaultMevTelemetryEnabled,
		MevTelemetryHosts:                   DefaultMevTelemetryHosts,
		MevTelemetryIdentifier:              DefaultMevTelemetryIdentifier,
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

	return result
}
