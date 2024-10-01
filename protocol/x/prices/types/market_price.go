package types

import (
	errorsmod "cosmossdk.io/errors"
)

// ValidateFromParam checks that the MarketPrice is valid and that it corresponds to the given MarketParam.
func (mp *MarketPrice) ValidateFromParam(marketParam MarketParam) error {
	if marketParam.Id != mp.Id {
		return errorsmod.Wrapf(
			ErrInvalidInput,
			"market param id %d does not match market price id %d",
			marketParam.Id,
			mp.Id,
		)
	}
	return nil
}
