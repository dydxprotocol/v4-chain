package client

import "time"

var (
	// SlinkyPriceServerConnectionTimeout controls the timeout of establishing a
	// grpc connection to the pricefeed server.
	SlinkyPriceServerConnectionTimeout = time.Second * 5
	// SlinkyPriceFetchDelay controls the frequency at which we pull prices from slinky and push
	// them to the pricefeed server.
	SlinkyPriceFetchDelay = time.Second * 2
	// SlinkyMarketParamFetchDelay is the frequency at which we query the x/price module to refresh mappings from
	// currency pair to x/price ID.
	SlinkyMarketParamFetchDelay = time.Millisecond * 1900
	SlinkySidecarCheckDelay     = time.Second * 60
)

const (
	// SlinkyClientDaemonModuleName is the module name used in logging.
	SlinkyClientDaemonModuleName                      = "slinky-client-daemon"
	SlinkyClientPriceFetcherDaemonModuleName          = "slinky-client-price-fetcher-daemon"
	SlinkyClientMarketPairFetcherDaemonModuleName     = "slinky-client-market-pair-fetcher-daemon"
	SlinkyClientSidecarVersionFetcherDaemonModuleName = "slinky-client-sidecar-version-fetcher-daemon"
	MinSidecarVersion                                 = "v1.0.12"
)
