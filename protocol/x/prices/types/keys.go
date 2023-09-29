package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "prices"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// MarketParamKeyPrefix is the prefix to retrieve all MarketParams
	MarketParamKeyPrefix = "market_param/"

	// MarketPriceKeyPrefix is the prefix to retrieve all MarketPrices
	MarketPriceKeyPrefix = "market_price/"
)
