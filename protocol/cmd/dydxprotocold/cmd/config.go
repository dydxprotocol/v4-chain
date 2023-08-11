package cmd

import (
	tmcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

const (
	// TODO(CORE-189): Support both USDC and DYDX for gas.
	// TODO(CORE-194): Update denom for IBC USDC.
	minGasPrice = "0ibc/usdc-placeholder"
)

// DydxAppConfig specifies dYdX app specific config.
type DydxAppConfig struct {
	serverconfig.Config
}

// TODO(DEC-1718): Audit tendermint and app config parameters for mainnet.

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
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
	srvCfg.MinGasPrices = minGasPrice

	appConfig := DydxAppConfig{
		Config: *srvCfg,
	}

	// Enable telemetry.
	appConfig.Telemetry.Enabled = true
	appConfig.Telemetry.PrometheusRetentionTime = 60

	// Enable API server (required for telemetry).
	appConfig.API.Enable = true
	appConfig.API.Address = "tcp://0.0.0.0:1317"

	// GRPC.
	appConfig.GRPC.Address = "0.0.0.0:9090"
	appConfig.GRPCWeb.Address = "0.0.0.0:9091"

	appTemplate := serverconfig.DefaultConfigTemplate

	return appTemplate, appConfig
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

	// Mempool config.
	cfg.Mempool.Version = "v1"
	cfg.Mempool.CacheSize = 20000
	cfg.Mempool.Size = 50000
	cfg.Mempool.TTLNumBlocks = 20 //nolint:staticcheck
	cfg.Mempool.KeepInvalidTxsInCache = true

	// Enable pex.
	cfg.P2P.PexReactor = true

	// Enable telemetry.
	cfg.Instrumentation.Prometheus = true

	return cfg
}
