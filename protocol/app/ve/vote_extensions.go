package ve

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	priceskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	// encoding and decoding vote extensions
	voteCodec codec.VoteExtensionCodec

	// fetching valid price updates and current markets
	pricesKeeper ExtendVotePricesKeeper

	// writing prices to the prices module store
	priceApplier VEPriceApplier
}

var (
	acceptResponse = &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}
	rejectResponse = &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT}
)

func NewVoteExtensionHandler(
	logger log.Logger,
	voteCodec codec.VoteExtensionCodec,
	pricesKeeper ExtendVotePricesKeeper,
	priceApplier VEPriceApplier,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:       logger,
		voteCodec:    voteCodec,
		pricesKeeper: pricesKeeper,
		priceApplier: priceApplier,
	}
}

// Returns a handler that extends pre-commit votes with the current
// prices pulled from the perpetually running price daemon
// In the case of an error, the handler will return an empty vote extension
// ensuring liveness in the case of a price daemon failure
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, request *abci.RequestExtendVote) (resp *abci.ResponseExtendVote, err error) {
		defer func() {
			if recovery := recover(); recovery != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", recovery,
				)
				resp = &abci.ResponseExtendVote{VoteExtension: []byte{}}
				err = ErrPanic{fmt.Errorf("%v", recovery)}
			}
		}()

		if request == nil {
			err = fmt.Errorf("nil request for extend vote")
			return nil, err
		}

		reqFinalizeBlock := &abci.RequestFinalizeBlock{
			Txs:    request.Txs,
			Height: request.Height,
		}

		// apply prices from prev block to ensure that the prices are up to date
		if _, err = h.priceApplier.ApplyPricesFromVE(ctx, reqFinalizeBlock); err != nil {
			h.logger.Error(
				"failed to aggregate oracle votes",
				"height", request.Height,
				"err", err,
			)
			err = PreBlockError{err}

			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

		veBytes, err := h.GetVEBytesFromCurrPrices(ctx)
		if err != nil {
			h.logger.Error(
				"failed to get vote extension bytes from current prices",
				"height", request.Height,
				"err", err,
			)
			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, err
		}

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

		if len(req.VoteExtension) == 0 {
			h.logger.Info(
				"empty vote extension",
				"height", req.Height,
			)

			return acceptResponse, nil
		}

		if err := ValidateVEMarketsAndPrices(
			ctx,
			h.pricesKeeper.(priceskeeper.Keeper),
			req.VoteExtension,
			h.voteCodec,
		); err != nil {
			h.logger.Error(
				"failed to decode and validate vote extension",
				"height", req.Height,
				"err", err,
			)
			return rejectResponse, err
		}

		return acceptResponse, nil
	}
}

func (h *VoteExtensionHandler) GetVEBytesFromCurrPrices(ctx sdk.Context) ([]byte, error) {
	priceUpdates := h.pricesKeeper.GetValidMarketPriceUpdates(ctx)

	if len(priceUpdates.MarketPriceUpdates) == 0 {
		return nil, fmt.Errorf("no valid median prices")
	}

	// turn prices from daemon into a VE
	voteExt, err := h.transformDaemonPricesToVE(priceUpdates.MarketPriceUpdates)
	if err != nil {
		return nil, err
	}

	veBytes, err := h.voteCodec.Encode(voteExt)
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
			continue
		}
		marketId := priceUpdate.GetMarketId()
		vePrices[marketId] = encodedPrice
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
