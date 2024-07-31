package ve

import (
	"fmt"
	"math/big"
	"sort"

	"cosmossdk.io/log"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	// encoding and decoding vote extensions
	voteCodec codec.VoteExtensionCodec

	// fetching valid price updates and current markets
	pricesKeeper PreBlockExecPricesKeeper

	// fetching last funding rates for price calc
	perpetualsKeeper ExtendVotePerpetualsKeeper

	// fetching mid price for price calc
	clobKeeper ExtendVoteClobKeeper

	// writing prices to the prices module store
	priceApplier VEPriceApplier
}

type VEPricePair struct {
	SpotPrice uint64
	PnlPrice  uint64
}

var (
	acceptResponse = &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}
	rejectResponse = &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT}
)

func NewVoteExtensionHandler(
	logger log.Logger,
	voteCodec codec.VoteExtensionCodec,
	pricesKeeper PreBlockExecPricesKeeper,
	perpetualsKeeper ExtendVotePerpetualsKeeper,
	clobKeeper ExtendVoteClobKeeper,
	priceApplier VEPriceApplier,
) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:           logger,
		voteCodec:        voteCodec,
		pricesKeeper:     pricesKeeper,
		perpetualsKeeper: perpetualsKeeper,
		clobKeeper:       clobKeeper,
		priceApplier:     priceApplier,
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
			DecidedLastCommit: abci.CommitInfo{
				Round: request.ProposedLastCommit.Round,
				Votes: []abci.VoteInfo{},
			},
		}

		// apply prices from prev block to ensure that the prices are up to date
		if err := h.priceApplier.ApplyPricesFromVE(ctx, reqFinalizeBlock); err != nil {
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
	) (resp *abci.ResponseVerifyVoteExtension, err error) {

		defer func() {
			if recovery := recover(); recovery != nil {
				h.logger.Error(
					"recovered from panic in VerifyVoteExtensionHandler",
					"err", recovery,
				)
				resp = rejectResponse
				err = ErrPanic{fmt.Errorf("%v", recovery)}

			}
		}()

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
			h.pricesKeeper,
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
	priceUpdates := h.getCurrentPrices(ctx)

	if len(priceUpdates) == 0 {
		return nil, fmt.Errorf("no valid prices")
	}

	// turn prices from daemon into a VE
	voteExt, err := h.transformDaemonPricesToVE(priceUpdates)
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
	priceupdates map[uint32]VEPricePair,
) (types.DaemonVoteExtension, error) {
	vePrices := make(map[uint32]*types.DaemonVoteExtension_PricePair)

	for marketId, priceUpdate := range priceupdates {
		// check if the marketId is valid
		encodedPricePair, err := h.GetEncodedPriceFromPriceUpdate(priceUpdate)
		if err != nil {
			continue
		}
		vePrices[marketId] = &encodedPricePair
	}

	return types.DaemonVoteExtension{
		Prices: vePrices,
	}, nil
}

func (h *VoteExtensionHandler) GetEncodedPriceFromPriceUpdate(
	priceUpdate VEPricePair,
) (types.DaemonVoteExtension_PricePair, error) {
	spotPrice := new(big.Int).SetUint64(priceUpdate.SpotPrice)
	pnlPrice := new(big.Int).SetUint64(priceUpdate.PnlPrice)

	encodedSpotPrice, err := veutils.GetVEEncodedPrice(spotPrice)
	if err != nil {
		return types.DaemonVoteExtension_PricePair{}, err
	}

	encodedPnlPrice, err := veutils.GetVEEncodedPrice(pnlPrice)
	if err != nil {
		// TODO: If PNL price can't be encoded should we default to spot price for pnl price
		// or do we allow for a price to be nil in a PricePair object
		return types.DaemonVoteExtension_PricePair{}, err
	}

	return types.DaemonVoteExtension_PricePair{
		SpotPrice: encodedSpotPrice,
		PnlPrice:  encodedPnlPrice,
	}, nil
}

func (h *VoteExtensionHandler) getCurrentPrices(
	ctx sdk.Context,
) map[uint32]VEPricePair {
	vePrices := make(map[uint32]VEPricePair)

	indexPrices := h.pricesKeeper.GetValidMarketPriceUpdates(ctx)

	for _, market := range indexPrices.MarketPriceUpdates {

		clobMidPrice := h.getClobMidPrice(ctx, market.MarketId)
		if clobMidPrice == nil {
			vePrices[market.MarketId] = VEPricePair{
				SpotPrice: market.Price,
				PnlPrice:  market.Price,
			}
			continue
		}
		smoothedPrice := h.getSmoothedPrice(market.MarketId)
		if smoothedPrice == nil {
			vePrices[market.MarketId] = VEPricePair{
				SpotPrice: market.Price,
				PnlPrice:  market.Price,
			}
			continue
		}
		lastFundingRate := h.getLastFundingRate(ctx, market.MarketId)
		if lastFundingRate == nil {
			vePrices[market.MarketId] = VEPricePair{
				SpotPrice: market.Price,
				PnlPrice:  market.Price,
			}
			continue
		}

		medianPrice := h.getMedianPrice(
			new(big.Int).SetUint64(market.Price),
			clobMidPrice,
			smoothedPrice,
			lastFundingRate,
		)

		vePrices[market.MarketId] = VEPricePair{
			SpotPrice: market.Price,
			PnlPrice:  medianPrice.Uint64(),
		}
	}

	return vePrices
}

func (h *VoteExtensionHandler) getMedianPrice(
	indexPrice *big.Int,
	clobMidPrice *big.Int,
	smoothedPrice *big.Int,
	lastFundingRate *big.Int,
) *big.Int {

	fundingWeightedPrice := h.getFundingWeightedIndexPrice(indexPrice, lastFundingRate)
	prices := []*big.Int{clobMidPrice, smoothedPrice, fundingWeightedPrice}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Cmp(prices[j]) < 0
	})

	return prices[1]
}

func (h *VoteExtensionHandler) getFundingWeightedIndexPrice(
	indexPrice *big.Int,
	lastFundingRate *big.Int,
) *big.Int {
	ppmFactor := new(big.Int).SetInt64(1000000)
	adjustedFundingRate := new(big.Int).Add(lastFundingRate, ppmFactor)
	fundingWeightedPrice := new(big.Int).Mul(indexPrice, adjustedFundingRate)
	fundingWeightedPrice = fundingWeightedPrice.Div(fundingWeightedPrice, ppmFactor)
	return fundingWeightedPrice
}

func (h *VoteExtensionHandler) getClobMidPrice(
	ctx sdk.Context,
	marketId uint32,
) *big.Int {
	clobPair, found := h.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(marketId))

	if !found {
		return nil
	}

	clobMetadata := h.clobKeeper.GetSingleMarketClobMetadata(ctx, clobPair)

	if clobMetadata.MidPrice == 0 {
		return nil
	}

	midPrice := clobMetadata.MidPrice.ToBigInt()
	subticksPerTick := new(big.Int).SetUint64(uint64(clobPair.SubticksPerTick))
	return new(big.Int).Div(midPrice, subticksPerTick)
}

func (h *VoteExtensionHandler) getSmoothedPrice(
	marketId uint32,
) *big.Int {
	smoothedPrice, exists := h.pricesKeeper.GetSmoothedPrice(marketId)
	if !exists || smoothedPrice == 0 {
		return nil
	}

	return new(big.Int).SetUint64(smoothedPrice)
}

func (h *VoteExtensionHandler) getLastFundingRate(
	ctx sdk.Context,
	marketId uint32,
) *big.Int {
	perpetual, err := h.perpetualsKeeper.GetPerpetual(ctx, marketId)
	if err != nil {
		return nil
	}

	return perpetual.LastFundingRate.BigInt()
}
