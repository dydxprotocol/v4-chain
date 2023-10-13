package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/json"
)

// Validate checks that the MarketParam is valid.
func (mp *MarketParam) Validate() error {
	// Validate pair.
	if mp.Pair == "" {
		return errorsmod.Wrap(ErrInvalidInput, "Pair cannot be empty")
	}

	if mp.MinExchanges == 0 {
		return ErrZeroMinExchanges
	}

	// Validate min price change.
	if mp.MinPriceChangePpm == 0 || mp.MinPriceChangePpm >= lib.MaxPriceChangePpm {
		return errorsmod.Wrapf(
			ErrInvalidInput,
			"Min price change in parts-per-million must be greater than 0 and less than %d",
			lib.MaxPriceChangePpm)
	}

	if err := json.IsValidJSON(mp.ExchangeConfigJson); err != nil {
		return errorsmod.Wrapf(
			ErrInvalidInput,
			"ExchangeConfigJson string is not valid: err=%v, input=%v",
			err,
			mp.ExchangeConfigJson,
		)
	}

	return nil
}
