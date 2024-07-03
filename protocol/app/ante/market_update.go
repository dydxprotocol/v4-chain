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

	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

var ErrNoCrossMarketUpdates = errors.New("cannot call MsgUpdateMarkets or MsgUpsertMarkets " +
	"on a market listed as cross margin")

type ValidateMarketUpdateDecorator struct {
	perpKeeper  perpetualstypes.PerpetualsKeeper
	priceKeeper pricestypes.PricesKeeper
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
) ValidateMarketUpdateDecorator {
	return ValidateMarketUpdateDecorator{
		perpKeeper:  perpKeeper,
		priceKeeper: priceKeeper,
		cache:       make(map[string]perpetualstypes.PerpetualMarketType),
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
	var msg = msgs[0]

	switch msg := msg.(type) {
	case *mmtypes.MsgUpdateMarkets:
		if contains := d.doMarketsContainCrossMarket(ctx, msg.UpdateMarkets); contains {
			return ctx, ErrNoCrossMarketUpdates
		}

	case *mmtypes.MsgUpsertMarkets:
		if contains := d.doMarketsContainCrossMarket(ctx, msg.Markets); contains {
			return ctx, ErrNoCrossMarketUpdates
		}
	default:
		return ctx, fmt.Errorf("unrecognized message type: %T", msg)
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
