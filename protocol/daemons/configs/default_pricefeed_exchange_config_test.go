package configs_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/configs"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	pfconstants "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"

	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/stretchr/testify/require"
)

var (
	binanceId = exchange_common.EXCHANGE_ID_BINANCE
	filePath  = fmt.Sprintf("config/%v", constants.PricefeedExchangeConfigFileName)
)

const (
	tomlString = `# This is a TOML config file.
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
	# cannot be 0. For multi-market API exchanges, the behavior will default to 1.
	[[exchanges]]
	ExchangeId = "Binance"
	IntervalMs = 2500
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "BinanceUS"
	IntervalMs = 2500
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Bitfinex"
	IntervalMs = 2500
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Bitstamp"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Bybit"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "CoinbasePro"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 3
	[[exchanges]]
	ExchangeId = "CryptoCom"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Gate"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Huobi"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Kraken"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Kucoin"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Mexc"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "Okx"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 1
	[[exchanges]]
	ExchangeId = "TestFixedPriceExchange"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 3
	[[exchanges]]
	ExchangeId = "TestVolatileExchange"
	IntervalMs = 2000
	TimeoutMs = 3000
	MaxQueries = 3
`
)

func TestGenerateDefaultExchangeTomlString(t *testing.T) {
	defaultConfigStringBuffer := configs.GenerateDefaultExchangeTomlString()
	require.Equal(
		t,
		tomlString,
		defaultConfigStringBuffer.String(),
	)
}

func TestWriteDefaultPricefeedExchangeToml(t *testing.T) {
	err := os.Mkdir("config", 0700)
	require.NoError(t, err)
	configs.WriteDefaultPricefeedExchangeToml("")

	buffer, err := os.ReadFile(filePath)
	require.NoError(t, err)

	require.Equal(t, tomlString, string(buffer[:]))
	os.RemoveAll("config")
}

func TestWriteDefaultPricefeedExchangeToml_FileExists(t *testing.T) {
	helloWorld := "Hello World"

	err := os.Mkdir("config", 0700)
	require.NoError(t, err)

	tmos.MustWriteFile(filePath, bytes.NewBuffer([]byte(helloWorld)).Bytes(), 0644)
	configs.WriteDefaultPricefeedExchangeToml("")

	buffer, err := os.ReadFile(filePath)
	require.NoError(t, err)

	require.Equal(t, helloWorld, string(buffer[:]))
	os.RemoveAll("config")
}

func TestReadExchangeStartupConfigFile(t *testing.T) {
	pwd, _ := os.Getwd()

	tests := map[string]struct {
		// parameters
		exchangeConfigSourcePath string
		doNotWriteFile           bool

		// expectations
		expectedExchangeId         types.ExchangeId
		expectedIntervalMsExchange uint32
		expectedTimeoutMsExchange  uint32
		expectedMaxQueries         uint32
		expectedPanic              error
	}{
		"valid": {
			exchangeConfigSourcePath:   "test_data/valid_test.toml",
			expectedExchangeId:         binanceId,
			expectedIntervalMsExchange: pfconstants.StaticExchangeQueryConfig[binanceId].IntervalMs,
			expectedTimeoutMsExchange:  pfconstants.StaticExchangeQueryConfig[binanceId].TimeoutMs,
			expectedMaxQueries:         pfconstants.StaticExchangeQueryConfig[binanceId].MaxQueries,
		},
		"config file cannot be found": {
			exchangeConfigSourcePath: "test_data/notexisting_test.toml",
			doNotWriteFile:           true,
			expectedPanic: fmt.Errorf(
				"open %s%s: no such file or directory",
				pwd+"/config/",
				constants.PricefeedExchangeConfigFileName,
			),
		},
		"config file cannot be unmarshalled": {
			exchangeConfigSourcePath: "test_data/broken_test.toml",
			expectedPanic:            errors.New("(1, 12): was expecting token [[, but got unclosed table array key instead"),
		},
		"config file has malformed values": {
			exchangeConfigSourcePath: "test_data/missingvals_test.toml",
			expectedPanic: errors.New(
				"One or more query config values are unset or are set to zero for exchange with id: 'BinanceUS'",
			),
		},
		"config file has incorrect values": {
			exchangeConfigSourcePath: "test_data/wrongvaltype_test.toml",
			expectedPanic: errors.New(
				"(3, 1): Can't convert a(string) to uint32",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if !tc.doNotWriteFile {
				err := os.Mkdir("config", 0700)
				require.NoError(t, err)

				file, err := os.Open(tc.exchangeConfigSourcePath)
				require.NoError(t, err)

				config, err := os.Create(filepath.Join("config", constants.PricefeedExchangeConfigFileName))
				require.NoError(t, err)
				_, err = config.ReadFrom(file)
				require.NoError(t, err)
			}

			if tc.expectedPanic != nil {
				require.PanicsWithError(
					t,
					tc.expectedPanic.Error(),
					func() { configs.ReadExchangeQueryConfigFile(pwd) },
				)

				os.RemoveAll("config")
				return
			}

			exchangeStartupConfigMap := configs.ReadExchangeQueryConfigFile(pwd)

			require.Equal(
				t,
				&types.ExchangeQueryConfig{
					ExchangeId: tc.expectedExchangeId,
					IntervalMs: tc.expectedIntervalMsExchange,
					TimeoutMs:  tc.expectedTimeoutMsExchange,
					MaxQueries: tc.expectedMaxQueries,
				},
				exchangeStartupConfigMap[tc.expectedExchangeId],
			)

			os.RemoveAll("config")
		})
	}

	// In case tests fail and the path was never removed.
	os.RemoveAll("config")
}
