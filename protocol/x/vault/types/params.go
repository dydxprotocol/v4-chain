package types

import (
	"math"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
)

// DefaultQuotingParams returns a default set of `x/vault` parameters.
func DefaultQuotingParams() QuotingParams {
	return QuotingParams{
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
func (p QuotingParams) Validate() error {
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
	// Skew factor times order_size_pct must be less than 2 to avoid skewing over the spread
	if uint64(p.SkewFactorPpm)*uint64(p.OrderSizePctPpm) >= 2_000_000*1_000_000 {
		return ErrInvalidSkewFactor
	}

	return nil
}

// Validate validates individual vault parameters.
func (v VaultParams) Validate() error {
	// Validate status.
	if v.Status == VaultStatus_VAULT_STATUS_UNSPECIFIED {
		return ErrUnspecifiedVaultStatus
	}
	// Validate quoting params.
	if v.QuotingParams != nil {
		if err := v.QuotingParams.Validate(); err != nil {
			return err
		}
	}

	return nil
}
