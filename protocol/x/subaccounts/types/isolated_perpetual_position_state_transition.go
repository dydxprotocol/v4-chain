package types

import "github.com/dydxprotocol/v4-chain/protocol/lib/int256"

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

// Represents a state transition for an isolated perpetual position in a subaccount.
type IsolatedPerpetualPositionStateTransition struct {
	SubaccountId *SubaccountId
	PerpetualId  uint32
	// TODO(DEC-715): Support non-USDC assets.
	// Quote quantums of collateral to transfer as a result of the state transition.
	QuoteQuantums *int256.Int
	// The state transition that occurred for the isolated perpetual positions.
	Transition PositionStateTransition
}
