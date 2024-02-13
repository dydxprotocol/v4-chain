package vote_extensions

import (
	"context"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	oracleclient "github.com/skip-mev/slinky/service/clients/oracle"
	oracletypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	types2 "github.com/skip-mev/slinky/x/oracle/types"
	"google.golang.org/grpc"

	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
)

type OracleClient struct {
	Slinky                 oracleclient.OracleClient
	PricesKeeper           priceskeeper.Keeper
	PriceFeedServiceClient api.PriceFeedServiceClient
	grpcClient             daemontypes.GrpcClient
	pricesSocket           string
	pricesConn             *grpc.ClientConn
}

func NewOracleClient(slinky oracleclient.OracleClient, pricesKeeper priceskeeper.Keeper, grpcClient daemontypes.GrpcClient, pricesSocket string) *OracleClient {
	return &OracleClient{
		Slinky:       slinky,
		PricesKeeper: pricesKeeper,
		grpcClient:   grpcClient,
		pricesSocket: pricesSocket,
	}
}

func (o *OracleClient) Start(ctx context.Context) error {
	cancelCtx, cf := context.WithTimeout(ctx, time.Second*5)
	defer cf()
	pricesConn, err := o.grpcClient.NewGrpcConnection(cancelCtx, o.pricesSocket)
	if err != nil {
		return err
	}
	o.pricesConn = pricesConn
	o.PriceFeedServiceClient = api.NewPriceFeedServiceClient(o.pricesConn)
	return o.Slinky.Start(ctx)
}

func (o *OracleClient) Stop() error {
	if o.pricesConn != nil {
		_ = o.grpcClient.CloseConnection(o.pricesConn)
	}
	return o.Slinky.Stop()
}

func (o *OracleClient) Prices(ctx context.Context, in *oracletypes.QueryPricesRequest, opts ...grpc.CallOption) (*oracletypes.QueryPricesResponse, error) {
	sdkCtx, ok := ctx.(sdk.Context)
	if !ok {
		return nil, fmt.Errorf("oracle client was passed on non-sdk context object")
	}
	slinkyResponse, err := o.Slinky.Prices(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	// Update the prices keeper w/ the most recent prices for the relevant markets
	var updates []*api.MarketPriceUpdate
	for currencyPairString, priceString := range slinkyResponse.Prices {
		currencyPair, err := types2.CurrencyPairFromString(currencyPairString)
		if err != nil {
			return nil, err
		}
		sdkCtx.Logger().Info("turned price pair to currency pair", "string", currencyPairString, "currency pair", currencyPair.String())
		id, found := o.PricesKeeper.GetIDForCurrencyPair(sdkCtx, currencyPair)
		if !found {
			sdkCtx.Logger().Info("slinky client returned currency pair not found in prices keeper", "currency pair", currencyPairString)
			continue
		}
		price, err := strconv.ParseUint(priceString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("slinky client returned price %s not parsable as uint64", priceString)
		}
		sdkCtx.Logger().Info("parsed update for", "market id", uint32(id), "price", price)
		updates = append(updates, &api.MarketPriceUpdate{
			MarketId: uint32(id),
			ExchangePrices: []*api.ExchangePrice{
				{
					ExchangeId:     "slinky",
					Price:          price,
					LastUpdateTime: &slinkyResponse.Timestamp,
				},
			},
		})
	}
	if o.PriceFeedServiceClient == nil {
		sdkCtx.Logger().Error("nil price feed service client")
	}
	_, err = o.PriceFeedServiceClient.UpdateMarketPrices(ctx, &api.UpdateMarketPricesRequest{MarketPriceUpdates: updates})
	if err != nil {
		sdkCtx.Logger().Error(err.Error())
		return nil, err
	}
	// Get valid price updates
	validUpdates := o.PricesKeeper.GetValidMarketPriceUpdates(sdkCtx)
	if validUpdates == nil {
		sdkCtx.Logger().Info("prices keeper returned no valid market price updates")
		return nil, nil
	}
	sdkCtx.Logger().Info("prices keeper returned valid updates", "length", len(validUpdates.MarketPriceUpdates))
	// Translate price updates into oracle response
	var outputResponse = &oracletypes.QueryPricesResponse{
		Prices:    make(map[string]string),
		Timestamp: slinkyResponse.Timestamp,
	}
	for _, update := range validUpdates.MarketPriceUpdates {
		mappedPair, found := o.PricesKeeper.GetCurrencyPairFromID(sdkCtx, uint64(update.GetMarketId()))
		if found {
			sdkCtx.Logger().Info("added currency pair", "pair", mappedPair.String())
			outputResponse.Prices[mappedPair.String()] = strconv.FormatUint(update.Price, 10)
		} else {
			sdkCtx.Logger().Info("failed to add currency pair", "pair", mappedPair.String())
		}
	}
	return outputResponse, nil
}
