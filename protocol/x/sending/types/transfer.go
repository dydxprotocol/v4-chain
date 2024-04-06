package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetSenderSubaccountUpdate returns the sender subaccount update
// for this transfer. Currently only supports quote balance update.
func (t *Transfer) GetSenderSubaccountUpdate() (update types.Update) {
	return types.Update{
		SubaccountId: t.Sender,
		AssetUpdates: []types.AssetUpdate{
			{
				AssetId:       t.AssetId,
				QuantumsDelta: new(int256.Int).Neg(t.GetQuantums()),
			},
		},
	}
}

// GetRecipientSubaccountUpdate returns the recipient subaccount update
// for this transfer. Currently only supports quote balance update.
func (t *Transfer) GetRecipientSubaccountUpdate() (update types.Update) {
	return types.Update{
		SubaccountId: t.Recipient,
		AssetUpdates: []types.AssetUpdate{
			{
				AssetId:       t.AssetId,
				QuantumsDelta: t.GetQuantums(),
			},
		},
	}
}

// GetBigQuantums returns the amount of the transfer in big notional.
// Currently only supports quote balance update.
func (t *Transfer) GetQuantums() *int256.Int {
	return int256.NewUnsignedInt(t.Amount)
}
