package types

// ExchangeConfig maintains an exchange specific configuration for the price feed client including:
// 1) `Markets` to query
// 2) `ExchangeStartupConfig` contains all information on how/how often the pricefeed client should
// query the exchange
type ExchangeConfig struct {
	Markets               []uint32
	ExchangeStartupConfig ExchangeStartupConfig
	IsMultiMarket         bool
}
