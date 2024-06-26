package ve

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/log"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	libtime "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/time"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	indexPriceCache *pricefeedtypes.MarketToExchangePrices

	veCodec codec.VoteExtensionCodec

	timeout time.Duration

	timeProvider libtime.TimeProvider

	pk *pk.Keeper
}

func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (resp *cometabci.ResponseExtendVote, err error) {
		defer func() {
			// catch panics if possible
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, ErrPanic{fmt.Errorf("%v", r)}
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
			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, fmt.Errorf("no valid median prices")
		}

		voteExt, err := h.transformDeamonPricesToVE(ctx, currPrices)
		if err != nil {
			h.logger.Error("failed to transform prices to vote extension", "height", req.Height, "err", err)
			// TODO: structure error
			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		bz, err := h.veCodec.Encode(voteExt)
		if err != nil {
			h.logger.Error("failed to encode vote extension", "height", req.Height, "err", err)
			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Debug("extending vote with deamon prices", "height", req.Height, "prices", len(currPrices))

		return &cometabci.ResponseExtendVote{VoteExtension: bz}, nil
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
