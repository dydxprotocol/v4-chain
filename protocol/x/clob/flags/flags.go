package flags

import (
	"strings"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type ClobFlags struct {
	MevTelemetryHosts      []string
	MevTelemetryIdentifier string
}

// List of CLI flags.
const (
	MevTelemetryHosts      = "mev-telemetry-hosts"
	MevTelemetryIdentifier = "mev-telemetry-identifier"
)

// Default values.

const (
	DefaultMevTelemetryHosts      = ""
	DefaultMevTelemetryIdentifier = ""
)

// AddFlagsToCmd adds flags to app initialization.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --non-validating-full-node true`.
func AddClobFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(
		MevTelemetryHosts,
		DefaultMevTelemetryHosts,
		"Sets the addresses (comma-delimited) to connect to the MEV Telemetry collection agents.",
	)
	cmd.Flags().String(
		MevTelemetryIdentifier,
		DefaultMevTelemetryIdentifier,
		"Sets the identifier to use for MEV Telemetry collection agents.",
	)
}

// GetFlagValuesFromOptions gets values from the `AppOptions` struct which contains values
// from the command-line flags.
func GetClobFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) ClobFlags {
	// Create default result.
	result := ClobFlags{
		MevTelemetryHosts:      []string{},
		MevTelemetryIdentifier: DefaultMevTelemetryIdentifier,
	}

	// Populate the flags if they exist.
	mevTelemetryHostString, ok := appOpts.Get(MevTelemetryHosts).(string)

	if !ok {
		return result
	}

	if mevTelemetryHostString == "" {
		return result
	}

	result.MevTelemetryHosts = strings.Split(mevTelemetryHostString, ",")

	if v, ok := appOpts.Get(MevTelemetryIdentifier).(string); ok {
		result.MevTelemetryIdentifier = v
	}

	return result
}
