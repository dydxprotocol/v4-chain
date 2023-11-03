package server

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	servertypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

// PriceFeedServer defines the fields required for price updates.
type PriceFeedServer struct {
	marketToExchange *pricefeedtypes.MarketToExchangePrices
}

// WithPriceFeedMarketToExchangePrices sets the `MarketToExchangePrices` field.
// This is used by the price feed service to communicate price updates
// to the main application.
func (server *Server) WithPriceFeedMarketToExchangePrices(
	marketToExchange *pricefeedtypes.MarketToExchangePrices,
) *Server {
	server.marketToExchange = marketToExchange
	return server
}

// ExpectPricefeedDaemon registers the pricefeed daemon with the server. This is required
// in order to ensure that the daemon service is called at least once during every
// maximumAcceptableUpdateDelay duration. It will cause the protocol to panic if the daemon does not
// respond within maximumAcceptableUpdateDelay duration.
func (server *Server) ExpectPricefeedDaemon(maximumAcceptableUpdateDelay time.Duration) {
	server.registerDaemon(servertypes.PricefeedDaemonServiceName, maximumAcceptableUpdateDelay)
}

// UpdateMarketPrices updates prices from exchanges for each market provided.
func (s *Server) UpdateMarketPrices(
	ctx context.Context,
	req *api.UpdateMarketPricesRequest,
) (*api.UpdateMarketPricesResponse, error) {
	// Measure latency in ingesting and handling gRPC price update.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedServer,
		time.Now(),
		metrics.PricefeedServerUpdatePrices,
		metrics.Latency,
	)

	// If the daemon is unable to report a response, there is either an error in the registration of
	// this daemon, or another one. In either case, the protocol should panic.
	if err := s.reportResponse(servertypes.PricefeedDaemonServiceName); err != nil {
		panic(err)
	}

	if s.marketToExchange == nil {
		panic(
			errorsmod.Wrapf(
				types.ErrServerNotInitializedCorrectly,
				"MarketToExchange not initialized",
			),
		)
	}

	if err := validateMarketPricesUpdatesMessage(req); err != nil {
		// Log if failure occurs during an update.
		s.logger.Error("Failed to validate price update message", "error", err)
		return nil, err
	}

	s.marketToExchange.UpdatePrices(req.MarketPriceUpdates)
	return &api.UpdateMarketPricesResponse{}, nil
}

// validateMarketPricesUpdatesMessage validates a `UpdateMarketPricesRequest`.
func validateMarketPricesUpdatesMessage(req *api.UpdateMarketPricesRequest) error {
	if len(req.MarketPriceUpdates) == 0 {
		return types.ErrPriceFeedMarketPriceUpdateEmpty
	}

	for _, mpu := range req.MarketPriceUpdates {
		if err := validateMarketPriceUpdate(mpu); err != nil {
			// Measure failure per market in validation.
			telemetry.IncrCounterWithLabels(
				[]string{
					metrics.PricefeedServer,
					metrics.PricefeedServerValidatePrices,
					metrics.Error,
				},
				1,
				[]gometrics.Label{
					pricefeedmetrics.GetLabelForMarketId(mpu.MarketId),
				},
			)

			return err
		}
	}

	return nil
}

// validateMarketPriceUpdate validates a single `MarketPriceUpdate`.
func validateMarketPriceUpdate(mpu *api.MarketPriceUpdate) error {
	for _, ep := range mpu.ExchangePrices {
		if ep.Price == constants.DefaultPrice {
			return generateSdkErrorForExchangePrice(
				types.ErrPriceFeedInvalidPrice,
				ep,
				mpu.MarketId,
			)
		}

		if ep.LastUpdateTime == nil {
			return generateSdkErrorForExchangePrice(
				types.ErrPriceFeedLastUpdateTimeNotSet,
				ep,
				mpu.MarketId,
			)
		}
	}
	return nil
}

// generateSdkErrorForExchangePrice generates an error for an invalid `ExchangePrice`.
func generateSdkErrorForExchangePrice(err error, ep *api.ExchangePrice, marketId uint32) error {
	return errorsmod.Wrapf(err, "ExchangePrice: %v and MarketId: %d", ep, marketId)
}
