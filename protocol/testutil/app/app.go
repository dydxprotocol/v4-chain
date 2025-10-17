package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	listingtypes "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	dbm "github.com/cosmos/cosmos-db"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/mempool"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	sdkproto "github.com/cosmos/gogoproto/proto"
	marketmapmoduletypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	appconstants "github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/cmd/dydxprotocold/cmd"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testlog "github.com/dydxprotocol/v4-chain/protocol/testutil/logger"
	aptypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	govplus "github.com/dydxprotocol/v4-chain/protocol/x/govplus/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	vesttypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

// localdydxprotocol Alice config/priv_validator_key.json.
const alicePrivValidatorKeyJson = `{
  "address": "124B880684400B4C0086BD4EE882DCC5B61CF7E3",
  "pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "YiARx8259Z+fGFUxQLrz/5FU2RYRT6f5yzvt7D7CrQM="
  },
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "65frslxv5ig0KSNKlJOHT2FKTkOzkb/66eDPsiBaNUtiIBHHzbn1n58YVTFAuvP/kVTZFhFPp/nLO+3sPsKtAw=="
  }
}
`

// localdydxprotocol Alice config/node_key.json.
const aliceNodeKeyJson = `{
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "8EGQBxfGMcRfH0C45UTedEG5Xi3XAcukuInLUqFPpskjp1Ny0c5XvwlKevAwtVvkwoeYYQSe0geQG/cF3GAcUA=="
  }
}
`

// MustMakeCheckTxOptions is a struct containing options for MustMakeCheckTx.* functions.
type MustMakeCheckTxOptions struct {
	// AccAddressForSigning is the account that's used to sign the transaction.
	AccAddressForSigning string
	// AccSequenceNumberForSigning is the account sequence number that's used to sign the transaction.
	AccSequenceNumberForSigning uint64
	// Amount of Gas for the transaction.
	Gas uint64
	// Gas fees offered for the transaction.
	FeeAmt sdk.Coins
}

// ValidateResponsePrepareProposal is a function that validates the response from the PrepareProposalHandler.
type ValidateResponsePrepareProposalFn func(sdk.Context, abcitypes.ResponsePrepareProposal) (haltChain bool)

// ValidateResponseProcessProposal is a function that validates the response from the ProcessProposalHandler.
type ValidateResponseProcessProposalFn func(sdk.Context, abcitypes.ResponseProcessProposal) (haltChain bool)

// ValidateFinalizeBlockFn is a function that validates the response from finalizing the block.
type ValidateFinalizeBlockFn func(
	ctx sdk.Context,
	request abcitypes.RequestFinalizeBlock,
	response abcitypes.ResponseFinalizeBlock,
) (haltchain bool)

// AdvanceToBlockOptions is a struct containing options for AdvanceToBlock.* functions.
type AdvanceToBlockOptions struct {
	// The time associated with the block. If left at the default value then block time will be left unchanged.
	BlockTime time.Time

	// Whether to increment the block time using linear interpolation among the blocks.
	// TODO(DEC-2156): Instead of an option, pass in a `BlockTimeFunc` to map each block to a
	// time giving user greater flexibility.
	LinearBlockTimeInterpolation bool

	// RequestPrepareProposalTxsOverride allows overriding the txs that gets passed into the
	// PrepareProposalHandler. This is useful for testing scenarios where unintended msg txs
	// end up in the mempool (i.e. CheckTx failed to filter bad msg txs out).
	RequestPrepareProposalTxsOverride [][]byte

	// RequestProcessProposalTxsOverride allows overriding the txs that gets passed into the
	// ProcessProposalHandler. This is useful for testing scenarios where bad validators end
	// up proposing an invalid block proposal.
	RequestProcessProposalTxsOverride [][]byte

	// DeliverTxsOverride allows overriding the TestApp from being the block proposer and
	// allows simulating transactions that were agreed to upon consensus to be delivered.
	// This skips the PrepareProposal and ProcessProposal phase.
	DeliverTxsOverride [][]byte

	ValidateRespPrepare   ValidateResponsePrepareProposalFn
	ValidateRespProcess   ValidateResponseProcessProposalFn
	ValidateFinalizeBlock ValidateFinalizeBlockFn
}

// DefaultTestApp creates an instance of app.App with default settings, suitable for unit testing. The app will be
// initialized with any specified flags as overrides, and with any specified base app options.
func DefaultTestApp(customFlags map[string]interface{}, baseAppOptions ...func(*baseapp.BaseApp)) *app.App {
	appOptions := appoptions.GetDefaultTestAppOptionsFromTempDirectory("", customFlags)
	logger, ok := appOptions.Get(testlog.LoggerInstanceForTest).(log.Logger)
	if !ok {
		logger, _ = testlog.TestLogger()
	}
	db := dbm.NewMemDB()
	dydxApp := app.New(
		logger,
		db,
		nil,
		true,
		appOptions,
		baseAppOptions...,
	)
	return dydxApp
}

// DefaultGenesis returns a genesis doc using configuration from the local net with a genesis time
// equivalent to unix epoch + 1 nanosecond. We specifically use non-zero because stateful orders
// validate that block time is non-zero (https://github.com/dydxprotocol/v4-chain/protocol/blob/
// 84a046554ab1b4725475500d94a0b3179fdd18c2/x/clob/keeper/stateful_order_state.go#L237).
func DefaultGenesis() (genesis types.GenesisDoc) {
	// NOTE: Tendermint uses a custom JSON decoder for GenesisDoc
	err := tmjson.Unmarshal([]byte(constants.GenesisState), &genesis)
	if err != nil {
		panic(err)
	}
	genesis.GenesisTime = time.Unix(0, 1)
	return genesis
}

// GenesisStates is a type constraint for all well known genesis state types.
type GenesisStates interface {
	authtypes.GenesisState |
		banktypes.GenesisState |
		perptypes.GenesisState |
		feetiertypes.GenesisState |
		stattypes.GenesisState |
		vesttypes.GenesisState |
		rewardstypes.GenesisState |
		blocktimetypes.GenesisState |
		clobtypes.GenesisState |
		pricestypes.GenesisState |
		satypes.GenesisState |
		assettypes.GenesisState |
		epochstypes.GenesisState |
		sendingtypes.GenesisState |
		delaymsgtypes.GenesisState |
		bridgetypes.GenesisState |
		govtypesv1.GenesisState |
		ratelimittypes.GenesisState |
		govplus.GenesisState |
		vaulttypes.GenesisState |
		revsharetypes.GenesisState |
		marketmapmoduletypes.GenesisState |
		aptypes.GenesisState
}

// UpdateGenesisDocWithAppStateForModule updates the supplied genesis doc using the provided function. The function
// type (any one of the well known GenesisStates) is used to derive which module will be updated.
// Will panic if there is an error in marshalling or unmarshalling the app state.
func UpdateGenesisDocWithAppStateForModule[T GenesisStates](genesisDoc *types.GenesisDoc, fn func(genesisState *T)) {
	var appState map[string]json.RawMessage
	err := json.Unmarshal(genesisDoc.AppState, &appState)
	if err != nil {
		panic(err)
	}

	var moduleName string
	var t T
	switch any(t).(type) {
	case authtypes.GenesisState:
		moduleName = authtypes.ModuleName
	case banktypes.GenesisState:
		moduleName = banktypes.ModuleName
	case blocktimetypes.GenesisState:
		moduleName = blocktimetypes.ModuleName
	case bridgetypes.GenesisState:
		moduleName = bridgetypes.ModuleName
	case delaymsgtypes.GenesisState:
		moduleName = delaymsgtypes.ModuleName
	case perptypes.GenesisState:
		moduleName = perptypes.ModuleName
	case clobtypes.GenesisState:
		moduleName = clobtypes.ModuleName
	case feetiertypes.GenesisState:
		moduleName = feetiertypes.ModuleName
	case pricestypes.GenesisState:
		moduleName = pricestypes.ModuleName
	case rewardstypes.GenesisState:
		moduleName = rewardstypes.ModuleName
	case vesttypes.GenesisState:
		moduleName = vesttypes.ModuleName
	case stattypes.GenesisState:
		moduleName = stattypes.ModuleName
	case satypes.GenesisState:
		moduleName = satypes.ModuleName
	case assettypes.GenesisState:
		moduleName = assettypes.ModuleName
	case epochstypes.GenesisState:
		moduleName = epochstypes.ModuleName
	case sendingtypes.GenesisState:
		moduleName = sendingtypes.ModuleName
	case govtypesv1.GenesisState:
		moduleName = govtypes.ModuleName
	case ratelimittypes.GenesisState:
		moduleName = ratelimittypes.ModuleName
	case govplus.GenesisState:
		moduleName = govplus.ModuleName
	case vaulttypes.GenesisState:
		moduleName = vaulttypes.ModuleName
	case revsharetypes.GenesisState:
		moduleName = revsharetypes.ModuleName
	case marketmapmoduletypes.GenesisState:
		moduleName = marketmapmoduletypes.ModuleName
	case listingtypes.GenesisState:
		moduleName = listingtypes.ModuleName
	case aptypes.GenesisState:
		moduleName = aptypes.ModuleName
	default:
		panic(fmt.Errorf("Unsupported type %T", t))
	}

	if protoMsg, ok := any(&t).(sdkproto.Message); ok {
		constants.TestEncodingCfg.Codec.MustUnmarshalJSON(appState[moduleName], protoMsg)
		fn(&t)
		appState[moduleName] = constants.TestEncodingCfg.Codec.MustMarshalJSON(protoMsg)
	} else {
		panic(fmt.Errorf("Unsupported type %T", t))
	}

	bz, err := json.Marshal(appState)
	if err != nil {
		panic(err)
	}
	genesisDoc.AppState = bz
}

// Used to instantiate new instances of the App.
type AppCreatorFn func() *app.App

// Used to instantiate new instances of the genesis doc.
type GenesisDocCreatorFn func() (genesis types.GenesisDoc)

// ExecuteCheckTxs is invoked once per block. Returning true will halt execution.
// The provided context will be a new CheckTx context using the last committed block height.
type ExecuteCheckTxs func(ctx sdk.Context, app *app.App) (stop bool)

// NewTestAppBuilder returns a new builder for TestApp.
//
// The default instance will return a builder with:
//   - DefaultGenesis
//   - no custom flags
//   - an ExecuteCheckTxs function that will stop on after the first block
//   - non-determinism checks enabled
//
// Note that the TestApp instance will have 3 non-determinism state checking apps:
//   - `parallelApp` is responsible for seeing all CheckTx requests, block proposals, blocks, and RecheckTx requests.
//     This allows it to detect state differences due to inconsistent in-memory structures (for example iteration order
//     in maps).
//   - `noCheckTxApp` is responsible for seeing all block proposals and blocks. This allows it to simulate a validator
//     that never received any of the CheckTx requests and that it will still accept blocks and arrive at the same
//     state hash.
//   - `crashingApp` is responsible for restarting before processing a block and sees all CheckTx requests, block
//     proposals, and blocks. This allows it to check that in memory state can be restored successfully on application
//     and that it will accept a block after a crash and arrive at the same state hash.
//
// Tests that rely on mutating internal application state directly (for example via keepers) will want to disable
// non-determinism checks via `WithNonDeterminismChecksEnabled(false)` otherwise the test will likely hit a
// non-determinism check that fails causing the test to fail. If possible, update the test instead to use genesis state
// to initialize state or `CheckTx` transactions to initialize the appropriate keeper state.
//
// Tests that rely on in-memory state to survive across block boundaries will want to disable crashing App CheckTx
// non-determinism checks via `WithCrashingAppCheckTxNonDeterminismChecksEnabled(false)` otherwise the test will likely
// hit a non-determinism check that fails causing the test to fail. For example unmatched short term
// orders in the memclob and order rate limits are only stored in memory and lost on application restart, and it would
// thus make sense to disable the crashing App CheckTx non-determinism check for tests that rely on this information
// surviving across block boundaries.
func NewTestAppBuilder(t testing.TB) TestAppBuilder {
	if t == nil {
		panic("t must not be nil")
	}
	return TestAppBuilder{
		genesisDocFn:                   DefaultGenesis,
		disableHealthMonitorForTesting: true,
		appOptions:                     make(map[string]interface{}),
		enableNonDeterminismChecks:     true,
		enableCrashingAppCheckTxNonDeterminismChecks: true,
		executeCheckTxs: func(ctx sdk.Context, app *app.App) (stop bool) {
			return true
		},
		t: t,
	}
}

// A builder containing fields necessary for the TestApp.
//
// Note that we specifically use value receivers for the With... methods because we want to make the builder instances
// immutable.
type TestAppBuilder struct {
	genesisDocFn                                 GenesisDocCreatorFn
	disableHealthMonitorForTesting               bool
	appOptions                                   map[string]interface{}
	executeCheckTxs                              ExecuteCheckTxs
	enableNonDeterminismChecks                   bool
	enableCrashingAppCheckTxNonDeterminismChecks bool
	t                                            testing.TB
}

// WithGenesisDocFn returns a builder like this one with specified function that will be used to create
// the genesis doc.
func (tApp TestAppBuilder) WithGenesisDocFn(fn GenesisDocCreatorFn) TestAppBuilder {
	tApp.genesisDocFn = fn
	return tApp
}

// WithHealthMonitorDisabledForTesting controls whether the daemon server health monitor is disabled for testing.
func (builder TestAppBuilder) WithHealthMonitorDisabledForTesting(disableHealthMonitorForTesting bool) TestAppBuilder {
	builder.disableHealthMonitorForTesting = disableHealthMonitorForTesting
	return builder
}

// WithNonDeterminismChecksEnabled controls whether non-determinism checks via distinct application instances
// state hash and CheckTx/ReCheckTx response comparisons.
//
// Tests that rely on mutating internal application state directly (for example via keepers) will want to disable
// non-determinism checks via `WithNonDeterminismChecksEnabled(false)` otherwise the test will likely hit a
// non-determinism check that fails causing the test to fail. If possible, update the test instead to use genesis state
// to initialize state or `CheckTx` transactions to initialize the appropriate keeper state.
func (builder TestAppBuilder) WithNonDeterminismChecksEnabled(enableNonDeterminismChecks bool) TestAppBuilder {
	builder.enableNonDeterminismChecks = enableNonDeterminismChecks
	return builder
}

// WithCrashingAppCheckTxNonDeterminismChecksEnabled controls whether the crashing App instance will ensure that
// the `CheckTx` result matches that of the main `App`.
//
// Tests that rely on in-memory state to survive across block boundaries will want to disable crashing App CheckTx
// non-determinism checks via `WithCrashingAppCheckTxNonDeterminismChecksEnabled(false)` otherwise the test will likely
// hit a non-determinism check that fails causing the test to fail. For example unmatched short term
// orders in the memclob and order rate limits are only stored in memory and lost on application restart, and it would
// thus make sense to disable the crashing App CheckTx non-determinism check for tests that rely on this information
// surviving across block boundaries.
func (builder TestAppBuilder) WithCrashingAppCheckTxNonDeterminismChecksEnabled(
	enableCrashingAppCheckTxNonDeterminismChecks bool) TestAppBuilder {
	builder.enableCrashingAppCheckTxNonDeterminismChecks = enableCrashingAppCheckTxNonDeterminismChecks
	return builder
}

// WithAppOptions returns a builder like this one with the specified app options.
func (builder TestAppBuilder) WithAppOptions(
	appOptions map[string]interface{},
) TestAppBuilder {
	builder.appOptions = appOptions
	return builder
}

// Build returns a new TestApp capable of being executed.
func (builder TestAppBuilder) Build() *TestApp {
	tApp := TestApp{
		builder: builder,
	}
	// Get the initial genesis state and initialize the chain and commit the results of the initialization.
	tApp.genesis = tApp.builder.genesisDocFn()
	if tApp.genesis.GenesisTime.UnixNano() <= time.UnixMilli(0).UnixNano() {
		tApp.builder.t.Fatal(fmt.Errorf(
			"Unable to start chain at time %v, must be greater than unix epoch.",
			tApp.genesis.GenesisTime,
		))
		return nil
	}

	// Launch the main instance of the application
	// TODO(CORE-721): Consolidate launch of apps into an abstraction since the logic is mostly repeated 4 times.
	{
		validatorHomeDir, err := prepareValidatorHomeDir(tApp.genesis)
		if err != nil {
			tApp.builder.t.Fatal(err)
			return nil
		}
		app, shutdownFn, err := launchValidatorInDir(validatorHomeDir, tApp.builder.appOptions)
		if err != nil {
			tApp.builder.t.Fatal(err)
			return nil
		}
		tApp.App = app

		tApp.builder.t.Cleanup(func() {
			doneErr := shutdownFn()

			// Clean-up the home directory.
			if err := os.RemoveAll(validatorHomeDir); err != nil {
				tApp.builder.t.Logf("Failed to clean-up temporary validator dir %s", validatorHomeDir)
			}

			if doneErr != nil {
				tApp.builder.t.Fatal(doneErr)
			}
		})
	}

	if tApp.builder.disableHealthMonitorForTesting {
		tApp.App.DisableHealthMonitorForTesting()
	}

	if tApp.builder.enableNonDeterminismChecks {
		// Filter out appOptions that shouldn't be shared to the App instances used for non-determinism checks.
		// TODO(CORE-720): Improve integration of in memory objects for e2e test framework that shouldn't be shared
		// across application instances.
		filteredAppOptions := make(map[string]interface{})
		for key, value := range tApp.builder.appOptions {
			if key != testlog.LoggerInstanceForTest && key != indexer.MsgSenderInstanceForTest {
				filteredAppOptions[key] = value
			}
		}

		// Launch the `parallelApp` instance.
		{
			validatorHomeDir, err := prepareValidatorHomeDir(tApp.genesis)
			if err != nil {
				tApp.builder.t.Fatal(err)
				return nil
			}
			app, shutdownFn, err := launchValidatorInDir(validatorHomeDir, filteredAppOptions)
			if err != nil {
				tApp.builder.t.Fatal(err)
				return nil
			}
			tApp.parallelApp = app

			tApp.builder.t.Cleanup(func() {
				doneErr := shutdownFn()

				// Clean-up the home directory.
				if err := os.RemoveAll(validatorHomeDir); err != nil {
					tApp.builder.t.Logf("Failed to clean-up temporary validator dir %s", validatorHomeDir)
				}

				if doneErr != nil {
					tApp.builder.t.Fatal(doneErr)
				}
			})
		}

		// Launch the `noCheckTx` instance.
		{
			validatorHomeDir, err := prepareValidatorHomeDir(tApp.genesis)
			if err != nil {
				tApp.builder.t.Fatal(err)
				return nil
			}
			app, shutdownFn, err := launchValidatorInDir(validatorHomeDir, filteredAppOptions)
			if err != nil {
				tApp.builder.t.Fatal(err)
				return nil
			}
			tApp.noCheckTxApp = app

			tApp.builder.t.Cleanup(func() {
				doneErr := shutdownFn()

				// Clean-up the home directory.
				if err := os.RemoveAll(validatorHomeDir); err != nil {
					tApp.builder.t.Logf("Failed to clean-up temporary validator dir %s", validatorHomeDir)
				}

				if doneErr != nil {
					tApp.builder.t.Fatal(doneErr)
				}
			})
		}

		// Launch the `crashingApp` instance.
		{
			validatorHomeDir, err := prepareValidatorHomeDir(tApp.genesis)
			if err != nil {
				tApp.builder.t.Fatal(err)
				return nil
			}
			app, shutdownFn, err := launchValidatorInDir(validatorHomeDir, filteredAppOptions)
			if err != nil {
				tApp.builder.t.Fatal(err)
				return nil
			}
			tApp.crashingApp = app

			tApp.builder.t.Cleanup(func() {
				doneErr := shutdownFn()

				// Clean-up the home directory.
				if err := os.RemoveAll(validatorHomeDir); err != nil {
					tApp.builder.t.Logf("Failed to clean-up temporary validator dir %s", validatorHomeDir)
				}

				if doneErr != nil {
					tApp.builder.t.Fatal(doneErr)
				}
			})

			tApp.restartCrashingApp = func() {
				// We shutdown the instance of the existing crashingApp.
				doneOrRestartErr := shutdownFn()
				tApp.crashingApp = nil

				if err == nil {
					app, shutdownFn, doneOrRestartErr = launchValidatorInDir(validatorHomeDir, filteredAppOptions)
				}

				// If we errored shutting down or relaunching then update the shutdownFn to return this error
				// and fatal the test.
				if err != nil {
					shutdownFn = func() error {
						return doneOrRestartErr
					}
					tApp.builder.t.Fatal(doneOrRestartErr)
					return
				}

				// Update the crashingApp pointer to the new instance of the application.
				tApp.crashingApp = app
			}
		}
	}

	return &tApp
}

// A TestApp used to executed ABCI++ flows. Note that callers should invoke `TestApp.CheckTx` over `TestApp.App.CheckTx`
// to ensure that the transaction is added to a "mempool" that will be considered during the Prepare/Process proposal
// phase.
//
// Note that the TestApp instance has 3 non-determinism state checking apps:
//   - `parallelApp` is responsible for seeing all CheckTx requests, block proposals, blocks, and RecheckTx requests.
//     This allows it to detect state differences due to inconsistent in-memory structures (for example iteration order
//     in maps).
//   - `noCheckTxApp` is responsible for seeing all block proposals and blocks. This allows it to simulate a validator
//     that never received any of the CheckTx requests and that it will still accept blocks and arrive at the same
//     state hash.
//   - `crashingApp` is responsible for restarting before processing a block and sees all CheckTx requests, block
//     proposals, and blocks. This allows it to check that in memory state can be restored successfully on application
//     and that it will accept a block after a crash and arrive at the same state hash.
//
// Note that TestApp.CheckTx is thread safe. All other methods are not thread safe.
type TestApp struct {
	// Should only be used to fetch read only state, all mutations should preferably happen through Genesis state,
	// TestApp.CheckTx, and block proposals.
	// TODO(CLOB-545): Hide App and copy the pointers to keepers to be prevent incorrect usage of App.CheckTx over
	// TestApp.CheckTx.
	App                *app.App
	parallelApp        *app.App
	noCheckTxApp       *app.App
	crashingApp        *app.App
	restartCrashingApp func()
	builder            TestAppBuilder
	genesis            types.GenesisDoc
	header             tmproto.Header
	passingCheckTxs    [][]byte
	passingCheckTxsMtx sync.Mutex
	initialized        bool
	halted             bool
	// mtx is used to enable writing concurrent tests that invoke AdvanceToBlock and CheckTx concurrently.
	// Note that AdvanceToBlock requires an exclusive lock similar to what is performed via CometBFT/Cosmos SDK
	// while CheckTx only requires a read lock since it invokes CheckTx across multiple instances of the application.
	// This allows for determinism invariant testing across these multiple instances of the application to occur.
	mtx sync.RWMutex
}

func (tApp *TestApp) Builder() TestAppBuilder {
	return tApp.builder
}

// InitChain initializes the chain. Will panic if initialized more than once.
func (tApp *TestApp) InitChain() sdk.Context {
	tApp.mtx.Lock()
	defer tApp.mtx.Unlock()

	if tApp.initialized {
		panic(errors.New("Cannot initialize chain that has been initialized already."))
	}
	tApp.initChainIfNeeded()
	return tApp.App.NewContextLegacy(true, tApp.header)
}

func (tApp *TestApp) initChainIfNeeded() {
	if tApp.initialized {
		return
	}

	tApp.initialized = true

	consensusParamsProto := tApp.genesis.ConsensusParams.ToProto()
	initChainRequest := abcitypes.RequestInitChain{
		InitialHeight:   tApp.genesis.InitialHeight,
		AppStateBytes:   tApp.genesis.AppState,
		ChainId:         tApp.genesis.ChainID,
		ConsensusParams: &consensusParamsProto,
		Time:            tApp.genesis.GenesisTime,
	}
	initChainResponse, err := tApp.App.InitChain(&initChainRequest)
	if err != nil {
		tApp.builder.t.Fatalf("Failed to initialize chain %+v, err %+v", initChainResponse, err)
	}

	if tApp.builder.enableNonDeterminismChecks {
		initChain(tApp.builder.t, tApp.parallelApp, initChainRequest, initChainResponse.AppHash)
		initChain(tApp.builder.t, tApp.noCheckTxApp, initChainRequest, initChainResponse.AppHash)
		initChain(tApp.builder.t, tApp.crashingApp, initChainRequest, initChainResponse.AppHash)
	}

	finalizeBlockRequest := abcitypes.RequestFinalizeBlock{
		Hash:   initChainResponse.AppHash,
		Height: 1,
		Time:   tApp.genesis.GenesisTime,
	}
	finalizeBlockResponse, err := tApp.App.FinalizeBlock(&finalizeBlockRequest)
	if err != nil {
		tApp.builder.t.Fatalf("Failed to finalize block %+v, err %+v", finalizeBlockResponse, err)
	}

	_, err = tApp.App.Commit()
	require.NoError(tApp.builder.t, err)
	if tApp.builder.enableNonDeterminismChecks {
		finalizeBlockAndCommit(tApp.builder.t, tApp.parallelApp, finalizeBlockRequest, tApp.App)
		finalizeBlockAndCommit(tApp.builder.t, tApp.noCheckTxApp, finalizeBlockRequest, tApp.App)
		finalizeBlockAndCommit(tApp.builder.t, tApp.crashingApp, finalizeBlockRequest, tApp.App)
	}

	tApp.header = tmproto.Header{
		ChainID:            tApp.genesis.ChainID,
		ProposerAddress:    constants.AliceAccAddress,
		Height:             tApp.App.LastBlockHeight(),
		Time:               tApp.genesis.GenesisTime,
		LastCommitHash:     tApp.App.LastCommitID().Hash,
		NextValidatorsHash: tApp.App.LastCommitID().Hash,
	}
}

// AdvanceToBlock advances the chain to the specified block and block time.
// If the specified block is the current block, simply returns the current context.
// For example, block = 10, t = 90 will advance to a block with height 10 and with a time of 90.
func (tApp *TestApp) AdvanceToBlock(
	block uint32,
	options AdvanceToBlockOptions,
) sdk.Context {
	tApp.mtx.Lock()
	defer tApp.mtx.Unlock()

	tApp.panicIfChainIsHalted()
	tApp.initChainIfNeeded()

	if options.BlockTime.IsZero() { // if time is not specified, use the current block time.
		options.BlockTime = tApp.header.Time
	}
	if int64(block) <= tApp.header.Height {
		panic(fmt.Errorf("Expected block height (%d) > current block height (%d).", block, tApp.header.Height))
	}
	if options.BlockTime.UnixNano() < tApp.header.Time.UnixNano() {
		panic(fmt.Errorf("Expected time (%v) >= current block time (%v).", options.BlockTime, tApp.header.Time))
	}
	if int64(block) == tApp.header.Height {
		return tApp.App.NewContextLegacy(true, tApp.header)
	}

	// Ensure that we grab the lock so that we can read and write passingCheckTxs correctly.
	tApp.passingCheckTxsMtx.Lock()
	defer tApp.passingCheckTxsMtx.Unlock()

	// Advance to the requested block using the requested block time.
	for tApp.App.LastBlockHeight() < int64(block) {
		tApp.panicIfChainIsHalted()
		tApp.header.Height = tApp.App.LastBlockHeight() + 1
		if tApp.header.Height == int64(block) {
			// By default, only update block time at the requested block.
			tApp.header.Time = options.BlockTime
		} else if options.LinearBlockTimeInterpolation {
			remainingDuration := options.BlockTime.Sub(tApp.header.Time)
			nextBlockDuration := remainingDuration / time.Duration(int64(block)-tApp.App.LastBlockHeight())
			tApp.header.Time = tApp.header.Time.Add(nextBlockDuration)
		}
		tApp.header.LastCommitHash = tApp.App.LastCommitID().Hash
		tApp.header.NextValidatorsHash = tApp.App.LastCommitID().Hash

		deliverTxs := options.DeliverTxsOverride
		if deliverTxs == nil {
			// Prepare the proposal and process it.
			prepareRequest := abcitypes.RequestPrepareProposal{
				Txs:                tApp.passingCheckTxs,
				MaxTxBytes:         math.MaxInt64,
				Height:             tApp.header.Height,
				Time:               tApp.header.Time,
				NextValidatorsHash: tApp.header.NextValidatorsHash,
				ProposerAddress:    tApp.header.ProposerAddress,
			}
			if options.RequestPrepareProposalTxsOverride != nil {
				prepareRequest.Txs = options.RequestPrepareProposalTxsOverride
			}
			prepareResponse, prepareErr := tApp.App.PrepareProposal(&prepareRequest)

			if options.ValidateRespPrepare != nil {
				haltChain := options.ValidateRespPrepare(
					tApp.App.NewContextLegacy(true, tApp.header),
					*prepareResponse,
				)
				tApp.halted = haltChain
				if tApp.halted {
					return tApp.App.NewContextLegacy(true, tApp.header)
				}
			}

			require.NoError(
				tApp.builder.t,
				prepareErr,
				"Expected prepare proposal request %+v to succeed, but failed with %+v and err %+v.",
				prepareRequest,
				prepareResponse,
				prepareErr,
			)

			// Pass forward any transactions that were modified through the preparation phase and process them.
			if options.RequestProcessProposalTxsOverride != nil {
				prepareResponse.Txs = options.RequestProcessProposalTxsOverride
			}
			processRequest := abcitypes.RequestProcessProposal{
				Txs:                prepareResponse.Txs,
				Hash:               tApp.header.AppHash,
				Height:             tApp.header.Height,
				Time:               tApp.header.Time,
				NextValidatorsHash: tApp.header.NextValidatorsHash,
				ProposerAddress:    tApp.header.ProposerAddress,
			}
			processResponse, processErr := tApp.App.ProcessProposal(&processRequest)

			if options.ValidateRespProcess != nil {
				haltChain := options.ValidateRespProcess(
					tApp.App.NewContextLegacy(true, tApp.header),
					*processResponse,
				)
				tApp.halted = haltChain
				if tApp.halted {
					return tApp.App.NewContextLegacy(true, tApp.header)
				}
			}

			require.Truef(
				tApp.builder.t,
				processErr == nil && processResponse.IsAccepted(),
				"Expected process proposal request %+v to be accepted, but failed with %+v and err %+v.",
				processRequest,
				processResponse,
				processErr,
			)

			// Check that all instances of the application can process the proposoal and come to the same result.
			if tApp.builder.enableNonDeterminismChecks {
				parallelProcessResponse, parallelProcessErr := tApp.parallelApp.ProcessProposal(&processRequest)
				require.Truef(
					tApp.builder.t,
					parallelProcessErr == nil && parallelProcessResponse.IsAccepted(),
					"Non-determinism detected, expected process proposal request %+v to be accepted, but failed with %+v and err %+v.",
					processRequest,
					parallelProcessResponse,
					parallelProcessErr,
				)
				noCheckTxProcessResponse, noCheckTxProcessErr := tApp.noCheckTxApp.ProcessProposal(&processRequest)
				require.Truef(
					tApp.builder.t,
					noCheckTxProcessErr == nil && noCheckTxProcessResponse.IsAccepted(),
					"Non-determinism detected, expected process proposal request %+v to be accepted, but failed with %+v and err %+v.",
					processRequest,
					noCheckTxProcessResponse,
					noCheckTxProcessErr,
				)
				crashingProcessResponse, crashingProcessErr := tApp.crashingApp.ProcessProposal(&processRequest)
				require.Truef(
					tApp.builder.t,
					crashingProcessErr == nil && crashingProcessResponse.IsAccepted(),
					"Non-determinism detected, expected process proposal request %+v to be accepted, but failed with %+v and err %+v.",
					processRequest,
					crashingProcessResponse,
					crashingProcessErr,
				)
			}

			deliverTxs = prepareResponse.Txs
		}

		txsNotInLastProposal := make([][]byte, 0)
		for _, tx := range tApp.passingCheckTxs {
			if !slices.ContainsFunc(deliverTxs, func(tx2 []byte) bool {
				return bytes.Equal(tx, tx2)
			}) {
				txsNotInLastProposal = append(txsNotInLastProposal, tx)
			}
		}
		tApp.passingCheckTxs = txsNotInLastProposal

		// Restart the crashingApp instance before processing the block.
		if tApp.builder.enableNonDeterminismChecks {
			tApp.restartCrashingApp()
		}

		// Finalize the block
		finalizeBlockRequest := abcitypes.RequestFinalizeBlock{
			Txs:                deliverTxs,
			Hash:               tApp.header.AppHash,
			Height:             tApp.header.Height,
			Time:               tApp.header.Time,
			NextValidatorsHash: tApp.header.NextValidatorsHash,
			ProposerAddress:    tApp.header.ProposerAddress,
		}
		finalizeBlockResponse, finalizeBlockErr := tApp.App.FinalizeBlock(&finalizeBlockRequest)

		if options.ValidateFinalizeBlock != nil {
			tApp.halted = options.ValidateFinalizeBlock(
				tApp.App.NewContextLegacy(false, tApp.header),
				finalizeBlockRequest,
				*finalizeBlockResponse,
			)
			if tApp.halted {
				return tApp.App.NewContextLegacy(true, tApp.header)
			}
		} else {
			require.NoErrorf(
				tApp.builder.t,
				finalizeBlockErr,
				"Expected block finalization to succeed but failed %+v with err %+v.",
				finalizeBlockResponse,
				finalizeBlockErr,
			)
			for i, txResult := range finalizeBlockResponse.TxResults {
				require.Conditionf(
					tApp.builder.t,
					txResult.IsOK,
					"Failed to deliver transaction %d that was accepted: %+v. Response: %+v",
					i,
					txResult,
					finalizeBlockResponse,
				)
			}
		}

		// Commit the block.
		_, err := tApp.App.Commit()
		require.NoError(tApp.builder.t, err)

		// Finalize and commit all the blocks for the non-determinism checkers.
		if tApp.builder.enableNonDeterminismChecks {
			finalizeBlockAndCommit(tApp.builder.t, tApp.parallelApp, finalizeBlockRequest, tApp.App)
			finalizeBlockAndCommit(tApp.builder.t, tApp.noCheckTxApp, finalizeBlockRequest, tApp.App)
			finalizeBlockAndCommit(tApp.builder.t, tApp.crashingApp, finalizeBlockRequest, tApp.App)
		}

		// Recheck the remaining transactions in the mempool pruning any that have failed during recheck.
		passingRecheckTxs := make([][]byte, 0)
		for _, passingCheckTx := range tApp.passingCheckTxs {
			recheckTxRequest := abcitypes.RequestCheckTx{
				Tx:   passingCheckTx,
				Type: abcitypes.CheckTxType_Recheck,
			}
			recheckTxResponse, recheckTxErr := tApp.App.CheckTx(&recheckTxRequest)
			if recheckTxErr == nil && recheckTxResponse.IsOK() {
				passingRecheckTxs = append(passingRecheckTxs, passingCheckTx)
			}

			if tApp.builder.enableNonDeterminismChecks {
				parallelRecheckTxResponse, parallelRecheckTxErr := tApp.parallelApp.CheckTx(&recheckTxRequest)
				require.Truef(
					tApp.builder.t,
					recheckTxResponse.Code == parallelRecheckTxResponse.Code &&
						((recheckTxErr == nil && parallelRecheckTxErr == nil) ||
							(recheckTxErr != nil && parallelRecheckTxErr != nil)),
					"Non-determinism detected during RecheckTx, expected %+v with err %+v, got %+v with err %+v.",
					recheckTxResponse,
					recheckTxErr,
					parallelRecheckTxResponse,
					parallelRecheckTxErr,
				)

				// None of the transactions should be rechecked in `noCheckTxApp` since the transaction will only
				// process block proposals and blocks. Also, none of the transactions should be rechecked for
				// tApp.crashingApp since the mempool should be discarded on each crash.
			}
		}
		tApp.passingCheckTxs = passingRecheckTxs
	}

	return tApp.App.NewContextLegacy(true, tApp.header)
}

// initChain initializes the chain and ensures that it did successfully and also that it reached the expected
// app hash.
func initChain(t testing.TB, app *app.App, request abcitypes.RequestInitChain, expectedAppHash []byte) {
	response, err := app.InitChain(&request)
	if err != nil {
		t.Fatalf("Failed to initialize chain %+v, err %+v", response, err)
	}
	require.Equal(t, expectedAppHash, response.AppHash)
}

// finalizeBlockAndCommit finalizes the block and commits the chain verifying that the chain matches the
// expected commit id.
func finalizeBlockAndCommit(
	t testing.TB,
	app *app.App,
	request abcitypes.RequestFinalizeBlock,
	expectedApp *app.App,
) {
	_, err := app.FinalizeBlock(&request)
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	diffs := make([]string, 0)
	if !bytes.Equal(app.LastCommitID().Hash, expectedApp.LastCommitID().Hash) {
		rootMulti := app.CommitMultiStore().(*rootmulti.Store)
		expectedRootMulti := expectedApp.CommitMultiStore().(*rootmulti.Store)

		commitInfo, err := rootMulti.GetCommitInfo(app.LastBlockHeight())
		require.NoError(t, err)
		expectedCommitInfo, err := expectedRootMulti.GetCommitInfo(app.LastBlockHeight())
		require.NoError(t, err)
		storeInfos := make(map[string]storetypes.StoreInfo)
		for _, storeInfo := range commitInfo.StoreInfos {
			storeInfos[storeInfo.Name] = storeInfo
		}

		storeDiffs := make([]string, 0)
		for _, storeInfo := range expectedCommitInfo.StoreInfos {
			if !bytes.Equal(storeInfos[storeInfo.Name].GetHash(), storeInfo.GetHash()) {
				diffs = append(
					diffs,
					fmt.Sprintf("Expected store '%s' hashes to match.", storeInfo.Name),
				)
				storeDiffs = append(storeDiffs, storeInfo.Name)
			}
		}

		for _, storeName := range storeDiffs {
			store := rootMulti.GetCommitKVStore(rootMulti.StoreKeysByName()[storeName])
			expectedStore := expectedRootMulti.GetCommitKVStore(expectedRootMulti.StoreKeysByName()[storeName])

			itr := store.Iterator(nil, nil)
			defer itr.Close()
			expectedItr := expectedStore.Iterator(nil, nil)
			defer expectedItr.Close()

			for ; itr.Valid(); itr.Next() {
				if !expectedItr.Valid() {
					diffs = append(
						diffs,
						fmt.Sprintf(
							"Expected key/value ('%s', '%s') to exist in store '%s'.",
							itr.Key(),
							itr.Value(),
							storeName,
						),
					)
					continue
				} else if !bytes.Equal(itr.Key(), expectedItr.Key()) {
					diffs = append(
						diffs,
						fmt.Sprintf(
							"Expected key '%s' to exist in store '%s'.",
							itr.Key(),
							storeName,
						),
					)
				} else if !bytes.Equal(itr.Value(), expectedItr.Value()) {
					diffs = append(
						diffs,
						fmt.Sprintf(
							"Found key '%s' with different value '%s' which "+
								"differs from original value '%s' in store '%s'.",
							itr.Key(),
							expectedItr.Value(),
							itr.Value(),
							storeName,
						),
					)
				}

				expectedItr.Next()
			}

			for ; expectedItr.Valid(); expectedItr.Next() {
				diffs = append(
					diffs,
					fmt.Sprintf(
						"Expected key/value (%s, %s) to not exist in store '%s'",
						expectedItr.Key(),
						expectedItr.Value(),
						storeName,
					),
				)
			}
		}
	}

	if len(diffs) != 0 {
		t.Errorf(
			"Expected no differences in stores but found %d difference(s):\n  %s",
			len(diffs),
			strings.Join(diffs, "\n  "),
		)
		t.FailNow()
	}

	// Ensure that all instances after committing the block came to the same commit hash.
	require.Equalf(t,
		expectedApp.LastCommitID(),
		app.LastCommitID(),
		"Non-determinism in state detected, expected LastCommitID to match.",
	)
}

// GetHeader fetches the current header of the test app.
func (tApp *TestApp) GetHeader() tmproto.Header {
	tApp.mtx.RLock()
	defer tApp.mtx.RUnlock()
	return tApp.header
}

// GetBlockHeight fetches the current block height of the test app.
func (tApp *TestApp) GetBlockHeight() int64 {
	tApp.mtx.RLock()
	defer tApp.mtx.RUnlock()
	return tApp.header.Height
}

// GetHalted fetches the halted flag.
func (tApp *TestApp) GetHalted() bool {
	tApp.mtx.RLock()
	defer tApp.mtx.RUnlock()
	return tApp.halted
}

// newTestingLogger returns a logger that will write to stdout if testing is verbose. This method replaces
// cometbft's log.TestingLogger, which re-uses the same logger for all tests, which can cause race test false positives
// when accessed by concurrent go routines in the same test.
func newTestingLogger() cmtlog.Logger {
	if testing.Verbose() {
		return cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	} else {
		return cmtlog.NewNopLogger()
	}
}

// CheckTx adds the transaction to a test specific "mempool" that will be used to deliver the transaction during
// Prepare/Process proposal. Note that this must be invoked over TestApp.App.CheckTx as the transaction will not
// be added to the "mempool" causing the transaction to not be supplied during the Prepare/Process proposal phase.
//
// This method is thread-safe.
func (tApp *TestApp) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	tApp.mtx.RLock()
	defer tApp.mtx.RUnlock()

	tApp.panicIfChainIsHalted()
	res, err := tApp.App.CheckTx(&req)
	// Note that the dYdX fork of CometBFT explicitly excludes place and cancel order messages. See
	// https://github.com/dydxprotocol/cometbft/blob/5e6c4b6/mempool/clist_mempool.go#L441
	if err == nil && res.IsOK() && !mempool.IsShortTermClobOrderTransaction(req.Tx, newTestingLogger()) {
		// We want to ensure that we hold the lock only for updating passingCheckTxs so that App.CheckTx can execute
		// concurrently.
		tApp.passingCheckTxsMtx.Lock()
		defer tApp.passingCheckTxsMtx.Unlock()
		tApp.passingCheckTxs = append(tApp.passingCheckTxs, req.Tx)
	}

	if tApp.builder.enableNonDeterminismChecks {
		// We expect the parallel app to always produce the same result since all in memory state should be
		// consistent with tApp.App and produce the same result.
		parallelRes, parallelErr := tApp.parallelApp.CheckTx(&req)
		require.Truef(
			tApp.builder.t,
			res.Code == parallelRes.Code && ((err == nil && parallelErr == nil) || (err != nil && parallelErr != nil)),
			"Parallel app non-determinism detected during CheckTx, expected %+v with err %+v, got %+v with err %+v.",
			res,
			err,
			parallelRes,
			parallelErr,
		)

		// The crashing app may or may not be able to get to a recoverable state that would produce equivalent
		// results. For example short-term orders and cancellations will be lost from in-memory state.
		crashingRes, crashingErr := tApp.crashingApp.CheckTx(&req)
		if tApp.builder.enableCrashingAppCheckTxNonDeterminismChecks {
			require.Truef(
				tApp.builder.t,
				res.Code == crashingRes.Code && ((err == nil && crashingErr == nil) || (err != nil && crashingErr != nil)),
				"Crashing app non-determinism detected during CheckTx, expected %+v with err %+v, got %+v with err %+v.",
				res,
				err,
				crashingRes,
				crashingErr,
			)
		}
	}
	return *res
}

// panicIfChainIsHalted panics if the chain is halted.
func (tApp *TestApp) panicIfChainIsHalted() {
	if tApp.halted {
		panic("Chain is halted")
	}
}

// PrepareProposal creates an abci `RequestPrepareProposal` using the current state of the chain
// and calls the PrepareProposal handler to return an abci `ResponsePrepareProposal`.
func (tApp *TestApp) PrepareProposal() (*abcitypes.ResponsePrepareProposal, error) {
	tApp.mtx.Lock()
	defer tApp.mtx.Unlock()

	return tApp.App.PrepareProposal(&abcitypes.RequestPrepareProposal{
		Txs:                tApp.passingCheckTxs,
		MaxTxBytes:         math.MaxInt64,
		Height:             tApp.header.Height,
		Time:               tApp.header.Time,
		NextValidatorsHash: tApp.header.NextValidatorsHash,
		ProposerAddress:    tApp.header.ProposerAddress,
	})
}

// GetProposedOperations returns the operations queue that would be proposed if a proposal was generated by the
// application for the current block height. This is helpful for testcases where we want to use DeliverTxsOverride
// to insert new transactions, but preserve the operations that would have been proposed.
func (tApp *TestApp) GetProposedOperationsTx() []byte {
	tApp.mtx.Lock()
	defer tApp.mtx.Unlock()

	request := abcitypes.RequestPrepareProposal{
		Txs:                tApp.passingCheckTxs,
		MaxTxBytes:         math.MaxInt64,
		Height:             tApp.header.Height,
		Time:               tApp.header.Time,
		NextValidatorsHash: tApp.header.NextValidatorsHash,
		ProposerAddress:    tApp.header.ProposerAddress,
	}
	response, err := tApp.App.PrepareProposal(&request)
	require.NoError(
		tApp.builder.t,
		err,
		"Expected prepare proposal request %+v to succeed, but failed with %+v and err %+v.",
		request,
		response,
		err,
	)
	return response.Txs[0]
}

// prepareValidatorHomeDir launches a validator using the `start` command with the specified genesis doc and application
// options. `shutdownFn` must be invoked to cancel the execution of the app. It will block till the application
// shuts down.
func prepareValidatorHomeDir(
	genesis types.GenesisDoc,
) (validatorHomeDir string, err error) {
	// Create the validators home directory as a temporary directory and fill it with:
	//  - config/priv_validator_key.json
	//  - config/node_key.json
	//  - config/genesis.json
	validatorHomeDir = filepath.Join(os.TempDir(), fmt.Sprint(time.Now().UnixNano()))
	if err = os.MkdirAll(fmt.Sprintf("%s/config/", validatorHomeDir), 0755); err != nil {
		return "", err
	}
	if err = os.WriteFile(
		filepath.Join(validatorHomeDir, "config", "priv_validator_key.json"),
		[]byte(alicePrivValidatorKeyJson),
		0755,
	); err != nil {
		return "", err
	}
	if err = os.WriteFile(
		filepath.Join(validatorHomeDir, "config", "node_key.json"),
		[]byte(aliceNodeKeyJson),
		0755,
	); err != nil {
		return "", err
	}
	if err = genesis.SaveAs(filepath.Join(validatorHomeDir, "config", "genesis.json")); err != nil {
		return "", err
	}
	return validatorHomeDir, err
}

func launchValidatorInDir(
	validatorHomeDir string,
	appOptions map[string]interface{},
) (a *app.App, shutdownFn func() error, err error) {
	// Create a context that can be cancelled to stop the Cosmos App.
	done := make(chan error, 1)
	parentCtx, cancelFn := context.WithCancel(context.Background())

	appCaptor := make(chan *app.App, 1)
	// Set up the root command using https://github.com/dydxprotocol/v4-chain/blob/
	// 1fa21ed5d848ed7cc6a98053838cadb68422079f/protocol/cmd/dydxprotocold/main.go#L12 as a basis.
	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmdWithInterceptors(
		option,
		validatorHomeDir,
		// Inject the app options and logger
		func(serverCtxPtr *server.Context) {
			for key, value := range appOptions {
				serverCtxPtr.Viper.Set(key, value)
			}

			// Set the test logger instance based upon AppOptions.
			if logger, ok := appOptions[testlog.LoggerInstanceForTest]; ok {
				serverCtxPtr.Logger = logger.(log.Logger)
			}
		},
		// Override the addresses to use domain sockets to avoid port conflicts.
		func(s string, appConfig *cmd.DydxAppConfig) (string, *cmd.DydxAppConfig) {
			// Note that the domain sockets need to typically be ~100 bytes or fewer otherwise they will fail to be
			// created. The actual limit is OS specific.
			apiSocketPath := filepath.Join(validatorHomeDir, "api_socket")
			grpcSocketPath := filepath.Join(validatorHomeDir, "grpc_socket")
			appConfig.API.Address = fmt.Sprintf("unix://%s", apiSocketPath)
			appConfig.GRPC.Address = fmt.Sprintf("unix://%s", grpcSocketPath)

			// TODO(CORE-29): This disables launching the daemons since not all daemons currently shutdown as needed.
			appConfig.API.Enable = false
			// We disable telemetry since multiple instances of the application fail to register causing a panic.
			appConfig.Telemetry.Enabled = false
			appConfig.Oracle.MetricsEnabled = false
			return s, appConfig
		},
		// Capture the application instance.
		func(app *app.App) *app.App {
			appCaptor <- app
			return app
		},
	)

	// Specify the start-up flags.
	// TODO(CLOB-930): Allow for these flags to be overridden.
	rootCmd.SetArgs([]string{
		"start",
		// Do not start tendermint.
		"--grpc-only",
		"true",
		"--home",
		validatorHomeDir,
		// TODO(CORE-29): Allow the daemons to be launched and cleaned-up successfully by default.
		"--price-daemon-enabled",
		"false",
		"--bridge-daemon-enabled",
		"false",
		"--liquidation-daemon-enabled",
		"false",
		"--bridge-daemon-eth-rpc-endpoint",
		"https://eth-sepolia.g.alchemy.com/v2/demo",
		"--oracle.enabled=false",
		"--oracle.metrics_enabled=false",
		"--log_level=error",
		// TODO(CT-1329): Currently, the TestApp framework does not work well with OE,
		// since the non-deterministic test app instances does not handle go routines.
		"--optimistic-execution-enabled=false",
	})

	ctx := svrcmd.CreateExecuteContext(parentCtx)
	rootCmd.PersistentFlags().String(
		flags.FlagLogLevel,
		tmcfg.DefaultLogLevel,
		"The logging level (trace|debug|info|warn|error|fatal|panic)",
	)
	rootCmd.PersistentFlags().String(
		flags.FlagLogFormat,
		tmcfg.LogFormatPlain,
		"The logging format (json|plain)",
	)
	executor := tmcli.PrepareBaseCmd(rootCmd, appconstants.AppDaemonName, app.DefaultNodeHome)
	// We need to launch the root command in a separate go routine since it only returns once the app is shutdown.
	// So we wait for either the app to be captured representing a successful start or capture an error.
	go func() {
		// ExecuteContext will block and will only return if interrupted.
		err := executor.ExecuteContext(ctx)
		done <- err
	}()
	select {
	case a = <-appCaptor:
		shutdownFn = func() error {
			cancelFn()
			return <-done
		}
		return a, shutdownFn, nil
	case err = <-done:
		// Send the error to done channel so that `Cleanup` function will not block.
		cancelFn()
		done <- err
		return nil, nil, err
	}
}

// MustMakeCheckTxsWithClobMsg creates one signed RequestCheckTx for each msg passed in.
// The messsage must use one of the hard-coded well known subaccount owners otherwise this will panic.
func MustMakeCheckTxsWithClobMsg[
	T clobtypes.MsgPlaceOrder |
		clobtypes.MsgCancelOrder |
		clobtypes.MsgBatchCancel |
		clobtypes.MsgUpdateLeverage](
	ctx sdk.Context,
	app *app.App,
	messages ...T,
) []abcitypes.RequestCheckTx {
	sdkMessages := make([]sdk.Msg, len(messages))
	var signerAddress string
	for i, msg := range messages {
		var m sdk.Msg
		switch v := any(msg).(type) {
		case clobtypes.MsgPlaceOrder:
			signerAddress = v.Order.OrderId.SubaccountId.Owner
			m = &v
		case clobtypes.MsgCancelOrder:
			signerAddress = v.OrderId.SubaccountId.Owner
			m = &v
		case clobtypes.MsgBatchCancel:
			signerAddress = v.SubaccountId.Owner
			m = &v
		case clobtypes.MsgUpdateLeverage:
			signerAddress = v.SubaccountId.Owner
			m = &v
		default:
			panic(fmt.Errorf("MustMakeCheckTxsWithClobMsg: Unknown message type %T", msg))
		}

		sdkMessages[i] = m
	}

	return MustMakeCheckTxsWithSdkMsg(
		ctx,
		app,
		MustMakeCheckTxOptions{
			AccAddressForSigning: signerAddress,
		},
		sdkMessages...,
	)
}

// MustMakeCheckTxsWithSdkMsg creates one signed RequestCheckTx for each msg passed in.
// The messsage must use one of the hard-coded well known subaccount owners otherwise this will panic.
func MustMakeCheckTxsWithSdkMsg(
	ctx sdk.Context,
	app *app.App,
	options MustMakeCheckTxOptions,
	messages ...sdk.Msg,
) (checkTxs []abcitypes.RequestCheckTx) {
	for _, msg := range messages {
		checkTxs = append(checkTxs, MustMakeCheckTx(ctx, app, options, msg))
	}

	return checkTxs
}

// MustMakeCheckTx creates a signed RequestCheckTx for the provided message.
// The message must use one of the hard-coded well known subaccount owners otherwise this will panic.
func MustMakeCheckTx(
	ctx sdk.Context,
	app *app.App,
	options MustMakeCheckTxOptions,
	messages ...sdk.Msg,
) abcitypes.RequestCheckTx {
	return MustMakeCheckTxWithPrivKeySupplier(
		ctx,
		app,
		options,
		constants.GetPrivateKeyFromAddress,
		messages...,
	)
}

// MustMakeCheckTxWithPrivKeySupplier creates a signed RequestCheckTx for the provided message. The `privKeySupplier`
// must provide a valid private key that matches the supplied account address.
func MustMakeCheckTxWithPrivKeySupplier(
	ctx sdk.Context,
	app *app.App,
	options MustMakeCheckTxOptions,
	privKeySupplier func(accAddress string) cryptotypes.PrivKey,
	messages ...sdk.Msg,
) abcitypes.RequestCheckTx {
	accAddress := sdk.MustAccAddressFromBech32(options.AccAddressForSigning)
	privKey := privKeySupplier(options.AccAddressForSigning)
	if !app.AccountKeeper.HasAccount(ctx, accAddress) {
		panic("Account not found")
	}
	account := app.AccountKeeper.GetAccount(ctx, accAddress)
	sequenceNumber := account.GetSequence()
	if options.AccSequenceNumberForSigning > 0 {
		sequenceNumber = options.AccSequenceNumberForSigning
	}
	checkTx, err := sims.GenSignedMockTx(
		rand.New(rand.NewSource(42)),
		app.TxConfig(),
		messages,
		options.FeeAmt,
		options.Gas,
		ctx.ChainID(),
		[]uint64{account.GetAccountNumber()},
		[]uint64{sequenceNumber},
		privKey,
	)
	if err != nil {
		panic(err)
	}
	bytes, err := app.TxConfig().TxEncoder()(checkTx)
	if err != nil {
		panic(err)
	}
	return abcitypes.RequestCheckTx{
		Tx:   bytes,
		Type: abcitypes.CheckTxType_New,
	}
}
