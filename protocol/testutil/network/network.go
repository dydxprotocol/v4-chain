package network

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"cosmossdk.io/log"
	pruningtypes "cosmossdk.io/store/pruning/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	v4module "github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib/marketmap"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/ci"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/gofrs/flock"
	"github.com/stretchr/testify/require"
)

type (
	Network = network.Network
	Config  = network.Config
)

var fileLock = flock.New("/tmp/test-cosmos-network.lock")

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t testing.TB, configs ...network.Config) *network.Network {
	// This is a workaround for an issue in the cosmos-sdk `testutil/network` package.
	// Specifically, the `testutil/network` package attempts to use a package-level lock to ensure that only one
	// test network is running at a time. This is problematic when running tests in parallel, as `go test`
	// will utilize a separate instance of `testutil/network` package for each test, and therefore the package-level lock
	// is not shared between the tests.
	// This is compounded by the fact that the `testutil/network` package uses a method called `FreeTCPAddr`,
	// which will attempt to start a listener on a random open port on the local machine.
	// It does this by calling `net.Listen` on `localhost:0`, and then immediately closing the listener.
	// It then uses the resulting ports to start multiple servers on the test network.
	// When multiple `Network.New` calls are made in parallel, it is possible that the `FreeTCPAddr` method will
	// return the same port for multiple calls, as the port remains open between the call to `FreeTCPAddr` and starting
	// up the gRPC server. The port contention is generally uncommon on Mac OS as the OS will assign monotonically
	// increasing ports to new listeners, but the implementation on Alpine Linux is different, and more commonly
	// returns the same ports open ports when invoked repeatedly.
	// For this reason, we use a file lock to ensure that only one `Network.New` call is made at a time in CI to help
	// prevent flaky tests. We do this only in CI, as it is not necessary when running tests locally on Mac OS.
	if ci.IsRunningOnGithubActions() {
		err := fileLock.Lock()
		require.NoError(t, err)
		defer func() {
			err := fileLock.Unlock()
			require.NoError(t, err)
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			for range c {
				fileLock.Close()
				return
			}
		}()
	}

	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = network.DefaultConfig(NewTestNetworkFixture)
	} else {
		cfg = configs[0]
	}

	// TODO(CORE-682): Remove shutdown override hook once Cosmos SDK invokes it as part of network#Cleanup.
	appConstructor := cfg.AppConstructor
	cfg.AppConstructor = func(val network.ValidatorI) servertypes.Application {
		app := appConstructor(val)
		t.Cleanup(func() {
			if err := app.Close(); err != nil {
				panic(err)
			}
		})
		return app
	}

	net, err := network.New(t, t.TempDir(), cfg)
	if err != nil {
		panic(err)
	}

	t.Cleanup(net.Cleanup)
	return net
}

// NetworkConfigOptions represent configuration options that can be passed to the `DefaultConfig` method
// when creating a test network.
type NetworkConfigOptions struct {
	// AppOptions allows for overriding the default AppOptions in the configuration.
	// This is useful for overriding things like command line flags which normally exist in this struct.
	AppOptions servertypes.AppOptions
	// OnNewApp is a function invoked immediately before the test network calls `app.New` to initialize
	// the application. It is called with the in-process validator running the application.
	OnNewApp func(val network.ValidatorI)
}

// DefaultConfig will initialize config for the network with custom application,
// genesis and single validator. All other parameters are inherited from cosmos-sdk/testutil/network.DefaultConfig
// If options are provided as nil, no custom options will be utilized.
func DefaultConfig(options *NetworkConfigOptions) network.Config {
	var appOptions servertypes.AppOptions
	var onNewApp func(val network.ValidatorI)

	// Parse options.
	if options == nil {
		appOptions = appoptions.GetDefaultTestAppOptions("", nil)
		onNewApp = func(val network.ValidatorI) {}
	} else {
		appOptions = options.AppOptions
		onNewApp = options.OnNewApp
	}

	encoding := app.GetEncodingConfig()

	// Inject default market map genesis based off of prices genesis
	cdc := codec.NewProtoCodec(v4module.InterfaceRegistry)
	genesisState := basic_manager.ModuleBasics.DefaultGenesis(encoding.Codec)

	var pricesGenesisState pricestypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[pricestypes.ModuleName], &pricesGenesisState)
	marketMap, err := marketmap.ConstructMarketMapFromParams(pricesGenesisState.MarketParams)
	if err != nil {
		panic(err)
	}
	marketmapGenesis := marketmaptypes.GenesisState{
		MarketMap: marketMap,
		Params:    marketmaptypes.DefaultParams(),
	}
	genesisState[marketmaptypes.ModuleName] = cdc.MustMarshalJSON(&marketmapGenesis)

	return network.Config{
		Codec:             encoding.Codec,
		TxConfig:          encoding.TxConfig,
		LegacyAmino:       encoding.Amino,
		InterfaceRegistry: encoding.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.ValidatorI) servertypes.Application {
			onNewApp(val)
			if appOptions.Get(flags.FlagHome) == "" || appOptions.Get(flags.FlagHome) == nil {
				appOptions.(*appoptions.FakeAppOptions).Set(flags.FlagHome, val.GetCtx().Config.RootDir)
			}

			return app.New(
				val.GetCtx().Logger,
				dbm.NewMemDB(),
				nil,
				true,
				appOptions,
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
				baseapp.SetChainID("dydxprotocol"),
			)
		},
		GenesisState:    genesisState,
		TimeoutCommit:   2 * time.Second,
		ChainID:         "dydxprotocol",
		NumValidators:   1,
		BondDenom:       sdk.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: pruningtypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}
}

// NewTestNetworkFixture returns a new simapp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	appOptions := appoptions.GetDefaultTestAppOptionsFromTempDirectory("", nil)
	dydxApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		appOptions,
	)

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return app.New(
			val.GetCtx().Logger,
			dbm.NewMemDB(),
			nil,
			true,
			appOptions,
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   dydxApp.DefaultGenesis(),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: dydxApp.InterfaceRegistry(),
			Codec:             dydxApp.AppCodec(),
			TxConfig:          dydxApp.TxConfig(),
			Amino:             dydxApp.LegacyAmino(),
		},
	}
}
