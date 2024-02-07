package prices

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	"github.com/dydxprotocol/v4-chain/protocol/app/process/errors"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	slinkyabci "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/abci/ve"
)

// SlinkyMarketPriceDecoder wraps an existing MarketPriceDecoder with logic to verify that the MarketPriceUpdateTx was indeed
// derived from vote-extensions injected into the block.
type SlinkyMarketPriceDecoder struct {
	// underlying UpdateMarketPriceTxDecoder
	decoder UpdateMarketPriceTxDecoder

	// underlying vote-extension aggregator
	agg prices.PriceUpdateGenerator
}

// NewSlinkyMarketPriceDecoder returns a new SlinkyMarketPriceDecoder
func NewSlinkyMarketPriceDecoder(decoder UpdateMarketPriceTxDecoder, agg prices.PriceUpdateGenerator) *SlinkyMarketPriceDecoder {
	return &SlinkyMarketPriceDecoder{
		decoder: decoder,
		agg:     agg,
	}
}

// DecodeUpdateMarketPricesTx returns a new `UpdateMarketPricesTx` after validating the following:
//   - the underlying decoder decodes successfully
//   - the UpdateMarketPricesTx follows correctly from the vote-extensions
//   - vote-extensions are enabled: each price per market-id is derived from the injected extended commit
//   - vote-extensions are disabled: no price updates are proposed
//   - vote-extensions are inserted into the block if necessary
//
// If error occurs during any of the checks, returns error.
func (mpd *SlinkyMarketPriceDecoder) DecodeUpdateMarketPricesTx(ctx sdk.Context, txs [][]byte) (*UpdateMarketPricesTx, error) {
	expectedMsg := &pricestypes.MsgUpdateMarketPrices{}

	// check if vote-extensions are enabled
	if ve.VoteExtensionsEnabled(ctx) {
		// if there isn't a vote-extension in the block when there should be, fail
		if len(txs) < slinkyabci.NumInjectedTxs {
			return nil, errors.GetDecodingError(msgUpdateMarketPricesType, fmt.Errorf("expected %v txs, got %v", slinkyabci.NumInjectedTxs, len(txs)))
		}

		// get the expected message from the injected vote-extensions
		var err error
		expectedMsg, err = mpd.agg.GetValidMarketPriceUpdates(ctx, txs[slinkyabci.OracleInfoIndex])
		if err != nil {
			return nil, errors.GetDecodingError(msgUpdateMarketPricesType, err)
		}
	}

	// use the underlying decoder to get the UpdateMarketPricesTx
	updateMarketPrices, err := mpd.decoder.DecodeUpdateMarketPricesTx(ctx, txs)
	if err != nil {
		return nil, err
	}

	updateMarketPricesMsg, ok := updateMarketPrices.GetMsg().(*pricestypes.MsgUpdateMarketPrices)
	if !ok {
		return nil, errors.GetDecodingError(msgUpdateMarketPricesType, fmt.Errorf("expected %T, got %T", expectedMsg, updateMarketPricesMsg))
	}

	// check that the UpdateMarketPricesTx matches the expected message
	if err := checkEqualityOfMarketPriceUpdate(expectedMsg, updateMarketPricesMsg); err != nil {
		return nil, err
	}

	return updateMarketPrices, nil
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

// checkEqualityOfMarketPriceUpdate checks that the given market-price updates are equivalent. Notice,
// this method only checks that the prices per market-id are correct, not that the market-ids are in the
// same order,
func checkEqualityOfMarketPriceUpdate(expectedMsg, actualMsg *pricestypes.MsgUpdateMarketPrices) error {
	// assert len is correct
	if len(expectedMsg.MarketPriceUpdates) != len(actualMsg.MarketPriceUpdates) {
		return IncorrectNumberUpdatesError(len(expectedMsg.MarketPriceUpdates), len(actualMsg.MarketPriceUpdates))
	}

	// map market-id to price-map
	expectedPricesPerMarketID := make(map[uint32]uint64)
	for _, marketPriceUpdate := range expectedMsg.MarketPriceUpdates {
		expectedPricesPerMarketID[marketPriceUpdate.MarketId] = marketPriceUpdate.Price
	}

	// check that the actual prices match the expected prices
	for _, marketPriceUpdate := range actualMsg.MarketPriceUpdates {
		expectedPrice, ok := expectedPricesPerMarketID[marketPriceUpdate.MarketId]
		if !ok {
			return MissingPriceUpdateForMarket(marketPriceUpdate.MarketId)
		}

		if expectedPrice != marketPriceUpdate.Price {
			return IncorrectPriceUpdateForMarket(marketPriceUpdate.MarketId, expectedPrice, marketPriceUpdate.Price)
		}
	}

	return nil
}
