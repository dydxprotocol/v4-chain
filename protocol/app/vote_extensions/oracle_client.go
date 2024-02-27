package vote_extensions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	oracleservicetypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	"google.golang.org/grpc"
)

// OracleClient is a wrapper around the default Slinky OracleClient interface. This object is responsible for requesting
// prices from the x/prices module (originally sent via the sidecar to the price-feed service), and
// injecting those prices into the vote-extension.
type OracleClient struct {
	PricesKeeper PricesKeeper
}

// NewOracleClient returns a new OracleClient object.
func NewOracleClient(pricesKeeper PricesKeeper) *OracleClient {
	return &OracleClient{
		PricesKeeper: pricesKeeper,
	}
}

// Start starts the OracleClient.
func (o *OracleClient) Start(ctx context.Context) error {
	return nil
}

// Stop stops the OracleClient.
func (o *OracleClient) Stop() error {
	return nil
}

// Prices is a wrapper around the Slinky OracleClient's Prices method. This method is responsible for doing the following:
//  1. Get the latest prices from the x/prices module's indexPriceCache via GetValidMarketPriceUpdates
//  2. Translate the response from x/prices into a QueryPricesResponse, and return it.
//
// This method fails if:
//   - The passed in context is not an sdk.Context
func (o *OracleClient) Prices(ctx context.Context, in *oracleservicetypes.QueryPricesRequest, opts ...grpc.CallOption) (*oracleservicetypes.QueryPricesResponse, error) {
	sdkCtx, ok := ctx.(sdk.Context)
	if !ok {
		return nil, fmt.Errorf("oracle client was passed on non-sdk context object")
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
		Timestamp: time.Now(),
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
