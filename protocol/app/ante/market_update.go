package ante

import (
	"errors"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"

	slinkylibs "github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

var ErrNoCrossMarketUpdates = errors.New("cannot call MsgUpdateMarkets or MsgUpsertMarkets " +
	"on a market listed as cross margin")

type MarketMapKeeper interface {
	GetAllMarkets(ctx sdk.Context) (map[string]mmtypes.Market, error)
}

type ValidateMarketUpdateDecorator struct {
	perpKeeper      perpetualstypes.PerpetualsKeeper
	priceKeeper     pricestypes.PricesKeeper
	marketMapKeeper MarketMapKeeper
	// write only cache for mapping slinky ticker strings to market types
	// only evicted on node restart
	cache map[string]perpetualstypes.PerpetualMarketType
}

// NewValidateMarketUpdateDecorator returns an AnteDecorator that is able to check for x/marketmap update messages
// and reject them if they are updating cross margin markets.
//
// NOTE: this is a stop-gap solution before more general functionality is added to x/marketmap to delay and gate
// certain update operations.
func NewValidateMarketUpdateDecorator(
	perpKeeper perpetualstypes.PerpetualsKeeper,
	priceKeeper pricestypes.PricesKeeper,
	marketMapKeeper MarketMapKeeper,
) ValidateMarketUpdateDecorator {
	return ValidateMarketUpdateDecorator{
		perpKeeper:      perpKeeper,
		priceKeeper:     priceKeeper,
		marketMapKeeper: marketMapKeeper,
		cache:           make(map[string]perpetualstypes.PerpetualMarketType),
	}
}

// AnteHandle performs the following checks:
// - check if tx contains x/marketmap/MsgUpdateMarkets or x/marketmap/MsgUpsertMarkets
// - check if the given Tx has more than one message if it has x/marketmap updates, reject if so
// - check if the x/marketmap update affects markets that are registered as cross margin
// in x/perpetuals, reject if so.
func (d ValidateMarketUpdateDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	// Ensure that if this is a market update message then that there is only one.
	// If it isn't a market update message then pass to the next AnteHandler.
	isSingleMarketUpdate, err := IsMarketUpdateTx(tx)
	if err != nil {
		return ctx, err
	}
	if !isSingleMarketUpdate {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	var (
		msg     = msgs[0]
		markets []mmtypes.Market
	)

	switch msg := msg.(type) {
	case *mmtypes.MsgUpdateMarkets:
		markets = msg.UpdateMarkets
	case *mmtypes.MsgUpsertMarkets:
		markets = msg.Markets
	default:
		return ctx, fmt.Errorf("unrecognized message type: %T", msg)
	}

	if contains := d.doMarketsContainCrossMarket(ctx, markets); contains {
		return ctx, ErrNoCrossMarketUpdates
	}

	// check if the market updates are safe
	if err := d.doMarketsUpdateEnabledValues(ctx, markets); err != nil {
		return ctx, errorsmod.Wrap(err, "market update is not safe")
	}

	return next(ctx, tx, simulate)
}

func (d ValidateMarketUpdateDecorator) doMarketsContainCrossMarket(ctx sdk.Context, markets []mmtypes.Market) bool {
	perps := d.perpKeeper.GetAllPerpetuals(ctx)

	for _, market := range markets {
		ticker := market.Ticker.CurrencyPair.String()

		marketType, found := d.cache[ticker]
		if !found {
			// search for market if we cannot find in cache
			for _, perp := range perps {
				params, exists := d.priceKeeper.GetMarketParam(ctx, perp.Params.MarketId)
				if !exists {
					return false
				}

				if MatchPairToSlinkyTicker(params.Pair, market.Ticker.CurrencyPair) {
					// populate cache
					marketType = perp.Params.MarketType
					d.cache[ticker] = marketType
				}
			}
		}
		if marketType == perpetualstypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS {
			return true
		}
	}

	return false
}

// doMarketsUpdateEnabledValues checks if the given markets updates are safe, specifically:
// 1. If a newly added market (market does not exist in x/prices) is added, it should be disabled in the market-map
// 2. If an existing market is updated, the market-update should not change the enabled value
func (d ValidateMarketUpdateDecorator) doMarketsUpdateEnabledValues(ctx sdk.Context, markets []mmtypes.Market) error {
	// get all market-params
	mps := d.priceKeeper.GetAllMarketParams(ctx)
	mm, err := d.marketMapKeeper.GetAllMarkets(ctx)
	if err != nil {
		return err
	}

	// convert to map for easy lookup
	mpMap, err := marketParamsSliceToMap(mps)
	if err != nil {
		return err
	}

	// check validity of incoming market-updates
	for _, market := range markets {
		_, exists := mpMap[market.Ticker.CurrencyPair.String()]
		if !exists {
			// if market does not exist in x/prices, it should be disabled
			if market.Ticker.Enabled {
				return errors.New("newly added market should be disabled")
			}
		} else {
			// find the market in the market-map
			mmMarket, exists := mm[market.Ticker.CurrencyPair.String()]
			if !exists {
				return errors.New("market does not exist in market-map")
			}

			// if market exists, it should not change the enabled value
			if mmMarket.Ticker.Enabled != market.Ticker.Enabled {
				return fmt.Errorf(
					"market should not change enabled value from %t to %t",
					mmMarket.Ticker.Enabled, market.Ticker.Enabled,
				)
			}
		}
	}

	return nil
}

func marketParamsSliceToMap(mps []pricestypes.MarketParam) (map[string]pricestypes.MarketParam, error) {
	mpMap := make(map[string]pricestypes.MarketParam)

	// create entry for each market-param
	for _, mp := range mps {
		// index will be the slinky-style ticker
		idx, err := slinkylibs.MarketPairToCurrencyPair(mp.Pair)
		if err != nil {
			return nil, err
		}

		// check for duplicate entries
		if _, exists := mpMap[idx.String()]; exists {
			return nil, errors.New("duplicate market-param entry")
		}

		mpMap[idx.String()] = mp
	}

	return mpMap, nil
}

// IsMarketUpdateTx returns `true` if the supplied `tx` consists of a single
// MsgUpdateMarkets or MsgUpsertMarkets
func IsMarketUpdateTx(tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	var hasMessage = false

	for _, msg := range msgs {
		switch msg.(type) {
		case *mmtypes.MsgUpdateMarkets, *mmtypes.MsgUpsertMarkets:
			hasMessage = true
		}

		if hasMessage {
			break
		}
	}

	if !hasMessage {
		return false, nil
	}

	numMsgs := len(msgs)
	if numMsgs > 1 {
		return false, errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"a transaction containing MsgUpdateMarkets or MsgUpsertMarkets may not contain more than one message",
		)
	}

	return true, nil
}

// MatchPairToSlinkyTicker matches a market params string of form "BTC-USD"
// to a slinky currency pair struct consisting of a BASE and QUOTE.
func MatchPairToSlinkyTicker(pair string, ticker slinkytypes.CurrencyPair) bool {
	pairSplit := strings.Split(pair, "-")
	if len(pairSplit) != 2 {
		return false
	}

	if pairSplit[0] == ticker.Base && pairSplit[1] == ticker.Quote {
		return true
	}

	return false
}
