package types

import (
	"testing"

	"cosmossdk.io/log"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExchangeToPrices_IsEmpty(t *testing.T) {
	etp := NewExchangeToPrice(0)

	require.Empty(t, etp.exchangeToPriceTimestamp)
}

func TestUpdatePrices_SingleExchangeSingleUpdate(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
		})

	require.Len(t, etp.exchangeToPriceTimestamp, 1)
	priceTimestamp, ok := etp.exchangeToPriceTimestamp[constants.ExchangeId1]
	require.True(t, ok)
	require.Equal(t, priceTimestamp.Price, constants.Price1)
	require.Equal(t, priceTimestamp.LastUpdateTime, constants.TimeT)
}

func TestUpdatePrices_SingleExchangeMultiUpdate(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange1_Price2_AfterTimeT,
		})

	// Update with greater timestamp overwrites
	require.Len(t, etp.exchangeToPriceTimestamp, 1)
	priceTimestamp := etp.exchangeToPriceTimestamp[constants.ExchangeId1]
	require.Equal(t, priceTimestamp.Price, constants.Price2)
	require.Equal(t, priceTimestamp.LastUpdateTime, constants.TimeTPlusThreshold)
}

func TestUpdatePrices_MultiExchangeSingleUpdate(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange2_Price2_TimeT,
		})

	require.Len(t, etp.exchangeToPriceTimestamp, 2)
	priceTimestamp1 := etp.exchangeToPriceTimestamp[constants.ExchangeId1]
	priceTimestamp2 := etp.exchangeToPriceTimestamp[constants.ExchangeId2]
	require.Equal(t, priceTimestamp1.Price, constants.Price1)
	require.Equal(t, priceTimestamp1.LastUpdateTime, constants.TimeT)
	require.Equal(t, priceTimestamp2.Price, constants.Price2)
	require.Equal(t, priceTimestamp2.LastUpdateTime, constants.TimeT)
}

func TestUpdatePrices_MultiExchangeMutliUpdate(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange1_Price2_AfterTimeT,
			constants.Exchange2_Price2_TimeT,
			constants.Exchange2_Price3_AfterTimeT,
		})

	// Update with greater timestamp overwrites
	require.Len(t, etp.exchangeToPriceTimestamp, 2)
	priceTimestamp1 := etp.exchangeToPriceTimestamp[constants.ExchangeId1]
	priceTimestamp2 := etp.exchangeToPriceTimestamp[constants.ExchangeId2]
	require.Equal(t, priceTimestamp1.Price, constants.Price2)
	require.Equal(t, priceTimestamp1.LastUpdateTime, constants.TimeTPlusThreshold)
	require.Equal(t, priceTimestamp2.Price, constants.Price3)
	require.Equal(t, priceTimestamp2.LastUpdateTime, constants.TimeTPlusThreshold)
}

func TestUpdatePrices_OldUpdateFails(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
		})
	priceTimestamp1 := etp.exchangeToPriceTimestamp[constants.ExchangeId1]

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
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
		})

	r := etp.GetValidPrices(log.NewNopLogger(), constants.TimeT)
	require.Len(t, r, 1)
	require.Equal(t, constants.Price1, r[0])
}

func TestGetValidPrices_Empty(t *testing.T) {
	etp := NewExchangeToPrice(0)

	r := etp.GetValidPrices(log.NewNopLogger(), constants.TimeT)
	require.Empty(t, r)
}

func TestGetValidPrices_OldPricesEmpty(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange2_Price2_TimeT,
		})

	r := etp.GetValidPrices(log.NewNopLogger(), constants.TimeTPlus1)
	require.Empty(t, r)
}

func TestGetValidPrices_ValidAndOldPrices(t *testing.T) {
	etp := NewExchangeToPrice(0)

	etp.UpdatePrices(
		[]*api.ExchangePrice{
			constants.Exchange1_Price1_TimeT,
			constants.Exchange2_Price3_AfterTimeT,
			constants.Exchange3_Price4_AfterTimeT,
		})

	// Exchange 1's Price is before cutoff, so it's ignored
	r := etp.GetValidPrices(log.NewNopLogger(), constants.TimeTPlus1)
	require.Len(t, r, 2)

	expected := []uint64{constants.Price3, constants.Price4}
	assert.ElementsMatch(t, expected, r)
}
