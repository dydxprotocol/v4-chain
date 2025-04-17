package client

import (
	"context"
	"strconv"

	"cosmossdk.io/log"
	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"
	oracleclient "github.com/dydxprotocol/slinky/service/clients/oracle"
	"github.com/dydxprotocol/slinky/service/servers/oracle/types"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
)

// PriceFetcher is responsible for pulling prices from the slinky sidecar and sending them to the pricefeed server.
type PriceFetcher interface {
	Start(ctx context.Context) error
	Stop()
	FetchPrices(ctx context.Context) error
}

// PriceFetcherImpl implements the PriceFetcher interface.
type PriceFetcherImpl struct {
	marketPairFetcher MarketPairFetcher
	indexPriceCache   *pricefeedtypes.MarketToExchangePrices
	slinky            oracleclient.OracleClient
	logger            log.Logger
}

// NewPriceFetcher creates a PriceFetcher.
func NewPriceFetcher(
	marketPairFetcher MarketPairFetcher,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	slinky oracleclient.OracleClient,
	logger log.Logger) PriceFetcher {
	return &PriceFetcherImpl{
		marketPairFetcher: marketPairFetcher,
		indexPriceCache:   indexPriceCache,
		slinky:            slinky,
		logger:            logger,
	}
}

// Start initializes the underlying connections of the PriceFetcher.
func (p *PriceFetcherImpl) Start(ctx context.Context) error {
	return p.slinky.Start(ctx)
}

// Stop closes all open connections.
func (p *PriceFetcherImpl) Stop() {
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
		currencyPair, err := slinkytypes.CurrencyPairFromString(currencyPairString)
		if err != nil {
			return err
		}

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
			p.logger.Error("slinky client returned a price not parsable as uint64", "price", priceString)
			continue
		}
		p.logger.Debug("Parsed Slinky price update",
			"market id", id,
			"price", price,
			"string", currencyPairString,
			"currency pair", currencyPair.String())

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

	p.logger.Info("Slinky returned valid market price updates", "count", len(updates), "updates", updates)

	// send the updates to the indexPriceCache
	if len(updates) == 0 {
		p.logger.Info("Slinky returned 0 valid market price updates")
		return nil
	}
	p.indexPriceCache.UpdatePrices(updates)
	return nil
}
