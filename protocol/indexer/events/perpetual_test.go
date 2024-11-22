package events

import (
	"testing"

	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	"github.com/stretchr/testify/require"
)

func TestNewUpdatePerpetualEvent_Success(t *testing.T) {
	updatePerpetualEvent := NewUpdatePerpetualEvent(
		5,
		"BTC-ETH",
		5,
		-8,
		2,
		perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
	)
	expectedUpdatePerpetualEventProto := &UpdatePerpetualEventV2{
		Id:               5,
		Ticker:           "BTC-ETH",
		MarketId:         5,
		AtomicResolution: -8,
		LiquidityTier:    2,
		MarketType:       perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
	}
	require.Equal(t, expectedUpdatePerpetualEventProto, updatePerpetualEvent)
}
