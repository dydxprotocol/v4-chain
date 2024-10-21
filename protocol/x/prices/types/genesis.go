package types

import (
	"fmt"
)

// DefaultGenesis returns the default Prices genesis state
func DefaultGenesis() *GenesisState {
	// TODO(CORE-430): Add all canonical markets.
	return &GenesisState{
		MarketParams: []MarketParam{
			{
				ExchangeConfigJson: "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"BTCUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"BTCUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tBTCUSD\"},{\"exchangeName\":\"Bitstamp\",\"ticker\":\"BTC/USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"CryptoCom\",\"ticker\":\"BTC_USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\"}]}", //nolint:lll
				Exponent:           -5,
				Id:                 0,
				MinExchanges:       1,
				MinPriceChangePpm:  1000,
				Pair:               "BTC-USD",
			},
			{
				ExchangeConfigJson: "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ETHUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ETHUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tETHUSD\"},{\"exchangeName\":\"Bitstamp\",\"ticker\":\"ETH/USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"CryptoCom\",\"ticker\":\"ETH_USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\"}]}", //nolint:lll
				Exponent:           -6,
				Id:                 1,
				MinExchanges:       1,
				MinPriceChangePpm:  1000,
				Pair:               "ETH-USD",
			},
		},
		MarketPrices: []MarketPrice{
			{
				Exponent: -5,
				Id:       0,
				Price:    2000000000,
			},
			{
				Exponent: -6,
				Id:       1,
				Price:    1500000000,
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated key for Markets.
	marketParamKeyMap := make(map[uint32]struct{})
	for _, marketParam := range gs.MarketParams {
		if _, exists := marketParamKeyMap[marketParam.Id]; exists {
			return fmt.Errorf("duplicated market param id")
		}
		marketParamKeyMap[marketParam.Id] = struct{}{}

		if err := marketParam.Validate(); err != nil {
			return err
		}
	}

	if len(gs.MarketParams) != len(gs.MarketPrices) {
		return fmt.Errorf("expected the same number of market prices and market params")
	}

	for i, marketPrice := range gs.MarketPrices {
		if err := marketPrice.ValidateFromParam(gs.MarketParams[i]); err != nil {
			return err
		}
	}

	return nil
}
