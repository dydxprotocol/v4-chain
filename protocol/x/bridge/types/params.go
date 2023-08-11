package types

import (
	"github.com/dydxprotocol/v4/lib"
)

func (m *EventParams) Validate() error {
	return nil
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
