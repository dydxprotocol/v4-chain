package ve

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	// used to encode prices before they are put into a VE
	indexPriceCache *pricefeedtypes.MarketToExchangePrices

	veCodec codec.VoteExtensionCodec

	pk pk.Keeper
}

func NewVoteExtensionHandler(
	logger log.Logger,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	vecodec codec.VoteExtensionCodec,
	pricekeeper pk.Keeper,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:          logger,
		indexPriceCache: indexPriceCache,
		veCodec:         vecodec,
		pk:              pricekeeper,
	}
}

// Returns a handler that extends pre-commit votes with the current
// prices pulled from the perpetually running price daemon
// In the case of an error, the handler will return an empty vote extension
// ensuring liveness in the case of a price daemon failure
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (resp *abci.ResponseExtendVote, err error) {
		defer func() {
			// catch panics if possible
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &abci.ResponseExtendVote{VoteExtension: []byte{}}, ErrPanic{fmt.Errorf("%v", r)}
			}
		}()

		if req == nil {
			ctx.Logger().Error("extend vote handler received a nil request")
			err = fmt.Errorf("nil request for extend vote")
			return nil, err
		}

		// TODO: call the method used to write prices to state, in the same way preBlocker does

		// TODO: does the daemon needs some time to warm up or can we include this in the first block
		pu := h.pk.GetValidMarketPriceUpdates(ctx)

		if len(pu.MarketPriceUpdates) == 0 {
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, fmt.Errorf("no valid median prices")
		}

		voteExt, err := h.transformDaemonPricesToVE(ctx, pu.MarketPriceUpdates)
		if err != nil {
			h.logger.Error("failed to transform prices to vote extension", "height", req.Height, "err", err)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		bz, err := h.veCodec.Encode(voteExt)
		if err != nil {
			h.logger.Error("failed to encode vote extension", "height", req.Height, "err", err)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Debug("extending vote with daemon prices", "height", req.Height, "prices", len(pu.MarketPriceUpdates))

		return &abci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(
		ctx sdk.Context,
		req *abci.RequestVerifyVoteExtension,
	) (_ *abci.ResponseVerifyVoteExtension, err error) {

		if req == nil {
			ctx.Logger().Error("extend vote handler received a nil request")
			err = fmt.Errorf("nil request for verify vote")
			return nil, err
		}

		// accept if vote extension is empty
		if len(req.VoteExtension) == 0 {
			h.logger.Info(
				"empty vote extension",
				"height", req.Height,
			)

			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		ve, err := h.veCodec.Decode(req.VoteExtension)
		if err != nil {
			h.logger.Error(
				"failed to decode vote extension",
				"height", req.Height,
				"err", err,
			)
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT}, err
		}

		if err := h.ValidateDaemonVE(ctx, ve); err != nil {
			h.logger.Error(
				"failed to validate vote extension",
				"height", req.Height,
				"err", err,
			)
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT}, err
		}

		h.logger.Debug(
			"validated vote extension",
			"height", req.Height,
			"size (bytes)", len(req.VoteExtension),
		)

		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}

// encode the prices from the daemon into VE data using GobEncode
func (h *VoteExtensionHandler) transformDaemonPricesToVE(
	ctx sdk.Context,
	priceupdates []*pricetypes.MarketPriceUpdates_MarketPriceUpdate,
) (types.DaemonVoteExtension, error) {

	vePrices := make(map[uint32][]byte)
	for _, pu := range priceupdates {
		// check if the marketId is valid
		var market pricetypes.MarketParam
		var ok bool

		marketId := pu.GetMarketId()
		price := pu.GetPrice()

		// Check if the marketId is valid
		// TODO: check if this is necessary given that we call GetValidMarketPriceUpdates in the ExtendVoteHandler
		if market, ok = h.pk.GetMarketParam(ctx, marketId); !ok {
			h.logger.Debug("market id not found", "marketId", marketId)
			continue
		}

		rawPrice := new(big.Int).SetUint64(price)

		encodedPrice, err := h.indexPriceCache.GetVEEncodedPrice(rawPrice)

		if err != nil {
			h.logger.Debug("failed to encode price", "price", price, "market", market.Pair, "err", err)
			continue
		}

		h.logger.Info("transformed daemon price", "market", market.Pair, "price", price)

		vePrices[marketId] = encodedPrice
	}
	h.logger.Info("transformed daemon prices", "prices", len(vePrices))
	return types.DaemonVoteExtension{
		Prices: vePrices,
	}, nil
}

func (h *VoteExtensionHandler) ValidateDaemonVE(
	ctx sdk.Context,
	ve types.DaemonVoteExtension,
) error {
	maxPairs := h.GetMaxPairs(ctx)
	if uint32(len(ve.Prices)) > maxPairs {
		return fmt.Errorf("too many prices in daemon vote extension: %d > %d", len(ve.Prices), maxPairs)
	}

	for _, bz := range ve.Prices {
		if len(bz) > constants.MaximumPriceSize {
			return fmt.Errorf("price bytes are too long: %d", len(bz))
		}
	}

	return nil
}

func (h *VoteExtensionHandler) GetMaxPairs(ctx sdk.Context) uint32 {
	markets := h.pk.GetAllMarketParams(ctx)
	// TODO: check how to handle this query in prepare / process proposal
	// given that pairs can be created/removed
	return uint32(len(markets))
}
