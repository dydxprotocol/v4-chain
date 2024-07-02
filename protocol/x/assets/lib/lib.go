package lib

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

// GetNetCollateralAndMarginRequirements returns the net collateral, initial margin, and maintenance margin
// that a given position (quantums) for a given assetId contributes to an account.
func GetNetCollateralAndMarginRequirements(
	id uint32,
	bigQuantums *big.Int,
) (
	risk margin.Risk,
	err error,
) {
	risk = margin.ZeroRisk()

	// Balance is zero.
	if bigQuantums.BitLen() == 0 {
		return risk, nil
	}

	// USDC.
	if id == types.AssetUsdc.Id {
		risk.NC = new(big.Int).Set(bigQuantums)
		return risk, nil
	}

	// Balance is positive.
	// TODO(DEC-581): add multi-collateral support.
	if bigQuantums.Sign() == 1 {
		return risk, types.ErrNotImplementedMulticollateral
	}

	// Balance is negative.
	// TODO(DEC-582): margin-trading
	return risk, types.ErrNotImplementedMargin
}
