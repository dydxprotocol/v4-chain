package client

import (
	"context"
	"fmt"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
	"strconv"

	"cosmossdk.io/log"
	"google.golang.org/grpc"

	oracleclient "github.com/skip-mev/slinky/service/clients/oracle"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
)

// PriceFetcher is responsible for pulling prices from the slinky sidecar and sending them to the pricefeed server.
type PriceFetcher interface {
	Start(ctx context.Context) error
	Stop()
	FetchPrices(ctx context.Context) error
}

// PriceFetcherImpl implements the PriceFetcher interface.
type PriceFetcherImpl struct {
	marketPairFetcher      MarketPairFetcher
	priceFeedServiceClient api.PriceFeedServiceClient
	grpcClient             daemontypes.GrpcClient
	pricesSocket           string
	pricesConn             *grpc.ClientConn
	slinky                 oracleclient.OracleClient
	logger                 log.Logger
}

// NewPriceFetcher creates a PriceFetcher.
func NewPriceFetcher(
	marketPairFetcher MarketPairFetcher,
	grpcClient daemontypes.GrpcClient,
	pricesSocket string,
	slinky oracleclient.OracleClient,
	logger log.Logger) PriceFetcher {
	return &PriceFetcherImpl{
		marketPairFetcher: marketPairFetcher,
		grpcClient:        grpcClient,
		pricesSocket:      pricesSocket,
		slinky:            slinky,
		logger:            logger,
	}
}

// Start initializes the underlying connections of the PriceFetcher.
func (p *PriceFetcherImpl) Start(ctx context.Context) error {
	cancelCtx, cf := context.WithTimeout(ctx, SlinkyPriceServerConnectionTimeout)
	defer cf()
	pricesConn, err := p.grpcClient.NewGrpcConnection(cancelCtx, p.pricesSocket)
	if err != nil {
		return err
	}
	p.pricesConn = pricesConn
	p.priceFeedServiceClient = api.NewPriceFeedServiceClient(p.pricesConn)
	return p.slinky.Start(ctx)
}

// Stop closes all open connections.
func (p *PriceFetcherImpl) Stop() {
	if p.pricesConn != nil {
		_ = p.grpcClient.CloseConnection(p.pricesConn)
	}
	_ = p.slinky.Stop()
}

// FetchPrices pulls prices from Slinky, translates the returned data format to dydx-compatible types,
// and sends the price updates to the index price cache via the pricefeed server.
// It uses the MarketPairFetcher to efficiently map between Slinky's CurrencyPair primary key and dydx's
// MarketParam (or MarketPrice) ID.
//
// The markets in the index price cache will only have a single index price (from slinky).
// This is because the sidecar pre-aggregates market data.
func (p *PriceFetcherImpl) FetchPrices(ctx context.Context) error {
	// get prices from slinky sidecar via GRPC
	slinkyResponse, err := p.slinky.Prices(ctx, &types.QueryPricesRequest{})
	if err != nil {
		return err
	}

	// update the prices keeper w/ the most recent prices for the relevant markets
	var updates []*api.MarketPriceUpdate
	for currencyPairString, priceString := range slinkyResponse.Prices {
		// convert currency-pair string (index) into currency-pair object
		currencyPair, err := oracletypes.CurrencyPairFromString(currencyPairString)
		if err != nil {
			return err
		}
		p.logger.Info("turned price pair to currency pair",
			"string", currencyPairString,
			"currency pair", currencyPair.String())

		// get the market id for the currency pair
		id, err := p.marketPairFetcher.GetIDForPair(currencyPair)
		if err != nil {
			p.logger.Info("slinky client returned currency pair not found in MarketPairFetcher",
				"currency pair", currencyPairString,
				"error", err)
			continue
		}

		// parse the price string into a uint64
		price, err := strconv.ParseUint(priceString, 10, 64)
		if err != nil {
			return fmt.Errorf("slinky client returned price %s not parsable as uint64", priceString)
		}
		p.logger.Info("parsed update for", "market id", id, "price", price)

		// append the update to the list of MarketPriceUpdates to be sent to the app's price-feed service
		updates = append(updates, &api.MarketPriceUpdate{
			MarketId: id,
			ExchangePrices: []*api.ExchangePrice{
				{
					ExchangeId:     "slinky",
					Price:          price,
					LastUpdateTime: &slinkyResponse.Timestamp,
				},
			},
		})
	}

	// send the updates to the app's price-feed service -> these will then be piped to the
	// x/prices indexPriceCache via the pricefeed service
	if p.priceFeedServiceClient == nil {
		p.logger.Error("nil price feed service client")
	}
	if len(updates) == 0 {
		p.logger.Info("Slinky returned 0 valid market price updates")
		return nil
	}
	_, err = p.priceFeedServiceClient.UpdateMarketPrices(ctx, &api.UpdateMarketPricesRequest{MarketPriceUpdates: updates})
	if err != nil {
		p.logger.Error(err.Error())
		return err
	}
	return nil
}
