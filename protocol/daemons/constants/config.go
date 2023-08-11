package constants

const (
	// PriceFeed configs
	// DefaultPrice is the default value for `Price` field in `UpdateMarketPricesRequest`.
	DefaultPrice = 0

	// Liquidation configs
	// LiquidationLoopDelayMs is the delay between each loop iteration.
	LiquidationLoopDelayMs = 1000
	// LiquidationGetSubaccountPageLimit defines the max number of subaccounts to be fetched
	// in a single paginated request.
	LiquidationGetSubaccountPageLimit = 1000

	PricefeedExchangeConfigFileName = "pricefeed_exchange_config.toml"
)
