package prices

import (
	"fmt"
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

// MarketPriceUpdateTx is the default implementation of the MarketPriceUpdateTx interface.
// It's Validate() method is responsible for validating the underlying msg in accordance with the dydx process-proposal
// logic pre vote-extensions
type UpdateMarketPricesTx struct {
	ctx          sdk.Context
	pricesKeeper PricesKeeper
	msg          *pricestypes.MsgUpdateMarketPrices
}

// Validate returns an error if:
// - the underlying msg fails `ValidateBasic`
// - the underlying msg values are not "valid" according to the index price.
func (umpt *UpdateMarketPricesTx) Validate() error {
	if err := umpt.msg.ValidateBasic(); err != nil {
		return fmt.Errorf("error validating msg: %v: %v", umpt.msg, err)
	}

	if err := umpt.pricesKeeper.PerformStatefulPriceUpdateValidation(umpt.ctx, umpt.msg, true); err != nil {
		return err
	}

	return nil
}

// GetMsg retrieves the MarketPriceUpdate msg from this tx
func (umpt *UpdateMarketPricesTx) GetMsg() sdk.Msg {
	return umpt.msg
}