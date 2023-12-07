package test_fixed_price_exchange

import (
	"fmt"
	"net/http"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

type FixedPriceTicker struct {
	Pair  string
	Price string
}

var _ price_function.Ticker = (*FixedPriceTicker)(nil)

func (t FixedPriceTicker) GetPair() string {
	return t.Pair
}

func (t FixedPriceTicker) GetAskPrice() string {
	return t.Price
}

func (t FixedPriceTicker) GetBidPrice() string {
	return t.Price
}

func (t FixedPriceTicker) GetLastPrice() string {
	return t.Price
}

func FixedExchangePriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	btcTicker := FixedPriceTicker{
		Pair:  "BTC-USD",
		Price: fmt.Sprintf("%f", TestFixedPriceExchangeParams.BTCUSDPrice),
	}
	ethTicker := FixedPriceTicker{
		Pair:  "ETH-USD",
		Price: fmt.Sprintf("%f", TestFixedPriceExchangeParams.ETHUSDPrice),
	}
	solTicker := FixedPriceTicker{
		Pair:  "SOL-USD",
		Price: fmt.Sprintf("%f", TestFixedPriceExchangeParams.SOLUSDPrice),
	}
	return price_function.GetMedianPricesFromTickers(
		[]FixedPriceTicker{btcTicker, ethTicker, solTicker},
		tickerToExponent,
		resolver,
	)
}
