package configs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	tmos "github.com/cometbft/cometbft/libs/os"
	daemonconstants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/pelletier/go-toml"
)

// Note: any changes to the comments/variables/mapstructure must be reflected in the appropriate
// struct in daemons/pricefeed/client/static_exchange_startup_config.go.
const (
	defaultTomlTemplate = `# This is a TOML config file.
	# StaticExchangeStartupConfig represents the mapping of exchanges to the parameters for
	# querying from them.
	#
	# ExchangeId - Unique string identifying an exchange.
	#
	# IntervalMs - Delays between sending API requests to get exchange market prices - cannot be 0.
	#
	# TimeoutMs - Max time to wait on an API call to an exchange - cannot be 0.
	#
	# MaxQueries - Max api calls to get market prices for an exchange to make in a task-loop -
	# cannot be 0. For multi-market API exchanges, the behavior will default to 1.{{ range $exchangeId, $element := .}}
	[[exchanges]]
	ExchangeId = "{{$element.ExchangeId}}"
	IntervalMs = {{$element.IntervalMs}}
	TimeoutMs = {{$element.TimeoutMs}}
	MaxQueries = {{$element.MaxQueries}}{{end}}
`
)

// GenerateDefaultExchangeTomlString creates the toml file string containing the default configs
// for querying each exchange.
func GenerateDefaultExchangeTomlString() bytes.Buffer {
	// Create the template for turning each `parsableExchangeStartupConfig` into a toml map config in
	// a stringified toml file.
	template, err := template.New("").Parse(defaultTomlTemplate)
	// Panic if failure occurs when parsing the template.
	if err != nil {
		panic(err)
	}

	// Encode toml string into `defaultExchangeToml` and return if successful. Otherwise, panic.
	var defaultExchangeToml bytes.Buffer
	err = template.Execute(&defaultExchangeToml, constants.StaticExchangeQueryConfig)
	if err != nil {
		panic(err)
	}
	return defaultExchangeToml
}

// WriteDefaultPricefeedExchangeToml reads in the toml string for the pricefeed client and
// writes said string to the config folder as a toml file if the config file does not exist.
func WriteDefaultPricefeedExchangeToml(homeDir string) {
	// Write file into config folder if file does not exist.
	configFilePath := getConfigFilePath(homeDir)
	if !tmos.FileExists(configFilePath) {
		buffer := GenerateDefaultExchangeTomlString()
		tmos.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
	}
}

// ReadExchangeQueryConfigFile gets a mapping of `exchangeIds` to `ExchangeQueryConfigs`
// where `ExchangeQueryConfig` for querying exchanges for market prices comes from parsing a TOML
// file in the config directory.
// NOTE: if the config file is not found for the price-daemon, return the static exchange query
// config.
func ReadExchangeQueryConfigFile(homeDir string) map[types.ExchangeId]*types.ExchangeQueryConfig {
	// Read file for exchange query configurations.
	tomlFile, err := os.ReadFile(getConfigFilePath(homeDir))
	if err != nil {
		panic(err)
	}

	// Unmarshal `tomlFile` into `exchanges` for `exchangeStartupConfigMap`.
	exchanges := map[string][]types.ExchangeQueryConfig{}
	if err = toml.Unmarshal(tomlFile, &exchanges); err != nil {
		panic(err)
	}

	// Populate configs for exchanges.
	exchangeStartupConfigMap := make(map[types.ExchangeId]*types.ExchangeQueryConfig, len(exchanges))
	for _, exchange := range exchanges["exchanges"] {
		// Zero is an invalid configuration value for all parameters. This could also point to the
		// configuration file being setup wrong with one or more exchange parameters unset.
		if exchange.IntervalMs == 0 ||
			exchange.TimeoutMs == 0 ||
			exchange.MaxQueries == 0 {
			panic(
				fmt.Errorf(
					"One or more query config values are unset or are set to zero for exchange with id: '%v'",
					exchange.ExchangeId,
				),
			)
		}

		// Insert Key-Value pair into `exchangeStartupConfigMap`.
		exchangeStartupConfigMap[exchange.ExchangeId] = &types.ExchangeQueryConfig{
			ExchangeId: exchange.ExchangeId,
			IntervalMs: exchange.IntervalMs,
			TimeoutMs:  exchange.TimeoutMs,
			MaxQueries: exchange.MaxQueries,
		}
	}

	return exchangeStartupConfigMap
}

// getConfigFilePath returns the path to the pricefeed exchange config file.
func getConfigFilePath(homeDir string) string {
	return filepath.Join(
		homeDir,
		"config",
		daemonconstants.PricefeedExchangeConfigFileName,
	)
}
