package vote_extensions

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	oracleclient "github.com/skip-mev/slinky/service/clients/oracle"
	oracletypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	types2 "github.com/skip-mev/slinky/x/oracle/types"
	"google.golang.org/grpc"
	"strconv"
)

type OracleClient struct {
	Slinky       oracleclient.OracleClient
	PricesKeeper priceskeeper.Keeper
	Context      *sdk.Context
}

func (o OracleClient) Start(ctx context.Context) error {
	return o.Slinky.Start(ctx)
}

func (o OracleClient) Stop() error {
	return o.Slinky.Stop()
}

func (o OracleClient) Prices(ctx context.Context, in *oracletypes.QueryPricesRequest, opts ...grpc.CallOption) (*oracletypes.QueryPricesResponse, error) {
	slinkyResponse, err := o.Slinky.Prices(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	// Update the prices keeper w/ the most recent prices for the relevant markets
	var updates []*types.MsgUpdateMarketPrices_MarketPrice
	for currencyPairString, priceString := range slinkyResponse.Prices {
		currencyPair, err := types2.CurrencyPairFromString(currencyPairString)
		if err != nil {
			return nil, err
		}
		o.Context.Logger().Info("turned price pair to currency pair", "string", currencyPairString, "currency pair", currencyPair.String())
		id, found := o.PricesKeeper.GetIDForCurrencyPair(*o.Context, currencyPair)
		if !found {
			o.Context.Logger().Info("slinky client returned currency pair not found in prices keeper", "currency pair", currencyPairString)
			continue
		}
		price, err := strconv.ParseUint(priceString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("slinky client returned price %s not parsable as uint64", priceString)
		}
		o.Context.Logger().Info("parsed update for", "market id", uint32(id), "price", price)
		updates = append(updates, &types.MsgUpdateMarketPrices_MarketPrice{
			MarketId: uint32(id),
			Price:    price,
		})
	}
	err = o.PricesKeeper.UpdateMarketPrices(*o.Context, updates)
	if err != nil {
		o.Context.Logger().Error(err.Error())
		return nil, err
	}
	// Get valid price updates
	validUpdates := o.PricesKeeper.GetValidMarketPriceUpdates(*o.Context)
	if validUpdates == nil {
		o.Context.Logger().Info("prices keeper returned no valid market price updates")
		return nil, nil
	}
	o.Context.Logger().Info("prices keeper returned valid updates", "length", len(validUpdates.MarketPriceUpdates))
	// Translate price updates into oracle response
	var outputResponse = &oracletypes.QueryPricesResponse{
		Prices:    make(map[string]string),
		Timestamp: slinkyResponse.Timestamp,
	}
	for _, update := range validUpdates.MarketPriceUpdates {
		mappedPair, found := o.PricesKeeper.GetCurrencyPairFromID(*o.Context, uint64(update.GetMarketId()))
		if found {
			o.Context.Logger().Info("added currency pair", "pair", mappedPair.String())
			outputResponse.Prices[mappedPair.String()] = strconv.FormatUint(update.Price, 10)
		} else {
			o.Context.Logger().Info("failed to add currency pair", "pair", mappedPair.String())
		}
	}
	return outputResponse, nil
}
