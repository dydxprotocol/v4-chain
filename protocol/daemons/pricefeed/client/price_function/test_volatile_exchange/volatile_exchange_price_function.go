package test_volatile_exchange

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	medianizer lib.Medianizer,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	PERCENTAGE_THROUGH_DAY := float64(time.Now().Unix()%SECONDS_IN_DAY) / float64(SECONDS_IN_DAY)
	RADIANS := PERCENTAGE_THROUGH_DAY * TestVolatileExchangeParams.Frequency * 2 * math.Pi
	WAVE_VALUE := math.Sin(RADIANS)
	PRICE := float64(TestVolatileExchangeParams.AveragePrice) * (1 + TestVolatileExchangeParams.Amplitude + WAVE_VALUE)

	volatile_exchange_ticker := VolatileExchangeTicker{
		Pair:  "TEST-USD",
		Price: fmt.Sprintf("%f", PRICE),
	}
	return price_function.GetMedianPricesFromTickers(
		[]VolatileExchangeTicker{volatile_exchange_ticker},
		tickerToExponent,
		medianizer,
	)
}
