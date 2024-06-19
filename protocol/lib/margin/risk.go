package margin

import (
	"math/big"
)

type Risk struct {
	MMR *big.Int // Maintenance Margin Requirement
	IMR *big.Int // Initial Margin Requirement
	NC  *big.Int // Net Collateral
}

// AddInPlace adds the values of b to a (in-place).
func (a *Risk) AddInPlace(b Risk) {
	a.MMR.Add(a.MMR, b.MMR)
	a.IMR.Add(a.IMR, b.IMR)
	a.NC.Add(a.NC, b.NC)
}

// IsInitialCollateralized returns true if the account has enough net collateral to meet the
// initial margin requirement.
func (a *Risk) IsInitialCollateralized() bool {
	return a.NC.Cmp(a.IMR) >= 0
}

// IsMaintenanceCollateralized returns true if the account has enough net collateral to meet the
// maintenance margin requirement.
func (a *Risk) IsMaintenanceCollateralized() bool {
	return a.NC.Cmp(a.MMR) >= 0
}

// IsLiquidatable returns true if the account is liquidatable given its maintenance margin requirement
// and net collateral.
//
// The account is liquidatable if both of the following are true:
// - The maintenance margin requirements are greater than zero (note that they can never be negative).
// - The maintenance margin requirements are greater than the account's net collateral.
func (a *Risk) IsLiquidatable() bool {
	return a.MMR.Sign() > 0 && a.MMR.Cmp(a.NC) > 0
}
