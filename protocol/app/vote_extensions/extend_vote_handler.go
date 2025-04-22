package vote_extensions

import (
	"fmt"
	"math/big"

	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// ExtendVoteHandler is a wrapper around the Slinky ExtendVoteHandler. This wrapper is responsible for
// applying the newest MarketPriceUpdates in a block so that the prices to be submitted in a vote extension are
// determined on the latest available information.
type ExtendVoteHandler struct {
	SlinkyExtendVoteHandler sdk.ExtendVoteHandler
	PricesTxDecoder         process.UpdateMarketPriceTxDecoder
	PricesKeeper            PricesKeeper
}

type NoopPriceApplier struct{}

func (n NoopPriceApplier) ApplyPricesFromVoteExtensions(
	_ sdk.Context, _ *cometabci.RequestFinalizeBlock) (map[slinkytypes.CurrencyPair]*big.Int, error) {
	return nil, nil
}
func (n NoopPriceApplier) GetPricesForValidator(_ sdk.ConsAddress) map[slinkytypes.CurrencyPair]*big.Int {
	return nil
}

// ExtendVoteHandler returns a sdk.ExtendVoteHandler, responsible for the following:
//  1. Decoding the x/prices MsgUpdateMarketPrices in the current block - fail on errors
//  2. Validating the proposed MsgUpdateMarketPrices in accordance with the ProcessProposal check
//  3. Updating the market prices in the PricesKeeper so that the GetValidMarketPriceUpdates function returns the
//     latest available market prices
//  4. Calling the Slinky ExtendVoteHandler to handle the rest of ExtendVote
//
// See:
// https://github.com/dydxprotocol/slinky/blob/a5b1d3d3a2723e4746b5d588c512d7cc052dc0ff/abci/ve/vote_extension.go#L77
// for the Slinky ExtendVoteHandler logic.
func (e *ExtendVoteHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (resp *cometabci.ResponseExtendVote, err error) {
		// Decode the x/prices txn in the current block
		updatePricesTx, err := e.PricesTxDecoder.DecodeUpdateMarketPricesTx(ctx, req.Txs)
		if err != nil {
			return nil, fmt.Errorf("DecodeMarketPricesTx failure %w", err)
		}

		// ensure that the proposed MsgUpdateMarketPrices is valid in accordance w/ stateful information
		// this check is equivalent to the check in ProcessProposal (indexPriceCache has not been updated)
		err = updatePricesTx.Validate()
		if err != nil {
			return nil, fmt.Errorf("updatePricesTx.Validate failure %w", err)
		}

		// Update the market prices in the PricesKeeper, so that the GetValidMarketPriceUpdates
		// function returns the latest available market prices.
		updateMarketPricesMsg, ok := updatePricesTx.GetMsg().(*prices.MsgUpdateMarketPrices)
		if !ok {
			return nil, fmt.Errorf("expected %s, got %T", "MsgUpdateMarketPrices", updateMarketPricesMsg)
		}

		// Update the market prices in the PricesKeeper
		err = e.PricesKeeper.UpdateMarketPrices(ctx, updateMarketPricesMsg.MarketPriceUpdates)
		if err != nil {
			return nil, fmt.Errorf("failed to update market prices in extend vote handler pre-slinky invocation %w", err)
		}

		// Call the Slinky ExtendVoteHandler
		return e.SlinkyExtendVoteHandler(ctx, req)
	}
}
