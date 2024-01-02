package types_test

import (
	storetypes "cosmossdk.io/store/types"
	"math"
	"testing"

	ante_types "github.com/dydxprotocol/v4-chain/protocol/app/ante/types"
	"github.com/stretchr/testify/require"
)

func TestInfiniteGasMeter(t *testing.T) {
	meter := ante_types.NewFreeInfiniteGasMeter()
	require.Equal(t, uint64(math.MaxUint64), meter.Limit())
	require.Equal(t, uint64(math.MaxUint64), meter.GasRemaining())
	require.Equal(t, uint64(0), meter.GasConsumed())
	require.Equal(t, uint64(0), meter.GasConsumedToLimit())
	meter.ConsumeGas(10, "consume 10")
	require.Equal(t, uint64(math.MaxUint64), meter.GasRemaining())
	require.Equal(t, uint64(0), meter.GasConsumed())
	require.Equal(t, uint64(0), meter.GasConsumedToLimit())
	meter.RefundGas(1, "refund 1")
	require.Equal(t, uint64(math.MaxUint64), meter.GasRemaining())
	require.Equal(t, uint64(0), meter.GasConsumed())
	require.False(t, meter.IsPastLimit())
	require.False(t, meter.IsOutOfGas())
	meter.ConsumeGas(storetypes.Gas(math.MaxUint64/2), "consume half max uint64")
	require.NotPanics(t, func() { meter.ConsumeGas(storetypes.Gas(math.MaxUint64/2)+2, "panic") })
	require.NotPanics(t, func() { meter.RefundGas(meter.GasConsumed()+1, "refund greater than consumed") })
	require.Equal(t, "FreeInfiniteGasMeter:\n  consumed: 0", meter.String())
}
