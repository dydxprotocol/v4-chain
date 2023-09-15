package types_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	prices_types "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	exchangeConfigInvalidExchangeName    = `{"exchangeName":"invalid"}`
	exchangeConfigEmptyTicker            = `{"exchangeName":"Coinbase"}`
	exchangeConfigInvalidAdjustByMarket  = `{"exchangeName":"Coinbase","ticker":"BTC-USD", "adjustByMarket":"invalid"}`
	exchangeConfigCoinbaseBtcAdjustByEth = `{"exchangeName":"Coinbase","ticker":"BTC-USD", "adjustByMarket":"ETH-USD"}`
	exchangeConfigBinanceBtc             = `{"exchangeName":"Binance","ticker":"BTCUSDT"}`
	exchangeConfigCoinbaseEth            = `{"exchangeName":"Coinbase","ticker":"ETH-USD"}`
	exchangeConfigBinanceEth             = `{"exchangeName":"Binance","ticker":"ETHUSDT"}`

	exchangeIdCoinbase = "Coinbase"
	exchangeIdBinance  = "Binance"
)

// newMockUpdatersForExchange returns a new mock ExchangeConfigUpdater for testing. These mocks
// are used to test the order in which the price feed mutable market configs update the updaters.
// A test that uses these mocks should fail if the fetcher is ever updated before the encoder.
func newMockUpdatersForExchange(exchangeId types.ExchangeId) (
	encoder *mocks.ExchangeConfigUpdater,
	fetcher *mocks.ExchangeConfigUpdater,
) {
	encoder = &mocks.ExchangeConfigUpdater{}
	encoder.On("GetExchangeId").Return(exchangeId)
	encoderUpdate := encoder.On("UpdateMutableExchangeConfig", mock.Anything, mock.Anything).Return(nil)

	fetcher = &mocks.ExchangeConfigUpdater{}
	fetcher.On("GetExchangeId").Return(exchangeId)
	fetcher.On("UpdateMutableExchangeConfig", mock.Anything, mock.Anything).Return(nil).NotBefore(encoderUpdate)

	return encoder, fetcher
}

// newTestPriceFeedMutableMarketConfigs returns a new PricefeedMutableMarketConfigs with exchange and
// market configurations for testing.
func newTestPriceFeedMutableMarketConfigs() (
	pfmmc *types.PricefeedMutableMarketConfigsImpl,
	encoder *mocks.ExchangeConfigUpdater,
	fetcher *mocks.ExchangeConfigUpdater,
	err error,
) {
	pfmmc = types.NewPriceFeedMutableMarketConfigs(
		[]types.ExchangeId{exchangeIdCoinbase, exchangeIdBinance},
	)
	for _, exchange := range []types.ExchangeId{exchangeIdCoinbase, exchangeIdBinance} {
		encoder, fetcher = newMockUpdatersForExchange(exchange)
		pfmmc.AddPriceEncoder(encoder)
		pfmmc.AddPriceFetcher(fetcher)
	}

	err = pfmmc.UpdateMarkets(constants.TestMarket7And8Params)

	return pfmmc, encoder, fetcher, err
}

func TestGetExchangeMarketConfigCopy_Mixed(t *testing.T) {
	tests := map[string]struct {
		Id            types.ExchangeId
		Expected      *types.MutableExchangeMarketConfig
		ExpectedError error
	}{
		"success: valid exchange id": {
			Id:       exchangeIdCoinbase,
			Expected: constants.CoinbaseMutableMarketConfig,
		},
		"failure: invalid exchange id": {
			Id:            "invalid",
			ExpectedError: fmt.Errorf("mutableExchangeMarketConfig not found for exchange 'invalid'"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pfmmc, _, _, err := newTestPriceFeedMutableMarketConfigs()
			require.NoError(t, err)
			actual, err := pfmmc.GetExchangeMarketConfigCopy(tc.Id)
			if tc.ExpectedError != nil {
				require.Nil(t, actual)
				require.Error(t, err, tc.ExpectedError.Error())
			} else {
				// Validate that this method returns a copy and not the original.
				require.NotSame(t, tc.Expected, actual)
				require.Equal(t, tc.Expected, actual)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetMarketConfigCopies(t *testing.T) {
	tests := map[string]struct {
		Ids           []types.MarketId
		Expected      []*types.MutableMarketConfig
		ExpectedError error
	}{
		"success: empty list of ids": {
			Ids:      []types.MarketId{},
			Expected: []*types.MutableMarketConfig{},
		},
		"success: 1 valid id": {
			Ids: []types.MarketId{8},
			Expected: []*types.MutableMarketConfig{
				constants.TestMutableMarketConfigs[8],
			},
		},
		"success: all valid ids": {
			Ids: []types.MarketId{7, 8},
			Expected: []*types.MutableMarketConfig{
				constants.TestMutableMarketConfigs[7],
				constants.TestMutableMarketConfigs[8],
			},
		},
		"failure: mix of valid and invalid": {
			Ids:           []types.MarketId{7, 8, 9},
			ExpectedError: fmt.Errorf("market 9 not found in mutableMarketToConfigs"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pfmmc, _, _, err := newTestPriceFeedMutableMarketConfigs()
			require.NoError(t, err)
			actual, err := pfmmc.GetMarketConfigCopies(tc.Ids)
			if tc.ExpectedError != nil {
				require.Nil(t, actual)
				require.Error(t, err, tc.ExpectedError.Error())
			} else {
				// Validate that this method returns a copy and not the original.
				require.NotSame(t, tc.Expected, actual)
				require.Equal(t, tc.Expected, actual)
				require.NoError(t, err)
			}
		})
	}
}

func validMarketParamWithExchangeConfig(exchangeConfig string) prices_types.MarketParam {
	return prices_types.MarketParam{
		Id:                 1,
		Exponent:           -2,
		Pair:               "BTC-USD",
		MinExchanges:       1,
		MinPriceChangePpm:  1,
		ExchangeConfigJson: exchangeConfig,
	}
}

func TestValidateAndTransformParams_Mixed(t *testing.T) {
	tests := map[string]struct {
		marketParams                   []prices_types.MarketParam
		expectedError                  error
		expectedMutableMarketConfigs   map[types.MarketId]*types.MutableMarketConfig
		expectedMutableExchangeConfigs map[types.ExchangeId]*types.MutableExchangeMarketConfig
	}{
		"Invalid: nil params": {
			marketParams:  nil,
			expectedError: errors.New("marketParams cannot be nil"),
		},
		"Invalid: invalid params (missing pair)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				MinExchanges:       1,
				MinPriceChangePpm:  1,
				ExchangeConfigJson: "{}",
			}},
			expectedError: errors.New("invalid market param 0: Pair cannot be empty: Invalid input"),
		},
		"Invalid: invalid exchangeConfigJson (empty, fails marketParams.Validate)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				MinExchanges:       1,
				MinPriceChangePpm:  1,
				ExchangeConfigJson: "",
			}},
			expectedError: errors.New("ExchangeConfigJson string is not valid"),
		},
		"Invalid: invalid exchangeConfigJson (not json, fails marketParams.Validate)": {
			marketParams: []prices_types.MarketParam{
				validMarketParamWithExchangeConfig("invalid"),
			},
			expectedError: errors.New("ExchangeConfigJson string is not valid"),
		},
		"Invalid: invalid exchangeConfigJson (does not conform to schema)": {
			marketParams: []prices_types.MarketParam{
				validMarketParamWithExchangeConfig(`{"exchanges":"invalid"}`),
			},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: json: cannot unmarshal string into Go struct " +
					"field ExchangeConfigJson.exchanges of type []types.ExchangeMarketConfigJson",
			),
		},
		"Invalid: invalid exchangeConfigJson (does not validate - empty exchanges)": {
			marketParams: []prices_types.MarketParam{
				validMarketParamWithExchangeConfig("{}"),
			},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: exchanges cannot be empty",
			),
		},
		"Invalid: invalid exchangeConfigJson (exchange name cannot be empty)": {
			marketParams: []prices_types.MarketParam{validMarketParamWithExchangeConfig(`{"exchanges":[{}]}`)},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: exchange name cannot be empty",
			),
		},
		"Invalid: invalid exchangeConfigJson (exchange name invalid)": {
			marketParams: []prices_types.MarketParam{
				validMarketParamWithExchangeConfig(fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigInvalidExchangeName)),
			},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: exchange name 'invalid' is " +
					"not valid",
			),
		},
		"Invalid: invalid exchangeConfigJson (ticker empty)": {
			marketParams: []prices_types.MarketParam{
				validMarketParamWithExchangeConfig(fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigEmptyTicker)),
			},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: ticker cannot be empty",
			),
		},
		"Invalid: invalid exchangeConfigJson (adjustment market invalid)": {
			marketParams: []prices_types.MarketParam{
				validMarketParamWithExchangeConfig(fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigInvalidAdjustByMarket)),
			},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: adjustment market " +
					"'invalid' is not valid"),
		},
		"Invalid: invalid params (duplicate ids)": {
			marketParams: []prices_types.MarketParam{
				{
					Id:                 1,
					Exponent:           -2,
					Pair:               "BTC-USD",
					MinExchanges:       1,
					MinPriceChangePpm:  1,
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigCoinbaseBtcAdjustByEth),
				},
				{
					Id:                 1,
					Exponent:           -3,
					Pair:               "ETH-USD",
					MinExchanges:       2,
					MinPriceChangePpm:  2,
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigCoinbaseEth),
				},
			},
			expectedError: errors.New("invalid market param 1: duplicate market id 1"),
		},
		"Valid: 2 markets, 2 exchanges, with adjust-by markets": {
			marketParams: []prices_types.MarketParam{
				{
					Id:                1,
					Pair:              "BTC-USD",
					Exponent:          -2,
					MinExchanges:      1,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						exchangeConfigCoinbaseBtcAdjustByEth,
						exchangeConfigBinanceBtc,
					),
				},
				{
					Id:                2,
					Pair:              "ETH-USD",
					Exponent:          -3,
					MinExchanges:      2,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						exchangeConfigCoinbaseEth,
						exchangeConfigBinanceEth,
					),
				},
			},
			expectedMutableMarketConfigs: map[types.MarketId]*types.MutableMarketConfig{
				1: {
					Id:           1,
					Exponent:     -2,
					Pair:         "BTC-USD",
					MinExchanges: 1,
				},
				2: {
					Id:           2,
					Exponent:     -3,
					Pair:         "ETH-USD",
					MinExchanges: 2,
				},
			},
			expectedMutableExchangeConfigs: map[types.ExchangeId]*types.MutableExchangeMarketConfig{
				exchangeIdCoinbase: {
					Id: exchangeIdCoinbase,
					MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
						1: {
							Ticker:         "BTC-USD",
							AdjustByMarket: newUint32WithValue(2),
						},
						2: {
							Ticker: "ETH-USD",
						},
					},
				},
				exchangeIdBinance: {
					Id: exchangeIdBinance,
					MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
						1: {
							Ticker: "BTCUSDT",
						},
						2: {
							Ticker: "ETHUSDT",
						},
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pfmmc, _, _, err := newTestPriceFeedMutableMarketConfigs()
			require.NoError(t, err)
			mutableExchangeConfigs, mutableMarketConfigs, err := pfmmc.ValidateAndTransformParams(tc.marketParams)
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMutableMarketConfigs, mutableMarketConfigs)
				require.Equal(t, tc.expectedMutableExchangeConfigs, mutableExchangeConfigs)
			}
		})
	}
}

// TestUpdatesEncoderAndFetcherInOrder tests that the price feed mutable market configs updates the encoder
// before the fetcher. This test is confirmed to fail if update order is switched.
func TestUpdatesEncoderAndFetcherInOrder(t *testing.T) {
	pfmmc, encoder, fetcher, err := newTestPriceFeedMutableMarketConfigs()
	require.NoError(t, err)
	require.NoError(t, pfmmc.UpdateMarkets(constants.TestMarket7And8Params))

	// Assert that an update happened. If it happened out of order, this test should fail due to the
	// mock call configurations.
	encoder.AssertCalled(t, "UpdateMutableExchangeConfig", mock.Anything, mock.Anything)
	fetcher.AssertCalled(t, "UpdateMutableExchangeConfig", mock.Anything, mock.Anything)
}
