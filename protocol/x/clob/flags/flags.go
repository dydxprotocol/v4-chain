package flags

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type ClobFlags struct {
	MaxLiquidationOrdersPerBlock uint32

	MevTelemetryHost       string
	MevTelemetryIdentifier string
}

// List of CLI flags.
const (
	// Liquidations.
	MaxLiquidationOrdersPerBlock = "max-liquidation-orders-per-block"

	// Mev.
	MevTelemetryHost       = "mev-telemetry-host"
	MevTelemetryIdentifier = "mev-telemetry-identifier"
)

// Default values.
const (
	DefaultMaxLiquidationOrdersPerBlock = 10

	DefaultMevTelemetryHost       = ""
	DefaultMevTelemetryIdentifier = ""
)

// AddFlagsToCmd adds flags to app initialization.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --non-validating-full-node true`.
func AddClobFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Uint32(
		MaxLiquidationOrdersPerBlock,
		DefaultMaxLiquidationOrdersPerBlock,
		"Sets the maximum number of liquidation orders to process per block.",
	)
	cmd.Flags().String(
		MevTelemetryHost,
		DefaultMevTelemetryHost,
		"Sets the address to connect to for the MEV Telemetry collection agent.",
	)
	cmd.Flags().String(
		MevTelemetryIdentifier,
		DefaultMevTelemetryIdentifier,
		"Sets the identifier to use for MEV Telemetry collection agent.",
	)
}

func GetDefaultClobFlags() ClobFlags {
	return ClobFlags{
		MaxLiquidationOrdersPerBlock: DefaultMaxLiquidationOrdersPerBlock,
		MevTelemetryHost:             DefaultMevTelemetryHost,
		MevTelemetryIdentifier:       DefaultMevTelemetryIdentifier,
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
	if v, ok := appOpts.Get(MevTelemetryHost).(string); ok {
		result.MevTelemetryHost = v
	}

	if v, ok := appOpts.Get(MevTelemetryIdentifier).(string); ok {
		result.MevTelemetryIdentifier = v
	}

	if v, ok := appOpts.Get(MaxLiquidationOrdersPerBlock).(uint32); ok {
		result.MaxLiquidationOrdersPerBlock = v
	}

	return result
}
