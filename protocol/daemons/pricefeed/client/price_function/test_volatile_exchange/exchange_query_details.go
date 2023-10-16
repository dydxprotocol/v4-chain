package test_volatile_exchange

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

const SECONDS_IN_DAY = 24 * 60 * 60

type VolatileExchangeParams struct {
	AveragePrice float64
	Amplitude    float64
	Frequency    float64
}

// Test Exchange used for testing purposes.
var (
	TestVolatileExchangeParams = VolatileExchangeParams{
		AveragePrice: 100,
		Amplitude:    0.95,
		Frequency:    1,
	}
	TestVolatileExchangeDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE,
		Url:           "https://jsonplaceholder.typicode.com/users",
		PriceFunction: VolatileExchangePriceFunction,
		IsMultiMarket: false,
	}
)
