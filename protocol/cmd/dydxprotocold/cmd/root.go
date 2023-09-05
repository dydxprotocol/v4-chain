package cmd

import (
	"cosmossdk.io/store/snapshots"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cosmossdk.io/simapp/params"
	rosettaCmd "cosmossdk.io/tools/rosetta/cmd"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	gometrics "github.com/armon/go-metrics"
	dbm "github.com/cometbft/cometbft-db"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	dydxapp "github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib/encoding"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

const (
	EnvPrefix = "DYDX"

	flagIAVLCacheSize = "iavl-cache-size"
)

// TODO(DEC-1097): improve `cmd/` by adding tests, custom app configs, custom init cmd, and etc.
// NewRootCmd creates a new root command for `dydxprotocold`. It is called once in the main function.
func NewRootCmd(option *RootCmdOption) *cobra.Command {
	encodingConfig := encoding.MakeEncodingConfig(basic_manager.ModuleBasics)
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(dydxapp.DefaultNodeHome).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   dydxapp.AppDaemonName,
		Short: "Start dydxprotocol app",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			if err := server.InterceptConfigsPreRunHandler(
				cmd,
				customAppTemplate,
				customAppConfig,
				customTMConfig,
			); err != nil {
				return err
			}

			return nil
		},
		SilenceUsage: true,
	}

	initRootCmd(rootCmd, option, encodingConfig)

	return rootCmd
}

// initRootCmd initializes the app's root command with useful commands.
func initRootCmd(rootCmd *cobra.Command, option *RootCmdOption, encodingConfig params.EncodingConfig) {
	gentxModule := basic_manager.ModuleBasics[genutiltypes.ModuleName].(genutil.AppModuleBasic)
	rootCmd.AddCommand(
		genutilcli.InitCmd(basic_manager.ModuleBasics, dydxapp.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, dydxapp.DefaultNodeHome, gentxModule.GenTxValidator),
		genutilcli.MigrateGenesisCmd(),
		genutilcli.GenTxCmd(
			basic_manager.ModuleBasics,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			dydxapp.DefaultNodeHome,
		),
		genutilcli.ValidateGenesisCmd(basic_manager.ModuleBasics),
		AddGenesisAccountCmd(dydxapp.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		config.Cmd(),
	)

	a := appCreator{encodingConfig}
	server.AddCommands(
		rootCmd,
		dydxapp.DefaultNodeHome,
		a.newApp,
		a.appExport,
		func(cmd *cobra.Command) {
			addModuleInitFlags(cmd)

			if option.startCmdCustomizer != nil {
				option.startCmdCustomizer(cmd)
			}
		},
	)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(dydxapp.DefaultNodeHome),
	)

	// add rosetta commands
	rootCmd.AddCommand(rosettaCmd.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Codec))
}

// addModuleInitFlags adds module specific init flags.
func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// queryCommand adds transaction and account querying commands.
func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	basic_manager.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// txCommand adds transaction signing, encoding / decoding, and broadcasting commands.
func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	basic_manager.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

// newApp initializes and returns a new app.
func (a appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}

		chainID = appGenesis.ChainID
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", server.GetAppDBBackend(appOpts), snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent)),
	)

	// Report app version and git commit in non-dev and non-staging environments.
	if !strings.Contains(chainID, "dev") && !strings.Contains(chainID, "staging") {
		version := version.NewInfo()
		telemetry.IncrCounterWithLabels(
			[]string{metrics.AppVersionAndGitCommit},
			1,
			[]gometrics.Label{
				metrics.GetLabelForStringValue(metrics.AppVersion, version.Version),
				metrics.GetLabelForStringValue(metrics.GitCommit, version.GitCommit),
			},
		)
	}

	return dydxapp.New(
		logger,
		db,
		traceStore,
		true,
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(int(cast.ToUint64(appOpts.Get(flagIAVLCacheSize)))),
		baseapp.SetChainID(chainID),
	)
}

// appExport creates and exports a new app, returns the state of the app for a genesis file.
//
// Deprecated: this feature relies on the use of known unstable, legacy export functionality
// from cosmos.
func (a appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	dydxApp := dydxapp.New(
		logger,
		db,
		traceStore,
		height == -1, // -1: no height provided
		appOpts,
	)

	if height != -1 {
		if err := dydxApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return dydxApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
