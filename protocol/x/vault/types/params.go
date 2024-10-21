package types

import (
	"math"
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	skewFactor := new(big.Int).SetUint64(uint64(p.SkewFactorPpm))
	orderSizePct := new(big.Int).SetUint64(uint64(p.OrderSizePctPpm))
	skewFactorOrderSizePctProduct := new(big.Int).Mul(skewFactor, orderSizePct)
	skewFactorOrderSizePctProductThreshold := big.NewInt(2_000_000 * 1_000_000)
	if skewFactorOrderSizePctProduct.Cmp(skewFactorOrderSizePctProductThreshold) >= 0 {
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

// Validate validates OperatorParams.
func (o OperatorParams) Validate() error {
	// Validate that operator is non-empty.
	if o.Operator == "" {
		return ErrEmptyOperator
	}

	return nil
}

// DefaultOperatorParams returns a default set of `x/vault` operator parameters.
func DefaultOperatorParams() OperatorParams {
	return OperatorParams{
		Operator: lib.GovModuleAddress.String(),
		Metadata: OperatorMetadata{
			Name:        "Governance",
			Description: "Governance Module Account",
		},
	}
}
