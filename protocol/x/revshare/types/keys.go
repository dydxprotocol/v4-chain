package types

// Module name and store keys.
const (
	// ModuleName defines the module name.
	ModuleName = "revshare"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)

// State
const (
	// Key for MarketMapperRevenueShareParams
	MarketMapperRevenueShareParamsKey = "MarketMapperRevenueShareParams"

	// Key prefix for storing MarketMapperRevShareDetails per market
	MarketMapperRevSharePrefix = "MarketMapperRevShare:"

	UnconditionalRevShareConfigKey = "UnconditionalRevShareConfig"

	// Key prefix for storing OrderRouterRevShareParams
	OrderRouterRevSharePrefix = "OrderRouterRevShare:"
)
