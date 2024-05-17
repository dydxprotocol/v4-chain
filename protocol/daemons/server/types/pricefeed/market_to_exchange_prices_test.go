package types

import (
	"testing"
	"time"

	"cosmossdk.io/log"

	pricefeed_types "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestNewMarketToExchangePrices_IsEmpty(t *testing.T) {
	mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

	require.Empty(t, mte.marketToExchangePrices)
}

func TestUpdatePrices_SingleUpdateSinglePrice(t *testing.T) {
	mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

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
	mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

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
	mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

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
	mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

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
		updatePriceInput           []*api.MarketPriceUpdate
		getPricesInputMarketParams []types.MarketParam
		getPricesInputTime         time.Time
	}{
		"No market specified": {
			updatePriceInput:           constants.AtTimeTPriceUpdate,
			getPricesInputMarketParams: []types.MarketParam{}, // No market specified.
			getPricesInputTime:         constants.TimeT,
		},
		"No valid price timestamps": {
			updatePriceInput:           constants.AtTimeTPriceUpdate,
			getPricesInputMarketParams: constants.AllMarketParamsMinExchanges2,
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
			getPricesInputMarketParams: []types.MarketParam{
				{
					Id:           constants.MarketId9,
					MinExchanges: 0, // Set to 0 to trigger median calc error
				},
			},
			getPricesInputTime: constants.TimeT,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)
			mte.UpdatePrices(tc.updatePriceInput)
			r := mte.GetValidMedianPrices(
				log.NewNopLogger(),
				tc.getPricesInputMarketParams,
				tc.getPricesInputTime,
			)

			require.Len(t, r, 0) // The result is empty.
		})
	}
}

func TestGetValidMedianPrices_MultiMarketSuccess(t *testing.T) {
	mte := NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

	mte.UpdatePrices(constants.MixedTimePriceUpdate)

	r := mte.GetValidMedianPrices(
		log.NewNopLogger(),
		constants.AllMarketParamsMinExchanges2,
		constants.TimeT,
	)

	require.Len(t, r, 3)
	require.Equal(t, uint64(2002), r[constants.MarketId9]) // Median of 1001, 2002, 3003
	require.Equal(t, uint64(2503), r[constants.MarketId8]) // Median of 2002, 3003
}
