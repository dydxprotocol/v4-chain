package appoptions

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server/config"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"os"

	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
)

// FakeAppOptions is a helper struct used for creating `servertypes.AppOptions` for simulator and end-to-end testing.
// This struct allows for customizing the `servertypes.AppOptions` value that is normally supplied from CLI arguments
// to `dydxprotocold`. The real concrete implementation for this interface is in the "viper" package which is used
// under the hood by "cobra", which is the CLI framework used by Cosmos SDK.
type FakeAppOptions struct {
	options map[string]interface{}
}

func NewFakeAppOptions() *FakeAppOptions {
	return &FakeAppOptions{
		options: make(map[string]interface{}),
	}
}

func (fao *FakeAppOptions) Set(option string, value interface{}) {
	fao.options[option] = value
}

// Get implements the `servertypes.AppOptions` interface.
func (fao *FakeAppOptions) Get(o string) interface{} {
	value, ok := fao.options[o]
	if !ok {
		return nil
	}

	return value
}

// GetDefaultTestAppOptions returns a default set of AppOptions with the daemons disabled for end-to-end
// and simulator testing.
func GetDefaultTestAppOptions(homePath string, customFlags map[string]interface{}) servertypes.AppOptions {
	fao := NewFakeAppOptions()

	fao.Set(flags.FlagHome, homePath)

	// Disable the Price Daemon for all end-to-end and integration tests by default.
	fao.Set(daemonflags.FlagPriceDaemonEnabled, false)

	// Disable the Slinky Daemon for all end-to-end and integration tests by default.
	fao.Set(daemonflags.FlagOracleEnabled, false)

	// Disable Slinky Metrics for all end-to-end and integration tests by default.
	fao.Set(daemonflags.FlagOracleMetricsEnabled, false)

	// Disable the Bridge Daemon for all end-to-end and integration tests by default.
	fao.Set(daemonflags.FlagBridgeDaemonEnabled, false)

	// Disable the Liquidation Daemon for all end-to-end and integration tests by default.
	fao.Set(daemonflags.FlagLiquidationDaemonEnabled, false)

	// Populate the default value for gRPC.
	fao.Set(appflags.GrpcAddress, config.DefaultGRPCAddress)

	for flag, value := range customFlags {
		fao.Set(flag, value)
	}
	return fao
}

func GetDefaultTestAppOptionsFromTempDirectory(
	homePath string,
	customFlags map[string]interface{},
) servertypes.AppOptions {
	dir, err := os.MkdirTemp(homePath, "testapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)
	return GetDefaultTestAppOptions(".", customFlags)
}
