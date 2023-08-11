package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateMarketPrices{}

func NewMarketPriceUpdate(id uint32, price uint64) *MsgUpdateMarketPrices_MarketPrice {
	return &MsgUpdateMarketPrices_MarketPrice{
		MarketId: id,
		Price:    price,
	}
}

func NewMsgUpdateMarketPrices(
	proposer string,
	updates []*MsgUpdateMarketPrices_MarketPrice,
) *MsgUpdateMarketPrices {
	return &MsgUpdateMarketPrices{
		Proposer:           proposer,
		MarketPriceUpdates: updates,
	}
}

func (msg *MsgUpdateMarketPrices) GetSigners() []sdk.AccAddress {
	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposer}
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
