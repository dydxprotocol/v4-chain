package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkproto "github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4/app"
	"github.com/dydxprotocol/v4/testutil/appoptions"
	"github.com/dydxprotocol/v4/testutil/constants"
	assettypes "github.com/dydxprotocol/v4/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4/x/epochs/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4/x/sending/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

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

// A TestApp used to executed ABCI++ flows.
type TestApp struct {
	builder TestAppBuilder
	app     *app.App
	genesis types.GenesisDoc
	header  tmproto.Header
}

func (tApp *TestApp) Builder() TestAppBuilder {
	return tApp.builder
}

// InitChain initializes the chain. Will panic if initialized more than once.
func (tApp *TestApp) InitChain() (sdk.Context, *app.App) {
	if tApp.app != nil {
		panic(errors.New("Cannot initialize chain that has been initialized already. Missing a Reset()?"))
	}
	tApp.initChainIfNeeded()
	return tApp.app.NewContext(true, tApp.header), tApp.app
}

func (tApp *TestApp) initChainIfNeeded() {
	if tApp.app != nil {
		return
	}

	// Get the initial genesis state and initialize the chain and commit the results of the initialization.
	tApp.genesis = tApp.builder.genesisDocFn()
	tApp.app = tApp.builder.appCreatorFn()
	baseapp.SetChainID(tApp.genesis.ChainID)(tApp.app.GetBaseApp())
	if tApp.genesis.GenesisTime.UnixNano() <= time.UnixMilli(0).UnixNano() {
		panic(fmt.Errorf(
			"Unable to start chain at time %v, must be greater than unix epoch.",
			tApp.genesis.GenesisTime,
		))
	}

	consensusParamsProto := tApp.genesis.ConsensusParams.ToProto()

	tApp.app.InitChain(abcitypes.RequestInitChain{
		InitialHeight:   tApp.genesis.InitialHeight,
		AppStateBytes:   tApp.genesis.AppState,
		ChainId:         tApp.genesis.ChainID,
		ConsensusParams: &consensusParamsProto,
		Time:            tApp.genesis.GenesisTime,
	})
	tApp.app.Commit()

	tApp.header = tmproto.Header{
		ChainID:            tApp.genesis.ChainID,
		ProposerAddress:    constants.AliceAccAddress,
		Height:             tApp.app.LastBlockHeight(),
		Time:               tApp.genesis.GenesisTime,
		LastCommitHash:     tApp.app.LastCommitID().Hash,
		NextValidatorsHash: tApp.app.LastCommitID().Hash,
	}
}

// AdvanceToBlockIfNecessary advances the chain to the specified block using the current block time.
// If the block is the same, then this function results in a no-op.
// Note that due to DEC-1248 the minimum block height is 2.
func (tApp *TestApp) AdvanceToBlockIfNecessary(block uint32) (sdk.Context, *app.App) {
	if int64(block) == tApp.GetBlockHeight() {
		return tApp.app.NewContext(true, tApp.header), tApp.app
	}

	return tApp.AdvanceToBlock(block)
}

// AdvanceToBlock advances the chain to the specified block using the current block time.
//
// Note that due to DEC-1248 the minimum block height is 2.
func (tApp *TestApp) AdvanceToBlock(block uint32) (sdk.Context, *app.App) {
	tApp.initChainIfNeeded()
	return tApp.AdvanceToBlockWithTime(block, tApp.header.Time)
}

// AdvanceToBlock advances the chain to the specified block and time.
//
// Note that due to DEC-1248 the minimum block height is 2.
func (tApp *TestApp) AdvanceToBlockWithTime(block uint32, t time.Time) (sdk.Context, *app.App) {
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

	// Advance to the requested block using the requested block time.
	for tApp.app.LastBlockHeight() < int64(block) {
		tApp.header.Height = tApp.app.LastBlockHeight() + 1
		tApp.header.Time = t
		tApp.header.LastCommitHash = tApp.app.LastCommitID().Hash
		tApp.header.NextValidatorsHash = tApp.app.LastCommitID().Hash

		// Prepare the proposal and process it.
		prepareProposalResponse := tApp.app.PrepareProposal(abcitypes.RequestPrepareProposal{
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
		processProposalResponse := tApp.app.ProcessProposal(processProposalRequest)
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

		// Start the next block
		tApp.app.BeginBlock(abcitypes.RequestBeginBlock{
			Header: tApp.header,
		})

		// Deliver the transaction from the previous block
		for _, bz := range prepareProposalResponse.Txs {
			deliverTxResponse := tApp.app.DeliverTx(abcitypes.RequestDeliverTx{Tx: bz})
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
		tApp.app.EndBlock(abcitypes.RequestEndBlock{Height: tApp.header.Height})
		tApp.app.Commit()
	}

	return tApp.app.NewContext(true, tApp.header), tApp.app
}

// Reset resets the chain such that it can be initialized and executed again.
func (tApp *TestApp) Reset() {
	tApp.app = nil
	tApp.genesis = types.GenesisDoc{}
	tApp.header = tmproto.Header{}
}

// GetBlockHeight fetches the current block height of the test app.
func (tApp *TestApp) GetBlockHeight() int64 {
	return tApp.header.Height
}

// MustMakeCheckTxs creates one signed RequestCheckTx for each msg passed in. The messsage must use one of the
// hard-coded well known subaccount owners otherwise this will panic.
func MustMakeCheckTxs[T clobtypes.MsgPlaceOrder | clobtypes.MsgCancelOrder](
	ctx sdk.Context,
	app *app.App,
	messages ...T,
) []abcitypes.RequestCheckTx {
	sdkMessages := make([]sdk.Msg, len(messages))
	for i, msg := range messages {
		var m sdk.Msg
		switch v := any(msg).(type) {
		case clobtypes.MsgPlaceOrder:
			m = &v
		case clobtypes.MsgCancelOrder:
			m = &v
		default:
			panic(fmt.Errorf("Unknown message type %T", msg))
		}

		sdkMessages[i] = m
	}

	return MustMakeChecksTxsWithSdkMsg(ctx, app, sdkMessages...)
}

// MustMakeChecksTxsWithSdkMsg creates one signed RequestCheckTx for each msg passed in. The messsage must use one of
// the hard-coded well known subaccount owners otherwise this will panic.
func MustMakeChecksTxsWithSdkMsg(
	ctx sdk.Context,
	app *app.App,
	messages ...sdk.Msg,
) (checkTxs []abcitypes.RequestCheckTx) {
	for _, msg := range messages {
		checkTxs = append(checkTxs, MustMakeCheckTx(ctx, app, msg))
	}

	return checkTxs
}

// MustMakeCheckTx creates a signed RequestCheckTx for the provided message. The message must use one of the
// hard-coded well known subaccount owners otherwise this will panic.
func MustMakeCheckTx(
	ctx sdk.Context,
	app *app.App,
	message sdk.Msg,
) abcitypes.RequestCheckTx {
	var subAccountOwner string
	switch v := any(message).(type) {
	case *clobtypes.MsgPlaceOrder:
		subAccountOwner = v.Order.OrderId.SubaccountId.Owner
	case *clobtypes.MsgCancelOrder:
		subAccountOwner = v.OrderId.SubaccountId.Owner
	default:
		panic(fmt.Errorf("Unknown message type %T", message))
	}

	var privKey cryptotypes.PrivKey
	var accAddress sdk.AccAddress
	switch subAccountOwner {
	case constants.AliceAccAddress.String():
		accAddress = constants.AliceAccAddress
		privKey = constants.AlicePrivateKey
	case constants.BobAccAddress.String():
		accAddress = constants.BobAccAddress
		privKey = constants.BobPrivateKey
	case constants.CarlAccAddress.String():
		accAddress = constants.CarlAccAddress
		privKey = constants.CarlPrivateKey
	case constants.DaveAccAddress.String():
		accAddress = constants.DaveAccAddress
		privKey = constants.DavePrivateKey
	}
	if !app.AccountKeeper.HasAccount(ctx, accAddress) {
		panic("Account not found")
	}
	account := app.AccountKeeper.GetAccount(ctx, accAddress)

	checkTx, err := sims.GenSignedMockTx(
		rand.New(rand.NewSource(42)),
		app.TxConfig(),
		[]sdk.Msg{message},
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
