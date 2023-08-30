package types

// ExchangeId is the unique id for an `Exchange` in the `Prices` module.
// The id will be matched against each exchange's `exchangeName` in the `MarketParam`'s `exchange_config_json`.
type ExchangeId = string
