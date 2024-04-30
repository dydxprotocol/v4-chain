package app_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	evidencetypes "cosmossdk.io/x/evidence/types"
	feegranttypes "cosmossdk.io/x/feegrant"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govsimulation "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	exportedtypes "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	ratelimitmodule "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	rewardsmodule "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vestmodule "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
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
func NewSimApp(t testing.TB, appCreator func() *app.App) *SimApp {
	simApp := &SimApp{appCreator(), nil}
	baseapp.SetChainID(simChainId)(simApp.GetBaseApp())

	// TODO(CORE-682): Remove shutdown override hook once Cosmos SDK invokes it as part of simapp.
	t.Cleanup(func() {
		if err := simApp.App.Close(); err != nil {
			t.Fatal(err)
		}
	})
	return simApp
}

// Returns the simulation manager.
func (app *SimApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

var genesisModuleOrder = []string{
	authtypes.ModuleName,
	banktypes.ModuleName,
	authz.ModuleName,
	capabilitytypes.ModuleName,
	feegranttypes.ModuleName,
	govtypes.ModuleName,
	stakingtypes.ModuleName,
	distributiontypes.ModuleName,
	slashingtypes.ModuleName,
	paramstypes.ModuleName,
	exportedtypes.ModuleName,
	evidencetypes.ModuleName,
	ratelimitmodule.ModuleName,
	ibctransfertypes.ModuleName,
	icatypes.ModuleName,
	pricestypes.ModuleName,
	assetstypes.ModuleName,
	perpetualstypes.ModuleName,
	satypes.ModuleName,
	clobtypes.ModuleName,
	sendingtypes.ModuleName,
	vestmodule.ModuleName,
	rewardsmodule.ModuleName,
	epochstypes.ModuleName,
	blocktimetypes.ModuleName,
}

var skippedGenesisModules = map[string]interface{}{
	// Skip adding the interchain accounts module since the modules simulation
	// https://github.com/cosmos/ibc-go/blob/2551dea/modules/apps/27-interchain-accounts/simulation/proposals.go#L23
	// adds both ICA host and controller messages while the app only supports host messages causing the
	// simulation to fail due to unroutable controller messages.
	icatypes.ModuleName: nil,
}

// WithRandomlyGeneratedOperationsSimulationManager uses the default weighted operations of each of
// the modules which are currently using randomness to generate operations for simulation.
func (app *SimApp) WithRandomlyGeneratedOperationsSimulationManager() {
	// Find all simulation modules and replace the auth one with one that is needed for simulation.
	simAppModules := []module.AppModuleSimulation{}
	for _, genesisModule := range genesisModuleOrder {
		if _, skipped := skippedGenesisModules[genesisModule]; skipped {
			continue
		}

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
	if len(simAppModules) != len(foundSimAppModules)-len(skippedGenesisModules) {
		panic(fmt.Sprintf(
			"Under specified AppModuleSimulation genesis order. "+
				"Genesis order is %s with skipped modules %s but found modules %s.",
			genesisModuleOrder,
			skippedGenesisModules,
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

	b.Cleanup(func() {
		require.NoError(b, os.RemoveAll(dir))
	})

	appOptions := defaultAppOptionsForSimulation()

	dydxApp := NewSimApp(
		b,
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
		AppStateFn(dydxApp.AppCodec(), dydxApp.SimulationManager()),
		simtypes.RandomAccounts,
		CustomSimulationOperations(dydxApp, dydxApp.AppCodec(), config),
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

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(dir))
	})

	appOptions := defaultAppOptionsForSimulation()

	dydxApp := NewSimApp(
		t,
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
		AppStateFn(dydxApp.AppCodec(), dydxApp.SimulationManager()),
		simtypes.RandomAccounts,
		CustomSimulationOperations(dydxApp, dydxApp.AppCodec(), config),
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
	numTimesToRunPerSeed := 3
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	appOptions := defaultAppOptionsForSimulation()

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simcli.FlagVerboseValue {
				logger = log.NewTestLogger(t)
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			dydxApp := NewSimApp(
				t,
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
				AppStateFn(dydxApp.AppCodec(), dydxApp.SimulationManager()),
				simtypes.RandomAccounts,
				CustomSimulationOperations(dydxApp, dydxApp.AppCodec(), config),
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
	appOptions[daemonflags.FlagLiquidationDaemonEnabled] = false
	return appOptions
}

// Note: this content comes from github.com/cosmos/cosmos-sdk/simapp/state.go:
// https://github.com/cosmos/cosmos-sdk/blob/1e8e923d3174cdfdb42454a96c27251ad72b6504/simapp/state.go

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
func AppStateFn(cdc codec.JSONCodec, simManager *module.SimulationManager) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		if simcli.FlagGenesisTimeValue == 0 {
			genesisTimestamp = simtypes.RandTimestamp(r)
		} else {
			genesisTimestamp = time.Unix(simcli.FlagGenesisTimeValue, 0)
		}

		chainID = config.ChainID
		switch {
		case config.ParamsFile != "" && config.GenesisFile != "":
			panic("cannot provide both a genesis file and a params file")

		case config.GenesisFile != "":
			// override the default chain-id from app to set it later to the config
			genesisDoc, accounts := AppStateFromGenesisFileFn(r, cdc, config.GenesisFile)

			if simcli.FlagGenesisTimeValue == 0 {
				// use genesis timestamp if no custom timestamp is provided (i.e no random timestamp)
				genesisTimestamp = genesisDoc.GenesisTime
			}

			appState = genesisDoc.AppState
			chainID = genesisDoc.ChainID
			simAccs = accounts

		case config.ParamsFile != "":
			appParams := make(simtypes.AppParams)
			bz, err := os.ReadFile(config.ParamsFile)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bz, &appParams)
			if err != nil {
				panic(err)
			}
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)

		default:
			appParams := make(simtypes.AppParams)
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)
		}

		rawState := make(map[string]json.RawMessage)
		err := json.Unmarshal(appState, &rawState)
		if err != nil {
			panic(err)
		}

		stakingStateBz, ok := rawState[stakingtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		stakingState := new(stakingtypes.GenesisState)
		err = cdc.UnmarshalJSON(stakingStateBz, stakingState)
		if err != nil {
			panic(err)
		}
		// compute not bonded balance
		notBondedTokens := math.ZeroInt()
		for _, val := range stakingState.Validators {
			if val.Status != stakingtypes.Unbonded {
				continue
			}
			notBondedTokens = notBondedTokens.Add(val.GetTokens())
		}
		notBondedCoins := sdk.NewCoin(stakingState.Params.BondDenom, notBondedTokens)
		// edit bank state to make it have the not bonded pool tokens
		bankStateBz, ok := rawState[banktypes.ModuleName]
		// TODO(ignore - from CosmosSDK): should we panic in this case
		if !ok {
			panic("bank genesis state is missing")
		}
		bankState := new(banktypes.GenesisState)
		err = cdc.UnmarshalJSON(bankStateBz, bankState)
		if err != nil {
			panic(err)
		}

		stakingAddr := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()
		var found bool
		for _, balance := range bankState.Balances {
			if balance.Address == stakingAddr {
				found = true
				break
			}
		}
		if !found {
			bankState.Balances = append(bankState.Balances, banktypes.Balance{
				Address: stakingAddr,
				Coins:   sdk.NewCoins(notBondedCoins),
			})
		}

		// change appState back
		rawState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingState)
		rawState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)

		// replace appstate
		appState, err = json.Marshal(rawState)
		if err != nil {
			panic(err)
		}
		return appState, simAccs, chainID, genesisTimestamp
	}
}

// AppStateRandomizedFn creates calls each module's GenesisState generator function
// and creates the simulation params
func AppStateRandomizedFn(
	simManager *module.SimulationManager, r *rand.Rand, cdc codec.JSONCodec,
	accs []simtypes.Account, genesisTimestamp time.Time, appParams simtypes.AppParams,
) (json.RawMessage, []simtypes.Account) {
	numAccs := int64(len(accs))
	// TODO(ignore - from CosmosSDK)
	// in case runtime.RegisterModules(...) is used, the genesis state of the module won't be reflected here
	genesisState := basic_manager.ModuleBasics.DefaultGenesis(cdc)

	// generate a random amount of initial stake coins and a random initial
	// number of bonded accounts
	var (
		numInitiallyBonded int64
		initialStake       math.Int
	)
	appParams.GetOrGenerate(
		simtestutil.StakePerAccount, &initialStake, r,
		func(r *rand.Rand) {
			// Since the stake token denom has 18 decimals, the initial stake balance needs to be at least
			// 1e18 to be considered valid. However, in the current implementation of auth simulation logic
			// (https://github.com/dydxprotocol/cosmos-sdk/blob/93454d9f/x/auth/simulation/genesis.go#L38),
			// `initialStake` is casted to an `int64` value (max_int64 ~= 9.22e18).
			// As such today the only valid range of values for `initialStake` is [1e18, max_int64]. Note
			// this only represents 1~9 full coins.
			// TODO(DEC-2132): Make this value more realistic by allowing larger values.
			initialStake = math.NewInt(r.Int63n(8.22e18) + 1e18)
		},
	)

	appParams.GetOrGenerate(
		simtestutil.InitiallyBondedValidators, &numInitiallyBonded, r,
		func(r *rand.Rand) { numInitiallyBonded = int64(r.Intn(299) + 1) },
	)

	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}

	fmt.Printf(
		`Selected randomly generated parameters for simulated genesis:
{
  stake_per_account: "%s",
  initially_bonded_validators: "%d"
}
`, initialStake, numInitiallyBonded,
	)

	simState := &module.SimulationState{
		AppParams:    appParams,
		Cdc:          cdc,
		Rand:         r,
		GenState:     genesisState,
		Accounts:     accs,
		InitialStake: initialStake,
		NumBonded:    numInitiallyBonded,
		GenTimestamp: genesisTimestamp,
		BondDenom:    sdk.DefaultBondDenom,
	}

	simManager.GenerateGenesisStates(simState)

	appState, err := json.Marshal(genesisState)
	if err != nil {
		panic(err)
	}

	return appState, accs
}

// AppStateFromGenesisFileFn util function to generate the genesis AppState
// from a genesis.json file.
func AppStateFromGenesisFileFn(
	r io.Reader,
	cdc codec.JSONCodec,
	genesisFile string,
) (tmtypes.GenesisDoc, []simtypes.Account) {
	bytes, err := os.ReadFile(genesisFile)
	if err != nil {
		panic(err)
	}

	var genesis tmtypes.GenesisDoc
	// NOTE: Tendermint uses a custom JSON decoder for GenesisDoc
	err = tmjson.Unmarshal(bytes, &genesis)
	if err != nil {
		panic(err)
	}

	var appState app.GenesisState
	err = json.Unmarshal(genesis.AppState, &appState)
	if err != nil {
		panic(err)
	}

	var authGenesis authtypes.GenesisState
	if appState[authtypes.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[authtypes.ModuleName], &authGenesis)
	}

	newAccs := make([]simtypes.Account, len(authGenesis.Accounts))
	for i, acc := range authGenesis.Accounts {
		// Pick a random private key, since we don't know the actual key
		// This should be fine as it's only used for mock Tendermint validators
		// and these keys are never actually used to sign by mock Tendermint.
		privkeySeed := make([]byte, 15)
		if _, err := r.Read(privkeySeed); err != nil {
			panic(err)
		}

		privKey := secp256k1.GenPrivKeyFromSecret(privkeySeed)

		a, ok := acc.GetCachedValue().(sdk.AccountI)
		if !ok {
			panic("expected account")
		}

		// create simulator accounts
		simAcc := simtypes.Account{PrivKey: privKey, PubKey: privKey.PubKey(), Address: a.GetAddress()}
		newAccs[i] = simAcc
	}

	return genesis, newAccs
}

// CustomSimulationOperations initializes the custom simulation params and
// returns all the modules weighted operations.
func CustomSimulationOperations(
	app runtime.AppI,
	cdc codec.JSONCodec,
	config simtypes.Config,
) []simtypes.WeightedOperation {
	simState := module.SimulationState{
		AppParams: make(simtypes.AppParams),
		Cdc:       cdc,
		TxConfig:  moduletestutil.MakeTestTxConfig(),
		BondDenom: sdk.DefaultBondDenom,
	}

	// Set the weight of MsgCancelProposal to zero.
	b, err := json.Marshal(0)
	if err != nil {
		panic("Failed to marshal operation weights")
	}
	simState.AppParams[govsimulation.OpWeightMsgCancelProposal] = b

	//nolint:staticcheck // used for legacy testing
	simState.LegacyProposalContents = app.SimulationManager().GetProposalContents(simState)
	simState.ProposalMsgs = app.SimulationManager().GetProposalMsgs(simState)
	return app.SimulationManager().WeightedOperations(simState)
}
