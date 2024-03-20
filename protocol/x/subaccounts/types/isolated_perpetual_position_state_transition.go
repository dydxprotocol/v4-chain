package types

import "math/big"

type PositionStateTransition uint

const (
	Opened PositionStateTransition = iota
	Closed
)

var positionStateTransitionStringMap = map[PositionStateTransition]string{
	Opened: "opened",
	Closed: "closed",
}

func (t PositionStateTransition) String() string {
	result, exists := positionStateTransitionStringMap[t]
	if !exists {
		return "UnexpectedStateTransitionError"
	}

	return result
}

// Represents a state transition for an isolated perpetual position.
type IsolatedPerpetualPositionStateTransition struct {
	PerpetualId uint32
	// TODO(DEC-715): Support non-USDC assets.
	// Quote quantums position size of the subaccount that has a state change for an isolated perpetual.
	QuoteQuantumsBeforeUpdate *big.Int
	// The state transition that occurred for the isolated perpetual positions.
	Transition PositionStateTransition
}
