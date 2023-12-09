package test_fixed_price_exchange

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

// Fixed prices for BTC-USD, ETH-USD, SOL-USD
const (
	BTC_USD_PRICE = 50000
	ETH_USD_PRICE = 4000
	SOL_USD_PRICE = 100
)

type FixedPriceExchangeParams struct {
	BTCUSDPrice float64
	ETHUSDPrice float64
	SOLUSDPrice float64
}

var (
	TestFixedPriceExchangeParams = FixedPriceExchangeParams{
		BTCUSDPrice: BTC_USD_PRICE,
		ETHUSDPrice: ETH_USD_PRICE,
		SOLUSDPrice: SOL_USD_PRICE,
	}
	TestFixedPriceExchangeDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE,
		Url:           "https://jsonplaceholder.typicode.com/users",
		PriceFunction: FixedExchangePriceFunction,
		IsMultiMarket: false,
	}
)
