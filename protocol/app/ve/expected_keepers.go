package ve

import (
	"math/big"

	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PreparePricesKeeper defines the expected Prices keeper used for `PrepareProposal`.
type PreBlockExecPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdate *pricestypes.MarketPriceUpdate,
	) (isSpotValid bool, isPnlValid bool)

	UpdateSmoothedSpotPrices(
		ctx sdk.Context,
		linearInterpolateFunc func(v0 uint64, v1 uint64, ppm uint32) (uint64, error),
	) error

	GetValidMarketSpotPriceUpdates(ctx sdk.Context) []*pricestypes.MarketSpotPriceUpdate

	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam

	GetMarketParam(ctx sdk.Context, id uint32) (market pricestypes.MarketParam, exists bool)

	GetSmoothedSpotPrice(markedId uint32) (uint64, bool)
}

type ExtendVoteClobKeeper interface {
	GetSingleMarketClobMetadata(ctx sdk.Context, clobPair clobtypes.ClobPair) clobtypes.ClobMetadata
	GetClobPair(ctx sdk.Context, id clobtypes.ClobPairId) (val clobtypes.ClobPair, found bool)
}

type ExtendVotePerpetualsKeeper interface {
	GetPerpetual(
		ctx sdk.Context,
		id uint32,
	) (val perptypes.Perpetual, err error)
}

type ExtendVoteDaemonPriceCache interface {
	GetVEEncodedPrice(price *big.Int) ([]byte, error)
}

type VEPriceApplier interface {
	ApplyPricesFromVE(ctx sdk.Context, req *abci.RequestFinalizeBlock, writeToCache bool) error
}
