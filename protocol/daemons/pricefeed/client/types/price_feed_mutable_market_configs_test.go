package types_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	prices_types "github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
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

// newTestPriceFeedMutableMarketConfigs returns a new PricefeedMutableMarketConfigs with exchange and
// market configurations for testing.
func newTestPriceFeedMutableMarketConfigs() (*types.PricefeedMutableMarketConfigsImpl, error) {
	pfmmc := types.NewPriceFeedMutableMarketConfigs(
		[]types.ExchangeId{exchangeIdCoinbase, exchangeIdBinance},
	)
	coinbaseFetcher := &mocks.ExchangeConfigUpdater{}
	coinbaseFetcher.On("GetExchangeId").Return(exchangeIdCoinbase)
	coinbaseFetcher.On("UpdateMutableExchangeConfig", mock.Anything, mock.Anything).Return(nil)

	binanceFetcher := &mocks.ExchangeConfigUpdater{}
	binanceFetcher.On("GetExchangeId").Return(exchangeIdBinance)
	binanceFetcher.On("UpdateMutableExchangeConfig", mock.Anything, mock.Anything).Return(nil)

	pfmmc.AddExchangeConfigUpdater(coinbaseFetcher)
	pfmmc.AddExchangeConfigUpdater(binanceFetcher)

	err := pfmmc.UpdateMarkets(constants.TestMarket7And8Params)
	return pfmmc, err
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
			pfmmc, err := newTestPriceFeedMutableMarketConfigs()
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
			pfmmc, err := newTestPriceFeedMutableMarketConfigs()
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
				Id:       1,
				Exponent: -2,
			}},
			expectedError: errors.New("invalid market param 0: pair cannot be empty"),
		},
		"Invalid: invalid exchangeConfigJson (empty)": {
			marketParams: []prices_types.MarketParam{{
				Id:       1,
				Exponent: -2,
				Pair:     "BTC-USD",
			}},
			expectedError: errors.New("invalid exchange config json for market param 0: unexpected end of JSON input"),
		},
		"Invalid: invalid exchangeConfigJson (not json)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: "invalid",
			}},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid character 'i' looking for beginning " +
					"of value"),
		},
		"Invalid: invalid exchangeConfigJson (does not conform to schema)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: `{"exchanges":"invalid"}`,
			}},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: json: cannot unmarshal string into Go struct " +
					"field ExchangeConfigJson.exchanges of type []types.ExchangeMarketConfigJson",
			),
		},
		"Invalid: invalid exchangeConfigJson (does not validate - empty exchanges)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: "{}",
			}},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: exchanges cannot be empty",
			),
		},
		"Invalid: invalid exchangeConfigJson (exchange name cannot be empty)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: `{"exchanges":[{}]}`,
			}},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: exchange name cannot be empty",
			),
		},
		"Invalid: invalid exchangeConfigJson (exchange name invalid)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigInvalidExchangeName),
			}},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: exchange name 'invalid' is " +
					"not valid",
			),
		},
		"Invalid: invalid exchangeConfigJson (ticker empty)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigEmptyTicker),
			}},
			expectedError: errors.New(
				"invalid exchange config json for market param 0: invalid exchange: ticker cannot be empty",
			),
		},
		"Invalid: invalid exchangeConfigJson (adjustment market invalid)": {
			marketParams: []prices_types.MarketParam{{
				Id:                 1,
				Exponent:           -2,
				Pair:               "BTC-USD",
				ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigInvalidAdjustByMarket),
			}},
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
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigCoinbaseBtcAdjustByEth),
				},
				{
					Id:                 1,
					Exponent:           -3,
					Pair:               "ETH-USD",
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, exchangeConfigCoinbaseEth),
				},
			},
			expectedError: errors.New("invalid market param 1: duplicate market id 1"),
		},
		"Valid: 2 markets, 2 exchanges": {
			marketParams: []prices_types.MarketParam{
				{
					Id:       1,
					Exponent: -2,
					Pair:     "BTC-USD",
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						exchangeConfigCoinbaseBtcAdjustByEth,
						exchangeConfigBinanceBtc,
					),
				},
				{
					Id:       2,
					Exponent: -3,
					Pair:     "ETH-USD",
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						exchangeConfigCoinbaseEth,
						exchangeConfigBinanceEth,
					),
				},
			},
			expectedMutableMarketConfigs: map[types.MarketId]*types.MutableMarketConfig{
				1: {
					Id:       1,
					Exponent: -2,
					Pair:     "BTC-USD",
				},
				2: {
					Id:       2,
					Exponent: -3,
					Pair:     "ETH-USD",
				},
			},
			expectedMutableExchangeConfigs: map[types.ExchangeId]*types.MutableExchangeMarketConfig{
				exchangeIdCoinbase: {
					Id: exchangeIdCoinbase,
					MarketToTicker: map[types.MarketId]string{
						1: "BTC-USD",
						2: "ETH-USD",
					},
				},
				exchangeIdBinance: {
					Id: exchangeIdBinance,
					MarketToTicker: map[types.MarketId]string{
						1: "BTCUSDT",
						2: "ETHUSDT",
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pfmmc, err := newTestPriceFeedMutableMarketConfigs()
			require.NoError(t, err)
			mutableExchangeConfigs, mutableMarketConfigs, err := pfmmc.ValidateAndTransformParams(tc.marketParams)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMutableMarketConfigs, mutableMarketConfigs)
				require.Equal(t, tc.expectedMutableExchangeConfigs, mutableExchangeConfigs)
			}
		})
	}
}
