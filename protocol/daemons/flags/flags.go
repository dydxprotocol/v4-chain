package flags

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// List of CLI flags for Server and Client.
const (
	// Flag names
	FlagUnixSocketAddress = "unix-socket-address"

	FlagPriceDaemonEnabled     = "price-daemon-enabled"
	FlagPriceDaemonLoopDelayMs = "price-daemon-loop-delay-ms"

	FlagBridgeDaemonEnabled        = "bridge-daemon-enabled"
	FlagBridgeDaemonLoopDelayMs    = "bridge-daemon-loop-delay-ms"
	FlagBridgeDaemonEthRpcEndpoint = "bridge-daemon-eth-rpc-endpoint"

	FlagLiquidationDaemonEnabled             = "liquidation-daemon-enabled"
	FlagLiquidationDaemonLoopDelayMs         = "liquidation-daemon-loop-delay-ms"
	FlagLiquidationDaemonSubaccountPageLimit = "liquidation-daemon-subaccount-page-limit"
	FlagLiquidationDaemonRequestChunkSize    = "liquidation-daemon-request-chunk-size"
)

type SharedFlags struct {
	SocketAddress string
}

type BridgeFlags struct {
	Enabled        bool
	LoopDelayMs    uint32
	EthRpcEndpoint string
}

type LiquidationFlags struct {
	Enabled             bool
	LoopDelayMs         uint32
	SubaccountPageLimit uint64
	RequestChunkSize    uint64
}

type PriceFlags struct {
	Enabled     bool
	LoopDelayMs uint32
}
type DaemonFlags struct {
	Shared      SharedFlags
	Bridge      BridgeFlags
	Liquidation LiquidationFlags
	Price       PriceFlags
}

var defaultDaemonFlags *DaemonFlags

// Returns the default values for the Daemon Flags using a singleton pattern.
func GetDefaultDaemonFlags() DaemonFlags {
	if defaultDaemonFlags == nil {
		defaultDaemonFlags = &DaemonFlags{
			Shared: SharedFlags{
				SocketAddress: "/tmp/daemons.sock",
			},
			Bridge: BridgeFlags{
				Enabled:        true,
				LoopDelayMs:    30_000,
				EthRpcEndpoint: "https://eth-sepolia.g.alchemy.com/v2/demo",
			},
			Liquidation: LiquidationFlags{
				Enabled:             true,
				LoopDelayMs:         1_600,
				SubaccountPageLimit: 1_000,
				RequestChunkSize:    500,
			},
			Price: PriceFlags{
				Enabled:     true,
				LoopDelayMs: 3_000,
			},
		}
	}
	return *defaultDaemonFlags
}

// AddSharedPriceFeedFlagsToCmd adds the required flags to instantiate a server and client for
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
		"Socket address for the price daemon to send updates to, if not set "+
			"will establish default location to ingest price updates from",
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
		FlagLiquidationDaemonSubaccountPageLimit,
		df.Liquidation.SubaccountPageLimit,
		"Limit on the number of subaccounts to fetch per query in the Liquidation Daemon task loop.",
	)
	cmd.Flags().Uint64(
		FlagLiquidationDaemonRequestChunkSize,
		df.Liquidation.RequestChunkSize,
		"Limit on the number of subaccounts per collateralization check in the Liquidation Daemon task loop.",
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
}

// GetDaemonFlagValuesFromOptions gets all daemon flag values from the `AppOptions` struct.
func GetDaemonFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) DaemonFlags {
	// Default value
	result := GetDefaultDaemonFlags()

	// Shared Flags
	if v, ok := appOpts.Get(FlagUnixSocketAddress).(string); ok {
		result.Shared.SocketAddress = v
	}

	// Bridge Daemon.
	if v, ok := appOpts.Get(FlagBridgeDaemonEnabled).(bool); ok {
		result.Bridge.Enabled = v
	}
	if v, ok := appOpts.Get(FlagBridgeDaemonLoopDelayMs).(uint32); ok {
		result.Bridge.LoopDelayMs = v
	}
	if v, ok := appOpts.Get(FlagBridgeDaemonEthRpcEndpoint).(string); ok {
		result.Bridge.EthRpcEndpoint = v
	}

	// Liquidation Daemon.
	if v, ok := appOpts.Get(FlagLiquidationDaemonEnabled).(bool); ok {
		result.Liquidation.Enabled = v
	}
	if v, ok := appOpts.Get(FlagLiquidationDaemonLoopDelayMs).(uint32); ok {
		result.Liquidation.LoopDelayMs = v
	}
	if v, ok := appOpts.Get(FlagLiquidationDaemonSubaccountPageLimit).(uint64); ok {
		result.Liquidation.SubaccountPageLimit = v
	}
	if v, ok := appOpts.Get(FlagLiquidationDaemonRequestChunkSize).(uint64); ok {
		result.Liquidation.RequestChunkSize = v
	}

	// Price Daemon.
	if v, ok := appOpts.Get(FlagPriceDaemonEnabled).(bool); ok {
		result.Price.Enabled = v
	}
	if v, ok := appOpts.Get(FlagPriceDaemonLoopDelayMs).(uint32); ok {
		result.Price.LoopDelayMs = v
	}

	return result
}
