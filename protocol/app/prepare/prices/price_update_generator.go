package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PriceUpdateGenerator is an interface to abstract the logic of retrieving a
// `MsgUpdateMarketPrices` for any block.
type PriceUpdateGenerator interface {
	GetValidMarketPriceUpdates(ctx sdk.Context, extCommitBz []byte) (*pricestypes.MsgUpdateMarketPrices, error)
}
