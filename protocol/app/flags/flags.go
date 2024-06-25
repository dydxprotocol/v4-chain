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
	DdErrorTrackingFormat bool

	// Existing flags
	GrpcAddress string
	GrpcEnable  bool

	// Grpc Streaming
	GrpcStreamingEnabled              bool
	GrpcStreamingFlushIntervalMs      uint32
	GrpcStreamingMaxBatchSize         uint32
	GrpcStreamingMaxChannelBufferSize uint32

	VEOracleEnabled bool // Slinky Vote Extensions
}

// List of CLI flags.
const (
	DdAgentHost               = "dd-agent-host"
	DdTraceAgentPort          = "dd-trace-agent-port"
	NonValidatingFullNodeFlag = "non-validating-full-node"
	DdErrorTrackingFormat     = "dd-error-tracking-format"

	// Cosmos flags below. These config values can be set as flags or in config.toml.
	GrpcAddress = "grpc.address"
	GrpcEnable  = "grpc.enable"

	// Grpc Streaming
	GrpcStreamingEnabled              = "grpc-streaming-enabled"
	GrpcStreamingFlushIntervalMs      = "grpc-streaming-flush-interval-ms"
	GrpcStreamingMaxBatchSize         = "grpc-streaming-max-batch-size"
	GrpcStreamingMaxChannelBufferSize = "grpc-streaming-max-channel-buffer-size"

	// Slinky VEs enabled
	VEOracleEnabled = "slinky-vote-extension-oracle-enabled"
)

// Default values.
const (
	DefaultDdAgentHost           = ""
	DefaultDdTraceAgentPort      = 8126
	DefaultNonValidatingFullNode = false
	DefaultDdErrorTrackingFormat = false

	DefaultGrpcStreamingEnabled              = false
	DefaultGrpcStreamingFlushIntervalMs      = 50
	DefaultGrpcStreamingMaxBatchSize         = 2000
	DefaultGrpcStreamingMaxChannelBufferSize = 2000

	DefaultVEOracleEnabled = true
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
	cmd.Flags().Bool(
		DdErrorTrackingFormat,
		DefaultDdErrorTrackingFormat,
		"Enable formatting of log error tags to datadog error tracking format",
	)
	cmd.Flags().Bool(
		GrpcStreamingEnabled,
		DefaultGrpcStreamingEnabled,
		"Whether to enable grpc streaming for full nodes",
	)
	cmd.Flags().Uint32(
		GrpcStreamingFlushIntervalMs,
		DefaultGrpcStreamingFlushIntervalMs,
		"Flush interval (in ms) for grpc streaming",
	)
	cmd.Flags().Uint32(
		GrpcStreamingMaxBatchSize,
		DefaultGrpcStreamingMaxBatchSize,
		"Maximum batch size before grpc streaming cancels all subscriptions",
	)
	cmd.Flags().Uint32(
		GrpcStreamingMaxChannelBufferSize,
		DefaultGrpcStreamingMaxChannelBufferSize,
		"Maximum per-subscription channel size before grpc streaming cancels a singular subscription",
	)
	cmd.Flags().Bool(
		VEOracleEnabled,
		DefaultVEOracleEnabled,
		"Whether to run on-chain oracle via slinky vote extensions",
	)
}

// Validate checks that the flags are valid.
func (f *Flags) Validate() error {
	// Validtors must have cosmos grpc services enabled.
	if !f.NonValidatingFullNode && !f.GrpcEnable {
		return fmt.Errorf("grpc.enable must be set to true - validating requires gRPC server")
	}

	// Grpc streaming
	if f.GrpcStreamingEnabled {
		if !f.GrpcEnable {
			return fmt.Errorf("grpc.enable must be set to true - grpc streaming requires gRPC server")
		}
		if f.GrpcStreamingMaxBatchSize == 0 {
			return fmt.Errorf("grpc streaming batch size must be positive number")
		}
		if f.GrpcStreamingFlushIntervalMs == 0 {
			return fmt.Errorf("grpc streaming flush interval must be positive number")
		}
		if f.GrpcStreamingMaxChannelBufferSize == 0 {
			return fmt.Errorf("grpc streaming channel size must be positive number")
		}
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
		DdErrorTrackingFormat: DefaultDdErrorTrackingFormat,

		// These are the default values from the Cosmos flags.
		GrpcAddress: config.DefaultGRPCAddress,
		GrpcEnable:  true,

		GrpcStreamingEnabled:              DefaultGrpcStreamingEnabled,
		GrpcStreamingFlushIntervalMs:      DefaultGrpcStreamingFlushIntervalMs,
		GrpcStreamingMaxBatchSize:         DefaultGrpcStreamingMaxBatchSize,
		GrpcStreamingMaxChannelBufferSize: DefaultGrpcStreamingMaxChannelBufferSize,

		VEOracleEnabled: true,
	}

	// Populate the flags if they exist.
	if option := appOpts.Get(NonValidatingFullNodeFlag); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.NonValidatingFullNode = v
		}
	}

	if option := appOpts.Get(DdAgentHost); option != nil {
		if v, err := cast.ToStringE(option); err == nil && len(v) > 0 {
			result.DdAgentHost = v
		}
	}

	if option := appOpts.Get(DdTraceAgentPort); option != nil {
		if v, err := cast.ToUint16E(option); err == nil {
			result.DdTraceAgentPort = v
		}
	}

	if option := appOpts.Get(DdErrorTrackingFormat); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.DdErrorTrackingFormat = v
		}
	}

	if option := appOpts.Get(GrpcAddress); option != nil {
		if v, err := cast.ToStringE(option); err == nil && len(v) > 0 {
			result.GrpcAddress = v
		}
	}

	if option := appOpts.Get(GrpcEnable); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.GrpcEnable = v
		}
	}

	if option := appOpts.Get(GrpcStreamingEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.GrpcStreamingEnabled = v
		}
	}

	if option := appOpts.Get(GrpcStreamingFlushIntervalMs); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.GrpcStreamingFlushIntervalMs = v
		}
	}

	if option := appOpts.Get(GrpcStreamingMaxBatchSize); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.GrpcStreamingMaxBatchSize = v
		}
	}

	if option := appOpts.Get(GrpcStreamingMaxChannelBufferSize); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.GrpcStreamingMaxChannelBufferSize = v
		}
	}

	if option := appOpts.Get(VEOracleEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.VEOracleEnabled = v
		}
	}

	return result
}
