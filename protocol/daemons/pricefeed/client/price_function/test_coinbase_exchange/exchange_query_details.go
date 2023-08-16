package test_coinbase_exchange

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/coinbase_pro"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

// Test Exchange used for testing purposes. We'll reuse the CoinbasePro price function.
var (
	TestCoinbaseExchangeHost    = "test.exchange"
	TestCoinbaseExchangePort    = "9888"
	TestCoinbaseExchangeDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_TEST_COINBASE_EXCHANGE,
		Url:           fmt.Sprintf("http://%s:%s/ticker?symbol=$", TestCoinbaseExchangeHost, TestCoinbaseExchangePort),
		PriceFunction: coinbase_pro.CoinbaseProPriceFunction,
	}
)
