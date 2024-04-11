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
		0,
		1000000,
	)
	expectedLiquidityTierUpsertEventProto := &LiquidityTierUpsertEventV2{
		Id:                     0,
		Name:                   "Large-Cap",
		InitialMarginPpm:       50000,
		MaintenanceFractionPpm: 600000,
		OpenInterestLowerCap:   0,
		OpenInterestUpperCap:   1000000,
	}
	require.Equal(t, expectedLiquidityTierUpsertEventProto, liquidityTierUpsertEvent)
}
