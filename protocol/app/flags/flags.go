package flags

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type Flags struct {
	DdAgentHost           string
	DdTraceAgentPort      uint16
	NonValidatingFullNode bool
}

// List of CLI flags.
const (
	DdAgentHost               = "dd-agent-host"
	DdTraceAgentPort          = "dd-trace-agent-port"
	NonValidatingFullNodeFlag = "non-validating-full-node"
)

// Default values.
const (
	DefaultDdAgentHost           = ""
	DefaultDdTraceAgentPort      = 8126
	DefaultNonValidatingFullNode = false
)

// AddFlagsToCmd adds flags to app initialization.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --non-validating-full-node true`.
func AddFlagsToCmd(cmd *cobra.Command) {
	cmd.
		Flags().
		Bool(
			NonValidatingFullNodeFlag,
			DefaultNonValidatingFullNode,
			"Whether to run in non-validating full-node mode. "+
				"This disables the pricing daemon and enables the full-node ProcessProposal logic. "+
				"Validators should _never_ use this mode.",
		)
	cmd.Flags().
		String(
			DdAgentHost,
			DefaultDdAgentHost,
			"Sets the address to connect to for the Datadog Agent.",
		)
	cmd.Flags().
		Uint16(
			DdTraceAgentPort,
			DefaultDdTraceAgentPort,
			"Sets the Datadog Agent port.",
		)
}

// GetFlagValuesFromOptions gets values from the `AppOptions` struct which contains values
// from the command-line flags.
func GetFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) Flags {
	// Create default result.
	result := Flags{
		NonValidatingFullNode: DefaultNonValidatingFullNode,
		DdAgentHost:           DefaultDdAgentHost,
		DdTraceAgentPort:      DefaultDdTraceAgentPort,
	}

	// Populate the flags if they exist.
	if v, ok := appOpts.Get(NonValidatingFullNodeFlag).(bool); ok {
		result.NonValidatingFullNode = v
	}

	if v, ok := appOpts.Get(DdAgentHost).(string); ok {
		result.DdAgentHost = v
	}

	if v, ok := appOpts.Get(DdTraceAgentPort).(uint16); ok {
		result.DdTraceAgentPort = v
	}

	return result
}
