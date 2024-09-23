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
		1000000,
		"0/1",
	)
	expectedUpdatePerpetualEventV1Proto := &UpdatePerpetualEventV1{
		Id:               5,
		Ticker:           "BTC-ETH",
		MarketId:         5,
		AtomicResolution: -8,
		LiquidityTier:    2,
		DangerIndexPpm:   1000000,
		PerpYieldIndex:   "0/1",
	}
	require.Equal(t, expectedUpdatePerpetualEventV1Proto, updatePerpetualEventV1)
}
