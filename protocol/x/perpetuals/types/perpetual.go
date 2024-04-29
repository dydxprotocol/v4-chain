package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/pkg/errors"
)

func (p *Perpetual) GetId() uint32 {
	return p.Params.Id
}

// Stateless validation on Perpetual params.
func (p *PerpetualParams) Validate() error {
	// Check if market type is valid
	if p.MarketType != PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS &&
		p.MarketType != PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
		return errorsmod.Wrap(
			ErrInvalidMarketType,
			fmt.Sprintf("market type %v", p.MarketType),
		)
	}

	// Validate `ticker`.
	if len(p.Ticker) == 0 {
		return errors.WithStack(ErrTickerEmptyString)
	}

	// Validate `defaultFundingPpm`
	defaultFundingPpm := lib.AbsInt32(p.DefaultFundingPpm)
	if defaultFundingPpm > MaxDefaultFundingPpmAbs {
		return errorsmod.Wrap(
			ErrDefaultFundingPpmMagnitudeExceedsMax,
			lib.IntToString(p.DefaultFundingPpm))
	}

	return nil
}
