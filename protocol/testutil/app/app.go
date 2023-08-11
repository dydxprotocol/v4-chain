package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/mempool"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdkproto "github.com/cosmos/gogoproto/proto"

	"github.com/dydxprotocol/v4/app"
	"github.com/dydxprotocol/v4/testutil/appoptions"
	"github.com/dydxprotocol/v4/testutil/constants"
	testtx "github.com/dydxprotocol/v4/testutil/tx"
	assettypes "github.com/dydxprotocol/v4/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4/x/epochs/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4/x/sending/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

// MustMakeCheckTxOptions is a struct containing options for MustMakeCheckTx.* functions.
type MustMakeCheckTxOptions struct {
	// AccAddressForSigning is the account that's used to sign the transaction.
	AccAddressForSigning string
}

// Create an instance of app.App with default settings, suitable for unit testing,
// with the option to override specific flags.
func DefaultTestApp(customFlags map[string]interface{}) *app.App {
	appOptions := appoptions.GetDefaultTestAppOptionsFromTempDirectory("", customFlags)
	logger := log.TestingLogger()
	db := dbm.NewMemDB()
	dydxApp := app.New(
		logger,
		db,
		nil,
		true,
		appOptions,
	)
	return dydxApp
}

// DefaultTestAppCreatorFn is a wrapper function around DefaultTestApp using the specified custom flags.
func DefaultTestAppCreatorFn(customFlags map[string]interface{}) AppCreatorFn {
	return func() *app.App {
		return DefaultTestApp(customFlags)
	}
}

// DefaultGenesis returns a genesis doc using configuration from the local net with a genesis time
// equivalent to unix epoch + 1 nanosecond. We specifically use non-zero because stateful orders
// validate that block time is non-zero (https://github.com/dydxprotocol/v4/blob/
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
		clobtypes.GenesisState |
		pricestypes.GenesisState |
		satypes.GenesisState |
		assettypes.GenesisState |
		epochstypes.GenesisState |
		sendingtypes.GenesisState
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
	case perptypes.GenesisState:
		moduleName = perptypes.ModuleName
	case clobtypes.GenesisState:
		moduleName = clobtypes.ModuleName
	case pricestypes.GenesisState:
		moduleName = pricestypes.ModuleName
	case satypes.GenesisState:
		moduleName = satypes.ModuleName
	case assettypes.GenesisState:
		moduleName = assettypes.ModuleName
	case epochstypes.GenesisState:
		moduleName = epochstypes.ModuleName
	case sendingtypes.GenesisState:
		moduleName = sendingtypes.ModuleName
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
// The default instance will return a builder using:
//   - DefaultGenesis
//   - DefaultTestAppCreatorFn with no custom flags
//   - an ExecuteCheckTxs function that will stop on after the first block
func NewTestAppBuilder() TestAppBuilder {
	return TestAppBuilder{
		genesisDocFn: DefaultGenesis,
		appCreatorFn: DefaultTestAppCreatorFn(nil),
		executeCheckTxs: func(ctx sdk.Context, app *app.App) (stop bool) {
			return true
		},
	}
}

// A builder containing fields necessary for the TestApp.
//
// Note that we specifically use value receivers for the With... methods because we want to make the builder instances
// immutable.
type TestAppBuilder struct {
	genesisDocFn    GenesisDocCreatorFn
	appCreatorFn    func() *app.App
	executeCheckTxs ExecuteCheckTxs
	t               *testing.T
}

// WithGenesisDocFn returns a builder like this one with specified function that will be used to create
// the genesis doc.
func (tApp TestAppBuilder) WithGenesisDocFn(fn GenesisDocCreatorFn) TestAppBuilder {
	tApp.genesisDocFn = fn
	return tApp
}

// WithAppCreatorFn returns a builder like this one with the specified function that will be used to create
// the application.
func (tApp TestAppBuilder) WithAppCreatorFn(fn AppCreatorFn) TestAppBuilder {
	tApp.appCreatorFn = fn
	return tApp
}

// WithTesting returns a builder like this one with the specified testing environment being specified.
func (tApp TestAppBuilder) WithTesting(t *testing.T) TestAppBuilder {
	tApp.t = t
	return tApp
}

// Build returns a new TestApp capable of being executed.
func (tApp TestAppBuilder) Build() TestApp {
	return TestApp{
		builder: tApp,
	}
}

// A TestApp used to executed ABCI++ flows. Note that callers should invoke `TestApp.CheckTx` over `TestApp.App.CheckTx`
// to ensure that the transaction is added to a "mempool" that will be considered during the Prepare/Process proposal
// phase.
//
// Note that TestApp.CheckTx is thread safe. All other methods are not thread safe.
type TestApp struct {
	// Should only be used to fetch read only state, all mutations should preferrably happen through Genesis state,
	// TestApp.CheckTx, and block proposals.
	// TODO(CLOB-545): Hide App and copy the pointers to keepers to be prevent incorrect usage of App.CheckTx over
	// TestApp.CheckTx.
	App                *app.App
	builder            TestAppBuilder
	genesis            types.GenesisDoc
	header             tmproto.Header
	passingCheckTxs    [][]byte
	passingCheckTxsMtx sync.Mutex
}

func (tApp *TestApp) Builder() TestAppBuilder {
	return tApp.builder
}

// InitChain initializes the chain. Will panic if initialized more than once.
func (tApp *TestApp) InitChain() sdk.Context {
	if tApp.App != nil {
		panic(errors.New("Cannot initialize chain that has been initialized already. Missing a Reset()?"))
	}
	tApp.initChainIfNeeded()
	return tApp.App.NewContext(true, tApp.header)
}

func (tApp *TestApp) initChainIfNeeded() {
	if tApp.App != nil {
		return
	}

	fmt.Printf("!!!!!! %+v\n", tApp.genesis.AppState)
	// Get the initial genesis state and initialize the chain and commit the results of the initialization.
	tApp.genesis = tApp.builder.genesisDocFn()
	tApp.App = tApp.builder.appCreatorFn()
	baseapp.SetChainID(tApp.genesis.ChainID)(tApp.App.GetBaseApp())
	if tApp.genesis.GenesisTime.UnixNano() <= time.UnixMilli(0).UnixNano() {
		panic(fmt.Errorf(
			"Unable to start chain at time %v, must be greater than unix epoch.",
			tApp.genesis.GenesisTime,
		))
	}

	consensusParamsProto := tApp.genesis.ConsensusParams.ToProto()

	tApp.App.InitChain(abcitypes.RequestInitChain{
		InitialHeight:   tApp.genesis.InitialHeight,
		AppStateBytes:   tApp.genesis.AppState,
		ChainId:         tApp.genesis.ChainID,
		ConsensusParams: &consensusParamsProto,
		Time:            tApp.genesis.GenesisTime,
	})
	tApp.App.Commit()

	tApp.header = tmproto.Header{
		ChainID:            tApp.genesis.ChainID,
		ProposerAddress:    constants.AliceAccAddress,
		Height:             tApp.App.LastBlockHeight(),
		Time:               tApp.genesis.GenesisTime,
		LastCommitHash:     tApp.App.LastCommitID().Hash,
		NextValidatorsHash: tApp.App.LastCommitID().Hash,
	}
}

// AdvanceToBlockIfNecessary advances the chain to the specified block using the current block time.
// If the block is the same, then this function results in a no-op.
// Note that due to DEC-1248 the minimum block height is 2.
func (tApp *TestApp) AdvanceToBlockIfNecessary(block uint32) sdk.Context {
	if int64(block) == tApp.GetBlockHeight() {
		return tApp.App.NewContext(true, tApp.header)
	}

	return tApp.AdvanceToBlock(block)
}

// AdvanceToBlock advances the chain to the specified block using the current block time.
// For example, block = 10 will advance to a block with height 10 without changing the time.
//
// Note that due to DEC-1248 the minimum block height is 2.
func (tApp *TestApp) AdvanceToBlock(block uint32) sdk.Context {
	tApp.initChainIfNeeded()
	return tApp.AdvanceToBlockWithTime(block, tApp.header.Time)
}

// AdvanceToBlockWithTime advances the chain to the specified block and block time.
// For example, block = 10, t = 90 will advance to a block with height 10 and with a time of 90.
//
// Note that due to DEC-1248 the minimum block height is 2.
func (tApp *TestApp) AdvanceToBlockWithTime(block uint32, t time.Time) sdk.Context {
	if int64(block) <= tApp.header.Height {
		panic(fmt.Errorf("Expected block height (%d) > current block height (%d).", block, tApp.header.Height))
	}
	if t.UnixNano() < tApp.header.Time.UnixNano() {
		panic(fmt.Errorf("Expected time (%v) >= current block time (%v).", t, tApp.header.Time))
	}

	tApp.initChainIfNeeded()

	// First advance to the prior block using the current block time. This ensures that we only update the time on
	// the requested block.
	if int64(block)-tApp.header.Height > 1 && t != tApp.header.Time {
		tApp.AdvanceToBlock(block - 1)
	}

	// Ensure that we grab the lock so that we can read and write passingCheckTxs correctly.
	tApp.passingCheckTxsMtx.Lock()
	defer tApp.passingCheckTxsMtx.Unlock()

	// Advance to the requested block using the requested block time.
	for tApp.App.LastBlockHeight() < int64(block) {
		tApp.header.Height = tApp.App.LastBlockHeight() + 1
		tApp.header.Time = t
		tApp.header.LastCommitHash = tApp.App.LastCommitID().Hash
		tApp.header.NextValidatorsHash = tApp.App.LastCommitID().Hash

		// Prepare the proposal and process it.
		prepareProposalResponse := tApp.App.PrepareProposal(abcitypes.RequestPrepareProposal{
			Txs:                tApp.passingCheckTxs,
			MaxTxBytes:         math.MaxInt64,
			Height:             tApp.header.Height,
			Time:               tApp.header.Time,
			NextValidatorsHash: tApp.header.NextValidatorsHash,
			ProposerAddress:    tApp.header.ProposerAddress,
		})
		// Pass forward any transactions that were modified through the preparation phase and process them.
		processProposalRequest := abcitypes.RequestProcessProposal{
			Txs:                prepareProposalResponse.Txs,
			Hash:               tApp.header.AppHash,
			Height:             tApp.header.Height,
			Time:               tApp.header.Time,
			NextValidatorsHash: tApp.header.NextValidatorsHash,
			ProposerAddress:    tApp.header.ProposerAddress,
		}
		processProposalResponse := tApp.App.ProcessProposal(processProposalRequest)
		if tApp.builder.t == nil {
			if !processProposalResponse.IsAccepted() {
				panic(fmt.Errorf(
					"Expected process proposal request %+v to be accepted, but failed with %+v.",
					processProposalRequest,
					processProposalResponse,
				))
			}
		} else {
			require.Truef(
				tApp.builder.t,
				processProposalResponse.IsAccepted(),
				"Expected process proposal request %+v to be accepted, but failed with %+v.",
				processProposalRequest,
				processProposalResponse,
			)
		}
		txsNotInLastProposal := make([][]byte, 0)
		for _, tx := range tApp.passingCheckTxs {
			if !slices.ContainsFunc(prepareProposalResponse.Txs, func(tx2 []byte) bool {
				return bytes.Equal(tx, tx2)
			}) {
				txsNotInLastProposal = append(txsNotInLastProposal, tx)
			}
		}
		tApp.passingCheckTxs = txsNotInLastProposal

		// Start the next block
		tApp.App.BeginBlock(abcitypes.RequestBeginBlock{
			Header: tApp.header,
		})

		// Deliver the transaction from the previous block
		for _, bz := range prepareProposalResponse.Txs {
			deliverTxResponse := tApp.App.DeliverTx(abcitypes.RequestDeliverTx{Tx: bz})
			if tApp.builder.t == nil {
				if !deliverTxResponse.IsOK() {
					panic(fmt.Errorf(
						"Failed to deliver transaction that was accepted: %+v.",
						deliverTxResponse,
					))
				}
			} else {
				require.Truef(
					tApp.builder.t,
					deliverTxResponse.IsOK(),
					"Failed to deliver transaction that was accepted: %+v.",
					deliverTxResponse,
				)
			}
		}

		// End the block and commit it.
		tApp.App.EndBlock(abcitypes.RequestEndBlock{Height: tApp.header.Height})
		tApp.App.Commit()
	}

	return tApp.App.NewContext(true, tApp.header)
}

// Reset resets the chain such that it can be initialized and executed again.
func (tApp *TestApp) Reset() {
	tApp.App = nil
	tApp.genesis = types.GenesisDoc{}
	tApp.header = tmproto.Header{}
	tApp.passingCheckTxs = nil
}

// GetBlockHeight fetches the current block height of the test app.
func (tApp *TestApp) GetBlockHeight() int64 {
	return tApp.header.Height
}

// CheckTx adds the transaction to a test specific "mempool" that will be used to deliver the transaction during
// Prepare/Process proposal. Note that this must be invoked over TestApp.App.CheckTx as the transaction will not
// be added to the "mempool" causing the transaction to not be supplied during the Prepare/Process proposal phase.
//
// This method is thread-safe.
func (tApp *TestApp) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	res := tApp.App.CheckTx(req)
	// Note that the dYdX fork of CometBFT explicitly excludes place and cancel order messages. See
	// https://github.com/dydxprotocol/cometbft/blob/4d4d3b0/mempool/v0/clist_mempool.go#L416
	if res.IsOK() && !mempool.IsClobOrderTransaction(req.Tx, log.TestingLogger()) {
		// We want to ensure that we hold the lock only for updating passingCheckTxs so that App.CheckTx can execute
		// concurrently.
		tApp.passingCheckTxsMtx.Lock()
		defer tApp.passingCheckTxsMtx.Unlock()
		tApp.passingCheckTxs = append(tApp.passingCheckTxs, req.Tx)
	}
	return res
}

// MustMakeCheckTxsWithClobMsg creates one signed RequestCheckTx for each msg passed in.
// The messsage must use one of the hard-coded well known subaccount owners otherwise this will panic.
func MustMakeCheckTxsWithClobMsg[T clobtypes.MsgPlaceOrder | clobtypes.MsgCancelOrder](
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
			m = &v
		case clobtypes.MsgCancelOrder:
			m = &v
		default:
			panic(fmt.Errorf("MustMakeCheckTxsWithClobMsg: Unknown message type %T", msg))
		}

		msgSignerAddress := testtx.MustGetSignerAddress(m)
		if signerAddress == "" {
			signerAddress = msgSignerAddress
		}
		if signerAddress != msgSignerAddress {
			panic(fmt.Errorf("Input msgs must have the same owner/signer address"))
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

	checkTx, err := sims.GenSignedMockTx(
		rand.New(rand.NewSource(42)),
		app.TxConfig(),
		messages,
		sdk.Coins{},
		0,
		ctx.ChainID(),
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
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
