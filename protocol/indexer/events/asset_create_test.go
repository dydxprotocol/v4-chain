package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAssetCreateEvent_Success(t *testing.T) {
	assetCreateEvent := NewAssetCreateEvent(
		0,
		"BTC",
		true,
		0,
		8,
	)
	expectedAssetCreateEventProto := &AssetCreateEventV1{
		Id:               0,
		Symbol:           "BTC",
		HasMarket:        true,
		MarketId:         0,
		AtomicResolution: 8,
	}
	require.Equal(t, expectedAssetCreateEventProto, assetCreateEvent)
}
