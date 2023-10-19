package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "rewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "tmp_" + ModuleName
)

// State
const (
	// RewardShareKeyPrefix is the prefix to retrieve reward shares for all addresses.
	RewardShareKeyPrefix = "Shares:"

	// ParamsKey is the key for the params
	ParamsKey = "Params"
)

// Module accounts
const (
	// TreasuryAccountName defines the root string for the rewards treasury account address.
	TreasuryAccountName = "rewards_treasury"

	// VesterAccountName defines the root string for the rewards vester account address.
	VesterAccountName = "rewards_vester"
)
