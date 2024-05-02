package app

import (
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	custommodule "github.com/StreamFinance-Protocol/stream-chain/protocol/app/module"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/configs"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	cosmosflags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	"google.golang.org/grpc"

	// App
	appconstants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/middleware"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/prepare"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/process"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	timelib "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/time"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/rate_limit"

	// Mempool
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mempool"

	// Daemons
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	liquidationclient "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/liquidation/client"
	metricsclient "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/metrics/client"
	pricefeedclient "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/constants"
	pricefeed_types "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/types"
	daemonserver "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server"
	daemonservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types"
	liquidationtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/liquidations"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	daemontypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/types"

	// Modules
	assetsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets"
	assetsmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/keeper"
	assetsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	blocktimemodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime"
	blocktimemodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/keeper"
	blocktimemoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	clobmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob"
	clobflags "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/flags"
	clobmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	clobmodulememclob "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	clobmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	delaymsgmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg"
	delaymsgmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/keeper"
	delaymsgmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	epochsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs"
	epochsmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/keeper"
	epochsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/types"
	feetiersmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers"
	feetiersmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/keeper"
	feetiersmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"

	perpetualsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals"
	perpetualsmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	perpetualsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricesmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices"
	pricesmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricesmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimitmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit"
	ratelimitmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	ratelimitmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	rewardsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards"
	rewardsmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/keeper"
	rewardsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/types"
	sendingmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending"
	sendingmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/keeper"
	sendingmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	statsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats"
	statsmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/keeper"
	statsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
	subaccountsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts"
	subaccountsmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	vestmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest"
	vestmodulekeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest/keeper"
	vestmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest/types"

	// IBC
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" // nolint:staticcheck
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	// Indexer
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/msgsender"

	// Grpc Streaming
	streaming "github.com/StreamFinance-Protocol/stream-chain/protocol/streaming/grpc"
	streamingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/streaming/grpc/types"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {

		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+appconstants.AppName)

	// Set DefaultPowerReduction to 1e18 to avoid overflow whe calculating
	// consensus power.
	sdk.DefaultPowerReduction = lib.PowerReduction
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry
	event             runtime.EventService
	closeOnce         func() error

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	AuthzKeeper      authzkeeper.Keeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	// IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCKeeper             *ibckeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	RatelimitKeeper       ratelimitmodulekeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	PricesKeeper pricesmodulekeeper.Keeper

	AssetsKeeper assetsmodulekeeper.Keeper

	BlockTimeKeeper blocktimemodulekeeper.Keeper

	DelayMsgKeeper delaymsgmodulekeeper.Keeper

	FeeTiersKeeper feetiersmodulekeeper.Keeper

	PerpetualsKeeper *perpetualsmodulekeeper.Keeper

	VestKeeper vestmodulekeeper.Keeper

	RewardsKeeper rewardsmodulekeeper.Keeper

	StatsKeeper statsmodulekeeper.Keeper

	SubaccountsKeeper subaccountsmodulekeeper.Keeper

	ClobKeeper *clobmodulekeeper.Keeper

	SendingKeeper sendingmodulekeeper.Keeper

	EpochsKeeper epochsmodulekeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	ModuleManager *module.Manager
	ModuleBasics  module.BasicManager

	// module configurator
	configurator module.Configurator

	IndexerEventManager  indexer_manager.IndexerEventManager
	GrpcStreamingManager streamingtypes.GrpcStreamingManager
	Server               *daemonserver.Server

	// startDaemons encapsulates the logic that starts all daemons and daemon services. This function contains a
	// closure of all relevant data structures that are shared with various keepers. Daemon services startup is
	// delayed until after the gRPC service is initialized so that the gRPC service will be available and the daemons
	// can correctly operate.
	startDaemons func()

	PriceFeedClient    *pricefeedclient.Client
	LiquidationsClient *liquidationclient.Client

	DaemonHealthMonitor *daemonservertypes.HealthMonitor
}

// assertAppPreconditions assert invariants required for an application to start.
func assertAppPreconditions() {
	// Check that the default power reduction is set correctly.
	if sdk.DefaultPowerReduction.BigInt().Cmp(big.NewInt(1_000_000_000_000_000_000)) != 0 {
		panic("DefaultPowerReduction is not set correctly")
	}
}

// New returns a reference to an initialized blockchain app
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	assertAppPreconditions()

	// dYdX specific command-line flags.
	appFlags := flags.GetFlagValuesFromOptions(appOpts)
	logger.Info("Parsed App flags", "Flags", appFlags)
	// Panic if this is not a full node and gRPC is disabled.
	if err := appFlags.Validate(); err != nil {
		panic(err)
	}

	initDatadogProfiler(logger, appFlags.DdAgentHost, appFlags.DdTraceAgentPort)

	encodingConfig := GetEncodingConfig()

	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(appconstants.AppName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		authzkeeper.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		crisistypes.StoreKey,
		slashingtypes.StoreKey,
		paramstypes.StoreKey,
		consensusparamtypes.StoreKey,
		upgradetypes.StoreKey,
		feegrant.StoreKey,
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
		ratelimitmoduletypes.StoreKey,
		icacontrollertypes.StoreKey,
		icahosttypes.StoreKey,
		evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		pricesmoduletypes.StoreKey,
		assetsmoduletypes.StoreKey,
		blocktimemoduletypes.StoreKey,
		feetiersmoduletypes.StoreKey,
		perpetualsmoduletypes.StoreKey,
		satypes.StoreKey,
		statsmoduletypes.StoreKey,
		vestmoduletypes.StoreKey,
		rewardsmoduletypes.StoreKey,
		clobmoduletypes.StoreKey,
		sendingmoduletypes.StoreKey,
		delaymsgmoduletypes.StoreKey,
		epochsmoduletypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(
		paramstypes.TStoreKey,
		clobmoduletypes.TransientStoreKey,
		statsmoduletypes.TransientStoreKey,
		rewardsmoduletypes.TransientStoreKey,
		indexer_manager.TransientStoreKey,
	)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, clobmoduletypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}
	app.closeOnce = sync.OnceValue[error](
		func() error {
			if app.PriceFeedClient != nil {
				app.PriceFeedClient.Stop()
			}
			if app.Server != nil {
				app.Server.Stop()
			}
			return nil
		},
	)

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		lib.GovModuleAddress.String(),
		app.event,
	)
	bApp.SetParamStore(&app.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		lib.GovModuleAddress.String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedAddresses(),
		lib.GovModuleAddress.String(),
		logger,
	)
	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		lib.GovModuleAddress.String(),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.StakingKeeper,
		lib.GovModuleAddress.String(),
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		lib.GovModuleAddress.String(),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[feegrant.StoreKey]),
		app.AccountKeeper,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.SlashingKeeper.Hooks()),
	)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(cosmosflags.FlagHome))
	// set the governance module account as the authority for conducting upgrades
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		lib.GovModuleAddress.String(),
	)

	// grant capabilities for the ibc, ibc-transfer, ICAHostKeeper and ratelimit modules

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedIBCTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	// scopedRatelimitKeeper is not used as an input to any other module.
	app.CapabilityKeeper.ScopeToModule(ratelimitmoduletypes.ModuleName)

	app.CapabilityKeeper.Seal()

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.getSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
		lib.GovModuleAddress.String(),
	)

	// Create ICA Host Keeper
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		keys[icahosttypes.StoreKey], // key
		app.getSubspace(icahosttypes.SubModuleName), // paramSpace
		app.IBCKeeper.ChannelKeeper,                 // ics4Wrapper, may be replaced with middleware such as ics29 fee
		app.IBCKeeper.ChannelKeeper,                 // channelKeeper
		app.IBCKeeper.PortKeeper,                    // portKeeper
		app.AccountKeeper,                           // accountKeeper
		scopedICAHostKeeper,                         // scopedKeeper
		app.MsgServiceRouter(),                      // msgRouter
		lib.GovModuleAddress.String(),               // authority
	)

	app.BlockTimeKeeper = *blocktimemodulekeeper.NewKeeper(
		appCodec,
		keys[blocktimemoduletypes.StoreKey],
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	blockTimeModule := blocktimemodule.NewAppModule(appCodec, app.BlockTimeKeeper)

	app.RatelimitKeeper = *ratelimitmodulekeeper.NewKeeper(
		appCodec,
		keys[ratelimitmoduletypes.StoreKey],
		app.BankKeeper,
		app.BlockTimeKeeper,
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	rateLimitModule := ratelimitmodule.NewAppModule(appCodec, app.RatelimitKeeper)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.getSubspace(ibctransfertypes.ModuleName),
		app.RatelimitKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedIBCTransferKeeper,
		lib.GovModuleAddress.String(),
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	// Wrap the x/ratelimit middlware over the IBC Transfer module
	var transferStack ibcporttypes.IBCModule = transferIBCModule
	transferStack = ratelimitmodule.NewIBCMiddleware(app.RatelimitKeeper, transferStack)

	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	// Ordering of `AddRoute` does not matter.
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)

	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.StakingKeeper,
		app.SlashingKeeper,
		addresscodec.NewBech32Codec(sdk.Bech32PrefixAccAddr),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	/****  dYdX specific modules/setup ****/
	msgSender, indexerFlags := getIndexerFromOptions(appOpts, logger)
	app.IndexerEventManager = indexer_manager.NewIndexerEventManager(
		msgSender,
		tkeys[indexer_manager.TransientStoreKey],
		indexerFlags.SendOffchainData,
	)

	app.GrpcStreamingManager = getGrpcStreamingManagerFromOptions(appFlags, logger)

	timeProvider := &timelib.TimeProviderImpl{}

	app.EpochsKeeper = *epochsmodulekeeper.NewKeeper(
		appCodec,
		keys[epochsmoduletypes.StoreKey],
	)
	epochsModule := epochsmodule.NewAppModule(appCodec, app.EpochsKeeper)

	// Get Daemon Flags.
	daemonFlags := daemonflags.GetDaemonFlagValuesFromOptions(appOpts)
	logger.Info("Parsed Daemon flags", "Flags", daemonFlags)

	// Create server that will ingest gRPC messages from daemon clients.
	// Note that gRPC clients will block on new gRPC connection until the gRPC server is ready to
	// accept new connections.
	app.Server = daemonserver.NewServer(
		logger,
		grpc.NewServer(),
		&daemontypes.FileHandlerImpl{},
		daemonFlags.Shared.SocketAddress,
	)
	// Setup server for pricefeed messages. The server will wait for gRPC messages containing price
	// updates and then encode them into an in-memory cache shared by the prices module.
	// The in-memory data structure is shared by the x/prices module and PriceFeed daemon.
	indexPriceCache := pricefeedtypes.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)
	app.Server.WithPriceFeedMarketToExchangePrices(indexPriceCache)

	// Setup server for liquidation messages. The server will wait for gRPC messages containing
	// potentially liquidatable subaccounts and then encode them into an in-memory slice shared by
	// the liquidations module.
	// The in-memory data structure is shared by the x/clob module and liquidations daemon.
	daemonLiquidationInfo := liquidationtypes.NewDaemonLiquidationInfo()
	app.Server.WithDaemonLiquidationInfo(daemonLiquidationInfo)

	app.DaemonHealthMonitor = daemonservertypes.NewHealthMonitor(
		daemonservertypes.DaemonStartupGracePeriod,
		daemonservertypes.HealthCheckPollFrequency,
		app.Logger(),
		daemonFlags.Shared.PanicOnDaemonFailureEnabled,
	)
	// Create a closure for starting daemons and daemon server. Daemon services are delayed until after the gRPC
	// service is started because daemons depend on the gRPC service being available. If a node is initialized
	// with a genesis time in the future, then the gRPC service will not be available until the genesis time, the
	// daemons will not be able to connect to the cosmos gRPC query service and finish initialization, and the daemon
	// monitoring service will panic.
	app.startDaemons = func() {
		maxDaemonUnhealthyDuration := time.Duration(daemonFlags.Shared.MaxDaemonUnhealthySeconds) * time.Second
		// Start server for handling gRPC messages from daemons.
		go app.Server.Start()

		// Start liquidations client for sending potentially liquidatable subaccounts to the application.
		if daemonFlags.Liquidation.Enabled {
			app.LiquidationsClient = liquidationclient.NewClient(logger)
			go func() {
				app.RegisterDaemonWithHealthMonitor(app.LiquidationsClient, maxDaemonUnhealthyDuration)
				if err := app.LiquidationsClient.Start(
					// The client will use `context.Background` so that it can have a different context from
					// the main application.
					context.Background(),
					daemonFlags,
					appFlags,
					&daemontypes.GrpcClientImpl{},
				); err != nil {
					panic(err)
				}
			}()
		}

		// Non-validating full-nodes have no need to run the price daemon.
		if !appFlags.NonValidatingFullNode && daemonFlags.Price.Enabled {
			exchangeQueryConfig := configs.ReadExchangeQueryConfigFile(homePath)
			// Start pricefeed client for sending prices for the pricefeed server to consume. These prices
			// are retrieved via third-party APIs like Binance and then are encoded in-memory and
			// periodically sent via gRPC to a shared socket with the server.
			app.PriceFeedClient = pricefeedclient.StartNewClient(
				// The client will use `context.Background` so that it can have a different context from
				// the main application.
				context.Background(),
				daemonFlags,
				appFlags,
				logger,
				&daemontypes.GrpcClientImpl{},
				exchangeQueryConfig,
				constants.StaticExchangeDetails,
				&pricefeedclient.SubTaskRunnerImpl{},
			)
			app.RegisterDaemonWithHealthMonitor(app.PriceFeedClient, maxDaemonUnhealthyDuration)
		}

		// Start the Metrics Daemon.
		// The metrics daemon is purely used for observability. It should never bring the app down.
		// TODO(CLOB-960) Don't start this goroutine if telemetry is disabled
		// Note: the metrics daemon is such a simple go-routine that we don't bother implementing a health-check
		// for this service. The task loop does not produce any errors because the telemetry calls themselves are
		// not error-returning, so in effect this daemon would never become unhealthy.
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						"Metrics Daemon exited unexpectedly with a panic.",
						"panic",
						r,
						"stack",
						string(debug.Stack()),
					)
				}
			}()
			metricsclient.Start(
				// The client will use `context.Background` so that it can have a different context from
				// the main application.
				context.Background(),
				logger,
			)
		}()
	}

	app.PricesKeeper = *pricesmodulekeeper.NewKeeper(
		appCodec,
		keys[pricesmoduletypes.StoreKey],
		indexPriceCache,
		pricesmoduletypes.NewMarketToSmoothedPrices(pricesmoduletypes.SmoothedPriceTrackingBlockHistoryLength),
		timeProvider,
		app.IndexerEventManager,
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	pricesModule := pricesmodule.NewAppModule(appCodec, app.PricesKeeper, app.AccountKeeper, app.BankKeeper)

	app.AssetsKeeper = *assetsmodulekeeper.NewKeeper(
		appCodec,
		keys[assetsmoduletypes.StoreKey],
		app.PricesKeeper,
		app.IndexerEventManager,
	)
	assetsModule := assetsmodule.NewAppModule(appCodec, app.AssetsKeeper)

	app.DelayMsgKeeper = *delaymsgmodulekeeper.NewKeeper(
		appCodec,
		keys[delaymsgmoduletypes.StoreKey],
		bApp.MsgServiceRouter(),
		// Permit delayed messages to be signed by the following modules.
		[]string{
			lib.GovModuleAddress.String(),
		},
	)
	delayMsgModule := delaymsgmodule.NewAppModule(appCodec, app.DelayMsgKeeper)

	app.PerpetualsKeeper = perpetualsmodulekeeper.NewKeeper(
		appCodec,
		keys[perpetualsmoduletypes.StoreKey],
		app.PricesKeeper,
		app.EpochsKeeper,
		app.IndexerEventManager,
		// gov module and delayMsg module accounts are allowed to send messages to the bridge module.
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	perpetualsModule := perpetualsmodule.NewAppModule(appCodec, app.PerpetualsKeeper)

	app.StatsKeeper = *statsmodulekeeper.NewKeeper(
		appCodec,
		app.EpochsKeeper,
		keys[statsmoduletypes.StoreKey],
		tkeys[statsmoduletypes.TransientStoreKey],
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	statsModule := statsmodule.NewAppModule(appCodec, app.StatsKeeper)

	app.FeeTiersKeeper = *feetiersmodulekeeper.NewKeeper(
		appCodec,
		app.StatsKeeper,
		keys[feetiersmoduletypes.StoreKey],
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	feeTiersModule := feetiersmodule.NewAppModule(appCodec, app.FeeTiersKeeper)

	app.VestKeeper = *vestmodulekeeper.NewKeeper(
		appCodec,
		keys[vestmoduletypes.StoreKey],
		app.BankKeeper,
		app.BlockTimeKeeper,
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	vestModule := vestmodule.NewAppModule(appCodec, app.VestKeeper)

	app.RewardsKeeper = *rewardsmodulekeeper.NewKeeper(
		appCodec,
		keys[rewardsmoduletypes.StoreKey],
		tkeys[rewardsmoduletypes.TransientStoreKey],
		app.AssetsKeeper,
		app.BankKeeper,
		app.FeeTiersKeeper,
		app.PricesKeeper,
		app.IndexerEventManager,
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	rewardsModule := rewardsmodule.NewAppModule(appCodec, app.RewardsKeeper)

	app.SubaccountsKeeper = *subaccountsmodulekeeper.NewKeeper(
		appCodec,
		keys[satypes.StoreKey],
		app.AssetsKeeper,
		app.BankKeeper,
		app.PerpetualsKeeper,
		app.BlockTimeKeeper,
		app.IndexerEventManager,
	)
	subaccountsModule := subaccountsmodule.NewAppModule(
		appCodec,
		app.SubaccountsKeeper,
	)

	clobFlags := clobflags.GetClobFlagValuesFromOptions(appOpts)
	logger.Info("Parsed CLOB flags", "Flags", clobFlags)

	memClob := clobmodulememclob.NewMemClobPriceTimePriority(app.IndexerEventManager.Enabled())
	memClob.SetGenerateOrderbookUpdates(app.GrpcStreamingManager.Enabled())

	app.ClobKeeper = clobmodulekeeper.NewKeeper(
		appCodec,
		keys[clobmoduletypes.StoreKey],
		memKeys[clobmoduletypes.MemStoreKey],
		tkeys[clobmoduletypes.TransientStoreKey],
		// set the governance and delaymsg module accounts as the authority for conducting upgrades
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
		memClob,
		app.SubaccountsKeeper,
		app.AssetsKeeper,
		app.BlockTimeKeeper,
		app.BankKeeper,
		app.FeeTiersKeeper,
		app.PerpetualsKeeper,
		app.PricesKeeper,
		app.StatsKeeper,
		app.RewardsKeeper,
		app.IndexerEventManager,
		app.GrpcStreamingManager,
		txConfig.TxDecoder(),
		clobFlags,
		rate_limit.NewPanicRateLimiter[*clobmoduletypes.MsgPlaceOrder](),
		rate_limit.NewPanicRateLimiter[*clobmoduletypes.MsgCancelOrder](),
		daemonLiquidationInfo,
	)
	clobModule := clobmodule.NewAppModule(
		appCodec,
		app.ClobKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.SubaccountsKeeper,
	)
	app.PerpetualsKeeper.SetClobKeeper(app.ClobKeeper)

	app.SendingKeeper = *sendingmodulekeeper.NewKeeper(
		appCodec,
		keys[sendingmoduletypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.SubaccountsKeeper,
		app.IndexerEventManager,
		// gov module and delayMsg module accounts are allowed to send messages to the sending module.
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
	)
	sendingModule := sendingmodule.NewAppModule(
		appCodec,
		app.SendingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.SubaccountsKeeper,
	)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.ModuleManager = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.getSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.getSubspace(banktypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.getSubspace(crisistypes.ModuleName)),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.StakingKeeper,
			app.getSubspace(slashingtypes.ModuleName),
			app.interfaceRegistry,
		),

		staking.NewAppModule(
			appCodec,
			app.StakingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.getSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(app.UpgradeKeeper, addresscodec.NewBech32Codec(sdk.Bech32PrefixAccAddr)),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		transferModule,
		pricesModule,
		assetsModule,
		blockTimeModule,
		feeTiersModule,
		perpetualsModule,
		statsModule,
		vestModule,
		rewardsModule,
		subaccountsModule,
		clobModule,
		sendingModule,
		delayMsgModule,
		epochsModule,
		rateLimitModule,
	)

	app.ModuleManager.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.ModuleManager.SetOrderBeginBlockers(
		blocktimemoduletypes.ModuleName, // Must be first
		authz.ModuleName,                // Delete expired grants.
		epochsmoduletypes.ModuleName,
		capabilitytypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ratelimitmoduletypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		pricesmoduletypes.ModuleName,
		assetsmoduletypes.ModuleName,
		feetiersmoduletypes.ModuleName,
		perpetualsmoduletypes.ModuleName,
		statsmoduletypes.ModuleName,
		satypes.ModuleName,
		clobmoduletypes.ModuleName,
		vestmoduletypes.ModuleName,
		rewardsmoduletypes.ModuleName,
		sendingmoduletypes.ModuleName,
		delaymsgmoduletypes.ModuleName,
	)

	app.ModuleManager.SetOrderPrepareCheckStaters(
		clobmoduletypes.ModuleName,
	)

	app.ModuleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ratelimitmoduletypes.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		pricesmoduletypes.ModuleName,
		assetsmoduletypes.ModuleName,
		feetiersmoduletypes.ModuleName,
		perpetualsmoduletypes.ModuleName,
		statsmoduletypes.ModuleName,
		satypes.ModuleName,
		clobmoduletypes.ModuleName,
		sendingmoduletypes.ModuleName,
		vestmoduletypes.ModuleName,
		rewardsmoduletypes.ModuleName,
		epochsmoduletypes.ModuleName,
		delaymsgmoduletypes.ModuleName,
		authz.ModuleName,                // No-op.
		blocktimemoduletypes.ModuleName, // Must be last
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.ModuleManager.SetOrderInitGenesis(
		epochsmoduletypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		ibctransfertypes.ModuleName,
		ratelimitmoduletypes.ModuleName,
		feegrant.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		pricesmoduletypes.ModuleName,
		assetsmoduletypes.ModuleName,
		blocktimemoduletypes.ModuleName,
		feetiersmoduletypes.ModuleName,
		perpetualsmoduletypes.ModuleName,
		statsmoduletypes.ModuleName,
		satypes.ModuleName,
		clobmoduletypes.ModuleName,
		vestmoduletypes.ModuleName,
		rewardsmoduletypes.ModuleName,
		sendingmoduletypes.ModuleName,
		delaymsgmoduletypes.ModuleName,
		authz.ModuleName,
	)

	// NOTE: by default, set migration order here to be the same as init genesis order,
	// in case there are dependencies between modules.
	// x/auth is run last since it depends on the x/staking module.
	app.ModuleManager.SetOrderMigrations(
		epochsmoduletypes.ModuleName,
		capabilitytypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		ibctransfertypes.ModuleName,
		ratelimitmoduletypes.ModuleName,
		feegrant.ModuleName,
		consensusparamtypes.ModuleName,
		icatypes.ModuleName,
		pricesmoduletypes.ModuleName,
		assetsmoduletypes.ModuleName,
		blocktimemoduletypes.ModuleName,
		feetiersmoduletypes.ModuleName,
		perpetualsmoduletypes.ModuleName,
		statsmoduletypes.ModuleName,
		satypes.ModuleName,
		clobmoduletypes.ModuleName,
		vestmoduletypes.ModuleName,
		rewardsmoduletypes.ModuleName,
		sendingmoduletypes.ModuleName,
		delaymsgmoduletypes.ModuleName,
		authz.ModuleName,
		// Auth must be migrated after staking.
		authtypes.ModuleName,
	)

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err := app.ModuleManager.RegisterServices(app.configurator)
	app.ModuleBasics = module.NewBasicManagerFromManager(
		app.ModuleManager,
		map[string]module.AppModuleBasic{
			custommodule.SlashingModuleBasic{}.Name(): custommodule.SlashingModuleBasic{},
		},
	)
	if err != nil {
		panic(err)
	}

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.setAnteHandler(encodingConfig.TxConfig)
	app.SetMempool(mempool.NewNoOpMempool())
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetPrecommiter(app.Precommitter)
	app.SetPrepareCheckStater(app.PrepareCheckStater)

	// PrepareProposal setup.
	if appFlags.NonValidatingFullNode {
		app.SetPrepareProposal(prepare.FullNodePrepareProposalHandler())
	} else {
		app.SetPrepareProposal(
			prepare.PrepareProposalHandler(
				txConfig,
				app.ClobKeeper,
				app.PricesKeeper,
				app.PerpetualsKeeper,
			),
		)
	}

	// ProcessProposal setup.
	if appFlags.NonValidatingFullNode {
		// Note: If the command-line flag `--non-validating-full-node` is enabled, this node will use
		// an implementation of `ProcessProposal` which always returns `abci.ResponseProcessProposal_ACCEPT`.
		// Full-nodes do not participate in consensus, and therefore should not participate in voting / `ProcessProposal`.
		app.SetProcessProposal(
			process.FullNodeProcessProposalHandler(
				txConfig,
				app.ClobKeeper,
				app.StakingKeeper,
				app.PerpetualsKeeper,
				app.PricesKeeper,
			),
		)
	} else {
		app.SetProcessProposal(
			process.ProcessProposalHandler(
				txConfig,
				app.ClobKeeper,
				app.StakingKeeper,
				app.PerpetualsKeeper,
				app.PricesKeeper,
			),
		)
	}

	// Note that panics from out of gas errors won't get logged, since the `OutOfGasMiddleware` is added in front of this,
	// so error will get handled by that middleware and subsequent middlewares won't get executed.
	// Also note that `AddRunTxRecoveryHandler` adds the handler in reverse order, meaning that handlers that appear
	// earlier in the list will get executed later in the chain.
	// The chain of middlewares is shared between `DeliverTx` and `CheckTx`; in order to provide additional metadata
	// based on execution context such as the block proposer, the logger used by the logging middleware is
	// stored in a global variable and can be overwritten as necessary.
	middleware.Logger = logger
	app.AddRunTxRecoveryHandler(middleware.NewRunTxPanicLoggingMiddleware(app.ModuleBasics))

	// Set handlers and store loaders for upgrades.
	app.setupUpgradeHandlers()
	app.setupUpgradeStoreLoaders()

	// Currently the only case that exists where the app is _not_ started with loadLatest=true is when state is
	// loaded and then immediately exported to a file. In those cases, `LoadHeight` within `app.go` is called instead.
	// This behavior can be invoked via running `dydxprotocold export`, which exports the chain state to a JSON file.
	// In the export case, the memclob does not need to be hydrated, as it is never used.
	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Hydrate memStores used for caching state.
		app.hydrateMemStores()

		// Hydrate the `memclob` with all ordersbooks from state,
		// and hydrate the next `checkState` as well as the `memclob` with stateful orders.
		app.hydrateMemclobWithOrderbooksAndStatefulOrders()

		// Hydrate the keeper in-memory data structures.
		app.hydrateKeeperInMemoryDataStructures()
	}
	app.initializeRateLimiters()

	// Report out app version and git commit. This will be run when validators restart.
	version := version.NewInfo()
	app.Logger().Info(
		"App instantiated",
		metrics.AppVersion,
		version.Version,
		metrics.GitCommit,
		version.GitCommit,
	)

	return app
}

// RegisterDaemonWithHealthMonitor registers a daemon service with the update monitor, which will commence monitoring
// the health of the daemon. If the daemon does not register, the method will panic.
func (app *App) RegisterDaemonWithHealthMonitor(
	healthCheckableDaemon daemontypes.HealthCheckable,
	maxDaemonUnhealthyDuration time.Duration,
) {
	if err := app.DaemonHealthMonitor.RegisterService(healthCheckableDaemon, maxDaemonUnhealthyDuration); err != nil {
		app.Logger().Error(
			"Failed to register daemon service with update monitor",
			"error",
			err,
			"service",
			healthCheckableDaemon.ServiceName(),
			"maxDaemonUnhealthyDuration",
			maxDaemonUnhealthyDuration,
		)
		panic(err)
	}
}

// DisableHealthMonitorForTesting disables the health monitor for testing.
func (app *App) DisableHealthMonitorForTesting() {
	app.DaemonHealthMonitor.DisableForTesting()
}

// hydrateMemStores hydrates the memStores used for caching state.
func (app *App) hydrateMemStores() {
	// Create an `uncachedCtx` where the underlying MultiStore is the `rootMultiStore`.
	// We use this to hydrate the `memStore` state with values from the underlying `rootMultiStore`.
	uncachedCtx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	// Initialize memstore in clobKeeper with order fill amounts and stateful orders.
	app.ClobKeeper.InitMemStore(uncachedCtx)
}

// initializeRateLimiters initializes the rate limiters from state if the application is
// not started from genesis.
func (app *App) initializeRateLimiters() {
	// Create an `uncachedCtx` where the underlying MultiStore is the `rootMultiStore`.
	// We use this to hydrate the `orderRateLimiter` with values from the underlying `rootMultiStore`.
	uncachedCtx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	app.ClobKeeper.InitalizeBlockRateLimitFromStateIfExists(uncachedCtx)
}

// hydrateMemclobWithOrderbooksAndStatefulOrders hydrates the memclob with orderbooks and stateful orders
// from state.
func (app *App) hydrateMemclobWithOrderbooksAndStatefulOrders() {
	// Create a `checkStateCtx` where the underlying MultiStore is the `CacheMultiStore` for
	// the `checkState`. We do this to avoid performing any state writes to the `rootMultiStore`
	// directly.
	checkStateCtx := app.BaseApp.NewContext(true)

	// Initialize memclob in clobKeeper with orderbooks using `ClobPairs` in state.
	app.ClobKeeper.InitMemClobOrderbooks(checkStateCtx)
	// Initialize memclob with all existing stateful orders.
	// TODO(DEC-1348): Emit indexer messages to indicate that application restarted.
	app.ClobKeeper.InitStatefulOrders(checkStateCtx)
}

// hydrateKeeperInMemoryDataStructures hydrates the keeper with ClobPairId and PerpetualId mapping
// and untriggered conditional orders from state.
func (app *App) hydrateKeeperInMemoryDataStructures() {
	// Create a `checkStateCtx` where the underlying MultiStore is the `CacheMultiStore` for
	// the `checkState`. We do this to avoid performing any state writes to the `rootMultiStore`
	// directly.
	checkStateCtx := app.BaseApp.NewContext(true)

	// Initialize the untriggered conditional orders data structure with untriggered
	// conditional orders in state.
	app.ClobKeeper.HydrateClobPairAndPerpetualMapping(checkStateCtx)
	// Initialize the untriggered conditional orders data structure with untriggered
	// conditional orders in state.
	app.ClobKeeper.HydrateUntriggeredConditionalOrders(checkStateCtx)
}

// GetBaseApp returns the base app of the application
func (app *App) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// PreBlocker application updates before each begin block.
func (app *App) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.ModuleManager.PreBlock(ctx)
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	ctx = ctx.WithExecMode(lib.ExecModeBeginBlock)

	// Update the proposer address in the logger for the panic logging middleware.
	proposerAddr := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	middleware.Logger = ctx.Logger().With("proposer_cons_addr", proposerAddr.String())

	app.scheduleForkUpgrade(ctx)
	return app.ModuleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	ctx = ctx.WithExecMode(lib.ExecModeEndBlock)

	// Reset the logger for middleware.
	// Note that the middleware is only used by `CheckTx` and `DeliverTx`, and not `EndBlocker`.
	// Panics from `EndBlocker` will not be logged by the middleware and will lead to consensus failures.
	middleware.Logger = app.Logger()

	response, err := app.ModuleManager.EndBlock(ctx)
	if err != nil {
		return response, err
	}
	block := app.IndexerEventManager.ProduceBlock(ctx)
	app.IndexerEventManager.SendOnchainData(block)
	return response, err
}

// Precommitter application updates before the commital of a block after all transactions have been delivered.
func (app *App) Precommitter(ctx sdk.Context) {
	if err := app.ModuleManager.Precommit(ctx); err != nil {
		panic(err)
	}
}

// PrepareCheckStater application updates after commit and before any check state is invoked.
func (app *App) PrepareCheckStater(ctx sdk.Context) {
	ctx = ctx.WithExecMode(lib.ExecModePrepareCheckState)

	if err := app.ModuleManager.PrepareCheckState(ctx); err != nil {
		panic(err)
	}
}

// InitChainer application update at chain initialization.
func (app *App) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
	if err != nil {
		panic(err)
	}
	initResponse, err := app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)
	if err != nil {
		panic(err)
	}
	block := app.IndexerEventManager.ProduceBlock(ctx)
	app.IndexerEventManager.SendOnchainData(block)
	app.IndexerEventManager.ClearEvents(ctx)

	app.Logger().Info("Initialized chain", "blockHeight", ctx.BlockHeight())
	return initResponse, err
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns app's TxConfig
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *App) DefaultGenesis() map[string]json.RawMessage {
	return app.ModuleBasics.DefaultGenesis(app.appCodec)
}

// getSubspace returns a param subspace for a given module name.
func (app *App) getSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	app.ModuleBasics.RegisterGRPCGatewayRoutes(
		clientCtx,
		apiSvr.GRPCGatewayRouter,
	)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}

	// Now that the API server has been configured, start the daemons.
	app.startDaemons()
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService registers the node service.
func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// SimulationManager always returns nil.
func (app *App) SimulationManager() *module.SimulationManager {
	return nil
}

// buildAnteHandler builds an AnteHandler object configured for the app.
func (app *App) buildAnteHandler(txConfig client.TxConfig) sdk.AnteHandler {
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				FeegrantKeeper:  app.FeeGrantKeeper,
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			ClobKeeper: app.ClobKeeper,
		},
	)
	if err != nil {
		panic(err)
	}

	return anteHandler
}

// setAnteHandler creates a new AnteHandler and sets it on the base app and clob keeper.
func (app *App) setAnteHandler(txConfig client.TxConfig) {
	anteHandler := app.buildAnteHandler(txConfig)
	// Prevent a cycle between when we create the clob keeper and the ante handler.
	app.ClobKeeper.SetAnteHandler(anteHandler)
	app.SetAnteHandler(anteHandler)
}

// Close invokes an ordered shutdown of routines.
func (app *App) Close() error {
	app.BaseApp.Close()
	return app.closeOnce()
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key,
	tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)

	// register the key tables for legacy param subspaces
	keyTable := ibcclient.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	paramsKeeper.Subspace(icahosttypes.SubModuleName).WithKeyTable(icahosttypes.ParamKeyTable())
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName).WithKeyTable(icacontrollertypes.ParamKeyTable())

	return paramsKeeper
}

// getIndexerFromOptions returns an instance of a msgsender.IndexerMessageSender from the specified options.
// This function will default to try to use any instance that is configured for test execution followed by loading
// an instance from command line flags and finally returning a no-op instance.
func getIndexerFromOptions(
	appOpts servertypes.AppOptions,
	logger log.Logger,
) (msgsender.IndexerMessageSender, indexer.IndexerFlags) {
	v, ok := appOpts.Get(indexer.MsgSenderInstanceForTest).(msgsender.IndexerMessageSender)
	if ok {
		return v, indexer.IndexerFlags{
			SendOffchainData: true,
		}
	}

	indexerFlags := indexer.GetIndexerFlagValuesFromOptions(appOpts)
	logger.Info(
		"Parsed Indexer flags",
		"Flags", indexerFlags,
	)

	var indexerMessageSender msgsender.IndexerMessageSender
	if len(indexerFlags.KafkaAddrs) == 0 {
		indexerMessageSender = msgsender.NewIndexerMessageSenderNoop()
	} else {
		var err error
		indexerMessageSender, err = msgsender.NewIndexerMessageSenderKafka(
			indexerFlags,
			nil,
			logger,
		)
		if err != nil {
			panic(err)
		}
	}
	return indexerMessageSender, indexerFlags
}

// getGrpcStreamingManagerFromOptions returns an instance of a streamingtypes.GrpcStreamingManager from the specified
// options. This function will default to returning a no-op instance.
func getGrpcStreamingManagerFromOptions(
	appFlags flags.Flags,
	logger log.Logger,
) (manager streamingtypes.GrpcStreamingManager) {
	if appFlags.GrpcStreamingEnabled {
		logger.Info("GRPC streaming is enabled")
		return streaming.NewGrpcStreamingManager()
	}
	return streaming.NewNoopGrpcStreamingManager()
}
