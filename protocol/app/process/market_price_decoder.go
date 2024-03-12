package process

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// MarketPriceDecoder is an interface for decoding market price transactions, This interface is responsible
// for distinguishing between logic for unmarshalling MarketPriceUpdates, between MarketPriceUpdates
// determined by the proposer's price-cache, and from VoteExtensions.
type UpdateMarketPriceTxDecoder interface {
	// DecodeUpdateMarketPricesTx decodes the tx-bytes from the RequestProcessProposal and returns a MarketPriceUpdateTx.
	DecodeUpdateMarketPricesTx(ctx sdk.Context, txs [][]byte) (*UpdateMarketPricesTx, error)

	// GetTxOffset returns the offset that other injected txs should be placed with respect to their normally
	// expected indices. This method is used to account for injected vote-extensions, or any other injected
	// txs from dependencies.
	GetTxOffset(ctx sdk.Context) int
}

func NewUpdateMarketPricesTx(
	ctx sdk.Context, pk ProcessPricesKeeper, msg *pricestypes.MsgUpdateMarketPrices) *UpdateMarketPricesTx {
	return &UpdateMarketPricesTx{
		ctx:          ctx,
		pricesKeeper: pk,
		msg:          msg,
	}
}
