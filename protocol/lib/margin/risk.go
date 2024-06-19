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

func (a *Risk) IsInitialCollateralized() bool {
	return a.NC.Cmp(a.IMR) >= 0
}

func (a *Risk) IsMaintenanceCollateralized() bool {
	return a.NC.Cmp(a.MMR) >= 0
}
