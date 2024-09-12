package heap

import (
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// subaccountToDeleverage is a struct containing a subaccount ID and perpetual ID to deleverage.
// This struct is used as a return type for the LiquidateSubaccountsAgainstOrderbook and
// GetSubaccountsWithOpenPositionsInFinalSettlementMarkets called in PrepareCheckState.
type SubaccountToDeleverage struct {
	SubaccountId satypes.SubaccountId
	PerpetualId  uint32
}
