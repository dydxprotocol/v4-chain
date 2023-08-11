package flags

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type ClobFlags struct {
	MevTelemetryHost       string
	MevTelemetryIdentifier string
}

// List of CLI flags.
const (
	MevTelemetryHost       = "mev-telemetry-host"
	MevTelemetryIdentifier = "mev-telemetry-identifier"
)

// Default values.
const (
	DefaultMevTelemetryHost       = ""
	DefaultMevTelemetryIdentifier = ""
)

// AddFlagsToCmd adds flags to app initialization.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --non-validating-full-node true`.
func AddClobFlagsToCmd(cmd *cobra.Command) {
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

// GetFlagValuesFromOptions gets values from the `AppOptions` struct which contains values
// from the command-line flags.
func GetClobFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) ClobFlags {
	// Create default result.
	result := ClobFlags{
		MevTelemetryHost:       DefaultMevTelemetryHost,
		MevTelemetryIdentifier: DefaultMevTelemetryIdentifier,
	}

	// Populate the flags if they exist.
	if v, ok := appOpts.Get(MevTelemetryHost).(string); ok {
		result.MevTelemetryHost = v
	}

	if v, ok := appOpts.Get(MevTelemetryIdentifier).(string); ok {
		result.MevTelemetryIdentifier = v
	}

	return result
}
