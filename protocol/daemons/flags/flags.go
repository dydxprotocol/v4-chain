package flags

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// List of CLI flags for Server and Client.
const (
	// Flag names
	FlagUnixSocketAddress           = "unix-socket-address"
	FlagPanicOnDaemonFailureEnabled = "panic-on-daemon-failure-enabled"
	FlagMaxDaemonUnhealthySeconds   = "max-daemon-unhealthy-seconds"

	FlagPriceDaemonEnabled     = "price-daemon-enabled"
	FlagPriceDaemonLoopDelayMs = "price-daemon-loop-delay-ms"

	FlagBridgeDaemonEnabled        = "bridge-daemon-enabled"
	FlagBridgeDaemonLoopDelayMs    = "bridge-daemon-loop-delay-ms"
	FlagBridgeDaemonEthRpcEndpoint = "bridge-daemon-eth-rpc-endpoint"

	FlagLiquidationDaemonEnabled        = "liquidation-daemon-enabled"
	FlagLiquidationDaemonLoopDelayMs    = "liquidation-daemon-loop-delay-ms"
	FlagLiquidationDaemonQueryPageLimit = "liquidation-daemon-query-page-limit"

	FlagSlinkyDaemonEnabled = "slinky-daemon-enabled"
)

// Shared flags contains configuration flags shared by all daemons.
type SharedFlags struct {
	// SocketAddress is the location of the unix socket to communicate with the daemon gRPC service.
	SocketAddress string
	// PanicOnDaemonFailureEnabled toggles whether the daemon should panic on failure.
	PanicOnDaemonFailureEnabled bool
	// MaxDaemonUnhealthySeconds is the maximum allowable duration for which a daemon can be unhealthy.
	MaxDaemonUnhealthySeconds uint32
}

// BridgeFlags contains configuration flags for the Bridge Daemon.
type BridgeFlags struct {
	// Enabled toggles the bridge daemon on or off.
	Enabled bool
	// LoopDelayMs configures the update frequency of the bridge daemon.
	LoopDelayMs uint32
	// EthRpcEndpoint is the endpoint for the Ethereum node where bridge data is queried.
	EthRpcEndpoint string
}

// LiquidationFlags contains configuration flags for the Liquidation Daemon.
type LiquidationFlags struct {
	// Enabled toggles the liquidation daemon on or off.
	Enabled bool
	// LoopDelayMs configures the update frequency of the liquidation daemon.
	LoopDelayMs uint32
	// QueryPageLimit configures the pagination limit for fetching subaccounts.
	QueryPageLimit uint64
}

// PriceFlags contains configuration flags for the Price Daemon.
type PriceFlags struct {
	// Enabled toggles the price daemon on or off.
	Enabled bool
	// LoopDelayMs configures the update frequency of the price daemon.
	LoopDelayMs uint32
}

type SlinkyFlags struct {
	// Enabled toggles the slinky daemon on or off.
	Enabled bool
}

// DaemonFlags contains the collected configuration flags for all daemons.
type DaemonFlags struct {
	Shared      SharedFlags
	Bridge      BridgeFlags
	Liquidation LiquidationFlags
	Price       PriceFlags
	Slinky      SlinkyFlags
}

var defaultDaemonFlags *DaemonFlags

// GetDefaultDaemonFlags returns the default values for the Daemon Flags using a singleton pattern.
func GetDefaultDaemonFlags() DaemonFlags {
	if defaultDaemonFlags == nil {
		defaultDaemonFlags = &DaemonFlags{
			Shared: SharedFlags{
				SocketAddress:               "/tmp/daemons.sock",
				PanicOnDaemonFailureEnabled: true,
				MaxDaemonUnhealthySeconds:   5 * 60, // 5 minutes.
			},
			Bridge: BridgeFlags{
				Enabled:        true,
				LoopDelayMs:    30_000,
				EthRpcEndpoint: "",
			},
			Liquidation: LiquidationFlags{
				Enabled:        true,
				LoopDelayMs:    1_600,
				QueryPageLimit: 1_000,
			},
			Price: PriceFlags{
				Enabled:     true,
				LoopDelayMs: 3_000,
			},
			Slinky: SlinkyFlags{
				Enabled: false,
			},
		}
	}
	return *defaultDaemonFlags
}

// AddDaemonFlagsToCmd adds the required flags to instantiate a server and client for
// price updates. These flags should be applied to the `start` command V4 Cosmos application.
// E.g. `dydxprotocold start --price-daemon-enabled=true --unix-socket-address $(unix_socket_address)`
func AddDaemonFlagsToCmd(
	cmd *cobra.Command,
) {
	//
	df := GetDefaultDaemonFlags()

	// Shared Flags.
	cmd.Flags().String(
		FlagUnixSocketAddress,
		df.Shared.SocketAddress,
		"Socket address for the daemons to send updates to, if not set "+
			"will establish default location to ingest daemon updates from",
	)
	cmd.Flags().Bool(
		FlagPanicOnDaemonFailureEnabled,
		df.Shared.PanicOnDaemonFailureEnabled,
		"Enables panicking when a daemon fails.",
	)
	cmd.Flags().Uint32(
		FlagMaxDaemonUnhealthySeconds,
		df.Shared.MaxDaemonUnhealthySeconds,
		"Maximum allowable duration for which a daemon can be unhealthy.",
	)

	// Bridge Daemon.
	cmd.Flags().Bool(
		FlagBridgeDaemonEnabled,
		df.Bridge.Enabled,
		"Enable Bridge Daemon. Set to false for non-validator nodes.",
	)
	cmd.Flags().Uint32(
		FlagBridgeDaemonLoopDelayMs,
		df.Bridge.LoopDelayMs,
		"Delay in milliseconds between running the Bridge Daemon task loop.",
	)
	cmd.Flags().String(
		FlagBridgeDaemonEthRpcEndpoint,
		df.Bridge.EthRpcEndpoint,
		"Ethereum Node Rpc Endpoint",
	)

	// Liquidation Daemon.
	cmd.Flags().Bool(
		FlagLiquidationDaemonEnabled,
		df.Liquidation.Enabled,
		"Enable Liquidation Daemon. Set to false for non-validator nodes.",
	)
	cmd.Flags().Uint32(
		FlagLiquidationDaemonLoopDelayMs,
		df.Liquidation.LoopDelayMs,
		"Delay in milliseconds between running the Liquidation Daemon task loop.",
	)
	cmd.Flags().Uint64(
		FlagLiquidationDaemonQueryPageLimit,
		df.Liquidation.QueryPageLimit,
		"Limit on the number of items to fetch per query in the Liquidation Daemon task loop.",
	)

	// Price Daemon.
	cmd.Flags().Bool(
		FlagPriceDaemonEnabled,
		df.Price.Enabled,
		"Enable Price Daemon. Set to false for non-validator nodes.",
	)
	cmd.Flags().Uint32(
		FlagPriceDaemonLoopDelayMs,
		df.Price.LoopDelayMs,
		"Delay in milliseconds between sending price updates to the application.",
	)

	// Slinky Daemon.
	cmd.Flags().Bool(
		FlagSlinkyDaemonEnabled,
		df.Slinky.Enabled,
		"Enable Slinky Daemon. Set to false for non-validator nodes.",
	)
}

// GetDaemonFlagValuesFromOptions gets all daemon flag values from the `AppOptions` struct.
func GetDaemonFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) DaemonFlags {
	// Default value
	result := GetDefaultDaemonFlags()

	// Shared Flags
	if option := appOpts.Get(FlagUnixSocketAddress); option != nil {
		if v, err := cast.ToStringE(option); err == nil && len(v) > 0 {
			result.Shared.SocketAddress = v
		}
	}
	if option := appOpts.Get(FlagPanicOnDaemonFailureEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.Shared.PanicOnDaemonFailureEnabled = v
		}
	}
	if option := appOpts.Get(FlagMaxDaemonUnhealthySeconds); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.Shared.MaxDaemonUnhealthySeconds = v
		}
	}

	// Bridge Daemon.
	if option := appOpts.Get(FlagBridgeDaemonEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.Bridge.Enabled = v
		}
	}
	if option := appOpts.Get(FlagBridgeDaemonLoopDelayMs); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.Bridge.LoopDelayMs = v
		}
	}
	if option := appOpts.Get(FlagBridgeDaemonEthRpcEndpoint); option != nil {
		if v, err := cast.ToStringE(option); err == nil && len(v) > 0 {
			result.Bridge.EthRpcEndpoint = v
		}
	}

	// Liquidation Daemon.
	if option := appOpts.Get(FlagLiquidationDaemonEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.Liquidation.Enabled = v
		}
	}
	if option := appOpts.Get(FlagLiquidationDaemonLoopDelayMs); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.Liquidation.LoopDelayMs = v
		}
	}
	if option := appOpts.Get(FlagLiquidationDaemonQueryPageLimit); option != nil {
		if v, err := cast.ToUint64E(option); err == nil {
			result.Liquidation.QueryPageLimit = v
		}
	}

	// Price Daemon.
	if option := appOpts.Get(FlagPriceDaemonEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.Price.Enabled = v
		}
	}
	if option := appOpts.Get(FlagPriceDaemonLoopDelayMs); option != nil {
		if v, err := cast.ToUint32E(option); err == nil {
			result.Price.LoopDelayMs = v
		}
	}

	// Slinky Daemon.
	if option := appOpts.Get(FlagSlinkyDaemonEnabled); option != nil {
		if v, err := cast.ToBoolE(option); err == nil {
			result.Slinky.Enabled = v
		}
	}

	return result
}
