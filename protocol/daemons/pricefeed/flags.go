package pricefeed

import (
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/dydxprotocol/v4/daemons/constants"
	"github.com/dydxprotocol/v4/lib"
	"github.com/spf13/cobra"
)

const (
	// MaxPriceAge defines the duration in which a price update is valid for.
	MaxPriceAge = time.Duration(30_000_000_000) // 30 sec, duration uses nanoseconds.
)

// List of CLI flags for Server and Client.
const (
	FlagPriceFeedUnixSocketAddr          = "pricefeed-unixsocketaddress"
	FlagPriceFeedEnabled                 = "pricefeed-enabled"
	FlagPriceFeedPriceUpdaterLoopDelayMs = "pricefeed-price-updater-loop-delay-ms"

	DefaultFlagPriceFeedUnixSocketAddrValue          = constants.DaemonSocketAddr
	DefaultFlagPriceFeedEnabledValue                 = true
	DefaultFlagPriceFeedPriceUpdaterLoopDelayMsValue = 3000

	GrpcAddress = "grpc.address"
)

// AddSharedPriceFeedFlagsToCmd adds the required flags to instantiate a server and client for
// price updates. These flags should be applied to the `start` command V4 Cosmos application.
// E.g. `dydxprotocold start --pricefeed-enabled=true --pricefeed-unixsocketaddress $(pricefeed socket address)`
func AddSharedPriceFeedFlagsToCmd(cmd *cobra.Command) {
	cmd.
		Flags().
		String(
			FlagPriceFeedUnixSocketAddr,
			DefaultFlagPriceFeedUnixSocketAddrValue,
			"Socket address for the price daemon to send updates to, if not set "+
				"will establish default location to ingest price updates from",
		)
	cmd.
		Flags().
		Bool(
			FlagPriceFeedEnabled,
			DefaultFlagPriceFeedEnabledValue,
			"Enable client and server for ingesting pricefeed updates, set to false for non-validator nodes",
		)
}

// AddClientPriceFeedFlagsToCmd defines the required command line flags to instantiate the client
// for price updates that are not also required for a server. These flags should be passed to the `start`
// command of the Cosmos application.
// E.g. `dydxprotocold start --pricefeed-exchange-config-file $(pricefeed config file path)
// --pricefeed-price-updater-loop-delay-ms=3000
func AddClientPriceFeedFlagsToCmd(cmd *cobra.Command) {
	cmd.
		Flags().
		Int(
			FlagPriceFeedPriceUpdaterLoopDelayMs,
			DefaultFlagPriceFeedPriceUpdaterLoopDelayMsValue,
			"Delay in milliseconds between sending price updates to the application",
		)
}

// GetServerPricefeedFlagValuesFromOptions gets values for creating a price reporting environment
// struct from the `AppOptions` struct which contains values from price feed command-line flags.
func GetServerPricefeedFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) (
	pricefeedEnabled bool,
	pricefeedUnixSocketAddress string,
) {
	pricefeedEnabled, ok := appOpts.Get(FlagPriceFeedEnabled).(bool)
	if !ok {
		pricefeedEnabled = DefaultFlagPriceFeedEnabledValue
	}

	pricefeedUnixSocketAddress, ok = appOpts.Get(FlagPriceFeedUnixSocketAddr).(string)
	if !ok {
		pricefeedUnixSocketAddress = DefaultFlagPriceFeedUnixSocketAddrValue
	}

	return pricefeedEnabled,
		pricefeedUnixSocketAddress
}

// GetClientPricefeedFlagValuesFromOptions gets values for creating a price reporting environment
// client struct from the `AppOptions` struct which contains values from price feed command-line
// flags.
func GetClientPricefeedFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) (
	pricefeedEnabled bool,
	pricefeedUnixSocketAddress string,
	priceUpdaterLoopDelayMs uint32,
) {
	pricefeedEnabled, pricefeedUnixSocketAddress = GetServerPricefeedFlagValuesFromOptions(appOpts)

	unconvertedPriceUpdaterLoopDelayMs, ok := appOpts.Get(FlagPriceFeedPriceUpdaterLoopDelayMs).(int)
	if !ok {
		unconvertedPriceUpdaterLoopDelayMs = DefaultFlagPriceFeedPriceUpdaterLoopDelayMsValue
	}
	priceUpdaterLoopDelayMs = lib.MustConvertIntegerToUint32(unconvertedPriceUpdaterLoopDelayMs)

	return pricefeedEnabled,
		pricefeedUnixSocketAddress,
		priceUpdaterLoopDelayMs
}

// GetGrpcServerAddress gets the gRPC server host and port, which is used by daemons that need to
// query application state. If the value is not defined in the config, it will default to
// `DefaultGrpcEndpoint`.
func GetGrpcServerAddress(
	appOpts servertypes.AppOptions,
) string {
	grpcServerAddress, ok := appOpts.Get(GrpcAddress).(string)
	if !ok {
		grpcServerAddress = constants.DefaultGrpcEndpoint
	}

	return grpcServerAddress
}
