package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// 80% of the revenue share is the max allowed. Realistically this will be lower,
	// set it high so we don't need a protocol upgrade to change this
	kMaxOrderRouterRevSharePpm = lib.OneHundredThousand * 8
)

var _ types.Msg = &MsgSetOrderRouterRevShare{}

// ValidateBasic performs validation to check that order router rev share is under 80%
func (msg *MsgSetOrderRouterRevShare) ValidateBasic() error {
	// Maximum fee share is 800_000 ppm
	if msg.OrderRouterRevShare.SharePpm > kMaxOrderRouterRevSharePpm {
		return errorsmod.Wrapf(
			ErrInvalidRevenueSharePpm,
			"rev share safety violation: rev shares greater than or equal to allowed amount",
		)
	}
	return nil
}
