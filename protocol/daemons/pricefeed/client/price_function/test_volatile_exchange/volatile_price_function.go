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

// Ensure that VolatileExchangeTicker implements the Ticker interface at compile time.
var _ price_function.Ticker = (*VolatileExchangeTicker)(nil)

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

// VolatileExchangePriceFunction generates a time-based price value. The value follows a cosine wave
// function, but that includes jumps from the lowest value to the highest value (and vice versa)
// once per period. The general formula is written below.
// - PRICE = AVERAGE * (1 + AMPLITUDE * WAVE_VALUE)
// - WAVE_VALUE = math.Cos(RADIANS)
// - RADIANS = (PHASE <= 0.5 ? PHASE * 4 : PHASE * 4 - 1) * math.Pi
// - PHASE = (FREQUENCY * UNIX_SECONDS / SECONDS_IN_DAY) % 1
// The following values are parametrized in `VolatileExchangeParams`:
// - AVERAGE, AMPLITUDE, FREQUENCY
func VolatileExchangePriceFunction(
	response *http.Response,
	tickerToExponent map[string]int32,
	resolver types.Resolver,
) (tickerToPrice map[string]uint64, unavailableTickers map[string]error, err error) {
	// Calculate the phase, how far in the period we are.
	// The phase is a value that goes from 0 to 1 over the timespan of (1 day / frequency).
	phase := math.Mod(
		TestVolatileExchangeParams.Frequency*
			float64(time.Now().Unix())/
			float64(SECONDS_IN_DAY),
		1,
	)

	// Next we get the radians. Over each period, we want the price to "jump" from the max to the min
	// and vice versa. Otherwise the price should move smoothly between min and max.
	// So we want the  first half of the period to move from 0-2 pi radians,
	// and the second half of the period to move from 1-3 pi radians.
	radMultiplier := phase * float64(4)
	if phase > 0.5 {
		radMultiplier -= float64(1)
	}
	radians := radMultiplier * float64(math.Pi)

	// Next we get the final wave value in the range from -1 to 1 based on the radians.
	waveValue := math.Cos(radians)

	// The price value is centered around `AveragePrice` with an amplitude of `Amplitude`.
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
