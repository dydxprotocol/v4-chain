package types

import "math"

// DefaultParams returns a default set of `x/vault` parameters.
func DefaultParams() Params {
	return Params{
		Layers:                 2,       // 2 layers
		SpreadMinPpm:           3_000,   // 30 bps
		SpreadBufferPpm:        1_500,   // 15 bps
		SkewFactorPpm:          500_000, // 0.5
		OrderSizePpm:           100_000, // 10%
		OrderExpirationSeconds: 2,       // 2 seconds
	}
}

// Validate validates `x/vault` parameters.
func (p Params) Validate() error {
	// Layers must be less than or equal to MaxUint8.
	if p.Layers > math.MaxUint8 {
		return ErrInvalidLayers
	}
	// Spread min ppm must be positive.
	if p.SpreadMinPpm == 0 {
		return ErrInvalidSpreadMinPpm
	}
	// Order size must be positive.
	if p.OrderSizePpm == 0 {
		return ErrInvalidOrderSizePpm
	}
	// Order expiration seconds must be positive.
	if p.OrderExpirationSeconds == 0 {
		return ErrInvalidOrderExpirationSeconds
	}

	return nil
}
