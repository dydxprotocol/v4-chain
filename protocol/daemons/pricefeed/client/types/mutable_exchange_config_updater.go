package types

// ExchangeConfigUpdater is the interface that wraps the UpdateMutableExchangeConfig method.
// ExchangeConfigUpdater objects are keyed by exchange id and receive updates notifying
// them that the mutable exchange market configuration has been updated, along with the all new configs.
// This interface is added to avoid import loops that occur when importing the `PriceFetcher` type
// directly into `PricefeedMutableMarketConfigs`.
type ExchangeConfigUpdater interface {
	GetExchangeId() ExchangeId
	// UpdateMutableExchangeConfig notifies the object that the mutable exchange market configuration
	// for this object's exchange has been updated with a new configuration. It also provides
	// the current market configs for all supported markets on the exchange.
	UpdateMutableExchangeConfig(
		newExchangeConfig *MutableExchangeMarketConfig,
		newMarketConfigs []*MutableMarketConfig,
	) error
}
