package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/pkg/errors"
)

func (p *Perpetual) GetId() uint32 {
	return p.Params.Id
}

// Stateless validation on Perpetual params.
func (p *PerpetualParams) Validate() error {
	// Validate `ticker`.
	if len(p.Ticker) == 0 {
		return errors.WithStack(ErrTickerEmptyString)
	}

	// Validate `defaultFundingPpm`
	defaultFundingPpm := lib.AbsInt32(p.DefaultFundingPpm)
	if defaultFundingPpm > MaxDefaultFundingPpmAbs {
		return errorsmod.Wrap(
			ErrDefaultFundingPpmMagnitudeExceedsMax,
			lib.Int32ToString(p.DefaultFundingPpm))
	}

	return nil
}
