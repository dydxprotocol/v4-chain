package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	kMaxOrderRouterRevSharePpm = lib.OneHundredThousand * 5
)

var _ types.Msg = &MsgSetOrderRouterRevShare{}

// ValidateBasic performs validation to check that order router rev share is under 50%
func (msg *MsgSetOrderRouterRevShare) ValidateBasic() error {
	// Maximum fee share is 500_000 ppm
	if msg.OrderRouterRevShare.SharePpm > kMaxOrderRouterRevSharePpm {
		return errorsmod.Wrapf(
			ErrInvalidRevenueSharePpm,
			"rev share safety violation: rev shares greater than or equal to allowed amount",
		)
	}
	return nil
}
