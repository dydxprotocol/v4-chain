package types

// ExchangeQueryConfig contains configuration values for querying an exchange, passed in on startup.
// The configuration values include
//  1. `ExchangeId`
//  2. `IntervalMs` delay between task-loops where each task-loop sends API requests to an exchange
//  3. `TimeoutMs` max time to wait on an API call to an exchange
//  4. `MaxQueries` max number of API calls made to an exchange per task-loop. This parameter is used
//     for rate limiting requests to the exchange.
//
// For single-market API exchanges, the price fetcher will send approximately
// MaxQueries API responses into the exchange's buffered channel once every IntervalMs milliseconds.
// Note: the `ExchangeQueryConfig` will be used in the map of `{ exchangeId, `ExchangeQueryConfig` }`
// that dictates how the pricefeed client queries for market prices.
type ExchangeQueryConfig struct {
	ExchangeId ExchangeId `json:"exchange_id"`
	IntervalMs uint32     `json:"interval_ms"`
	TimeoutMs  uint32     `json:"timeout_ms"`
	MaxQueries uint32     `json:"max_queries"`
}
