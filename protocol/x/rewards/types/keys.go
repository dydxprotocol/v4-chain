package types

const (
	// ModuleName defines the module name
	ModuleName = "rewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "transient_" + ModuleName

	// TreasuryAccountName defines the root string for the rewards treasury account address.
	TreasuryAccountName = "rewards_treasury"

	// VesterAccountName defines the root string for the rewards vester account address.
	VesterAccountName = "rewards_vester"

	// RewardShareKeyPrefix is the prefix to retrieve reward shares for all addresses.
	RewardShareKeyPrefix = "reward_shares/"
)
