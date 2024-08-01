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
	MarketParamKeyPrefix = "Param:"

	// MarketPriceKeyPrefix is the prefix to retrieve all MarketPrices
	MarketPriceKeyPrefix = "Price:"

	// NextIDKey is the key for the next market ID
	NextMarketIDKey = "NextMarketID"
)
