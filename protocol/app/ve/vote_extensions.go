package ve

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	// encoding and decoding vote extensions
	veCodec codec.VoteExtensionCodec

	// fetching valid price updates and current markets
	pricesKeeper ExtendVotePricesKeeper

	// writing prices to the store
	// prices are written to the store here to ensure
	// GetValidMarketPriceUpdates returns the latest
	// accurate prices
	priceApplier VEPriceApplier
}

func NewVoteExtensionHandler(
	logger log.Logger,
	vecodec codec.VoteExtensionCodec,
	pricesKeeper ExtendVotePricesKeeper,
	priceApplier VEPriceApplier,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:       logger,
		veCodec:      vecodec,
		pricesKeeper: pricesKeeper,
		priceApplier: priceApplier,
	}
}

// Returns a handler that extends pre-commit votes with the current
// prices pulled from the perpetually running price daemon
// In the case of an error, the handler will return an empty vote extension
// ensuring liveness in the case of a price daemon failure
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (resp *abci.ResponseExtendVote, err error) {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &abci.ResponseExtendVote{VoteExtension: []byte{}}, ErrPanic{fmt.Errorf("%v", r)}
			}
		}()

		if req == nil {
			err = fmt.Errorf("nil request for extend vote")
			return nil, err
		}

		reqFinalizeBlock := &abci.RequestFinalizeBlock{
			Txs:    req.Txs,
			Height: req.Height,
		}

		// apply prices to ensure GetValidMarketPriceUpdates returns the latest prices
		if _, err = h.priceApplier.ApplyPricesFromVoteExtensions(ctx, reqFinalizeBlock); err != nil {
			h.logger.Error(
				"failed to aggregate oracle votes",
				"height", req.Height,
				"err", err,
			)
			err = PreBlockError{err}

			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		priceUpdates := h.pricesKeeper.GetValidMarketPriceUpdates(ctx)
		if len(priceUpdates.MarketPriceUpdates) == 0 {
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, fmt.Errorf("no valid median prices")
		}

		voteExt, err := h.transformDaemonPricesToVE(ctx, priceUpdates.MarketPriceUpdates)
		if err != nil {
			h.logger.Error("failed to transform prices to vote extension", "height", req.Height, "err", err)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		veBytes, err := h.veCodec.Encode(voteExt)
		if err != nil {
			h.logger.Error("failed to encode vote extension", "height", req.Height, "err", err)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Debug("extending vote with daemon prices", "height", req.Height, "prices", len(priceUpdates.MarketPriceUpdates))

		return &abci.ResponseExtendVote{VoteExtension: veBytes}, nil
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

		if err := h.ValidateDaemonVE(ctx, h.pricesKeeper, ve); err != nil {
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
		if market, ok = h.pricesKeeper.GetMarketParam(ctx, marketId); !ok {
			h.logger.Debug("market id not found", "marketId", marketId)
			continue
		}

		rawPrice := new(big.Int).SetUint64(price)

		encodedPrice, err := veutils.GetVEEncodedPrice(rawPrice)

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
	pricesKeeper ExtendVotePricesKeeper,
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
		// TODO: Should we check we can decode here
		// if _, err := pricesKeeper.GetMarketPriceUpdateFromBytes(id, bz); err != nil {
		// 	return fmt.Errorf("failed to decode price: %v", err)
		// }
	}

	return nil
}

func (h *VoteExtensionHandler) GetMaxPairs(ctx sdk.Context) uint32 {
	markets := h.pricesKeeper.GetAllMarketParams(ctx)
	// TODO: check how to handle this query in prepare / process proposal
	// given that pairs can be created/removed
	return uint32(len(markets))
}
