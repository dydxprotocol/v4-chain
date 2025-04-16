package ante

import (
	"errors"
	"fmt"

	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	mmtypes "github.com/dydxprotocol/slinky/x/marketmap/types"

	slinkylibs "github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

var ErrRestrictedMarketUpdates = errors.New("cannot call MsgUpdateMarkets or MsgUpsertMarkets " +
	"on a restricted market")

type MarketMapKeeper interface {
	GetAllMarkets(ctx sdk.Context) (map[string]mmtypes.Market, error)
}

var (
	cpUSDTUSD = slinkytypes.CurrencyPair{
		Base:  "USDT",
		Quote: "USD",
	}
)

type ValidateMarketUpdateDecorator struct {
	perpKeeper      perpetualstypes.PerpetualsKeeper
	priceKeeper     pricestypes.PricesKeeper
	marketMapKeeper MarketMapKeeper
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

	if contains, ticker := d.doMarketsContainRestrictedMarket(ctx, markets); contains {
		return ctx, fmt.Errorf("%w: %s", ErrRestrictedMarketUpdates, ticker)
	}

	// check if the market updates are safe
	if err := d.doMarketsUpdateEnabledValues(ctx, markets); err != nil {
		return ctx, pricestypes.ErrUnsafeMarketUpdate.Wrap(err.Error())
	}

	return next(ctx, tx, simulate)
}

// doMarketsContainRestrictedMarket checks if any of the given markets are restricted:
// 1. markets listed as CROSS perpetuals are restricted
// 2. the USDT/USD market is always restricted
func (d ValidateMarketUpdateDecorator) doMarketsContainRestrictedMarket(
	ctx sdk.Context,
	markets []mmtypes.Market,
) (bool, string) {
	// Grab all the perpetuals markets
	perps := d.perpKeeper.GetAllPerpetuals(ctx)
	restrictedMap := make(map[string]bool, len(perps))

	// Attempt to fetch the corresponding Prices market and map it to a currency pair
	for _, perp := range perps {
		params, exists := d.priceKeeper.GetMarketParam(ctx, perp.Params.MarketId)
		if !exists {
			continue
		}
		cp, err := slinkylibs.MarketPairToCurrencyPair(params.Pair)
		if err != nil {
			continue
		}
		restrictedMap[cp.String()] = perp.Params.MarketType == perpetualstypes.
			PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS
	}

	// add usdt/usd market to be restricted
	restrictedMap[cpUSDTUSD.String()] = true

	// Look in the mapped currency pairs to see if we have invalid updates
	for _, market := range markets {
		ticker := market.Ticker.CurrencyPair.String()

		restricted, found := restrictedMap[ticker]
		if found && restricted {
			return true, ticker
		}
	}

	return false, ""
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
				return pricestypes.ErrAdditionOfEnabledMarket
			}
		} else {
			// find the market in the market-map
			mmMarket, exists := mm[market.Ticker.CurrencyPair.String()]
			if !exists {
				return pricestypes.ErrMarketDoesNotExistInMarketMap
			}

			// if market exists, it should not change the enabled value
			if mmMarket.Ticker.Enabled != market.Ticker.Enabled {
				return pricestypes.ErrMarketUpdateChangesMarketMapEnabledValue.Wrapf(
					"market-map market: %t, incoming market update: %t", mmMarket.Ticker.Enabled, market.Ticker.Enabled,
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
