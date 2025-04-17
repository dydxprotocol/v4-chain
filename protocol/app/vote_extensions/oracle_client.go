package vote_extensions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	oracleservicetypes "github.com/dydxprotocol/slinky/service/servers/oracle/types"
	"google.golang.org/grpc"
)

// OraclePrices is an implementation of the Slinky OracleClient interface.
// This object is responsible for requesting prices from the x/prices module, and injecting those prices into the
// vote-extension.
// The
type OraclePrices struct {
	PricesKeeper PricesKeeper
}

// NewOraclePrices returns a new OracleClient object.
func NewOraclePrices(pricesKeeper PricesKeeper) *OraclePrices {
	return &OraclePrices{
		PricesKeeper: pricesKeeper,
	}
}

// Start is a no-op.
func (o *OraclePrices) Start(_ context.Context) error {
	return nil
}

// Stop is a no-op.
func (o *OraclePrices) Stop() error {
	return nil
}

// Prices is called in ExtendVoteHandler to determine which Prices are put into the extended commit.
// This method is responsible for doing the following:
//  1. Get the latest prices from the x/prices module's indexPriceCache via GetValidMarketPriceUpdates
//  2. Translate the response from x/prices into a QueryPricesResponse, and return it.
//
// This method fails if:
//   - The passed in context is not an sdk.Context
func (o *OraclePrices) Prices(ctx context.Context,
	_ *oracleservicetypes.QueryPricesRequest,
	_ ...grpc.CallOption) (*oracleservicetypes.QueryPricesResponse, error) {
	sdkCtx, ok := ctx.(sdk.Context)
	if !ok {
		return nil, fmt.Errorf("OraclePrices was passed on non-sdk context object")
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
