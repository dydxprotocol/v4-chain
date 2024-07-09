package ve

import (
	"math/big"

	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PreparePricesKeeper defines the expected Prices keeper used for `PrepareProposal`.
type PreparePricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *pricestypes.MarketPriceUpdates,
		performNonDeterministicValidation bool,
	) error

	GetValidMarketPriceUpdates(ctx sdk.Context) *pricestypes.MarketPriceUpdates
	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam
	GetMarketPriceUpdateFromBytes(id uint32, bz []byte) (*pricestypes.MarketPriceUpdates_MarketPriceUpdate, error)
}

type ExtendVotePricesKeeper interface {
	GetValidMarketPriceUpdates(ctx sdk.Context) *pricestypes.MarketPriceUpdates
	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam
	GetMarketPriceUpdateFromBytes(id uint32, bz []byte) (*pricestypes.MarketPriceUpdates_MarketPriceUpdate, error)
	GetMarketParam(ctx sdk.Context, id uint32) (market pricestypes.MarketParam, exists bool)
}

type ExtendVoteIndexPriceCache interface {
	GetVEEncodedPrice(price *big.Int) ([]byte, error)
}
