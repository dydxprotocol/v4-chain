package process

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/slinky/abci/ve"

	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// SlinkyMarketPriceDecoder wraps an existing MarketPriceDecoder with logic to verify that the MarketPriceUpdateTx
// was indeed derived from vote-extensions injected into the block.
type SlinkyMarketPriceDecoder struct {
	// underlying UpdateMarketPriceTxDecoder
	decoder UpdateMarketPriceTxDecoder

	// underlying vote-extension aggregator
	agg prices.PriceUpdateGenerator
}

// NewSlinkyMarketPriceDecoder returns a new SlinkyMarketPriceDecoder
func NewSlinkyMarketPriceDecoder(
	decoder UpdateMarketPriceTxDecoder, agg prices.PriceUpdateGenerator) *SlinkyMarketPriceDecoder {
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
func (mpd *SlinkyMarketPriceDecoder) DecodeUpdateMarketPricesTx(
	ctx sdk.Context, txs [][]byte) (*UpdateMarketPricesTx, error) {
	expectedMsg := &pricestypes.MsgUpdateMarketPrices{}

	// check if vote-extensions are enabled
	if ve.VoteExtensionsEnabled(ctx) {
		// if there isn't a vote-extension in the block when there should be, fail
		if len(txs) < constants.OracleVEInjectedTxs {
			return nil, getDecodingError(
				msgUpdateMarketPricesType,
				fmt.Errorf("expected %v txs, got %v", constants.OracleVEInjectedTxs, len(txs)))
		}

		// get the expected message from the injected vote-extensions
		var err error
		expectedMsg, err = mpd.agg.GetValidMarketPriceUpdates(ctx, txs[constants.OracleInfoIndex])
		if err != nil {
			return nil, getDecodingError(msgUpdateMarketPricesType, err)
		}
	}

	// use the underlying decoder to get the UpdateMarketPricesTx
	// If VE are not enabled with Slinky, then there should be no price updates
	updateMarketPrices, err := mpd.decoder.DecodeUpdateMarketPricesTx(ctx, txs)
	if err != nil {
		return nil, err
	}

	updateMarketPricesMsg, ok := updateMarketPrices.GetMsg().(*pricestypes.MsgUpdateMarketPrices)
	if !ok {
		return nil, getDecodingError(
			msgUpdateMarketPricesType, fmt.Errorf("expected %T, got %T", expectedMsg, updateMarketPricesMsg))
	}

	// check that the UpdateMarketPricesTx matches the expected message
	if err := checkEqualityOfMarketPriceUpdate(expectedMsg, updateMarketPricesMsg); err != nil {
		return nil, err
	}

	return updateMarketPrices, nil
}

// GetTxOffset returns the offset that other injected txs should be placed with respect to their normally
// expected indices. If vote-extensions are enabled, constants.OracleVEInjectedTxs is the expected offset,
// otherwise 0 is the expected offset.
func (mpd *SlinkyMarketPriceDecoder) GetTxOffset(ctx sdk.Context) int {
	if ve.VoteExtensionsEnabled(ctx) {
		return constants.OracleVEInjectedTxs
	}
	return 0
}

// checkEqualityOfMarketPriceUpdate checks that the given market-price updates are equivalent
// and both pass validate basic checks.
func checkEqualityOfMarketPriceUpdate(expectedMsg, actualMsg *pricestypes.MsgUpdateMarketPrices) error {
	// assert that the market-price updates are valid
	if err := expectedMsg.ValidateBasic(); err != nil {
		return InvalidMarketPriceUpdateError(err)
	}

	if err := actualMsg.ValidateBasic(); err != nil {
		return InvalidMarketPriceUpdateError(err)
	}
	// assert len is correct
	if len(expectedMsg.MarketPriceUpdates) != len(actualMsg.MarketPriceUpdates) {
		return IncorrectNumberUpdatesError(len(expectedMsg.MarketPriceUpdates), len(actualMsg.MarketPriceUpdates))
	}

	// check that the actual prices match the expected prices (both are sorted so markets are in the same order)
	for i, actualUpdate := range actualMsg.MarketPriceUpdates {
		expectedUpdate := expectedMsg.MarketPriceUpdates[i]

		if expectedUpdate.MarketId != actualUpdate.MarketId {
			return MissingPriceUpdateForMarket(expectedUpdate.MarketId)
		}

		if expectedUpdate.Price != actualUpdate.Price {
			return IncorrectPriceUpdateForMarket(expectedUpdate.MarketId, expectedUpdate.Price, actualUpdate.Price)
		}
	}

	return nil
}
