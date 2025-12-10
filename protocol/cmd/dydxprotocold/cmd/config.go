package cmd

import (
	"time"

	tmcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	oracleconfig "github.com/dydxprotocol/slinky/oracle/config"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

const (
	// `minGasPriceUusdc` is default minimum gas price in micro USDC.
	minGasPriceUusdc = "0.025" + assettypes.UusdcDenom
	// `minGasPriceStakeToken` is the default minimum gas price in stake token.
	// TODO(GENESIS): `adv4tnt` is a placeholder for the stake token of the dYdX chain.
	// Before this software is published for genesis, `adv4tnt` should be replaced with
	// the chain stake token. It's also recommended that the min gas price in stake token
	// is roughly the same in value as 0.025 micro USDC.
	minGasPriceStakeToken = "25000000000adv4tnt"
	// `minGasPrice` defines the default `minimum-gas-prices` attribute in validator's `app.toml` file.
	MinGasPrice = minGasPriceUusdc + "," + minGasPriceStakeToken
)

// DydxAppConfig specifies dYdX app specific config.
type DydxAppConfig struct {
	serverconfig.Config
	Oracle oracleconfig.AppConfig `mapstructure:"oracle"`
}

// TODO(DEC-1718): Audit tendermint and app config parameters for mainnet.

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, *DydxAppConfig) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In simapp, we set the min gas prices to 0.
	srvCfg.MinGasPrices = MinGasPrice

	appConfig := DydxAppConfig{
		Config: *srvCfg,
		Oracle: oracleconfig.AppConfig{
			Enabled:        true,
			OracleAddress:  "localhost:8080",
			ClientTimeout:  time.Second * 2,
			MetricsEnabled: false,
		},
	}

	// Enable telemetry.
	appConfig.Telemetry.Enabled = true
	appConfig.Telemetry.PrometheusRetentionTime = 60

	// Enable API server (required for telemetry).
	appConfig.API.Enable = true
	appConfig.API.Address = "tcp://0.0.0.0:1317"

	// GRPC.
	appConfig.GRPC.Address = "0.0.0.0:9090"

	appTemplate := serverconfig.DefaultConfigTemplate + oracleconfig.DefaultConfigTemplate

	return appTemplate, &appConfig
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// TODO(DEC-1716): Set default seeds.
	cfg.P2P.Seeds = ""

	// Expose the Tendermint RPC.
	cfg.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	cfg.RPC.CORSAllowedOrigins = []string{"*"}
	// goroutine profiling showed that we were using exactly 900 threads (the default) which was throttling
	// the maximum amount of load that the process could take. As of the last load test, at max QPS we were
	// seeing ~1700 threads being used.
	cfg.RPC.MaxOpenConnections = 8000
	cfg.RPC.GRPCMaxOpenConnections = 8000

	// Mempool config.
	// We specifically are using a number greater than max QPS (currently set at 5000) * ShortBlockWindow to prevent
	// a replay attack that is possible with short-term order placements and cancellations. The attack would consume
	// a users rate limit if the entry is evicted from the mempool cache as it would be possible for the transaction
	// to go through `CheckTx` again causing it to hit rate limit code against the users account.
	cfg.Mempool.CacheSize = 5000 * int(clobtypes.ShortBlockWindow)
	cfg.Mempool.Size = 100000
	cfg.Mempool.TTLNumBlocks = 40 //nolint:staticcheck
	cfg.Mempool.KeepInvalidTxsInCache = true

	// Enable pex.
	cfg.P2P.PexReactor = true

	// Enable telemetry.
	cfg.Instrumentation.Prometheus = true

	// Set default commit timeout to 500ms for faster block time.
	// Note: avoid using 1s since it's considered tne default Tendermint value
	// (https://github.com/dydxprotocol/tendermint/blob/dc03b21cf5d54c641e1d14b14fae5920fa7ba656/config/config.go#L982)
	// and will be overridden by `interceptConfigs` in `cosmos-sdk`.
	cfg.Consensus.TimeoutCommit = 500 * time.Millisecond

	return cfg
}
