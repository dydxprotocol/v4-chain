package server

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"
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

// UpdateMarketPrices updates prices from exchanges for each market provided.
func (s *Server) UpdateMarketPrices(
	ctx context.Context,
	req *api.UpdateMarketPricesRequest,
) (
	response *api.UpdateMarketPricesResponse,
	err error,
) {
	// Measure latency in ingesting and handling gRPC price update.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedServer,
		time.Now(),
		metrics.PricefeedServerUpdatePrices,
		metrics.Latency,
	)

	// This panic is an unexpected condition because we initialize the market price cache in app initialization before
	// starting the server or daemons.
	if s.marketToExchange == nil {
		panic(
			errorsmod.Wrapf(
				daemontypes.ErrServerNotInitializedCorrectly,
				"MarketToExchange not initialized",
			),
		)
	}

	if err = validateMarketPricesUpdatesMessage(req); err != nil {
		// Log if failure occurs during an update.
		s.logger.Error("Failed to validate price update message", "error", err)
		return nil, err
	}

	s.marketToExchange.UpdatePrices(req.MarketPriceUpdates)

	// Capture valid responses in metrics.
	s.reportValidResponse(types.PricefeedDaemonServiceName)

	return &api.UpdateMarketPricesResponse{}, nil
}

// validateMarketPricesUpdatesMessage validates a `UpdateMarketPricesRequest`.
func validateMarketPricesUpdatesMessage(req *api.UpdateMarketPricesRequest) error {
	if len(req.MarketPriceUpdates) == 0 {
		return daemontypes.ErrPriceFeedMarketPriceUpdateEmpty
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
				daemontypes.ErrPriceFeedInvalidPrice,
				ep,
				mpu.MarketId,
			)
		}

		if ep.LastUpdateTime == nil {
			return generateSdkErrorForExchangePrice(
				daemontypes.ErrPriceFeedLastUpdateTimeNotSet,
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
