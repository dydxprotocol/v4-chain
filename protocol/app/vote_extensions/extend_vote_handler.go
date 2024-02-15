package vote_extensions

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type ExtendVoteHandler struct {
	SlinkyExtendVoteHandler sdk.ExtendVoteHandler
	PricesTxDecoder         process.UpdateMarketPriceTxDecoder
	PricesKeeper            PricesKeeper
}

func (e *ExtendVoteHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (resp *cometabci.ResponseExtendVote, err error) {
		// Decode the x/prices txn in the current block
		updatePricesTx, err := e.PricesTxDecoder.DecodeUpdateMarketPricesTx(ctx, req.Txs)
		if err != nil {
			return nil, fmt.Errorf("DecodeMarketPricesTx failure %w", err)
		}
		// Apply the changes to the PricesKeeper so valid prices are proposed
		if err = e.PricesKeeper.UpdateSmoothedPrices(ctx, lib.Uint64LinearInterpolate); err != nil {
			return nil, fmt.Errorf("UpdateSmoothedPrices failure %w", err)
		}
		err = updatePricesTx.Validate()
		if err != nil {
			return nil, fmt.Errorf("updatePricesTx.Validate failure %w", err)
		}
		updateMarketPricesMsg, ok := updatePricesTx.GetMsg().(*prices.MsgUpdateMarketPrices)
		if !ok {
			return nil, fmt.Errorf("expected %s, got %T", "MsgUpdateMarketPrices", updateMarketPricesMsg)
		}
		err = e.PricesKeeper.UpdateMarketPrices(ctx, updateMarketPricesMsg.MarketPriceUpdates)
		if err != nil {
			return nil, fmt.Errorf("failed to update market prices in extend vote handler pre-slinky invocation %w", err)
		}
		// Call the Slinky ExtendVoteHandler
		return e.SlinkyExtendVoteHandler(ctx, req)
	}
}
