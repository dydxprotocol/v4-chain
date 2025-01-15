package events

import (
	"testing"

	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
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
		100, // defaultFundingPpm
	)
	expectedUpdatePerpetualEventProto := &UpdatePerpetualEventV3{
		Id:                   5,
		Ticker:               "BTC-ETH",
		MarketId:             5,
		AtomicResolution:     -8,
		LiquidityTier:        2,
		MarketType:           v1types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		DefaultFunding8HrPpm: 100,
	}
	require.Equal(t, expectedUpdatePerpetualEventProto, updatePerpetualEvent)
}
