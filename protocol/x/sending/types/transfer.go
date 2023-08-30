package types

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetSenderSubaccountUpdate returns the sender subaccount update
// for this transfer. Currently only supports quote balance update.
func (t *Transfer) GetSenderSubaccountUpdate() (update types.Update) {
	return types.Update{
		SubaccountId: t.Sender,
		AssetUpdates: []types.AssetUpdate{
			{
				AssetId:          t.AssetId,
				BigQuantumsDelta: new(big.Int).Neg(t.GetBigQuantums()),
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
				AssetId:          t.AssetId,
				BigQuantumsDelta: t.GetBigQuantums(),
			},
		},
	}
}

// GetBigQuantums returns the amount of the transfer in big notional.
// Currently only supports quote balance update.
func (t *Transfer) GetBigQuantums() (bigNotional *big.Int) {
	return new(big.Int).SetUint64(t.Amount)
}
