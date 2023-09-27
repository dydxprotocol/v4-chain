package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUpdatePerpetualEventV1_Success(t *testing.T) {
	updatePerpetualEventV1 := NewUpdatePerpetualEventV1(
		5,
		"BTC-ETH",
		5,
		-8,
		2,
	)
	expectedUpdatePerpetualEventV1Proto := &UpdatePerpetualEventV1{
		Id:               5,
		Ticker:           "BTC-ETH",
		MarketId:         5,
		AtomicResolution: -8,
		LiquidityTier:    2,
	}
	require.Equal(t, expectedUpdatePerpetualEventV1Proto, updatePerpetualEventV1)
}
