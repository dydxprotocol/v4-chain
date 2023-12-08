package flags

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type Flags struct {
	DdAgentHost           string
	DdTraceAgentPort      uint16
	NonValidatingFullNode bool

	// Existing flags
	GrpcAddress string
	GrpcEnable  bool
}

// List of CLI flags.
const (
	DdAgentHost               = "dd-agent-host"
	DdTraceAgentPort          = "dd-trace-agent-port"
	NonValidatingFullNodeFlag = "non-validating-full-node"

	// Cosmos flags below. These config values can be set as flags or in config.toml.
	GrpcAddress = "grpc.address"
	GrpcEnable  = "grpc.enable"
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
	cmd.Flags().Bool(
		NonValidatingFullNodeFlag,
		DefaultNonValidatingFullNode,
		"Whether to run in non-validating full-node mode. "+
			"This disables the pricing daemon and enables the full-node ProcessProposal logic. "+
			"Validators should _never_ use this mode.",
	)
	cmd.Flags().String(
		DdAgentHost,
		DefaultDdAgentHost,
		"Sets the address to connect to for the Datadog Agent.",
	)
	cmd.Flags().Uint16(
		DdTraceAgentPort,
		DefaultDdTraceAgentPort,
		"Sets the Datadog Agent port.",
	)
}

// Validate checks that the flags are valid.
func (f *Flags) Validate() error {
	// Validtors must have cosmos grpc services enabled.
	if !f.NonValidatingFullNode && !f.GrpcEnable {
		return fmt.Errorf("grpc.enable must be set to true - validating requires gRPC server")
	}
	return nil
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

		// These are the default values from the Cosmos flags.
		GrpcAddress: config.DefaultGRPCAddress,
		GrpcEnable:  true,
	}

	// Populate the flags if they exist.
	if option := appOpts.Get(NonValidatingFullNodeFlag); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.NonValidatingFullNode = v
		}
	}

	if option := appOpts.Get(DdAgentHost); option != nil {
		if v, err := cast.ToStringE(option); err == nil {
			result.DdAgentHost = v
		}
	}

	if option := appOpts.Get(DdTraceAgentPort); option != nil {
		if v, err := cast.ToUint16E(option); err == nil {
			result.DdTraceAgentPort = v
		}
	}

	if option := appOpts.Get(GrpcAddress); option != nil {
		if v, err := cast.ToStringE(option); err == nil {
			result.GrpcAddress = v
		}
	}

	if option := appOpts.Get(GrpcEnable); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.GrpcEnable = v
		}
	}

	return result
}
