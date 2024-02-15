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
	oracleservicetypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"google.golang.org/grpc"

	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
)

// OracleClient is a wrapper around the default Slinky OracleClient interface. This object is responsible for requesting
// prices from the sidecar, sending them to the x/prices module (via the price-feed service), and relaying what
// prices (according to the x/prices module) a node should inject into their vote-extension.
type OracleClient struct {
	Slinky                 oracleclient.OracleClient
	PricesKeeper           priceskeeper.Keeper
	PriceFeedServiceClient api.PriceFeedServiceClient
	grpcClient             daemontypes.GrpcClient
	pricesSocket           string
	pricesConn             *grpc.ClientConn
}

// NewOracleClient returns a new OracleClient object.
func NewOracleClient(slinky oracleclient.OracleClient, pricesKeeper priceskeeper.Keeper, grpcClient daemontypes.GrpcClient, pricesSocket string) *OracleClient {
	return &OracleClient{
		Slinky:       slinky,
		PricesKeeper: pricesKeeper,
		grpcClient:   grpcClient,
		pricesSocket: pricesSocket,
	}
}

// Start starts the OracleClient. This method is responsible for establishing a connection to the price-feed service
// and the sidecar, and then starting the Slinky OracleClient. This method will timeout after 5 seconds if no
// connection is established.
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

// Stop stops the OracleClient. This method is responsible for closing the connection to the sidecar and to
// the price-feed service.
func (o *OracleClient) Stop() error {
	if o.pricesConn != nil {
		_ = o.grpcClient.CloseConnection(o.pricesConn)
	}
	return o.Slinky.Stop()
}

// Prices is a wrapper around the Slinky OracleClient's Prices method. This method is responsible for doing the following:
//  1. Request the latest prices from the oracle-sidecar
//  2. Relay the latest prices to the price-feed service, which will then update the x/prices module's indexPriceCache
//  3. Get the latest prices from the x/prices module's indexPriceCache via GetValidMarketPriceUpdates
//  4. Translate the response from x/prices into a QueryPricesResponse, and return this
//
// This method fails if:
//   - The sidecar returns an error
//   - The sidecar returns a price that cannot be parsed as a uint64
//   - The price-feed service client returns an error
func (o *OracleClient) Prices(ctx context.Context, in *oracleservicetypes.QueryPricesRequest, opts ...grpc.CallOption) (*oracleservicetypes.QueryPricesResponse, error) {
	sdkCtx, ok := ctx.(sdk.Context)
	if !ok {
		return nil, fmt.Errorf("oracle client was passed on non-sdk context object")
	}

	// get prices from slinky sidecar via GRPC
	slinkyResponse, err := o.Slinky.Prices(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	// update the prices keeper w/ the most recent prices for the relevant markets
	var updates []*api.MarketPriceUpdate
	for currencyPairString, priceString := range slinkyResponse.Prices {
		// convert currency-pair string (index) into currency-pair object
		currencyPair, err := oracletypes.CurrencyPairFromString(currencyPairString)
		if err != nil {
			return nil, err
		}
		sdkCtx.Logger().Info("turned price pair to currency pair", "string", currencyPairString, "currency pair", currencyPair.String())

		// get the market id for the currency pair
		id, found := o.PricesKeeper.GetIDForCurrencyPair(sdkCtx, currencyPair)
		if !found {
			sdkCtx.Logger().Info("slinky client returned currency pair not found in prices keeper", "currency pair", currencyPairString)
			continue
		}

		// parse the price string into a uint64
		price, err := strconv.ParseUint(priceString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("slinky client returned price %s not parsable as uint64", priceString)
		}
		sdkCtx.Logger().Info("parsed update for", "market id", uint32(id), "price", price)

		// append the update to the list of MarketPriceUpdates to be sent to the app's price-feed service
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

	// send the updates to the app's price-feed service -> these will then be piped to the x/prices indexPriceCache via the price
	// feed service
	if o.PriceFeedServiceClient == nil {
		sdkCtx.Logger().Error("nil price feed service client")
	}
	_, err = o.PriceFeedServiceClient.UpdateMarketPrices(ctx, &api.UpdateMarketPricesRequest{MarketPriceUpdates: updates})
	if err != nil {
		sdkCtx.Logger().Error(err.Error())
		return nil, err
	}

	// get the final prices to include in the vote-extension from the x/prices module
	validUpdates := o.PricesKeeper.GetValidMarketPriceUpdates(sdkCtx)
	if validUpdates == nil {
		sdkCtx.Logger().Info("prices keeper returned no valid market price updates")
		return nil, nil
	}
	sdkCtx.Logger().Info("prices keeper returned valid updates", "length", len(validUpdates.MarketPriceUpdates))

	// translate price updates into oracle response
	var outputResponse = &oracleservicetypes.QueryPricesResponse{
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
