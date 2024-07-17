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

		// apply prices from prev block to ensure that the prices are up to date
		if _, err = h.priceApplier.ApplyPricesFromVE(ctx, reqFinalizeBlock); err != nil {
			h.logger.Error(
				"failed to aggregate oracle votes",
				"height", req.Height,
				"err", err,
			)
			err = PreBlockError{err}

			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		veBytes, err := h.GetVEBytesFromCurrPrices(ctx)
		if err != nil {
			h.logger.Error(
				"failed to get vote extension bytes from current prices",
				"height", req.Height,
				"err", err,
			)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		h.logger.Debug("extending vote with daemon prices", "height", req.Height)

		return &abci.ResponseExtendVote{VoteExtension: veBytes}, nil
	}
}

func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(
		ctx sdk.Context,
		req *abci.RequestVerifyVoteExtension,
	) (_ *abci.ResponseVerifyVoteExtension, err error) {

		if req == nil {
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

		if err := h.ValidateVEPriceByteSize(ctx, h.pricesKeeper, req.VoteExtension); err != nil {
			h.logger.Error(
				"failed to decode and validate vote extension",
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

func (h *VoteExtensionHandler) GetVEBytesFromCurrPrices(ctx sdk.Context) ([]byte, error) {
	priceUpdates := h.pricesKeeper.GetValidMarketPriceUpdates(ctx)

	if len(priceUpdates.MarketPriceUpdates) == 0 {
		return nil, fmt.Errorf("no valid median prices")
	}

	voteExt, err := h.transformDaemonPricesToVE(priceUpdates.MarketPriceUpdates)
	if err != nil {
		return nil, err
	}

	veBytes, err := h.veCodec.Encode(voteExt)
	if err != nil {
		return nil, err
	}

	return veBytes, nil
}

func (h *VoteExtensionHandler) transformDaemonPricesToVE(
	priceupdates []*pricetypes.MarketPriceUpdates_MarketPriceUpdate,
) (types.DaemonVoteExtension, error) {

	vePrices := make(map[uint32][]byte)

	for _, priceUpdate := range priceupdates {
		// check if the marketId is valid
		encodedPrice, err := h.GetEncodedPriceFromPriceUpdate(priceUpdate)

		if err != nil {
			h.logger.Debug(
				"failed to encode price",
				"price", priceUpdate.GetPrice(),
				"market", priceUpdate.GetMarketId(),
				"err", err,
			)
			continue
		}

		vePrices[priceUpdate.GetMarketId()] = encodedPrice
	}

	return types.DaemonVoteExtension{
		Prices: vePrices,
	}, nil
}

func (h *VoteExtensionHandler) GetEncodedPriceFromPriceUpdate(
	priceUpdate *pricetypes.MarketPriceUpdates_MarketPriceUpdate,
) ([]byte, error) {

	price := new(big.Int).SetUint64(priceUpdate.GetPrice())

	encodedPrice, err := veutils.GetVEEncodedPrice(price)
	if err != nil {
		return nil, err
	}

	return encodedPrice, nil

}

func (h *VoteExtensionHandler) ValidateVEPriceByteSize(
	ctx sdk.Context,
	pricesKeeper ExtendVotePricesKeeper,
	veBytes []byte,
) error {

	ve, err := h.veCodec.Decode(veBytes)
	if err != nil {
		return fmt.Errorf("failed to decode vote extension: %v", err)
	}

	maxPairs := h.GetMaxMarketPairs(ctx)
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

func (h *VoteExtensionHandler) GetMaxMarketPairs(ctx sdk.Context) uint32 {
	markets := h.pricesKeeper.GetAllMarketParams(ctx)
	return uint32(len(markets))
}
