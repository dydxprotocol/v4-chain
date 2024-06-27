package ve

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/log"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	libtime "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/time"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	indexPriceCache *pricefeedtypes.MarketToExchangePrices

	veCodec codec.VoteExtensionCodec

	timeout time.Duration

	timeProvider libtime.TimeProvider

	pk pk.Keeper
}

func NewVoteExtensionHandler(
	logger log.Logger,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	timeout time.Duration,
	vecodec codec.VoteExtensionCodec,
	timeProvider libtime.TimeProvider,
	pricekeeper pk.Keeper,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:          logger,
		timeout:         timeout,
		indexPriceCache: indexPriceCache,
		veCodec:         vecodec,
		timeProvider:    timeProvider,
		pk:              pricekeeper,
	}
}

// Returns a handler that extends pre-commit votes with the current
// prices pulled from the perpetually running price deamon
// In the case of an error, the handler will return an empty vote extension
// ensuring liveness in the case of a price deamon failure
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
			// TODO: return dynamic structured error with name of cometBFT request
			err = fmt.Errorf("nil request for extend vote")
			return nil, err
		}

		// TODO: call the method used to write prices to state, in the same way preBlocker does

		marketParams := h.pk.GetAllMarketParams(ctx)
		currPrices := h.indexPriceCache.GetValidMedianPrices(marketParams, h.timeProvider.Now())

		if len(currPrices) == 0 {
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, fmt.Errorf("no valid median prices")
		}

		voteExt, err := h.transformDeamonPricesToVE(ctx, currPrices)
		if err != nil {
			h.logger.Error("failed to transform prices to vote extension", "height", req.Height, "err", err)
			// TODO: structure error
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		bz, err := h.veCodec.Encode(voteExt)
		if err != nil {
			h.logger.Error("failed to encode vote extension", "height", req.Height, "err", err)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Debug("extending vote with deamon prices", "height", req.Height, "prices", len(currPrices))

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
			// TODO: return dynamic structured error with name of cometBFT request
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

		if err := h.ValidateDeamonVE(ctx, ve); err != nil {
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

// encode the prices from the deamon into VE data using GobEncode
func (h *VoteExtensionHandler) transformDeamonPricesToVE(
	ctx sdk.Context,
	prices map[uint32]uint64,
) (types.DeamonVoteExtension, error) {

	vePrices := make(map[uint32][]byte)
	for marketId, price := range prices {
		// check if the marketId is valid
		var market pricetypes.MarketParam
		var ok bool

		// Check if the marketId is valid
		if market, ok = h.pk.GetMarketParam(ctx, marketId); !ok {
			h.logger.Debug("market id not found", "marketId", marketId)
			continue
		}
		priceString := fmt.Sprintf("%d", price)
		rawPrice, converted := new(big.Int).SetString(priceString, 10)

		if !converted {
			// TODO: check if we just ignore price and continue or return error
			return types.DeamonVoteExtension{}, fmt.Errorf("failed to convert price string to big.Int: %s", priceString)
		}

		encodedPrice, err := h.indexPriceCache.GetEncodedPrice(rawPrice)

		if err != nil {
			h.logger.Debug("failed to encode price", "price", price, "market", market.Pair, "err", err)
			continue
		}

		h.logger.Info("transformed deamon price", "market", market.Pair, "price", price)

		vePrices[marketId] = encodedPrice
	}
	h.logger.Info("transformed deamon prices", "prices", len(vePrices))
	return types.DeamonVoteExtension{
		Prices: vePrices,
	}, nil
}

func (h *VoteExtensionHandler) ValidateDeamonVE(
	ctx sdk.Context,
	ve types.DeamonVoteExtension,
) error {
	maxPairs := h.GetMaxPairs(ctx)
	if uint32(len(ve.Prices)) > maxPairs {
		return fmt.Errorf("too many prices in deamon vote extension: %d > %d", len(ve.Prices), maxPairs)
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
