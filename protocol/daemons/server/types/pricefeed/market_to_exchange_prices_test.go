package types

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestNewMarketToExchangePrices_IsEmpty(t *testing.T) {
	mte := NewMarketToExchangePrices()

	require.Empty(t, mte.marketToExchangePrices)
}

func TestUpdatePrices_SingleUpdateSinglePrice(t *testing.T) {
	mte := NewMarketToExchangePrices()

	mte.UpdatePrices(
		[]*api.MarketPriceUpdate{
			{
				MarketId: constants.MarketId9,
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
				},
			},
		})

	require.Len(t, mte.marketToExchangePrices, 1)
	_, ok := mte.marketToExchangePrices[constants.MarketId9]
	require.True(t, ok)
}

func TestUpdatePrices_SingleUpdateMultiPrices(t *testing.T) {
	mte := NewMarketToExchangePrices()

	mte.UpdatePrices(
		[]*api.MarketPriceUpdate{
			{
				MarketId: constants.MarketId9,
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange2_Price2_TimeT,
				},
			},
		})

	require.Len(t, mte.marketToExchangePrices, 1)
	_, ok := mte.marketToExchangePrices[constants.MarketId9]
	require.True(t, ok)
}

func TestUpdatePrices_MultiUpdatesMultiPrices(t *testing.T) {
	mte := NewMarketToExchangePrices()

	mte.UpdatePrices(
		[]*api.MarketPriceUpdate{
			{
				MarketId: constants.MarketId9,
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange2_Price2_TimeT,
				},
			},
			{
				MarketId: constants.MarketId8,
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange2_Price2_TimeT,
				},
			},
		})

	require.Len(t, mte.marketToExchangePrices, 2)
	_, ok9 := mte.marketToExchangePrices[constants.MarketId9]
	require.True(t, ok9)
	_, ok8 := mte.marketToExchangePrices[constants.MarketId8]
	require.True(t, ok8)
}

func TestUpdatePrices_MultiUpdatesMultiPricesRepeated(t *testing.T) {
	mte := NewMarketToExchangePrices()

	mte.UpdatePrices(
		[]*api.MarketPriceUpdate{
			{
				MarketId: constants.MarketId9,
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange2_Price2_TimeT,
				},
			},
			{
				MarketId: constants.MarketId9, // Repeated market
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange3_Price4_AfterTimeT,
				},
			},
			{
				MarketId: constants.MarketId8,
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange2_Price2_TimeT,
				},
			},
			{
				MarketId: constants.MarketId8, // Repeated market
				ExchangePrices: []*api.ExchangePrice{
					constants.Exchange1_Price1_TimeT,
					constants.Exchange3_Price4_AfterTimeT,
				},
			},
		})

	require.Len(t, mte.marketToExchangePrices, 2)
	_, ok9 := mte.marketToExchangePrices[constants.MarketId9]
	require.True(t, ok9)
	_, ok8 := mte.marketToExchangePrices[constants.MarketId8]
	require.True(t, ok8)
}

func TestGetValidMedianPrices_EmptyResult(t *testing.T) {
	tests := map[string]struct {
		updatePriceInput      []*api.MarketPriceUpdate
		getPricesInputMarkets []types.Market
		getPricesInputTime    time.Time
	}{
		"No market specified": {
			updatePriceInput:      constants.AtTimeTPriceUpdate,
			getPricesInputMarkets: []types.Market{}, // No market specified.
			getPricesInputTime:    constants.TimeT,
		},
		"No exchange specified": {
			updatePriceInput: constants.AtTimeTPriceUpdate,
			getPricesInputMarkets: []types.Market{
				{
					Id:        constants.MarketId9,
					Exchanges: []uint32{}, // No exchanges specified
				},
				{
					Id:        constants.MarketId8,
					Exchanges: []uint32{}, // No exchanges specified
				},
				{
					Id:        constants.MarketId7,
					Exchanges: []uint32{}, // No exchanges specified
				},
			},
			getPricesInputTime: constants.TimeT,
		},
		"No valid exchange specified": {
			updatePriceInput: constants.AtTimeTPriceUpdate,
			getPricesInputMarkets: []types.Market{
				{
					Id:        constants.MarketId9,
					Exchanges: []uint32{11, 12, 13}, // Exchanges are invalid
				},
				{
					Id:        constants.MarketId8,
					Exchanges: []uint32{14, 15, 16}, // Exchanges are invalid
				},
				{
					Id:        constants.MarketId7,
					Exchanges: []uint32{17, 18, 19}, // Exchanges are invalid
				},
			},
			getPricesInputTime: constants.TimeT,
		},
		"No valid price timestamps": {
			updatePriceInput:      constants.AtTimeTPriceUpdate,
			getPricesInputMarkets: constants.AllMarketsMinExchanges2,
			// Updates @ timeT are invalid at this read time
			getPricesInputTime: constants.TimeTPlusThreshold.Add(time.Duration(1)),
		},
		"Empty prices does not throw": {
			updatePriceInput: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId9,
					ExchangePrices: []*api.ExchangePrice{
						constants.Exchange1_Price3_BeforeTimeT, // Invalid time
					},
				},
			},
			getPricesInputMarkets: []types.Market{
				{
					Id:           constants.MarketId9,
					Exchanges:    []uint32{constants.ExchangeFeedId1},
					MinExchanges: 0, // Set to 0 to trigger median calc error
				},
			},
			getPricesInputTime: constants.TimeT,
		},
		"Does not meet min exchanges": {
			updatePriceInput: constants.AtTimeTPriceUpdate,
			// MinExchanges is 3 for all markets, but updates are from 2 exchanges
			getPricesInputMarkets: constants.AllMarketsMinExchanges3,
			getPricesInputTime:    constants.TimeT,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mte := NewMarketToExchangePrices()
			mte.UpdatePrices(tc.updatePriceInput)
			r := mte.GetValidMedianPrices(
				tc.getPricesInputMarkets,
				tc.getPricesInputTime,
			)

			require.Len(t, r, 0) // The result is empty.
		})
	}
}

func TestGetValidMedianPrices_MultiMarketSuccess(t *testing.T) {
	mte := NewMarketToExchangePrices()

	mte.UpdatePrices(constants.MixedTimePriceUpdate)

	r := mte.GetValidMedianPrices(constants.AllMarketsMinExchanges2, constants.TimeT)

	require.Len(t, r, 2)
	require.Equal(t, uint64(2002), r[constants.MarketId9]) // Median of 1001, 2002, 3003
	require.Equal(t, uint64(2503), r[constants.MarketId8]) // Median of 2002, 3003
	// Market7 only has 1 valid price due to update time constraint,
	// but the min exchanges required is 2. Therefore, no median price.
}
