package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestNewMarketToPrice_IsEmpty(t *testing.T) {
	mtp := types.NewMarketToPrice()

	require.Empty(t, mtp.MarketToPriceTimestamp)
}

func TestUpdatePrice_Valid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market9_TimeT_Price1))

	require.Len(t, mtp.MarketToPriceTimestamp, 1)

	marketPriceTimestamp := mtp.GetAllPrices()[0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateValid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market9_TimeTMinusThreshold_Price2))
	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market9_TimeT_Price1))

	require.Len(t, mtp.MarketToPriceTimestamp, 1)

	marketPriceTimestamp := mtp.GetAllPrices()[0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateInvalid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market9_TimeT_Price1))
	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market9_TimeTMinusThreshold_Price2))

	require.Len(t, mtp.MarketToPriceTimestamp, 1)

	marketPriceTimestamp := mtp.GetAllPrices()[0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateForTwoMarketsValid(t *testing.T) {
	mtp := types.NewMarketToPrice()

	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market9_TimeT_Price1))
	mtp.UpdatePrice((*types.MarketPriceTimestamp)(constants.Market8_TimeTMinusThreshold_Price2))

	require.Len(t, mtp.MarketToPriceTimestamp, 2)

	marketPriceTimestamp := mtp.MarketToPriceTimestamp[constants.MarketId9]
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdateTime)

	marketPriceTimestamp2 := mtp.MarketToPriceTimestamp[constants.MarketId8]
	require.Equal(t, constants.Price2, marketPriceTimestamp2.Price)
	require.Equal(t, constants.TimeTMinusThreshold, marketPriceTimestamp2.LastUpdateTime)
}
