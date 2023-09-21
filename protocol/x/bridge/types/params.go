package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

func (m *EventParams) Validate() error {
	// TODO(CORE-601): More properly validate Ethereum address.
	if m.EthAddress == "" {
		return errorsmod.Wrap(ErrInvalidEthAddress, "Ethereum contract address cannot be empty")
	}
	return sdk.ValidateDenom(m.Denom)
}

func (m *ProposeParams) Validate() error {
	if m.ProposeDelayDuration < 0 {
		return ErrNegativeDuration
	}
	if m.SkipIfBlockDelayedByDuration < 0 {
		return ErrNegativeDuration
	}
	if m.SkipRatePpm > lib.OneMillion {
		return ErrRateOutOfBounds
	}
	return nil
}

func (m *SafetyParams) Validate() error {
	return nil
}
