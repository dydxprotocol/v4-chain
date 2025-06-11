package types

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// GetBuilderFee returns the fee amount for a builder given a fill amount.
// Returns 0 if the builder code is nil.
func (bc *BuilderCodeParameters) GetBuilderFee(fillAmount *big.Int) *big.Int {
	if bc == nil {
		return big.NewInt(0)
	}

	// Calculate fee using ppm (parts per million)
	return lib.BigMulPpm(fillAmount, lib.BigU(bc.FeePpm), true)
}
