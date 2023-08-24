package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLiquidityTierUpsertEvent_Success(t *testing.T) {
	liquidityTierUpsertEvent := NewLiquidityTierUpsertEvent(
		0,
		"Large-Cap",
		50000,
		600000,
		1000000000000,
	)
	expectedLiquidityTierUpsertEventProto := &LiquidityTierUpsertEventV1{
		Id:                     0,
		Name:                   "Large-Cap",
		InitialMarginPpm:       50000,
		MaintenanceFractionPpm: 600000,
		BasePositionNotional:   1000000000000,
	}
	require.Equal(t, expectedLiquidityTierUpsertEventProto, liquidityTierUpsertEvent)
}
