package test_volatile_exchange

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

// VolatileExchangeTicker is our representation of ticker information for a test market.
// It implements interface `Ticker` in util.go.
type VolatileExchangeTicker struct {
	Pair  string
	Price string
}

func (t VolatileExchangeTicker) GetPair() string {
	return t.Pair
}

func (t VolatileExchangeTicker) GetAskPrice() string {
	return t.Price
}

func (t VolatileExchangeTicker) GetBidPrice() string {
	return t.Price
}

func (t VolatileExchangeTicker) GetLastPrice() string {
	return t.Price
}

// VolatileExchangePriceFunction generates a time-based price value based off of the following
// function:
// - PRICE = AVERAGE * (1 + AMPLITUDE * WAVE_VALUE)
// - WAVE_VALUE = math.Sin(RADIANS)
// - RADIANS = PERCENTAGE_THROUGH_DAY * FREQUENCY * 2 * math.Pi
// - PERCENTAGE_THROUGH_DAY = (time.Now().Unix() % SECONDS_IN_A_DAY) / SECONDS_IN_A_DAY
// The following values are parametrized in `VolatileExchangeParams`:
// - AVERAGE, AMPLITUDE, FREQUENCY
func VolatileExchangePriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	percentageThroughDay := float64(time.Now().Unix()%SECONDS_IN_DAY) / float64(SECONDS_IN_DAY)
	radians := percentageThroughDay * TestVolatileExchangeParams.Frequency * 2 * math.Pi
	waveValue := math.Sin(radians)
	price := float64(TestVolatileExchangeParams.AveragePrice) * (1 + TestVolatileExchangeParams.Amplitude*waveValue)

	volatile_exchange_ticker := VolatileExchangeTicker{
		Pair:  "TEST-USD",
		Price: fmt.Sprintf("%f", price),
	}
	return price_function.GetMedianPricesFromTickers(
		[]VolatileExchangeTicker{volatile_exchange_ticker},
		tickerToExponent,
		resolver,
	)
}
