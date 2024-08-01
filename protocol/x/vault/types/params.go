package types

import (
	"math"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
)

// DefaultParams returns a default set of `x/vault` parameters.
func DefaultParams() Params {
	return Params{
		Layers:                           2,                            // 2 layers
		SpreadMinPpm:                     10_000,                       // 100 bps
		SpreadBufferPpm:                  1_500,                        // 15 bps
		SkewFactorPpm:                    2_000_000,                    // 2
		OrderSizePctPpm:                  100_000,                      // 10%
		OrderExpirationSeconds:           60,                           // 60 seconds
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000), // 1_000 USDC
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
	if p.OrderSizePctPpm == 0 {
		return ErrInvalidOrderSizePctPpm
	}
	// Order expiration seconds must be positive.
	if p.OrderExpirationSeconds == 0 {
		return ErrInvalidOrderExpirationSeconds
	}
	// Activation threshold quote quantums must be non-negative.
	if p.ActivationThresholdQuoteQuantums.Sign() < 0 {
		return ErrInvalidActivationThresholdQuoteQuantums
	}

	return nil
}

// Validate validates individual vault parameters.
func (v VaultParams) Validate() error {
	return nil
}
