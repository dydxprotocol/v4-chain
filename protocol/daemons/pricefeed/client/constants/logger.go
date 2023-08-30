package constants

import "time"

const (
	SubmoduleLogKey = "submodule"

	ExchangeIdLogKey = "exchangeId"
	MarketIdLogKey   = "marketId"
	PriceLogKey      = "Price"

	PricefeedDaemonModuleName       = "pricefeed-daemon"
	PriceFetcherSubmoduleName       = "price-fetcher"
	PriceEncoderSubmoduleName       = "price-encoder"
	PriceUpdaterSubmoduleName       = "price-updater"
	MarketParamUpdaterSubmoduleName = "market-param-updater"

	// PriceDaemonStartupErrorGracePeriod defines the amount of time the daemon waits before logging issues that are
	// intermittent on daemon startup as true errors. Examples of this includes price conversion failures due to
	// an uninitialized prices cache, and failures to fetch market param updates due to a delay on the protocol side
	// in starting the prices query service.
	// If the protocol is not started within this grace period, the daemon will report these errors as true errors.
	PriceDaemonStartupErrorGracePeriod = 120 * time.Second
)
