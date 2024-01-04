package types

import (
	"fmt"
	"math"

	storetypes "cosmossdk.io/store/types"
)

type freeInfiniteGasMeter struct {
}

// NewFreeInfiniteGasMeter returns a new gas meter without a limit, where gas is not consumed.
func NewFreeInfiniteGasMeter() storetypes.GasMeter {
	return &freeInfiniteGasMeter{}
}

// GasConsumed returns the gas consumed from the GasMeter, which is always 0 in freeInfiniteGasMeter
func (g *freeInfiniteGasMeter) GasConsumed() storetypes.Gas {
	return 0
}

// GasConsumedToLimit returns the gas consumed from the GasMeter since the gas is not confined to a limit.
// NOTE: This behaviour is only called when recovering from panic when BlockGasMeter consumes gas past the limit.
// This should never occur for the freeInfiniteGasMeter, which never increments gas consumed.
func (g *freeInfiniteGasMeter) GasConsumedToLimit() storetypes.Gas {
	return 0
}

// GasRemaining returns MaxUint64 since limit is not confined in freeInfiniteGasMeter.
func (g *freeInfiniteGasMeter) GasRemaining() storetypes.Gas {
	return math.MaxUint64
}

// Limit returns MaxUint64 since limit is not confined in freeInfiniteGasMeter.
func (g *freeInfiniteGasMeter) Limit() storetypes.Gas {
	return math.MaxUint64
}

// ConsumeGas is a no-op as no gas is consumed with the freeInfiniteGasMeter.
func (g *freeInfiniteGasMeter) ConsumeGas(amount storetypes.Gas, descriptor string) {
}

// RefundGas is a no-op as no gas is consumed with the freeInfiniteGasMeter.
func (g *freeInfiniteGasMeter) RefundGas(amount storetypes.Gas, descriptor string) {
}

// IsPastLimit returns false since the gas limit is not confined.
func (g *freeInfiniteGasMeter) IsPastLimit() bool {
	return false
}

// IsOutOfGas returns false since the gas limit is not confined.
func (g *freeInfiniteGasMeter) IsOutOfGas() bool {
	return false
}

// String returns the FreeInfiniteGasMeter's gas consumed.
func (g *freeInfiniteGasMeter) String() string {
	return fmt.Sprintf("FreeInfiniteGasMeter:\n  consumed: %d", 0)
}
