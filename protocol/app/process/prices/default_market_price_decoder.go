package prices

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process/errors"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"reflect"
)

var (
	msgUpdateMarketPricesType = reflect.TypeOf(pricestypes.MsgUpdateMarketPrices{})
)

const (
	UpdateMarketPricesTxLenOffset = -1
)

// DefaultUpdateMarketPriceTxDecoder is the default implementation of the MarketPriceDecoder interface.
// This implementation is expected to default MarketPriceUpdates in accordance with the dydx process-proposal
// logic pre vote-extensions
type DefaultUpdateMarketPriceTxDecoder struct {
	// pk is the expecte dependency on x/prices keeper, used for stateful validation of the returned MarketPriceUpdateTx
	pk PricesKeeper

	// tx decoder used for unmarshalling the market-price-update tx
	txDecoder sdk.TxDecoder
}

// DecodeUpdateMarketPricesTx returns a new `UpdateMarketPricesTx` after validating the following:
//   - decodes the given tx bytes
//   - checks the num of msgs in the tx matches expectations
//   - checks the msg is of expected type
//
// If error occurs during any of the checks, returns error.
func (mpd *DefaultUpdateMarketPriceTxDecoder) DecodeUpdateMarketPricesTx(ctx sdk.Context, txs [][]byte) (*UpdateMarketPricesTx, error) {
	tx, err := mpd.txDecoder(txs[len(txs)+UpdateMarketPricesTxLenOffset])
	if err != nil {
		return nil, errors.GetDecodingError(msgUpdateMarketPricesType, err)
	}

	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, fmt.Errorf("unexpected number of msgs: %v", len(msgs))
	}

	updateMarketPrices, ok := msgs[0].(*pricestypes.MsgUpdateMarketPrices)
	if !ok {
		return nil, fmt.Errorf("unexpected msg type: %v", msgs[0])
	}

	return &UpdateMarketPricesTx{
		ctx:          ctx,
		pricesKeeper: mpd.pk,
		msg:          updateMarketPrices,
	}, nil
}

// GetTxOffset returns the offset that other injected txs should be placed with respect to their normally
// expected indices. No offset is expected for the default implementation.
func (mpd DefaultUpdateMarketPriceTxDecoder) GetTxOffset(sdk.Context) int {
	return 0
}
