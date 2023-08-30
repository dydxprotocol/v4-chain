package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

var (
	marketId          = constants.MarketId0
	priceWithExponent = constants.Price5
	pair              = constants.BtcUsdPair
	minPriceChangePpm = uint32(50)
	exponent          = constants.BtcUsdExponent
)

func TestNewMarketPriceUpdateEvent_Success(t *testing.T) {
	priceUpdateEvent := events.NewMarketPriceUpdateEvent(marketId, priceWithExponent)
	expectedMarketEventProto := &events.MarketEventV1{
		MarketId: marketId,
		Event: &events.MarketEventV1_PriceUpdate{
			PriceUpdate: &events.MarketPriceUpdateEventV1{
				PriceWithExponent: priceWithExponent,
			},
		},
	}
	require.Equal(t, expectedMarketEventProto, priceUpdateEvent)
}

func TestNewMarketModifyEvent_Success(t *testing.T) {
	marketModifyEvent := events.NewMarketModifyEvent(marketId, pair, minPriceChangePpm)
	expectedMarketEventProto := &events.MarketEventV1{
		MarketId: marketId,
		Event: &events.MarketEventV1_MarketModify{
			MarketModify: &events.MarketModifyEventV1{
				Base: &events.MarketBaseEventV1{
					Pair:              pair,
					MinPriceChangePpm: minPriceChangePpm,
				},
			},
		},
	}
	require.Equal(t, expectedMarketEventProto, marketModifyEvent)
}

func TestNewMarketCreateEvent_Success(t *testing.T) {
	marketCreateEvent := events.NewMarketCreateEvent(
		marketId,
		pair,
		minPriceChangePpm,
		int32(exponent),
	)
	expectedMarketEventProto := &events.MarketEventV1{
		MarketId: marketId,
		Event: &events.MarketEventV1_MarketCreate{
			MarketCreate: &events.MarketCreateEventV1{
				Base: &events.MarketBaseEventV1{
					Pair:              pair,
					MinPriceChangePpm: minPriceChangePpm,
				},
				Exponent: int32(exponent),
			},
		},
	}
	require.Equal(t, expectedMarketEventProto, marketCreateEvent)
}
