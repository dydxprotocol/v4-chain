package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	"github.com/skip-mev/slinky/abci/ve"
	slinkyabci "github.com/skip-mev/slinky/abci/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process/errors"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"fmt"
)

// SlinkyMarketPriceDecoder wraps an existing MarketPriceDecoder with logic to verify that the MarketPriceUpdateTx was indeed
// derived from vote-extensions injected into the block.
type SlinkyMarketPriceDecoder struct {
	// underlying UpdateMarketPriceTxDecoder
	decoder UpdateMarketPriceTxDecoder

	// underlying vote-extension aggregator
	agg prices.PriceUpdateGenerator
}

// DecodeUpdateMarketPricesTx returns a new `UpdateMarketPricesTx` after validating the following:
//   - the underlying decoder decodes successfully
//   - the UpdateMarketPricesTx follows correctly from the vote-extensions
//   - vote-extensions are inserted into the block if necessary
//
// If error occurs during any of the checks, returns error.
func (mpd *SlinkyMarketPriceDecoder) DecodeUpdateMarketPricesTx(ctx sdk.Context, txs [][]byte) (*UpdateMarketPricesTx, error) {
	var extendedCommitBz []byte

	// check if vote-extensions are enabled
	if ve.VoteExtensionsEnabled(ctx) {
		// if there isn't a vote-extension in the block when there should be, fail
		if len(txs) < slinkyabci.NumInjectedTxs {
			return nil, errors.GetDecodingError(msgUpdateMarketPricesType, fmt.Errorf("expected %v txs, got %v", slinkyabci.NumInjectedTxs, len(txs)))
		}

		// get the expected extended commit bytes
		extendedCommitBz = txs[slinkyabci.NumInjectedTxs]
	}

	// use the underlying decoder to get the UpdateMarketPricesTx
	updateMarketPrices, err := mpd.decoder.DecodeUpdateMarketPricesTx(ctx, txs)
	if err != nil {
		return nil, err
	}

	// get the expected message from the injected vote-extensions
	expectedMsg, err := mpd.agg.GetValidMarketPriceUpdates(ctx, extendedCommitBz)
	if err != nil {
		return nil, err
	}

	// check that the UpdateMarketPricesTx matches the expected message
	return updateMarketPrices, checkEqualityOfMarketPriceUpdate(expectedMsg, updateMarketPrices.GetMsg())
}

// GetTxOffset returns the offset that other injected txs should be placed with respect to their normally
// expected indices. If vote-extensions are enabled, slinkyabci.NumInjectedTxs is the expected offset, 
// otherwise 0 is the expected offset.
func (mpd *SlinkyMarketPriceDecoder) GetTxOffset(ctx sdk.Context) int {
	if ve.VoteExtensionsEnabled(ctx) {
		return slinkyabci.NumInjectedTxs
	}
	return 0
}

// checkEqualityOfMarketPriceUpdate checks that the given market-price updates are equivalent
func checkEqualityOfMarketPriceUpdate(expectedMsgI, actualMsgI sdk.Msg) error {
	expectedMsg, ok := expectedMsgI.(*pricestypes.MsgUpdateMarketPrices)
	if !ok {
		return fmt.Errorf("expected message to be of type %T, got %T", expectedMsg, expectedMsgI)
	}

	actualMsg, ok := actualMsgI.(*pricestypes.MsgUpdateMarketPrices)
	if !ok {
		return fmt.Errorf("actual message to be of type %T, got %T", actualMsg, actualMsgI)
	}


	// assert len is correct
	if len(expectedMsg.MarketPriceUpdates) != len(actualMsg.MarketPriceUpdates) {
		return fmt.Errorf("expected %v market-price updates, got %v", len(expectedMsg.MarketPriceUpdates), len(actualMsg.MarketPriceUpdates))
	}

	// assert each market-price update is correct
	for i, expectedMarketPriceUpdate := range expectedMsg.MarketPriceUpdates {
		actualMarketPriceUpdate := actualMsg.MarketPriceUpdates[i]
		if expectedMarketPriceUpdate.MarketId != actualMarketPriceUpdate.MarketId {
			return fmt.Errorf("expected market id %v, got %v", expectedMarketPriceUpdate.MarketId, actualMarketPriceUpdate.MarketId)
		}
		if expectedMarketPriceUpdate.Price != actualMarketPriceUpdate.Price {
			return fmt.Errorf("expected price %v, got %v", expectedMarketPriceUpdate.Price, actualMarketPriceUpdate.Price)
		}
	}

	return nil
}
