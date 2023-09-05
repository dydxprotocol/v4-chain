package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateMarketPrices{}

func NewMarketPriceUpdate(id uint32, price uint64) *MsgUpdateMarketPrices_MarketPrice {
	return &MsgUpdateMarketPrices_MarketPrice{
		MarketId: id,
		Price:    price,
	}
}

func NewMsgUpdateMarketPrices(updates []*MsgUpdateMarketPrices_MarketPrice) *MsgUpdateMarketPrices {
	return &MsgUpdateMarketPrices{
		MarketPriceUpdates: updates,
	}
}

func (msg *MsgUpdateMarketPrices) GetSigners() []sdk.AccAddress {
	// Return empty slice because app-injected msg is not expected to be signed.
	return []sdk.AccAddress{}
}

// ValidateBasic performs stateless validations on the message. Specifically:
// - Update prices are non-zero.
// - Updates are sorted by market id in ascending order.
// - Updates do not contain duplicate markets.
func (msg *MsgUpdateMarketPrices) ValidateBasic() error {
	for i, marketPriceUpdate := range msg.MarketPriceUpdates {
		// Check price is not 0.
		if marketPriceUpdate.Price == uint64(0) {
			return sdkerrors.Wrapf(
				ErrInvalidMarketPriceUpdateStateless,
				"price cannot be 0 for market id (%d)",
				marketPriceUpdate.MarketId,
			)
		}

		// Check updates are sorted by market id and there are no duplicates.
		if i > 0 && msg.MarketPriceUpdates[i-1].MarketId >= marketPriceUpdate.MarketId {
			return sdkerrors.Wrap(
				ErrInvalidMarketPriceUpdateStateless,
				"market price updates must be sorted by market id in ascending order and cannot contain duplicates",
			)
		}
	}
	return nil
}
