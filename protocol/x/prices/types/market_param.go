package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// Validate checks that the MarketParam is valid.
func (mp *MarketParam) Validate() error {
	// Validate pair.
	if mp.Pair == "" {
		return errorsmod.Wrap(ErrInvalidInput, "Pair cannot be empty")
	}

	// Validate min price change.
	if mp.MinPriceChangePpm == 0 || mp.MinPriceChangePpm >= lib.MaxPriceChangePpm {
		return errorsmod.Wrapf(
			ErrInvalidInput,
			"Min price change in parts-per-million must be greater than 0 and less than %d",
			lib.MaxPriceChangePpm)
	}

	return nil
}
