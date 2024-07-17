package ve

import (
	"math/big"

	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
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
}

type ExtendVotePricesKeeper interface {
	GetValidMarketPriceUpdates(ctx sdk.Context) *pricestypes.MarketPriceUpdates
	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam
	GetMarketParam(ctx sdk.Context, id uint32) (market pricestypes.MarketParam, exists bool)
}

type ExtendVoteIndexPriceCache interface {
	GetVEEncodedPrice(price *big.Int) ([]byte, error)
}

type VEPriceApplier interface {
	ApplyPricesFromVoteExtensions(ctx sdk.Context, req *abci.RequestFinalizeBlock) (map[string]*big.Int, error)
}
