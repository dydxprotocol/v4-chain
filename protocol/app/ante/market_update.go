package ante

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var ErrnoCrossMarketUpdates = errors.New("cannot call MsgUpdateMarkets or MsgUpsertMarkets on a market listed as cross margin")

type ValidateMarketUpdateDecorator struct {
	pk    perpetualstypes.PerpetualsKeeper
	cache map[string]perpetualstypes.PerpetualMarketType
}

func NewValidateMarketUpdateDecorator(pk perpetualstypes.PerpetualsKeeper) ValidateMarketUpdateDecorator {
	return ValidateMarketUpdateDecorator{
		pk:    pk,
		cache: make(map[string]perpetualstypes.PerpetualMarketType),
	}
}

func (d *ValidateMarketUpdateDecorator) AnteHandle(
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
			return ctx, ErrnoCrossMarketUpdates
		}

	case *mmtypes.MsgUpsertMarkets:
		if contains := d.doMarketsContainCrossMarket(ctx, msg.Markets); contains {
			return ctx, ErrnoCrossMarketUpdates
		}
	default:
		return ctx, fmt.Errorf("unrecognized message type: %T", msg)
	}

	return next(ctx, tx, simulate)
}

func (d *ValidateMarketUpdateDecorator) doMarketsContainCrossMarket(ctx sdk.Context, markets []mmtypes.Market) bool {
	perps := d.pk.GetAllPerpetuals(ctx)

	for _, market := range markets {
		var (
			ticker     = market.Ticker.CurrencyPair.String()
			marketType perpetualstypes.PerpetualMarketType
			found      bool
		)

		if marketType, found = d.cache[ticker]; !found {
			// search for market if we cannot find
			for _, perp := range perps {
				if perp.Params.Ticker == ticker {
					// populate cache
					marketType := perp.Params.MarketType
					d.cache[ticker] = marketType
					break
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
