package app_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	feegranttypes "cosmossdk.io/x/feegrant"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	exportedtypes "github.com/cosmos/ibc-go/v7/modules/core/exported"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	rewardsmodule "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vestmodule "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/stretchr/testify/require"
)

var (
	simChainId                         = "v4-app-sim"
	_          runtime.AppI            = (*SimApp)(nil)
	_          servertypes.Application = (*SimApp)(nil)
)

func init() {
	simcli.GetSimulatorFlags()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

// An application that uses simulated operations.
type SimApp struct {
	*app.App

	// the simulation manager
	sm *module.SimulationManager
}

// Constructs an application capable of simulation using the appCreator
func NewSimApp(appCreator func() *app.App) *SimApp {
	simApp := &SimApp{appCreator(), nil}
	baseapp.SetChainID(simChainId)(simApp.GetBaseApp())
	return simApp
}

// Returns the simulation manager.
func (app *SimApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

var genesisModuleOrder = []string{
	authtypes.ModuleName,
	banktypes.ModuleName,
	capabilitytypes.ModuleName,
	feegranttypes.ModuleName,
	govtypes.ModuleName,
	stakingtypes.ModuleName,
	distributiontypes.ModuleName,
	slashingtypes.ModuleName,
	paramstypes.ModuleName,
	exportedtypes.ModuleName,
	ibctransfertypes.ModuleName,
	pricestypes.ModuleName,
	assetstypes.ModuleName,
	perpetualstypes.ModuleName,
	satypes.ModuleName,
	clobtypes.ModuleName,
	sendingtypes.ModuleName,
	vestmodule.ModuleName,
	rewardsmodule.ModuleName,
	epochstypes.ModuleName,
}

// WithRandomlyGeneratedOperationsSimulationManager uses the default weighted operations of each of
// the modules which are currently using randomness to generate operations for simulation.
func (app *SimApp) WithRandomlyGeneratedOperationsSimulationManager() {
	// Find all simulation modules and replace the auth one with one that is needed for simulation.
	simAppModules := []module.AppModuleSimulation{}
	for _, genesisModule := range genesisModuleOrder {
		if simAppModule, ok := app.ModuleManager.Modules[genesisModule].(module.AppModuleSimulation); ok {
			// Replace the auth module so that it generates some random accounts.
			if simAppModule.(module.AppModule).Name() == authtypes.ModuleName {
				authSubspace, _ := app.ParamsKeeper.GetSubspace(authtypes.ModuleName)
				simAppModules = append(simAppModules, auth.NewAppModule(
					app.AppCodec(),
					app.AccountKeeper,
					authsims.RandomGenesisAccounts,
					authSubspace,
				))
			} else {
				simAppModules = append(simAppModules, simAppModule)
			}
		} else {
			panic("Unable to find AppModuleSimulation " + genesisModule)
		}
	}
	foundSimAppModules := []string{}
	for _, appModule := range app.ModuleManager.Modules {
		if simAppModule, ok := appModule.(module.AppModuleSimulation); ok {
			foundSimAppModules = append(foundSimAppModules, simAppModule.(module.AppModuleBasic).Name())
		}
	}
	if len(simAppModules) != len(foundSimAppModules) {
		panic(fmt.Sprintf(
			"Under specified AppModuleSimulation genesis order. "+
				"Genesis order is %s but found modules %s.",
			genesisModuleOrder,
			foundSimAppModules,
		))
	}
	// Create the simulation manager and define the order of the modules for deterministic simulations.
	app.sm = module.NewSimulationManager(simAppModules...)
	app.sm.RegisterStoreDecoders()
}

// BenchmarkSimulation run the chain simulation.
// Copied from:
// https://github.com/cosmos/cosmos-sdk/blob/1e8e923d3174cdfdb42454a96c27251ad72b6504/simapp/sim_bench_test.go#L21
func BenchmarkFullAppSimulation(b *testing.B) {
	b.ReportAllocs()

	config := simcli.NewConfigFromFlags()
	config.ChainID = simChainId

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"goleveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	if err != nil {
		b.Fatalf("simulation setup failed: %s", err.Error())
	}

	if skip {
		b.Skip("skipping benchmark application simulation")
	}

	defer func() {
		require.NoError(b, db.Close())
		require.NoError(b, os.RemoveAll(dir))
	}()

	appOptions := defaultAppOptionsForSimulation()

	dydxApp := NewSimApp(
		func() *app.App {
			return app.New(
				logger,
				db,
				nil,
				true,
				appOptions,
				interBlockCacheOpt(),
			)
		})
	dydxApp.WithRandomlyGeneratedOperationsSimulationManager()

	// Note: While our app does not use the `vesting` module, the `auth` module still attempts to create
	// vesting accounts during simulation here:
	// https://github.com/dydxprotocol/cosmos-sdk/blob/dydx-fork-v0.47.0-rc2/x/auth/simulation/genesis.go#L26
	// For this reason, we need to register the `vesting` module interfaces so that the Genesis state of `auth` can be
	// marshaled properly.
	vestingtypes.RegisterInterfaces(dydxApp.InterfaceRegistry())

	// Run randomized simulations
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		dydxApp.GetBaseApp(),
		app.AppStateFn(dydxApp.AppCodec(), dydxApp.SimulationManager()),
		simtypes.RandomAccounts,
		simtestutil.SimulationOperations(dydxApp, dydxApp.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		dydxApp.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	if err = simtestutil.CheckExportSimulation(dydxApp, config, simParams); err != nil {
		b.Fatal(err)
	}

	if simErr != nil {
		b.Fatal(simErr)
	}

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

// TestFullAppSimulation was largely copied and modified from the `simapp` package in Cosmos SDK here:
// https://github.com/cosmos/cosmos-sdk/blob/f08ba9eafaacbc710d3211434a82c1828c57687b/simapp/sim_test.go#L68
func TestFullAppSimulation(t *testing.T) {
	config := simcli.NewConfigFromFlags()
	config.ChainID = simChainId

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"leveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	}()

	appOptions := defaultAppOptionsForSimulation()

	dydxApp := NewSimApp(
		func() *app.App {
			return app.New(
				logger,
				db,
				nil,
				true,
				appOptions,
			)
		})
	dydxApp.WithRandomlyGeneratedOperationsSimulationManager()
	require.Equal(t, "dydxprotocol", dydxApp.Name())

	// Note: While our app does not use the `vesting` module, the `auth` module still attempts to create
	// vesting accounts during simulation here:
	// https://github.com/dydxprotocol/cosmos-sdk/blob/dydx-fork-v0.47.0-rc2/x/auth/simulation/genesis.go#L26
	// For this reason, we need to register the `vesting` module interfaces so that the Genesis state of `auth` can be
	// marshaled properly.
	vestingtypes.RegisterInterfaces(dydxApp.InterfaceRegistry())

	// run randomized simulation
	_, _, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		dydxApp.GetBaseApp(),
		app.AppStateFn(dydxApp.AppCodec(), dydxApp.SimulationManager()),
		simtypes.RandomAccounts,
		simtestutil.SimulationOperations(dydxApp, dydxApp.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		dydxApp.AppCodec(),
	)
	require.NoError(t, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

// TestAppStateDeterminism was largely copied and modified from the `simapp` package in Cosmos SDK here:
// https://github.com/cosmos/cosmos-sdk/blob/1e8e923d3174cdfdb42454a96c27251ad72b6504/simapp/sim_test.go#L316
func TestAppStateDeterminism(t *testing.T) {
	if !simcli.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simcli.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = simChainId

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	appOptions := defaultAppOptionsForSimulation()

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simcli.FlagVerboseValue {
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			dydxApp := NewSimApp(
				func() *app.App {
					return app.New(
						logger,
						db,
						nil,
						true,
						appOptions,
						interBlockCacheOpt(),
					)
				})
			dydxApp.WithRandomlyGeneratedOperationsSimulationManager()

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			// Note: While our app does not use the `vesting` module, the `auth` module still attempts to create
			// vesting accounts during simulation here:
			// https://github.com/dydxprotocol/cosmos-sdk/blob/dydx-fork-v0.47.0-rc2/x/auth/simulation/genesis.go#L26
			// For this reason, we need to register the `vesting` module interfaces so that the Genesis state of `auth` can be
			// marshaled properly.
			vestingtypes.RegisterInterfaces(dydxApp.InterfaceRegistry())

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				dydxApp.GetBaseApp(),
				app.AppStateFn(dydxApp.AppCodec(), dydxApp.SimulationManager()),
				simtypes.RandomAccounts,
				simtestutil.SimulationOperations(dydxApp, dydxApp.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				dydxApp.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simtestutil.PrintStats(db)
			}

			appHash := dydxApp.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}

func defaultAppOptionsForSimulation() simtestutil.AppOptionsMap {
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue
	appOptions[daemonflags.FlagPriceDaemonEnabled] = false
	appOptions[daemonflags.FlagBridgeDaemonEnabled] = false
	return appOptions
}
