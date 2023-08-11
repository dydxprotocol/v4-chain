package types

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExchangeToPrices_IsEmpty(t *testing.T) {
	etp := NewExchangeToPrice()

	require.Empty(t, etp.exchangeToPriceTimestamp)
}

func TestUpdatePrices_SingleExchangeSingleUpdate(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
		})

	require.Len(t, etp.exchangeToPriceTimestamp, 1)
	priceTimestamp := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId1]
	require.Equal(t, priceTimestamp.Price, constants.Price1)
	require.Equal(t, priceTimestamp.LastUpdateTime, constants.TimeT)
}

func TestUpdatePrices_SingleExchangeMultiUpdate(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange1_Price2_AfterTimeT,
		})

	// Update with greater timestamp overwrites
	require.Len(t, etp.exchangeToPriceTimestamp, 1)
	priceTimestamp := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId1]
	require.Equal(t, priceTimestamp.Price, constants.Price2)
	require.Equal(t, priceTimestamp.LastUpdateTime, constants.TimeTPlusThreshold)
}

func TestUpdatePrices_MultiExchangeSingleUpdate(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange2_Price2_TimeT,
		})

	require.Len(t, etp.exchangeToPriceTimestamp, 2)
	priceTimestamp1 := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId1]
	priceTimestamp2 := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId2]
	require.Equal(t, priceTimestamp1.Price, constants.Price1)
	require.Equal(t, priceTimestamp1.LastUpdateTime, constants.TimeT)
	require.Equal(t, priceTimestamp2.Price, constants.Price2)
	require.Equal(t, priceTimestamp2.LastUpdateTime, constants.TimeT)
}

func TestUpdatePrices_MultiExchangeMutliUpdate(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange1_Price2_AfterTimeT,
			constants.Exchange2_Price2_TimeT,
			constants.Exchange2_Price3_AfterTimeT,
		})

	// Update with greater timestamp overwrites
	require.Len(t, etp.exchangeToPriceTimestamp, 2)
	priceTimestamp1 := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId1]
	priceTimestamp2 := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId2]
	require.Equal(t, priceTimestamp1.Price, constants.Price2)
	require.Equal(t, priceTimestamp1.LastUpdateTime, constants.TimeTPlusThreshold)
	require.Equal(t, priceTimestamp2.Price, constants.Price3)
	require.Equal(t, priceTimestamp2.LastUpdateTime, constants.TimeTPlusThreshold)
}

func TestUpdatePrices_OldUpdateFails(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
		})
	priceTimestamp1 := etp.exchangeToPriceTimestamp[constants.ExchangeFeedId1]

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price3_BeforeTimeT,
		})

	// Older timestamp does NOT update the prices.
	require.Len(t, etp.exchangeToPriceTimestamp, 1)
	require.Equal(t, priceTimestamp1.Price, constants.Price1)
	require.Equal(t, priceTimestamp1.LastUpdateTime, constants.TimeT)
}

func TestGetValidPrices(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
		})

	r := etp.GetValidPrices(constants.ValidExchanges1, constants.TimeT)
	require.Len(t, r, 1)
	require.Equal(t, constants.Price1, r[0])
}

func TestGetValidPrices_Empty(t *testing.T) {
	etp := NewExchangeToPrice()

	r := etp.GetValidPrices(constants.ValidExchangesAll, constants.TimeT)
	require.Empty(t, r)
}

func TestGetValidPrices_OldPricesEmpty(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange2_Price2_TimeT,
		})

	r := etp.GetValidPrices(constants.ValidExchangesAll, constants.TimeTPlus1)
	require.Empty(t, r)
}

func TestGetValidPrices_ValidAndOldPrices(t *testing.T) {
	etp := NewExchangeToPrice()

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange2_Price3_AfterTimeT,
			constants.Exchange3_Price4_AfterTimeT,
		})

	// Exchange 1's Price is before cutoff, so it's ignored
	r := etp.GetValidPrices(constants.ValidExchangesAll, constants.TimeTPlus1)
	require.Len(t, r, 2)

	expected := []uint64{constants.Price3, constants.Price4}
	assert.ElementsMatch(t, expected, r)
}

func TestGetValidPrices_MixedValidAndInvalid(t *testing.T) {
	tests := map[string]struct {
		input             []*api.ExchangePrice
		expectedResultLen int
		expectedPrices    []uint64
	}{
		"Valid Exchange + Invalid Time": {
			input:             []*api.ExchangePrice{constants.Exchange1_Price1_TimeT},
			expectedResultLen: 0,
			expectedPrices:    nil,
		},
		"Invalid Exchange + Invalid Time": {
			input:             []*api.ExchangePrice{constants.Exchange2_Price2_TimeT},
			expectedResultLen: 0,
			expectedPrices:    nil,
		},
		"Invalid Exchange + Valid Time": {
			input:             []*api.ExchangePrice{constants.Exchange3_Price4_AfterTimeT},
			expectedResultLen: 0,
			expectedPrices:    nil,
		},
		"Mixed: All Invalid": {
			input: []*api.ExchangePrice{
				constants.Exchange1_Price1_TimeT,      // valid exchange   + invalid time
				constants.Exchange2_Price2_TimeT,      // invalid exchange + invalid time
				constants.Exchange3_Price4_AfterTimeT, // invalid exchange + valid time
			},
			expectedResultLen: 0,
			expectedPrices:    []uint64{},
		},
		"Mixed: One valid, Rest invalid": {
			input: []*api.ExchangePrice{
				constants.Exchange1_Price2_AfterTimeT, // valid exchange   + valid time
				constants.Exchange2_Price2_TimeT,      // invalid exchange + invalid time
				constants.Exchange3_Price4_AfterTimeT, // invalid exchange + valid time
			},
			expectedResultLen: 1,
			expectedPrices:    []uint64{constants.Price2},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			etp := NewExchangeToPrice()
			etp.UpdatePrices(tc.input)
			r := etp.GetValidPrices(constants.ValidExchanges1, constants.TimeTPlus1)

			require.Len(t, r, tc.expectedResultLen)
			if tc.expectedResultLen > 0 {
				assert.ElementsMatch(t, tc.expectedPrices, r)
			}
		})
	}
}
