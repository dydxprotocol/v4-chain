package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "stats"

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "tmp_" + ModuleName

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// EpochStatsKeyPrefix is the prefix to retrieve the EpochStats for a given epoch
	EpochStatsKeyPrefix = "Epoch:"

	// UserStatsKeyPrefix is the prefix to retrieve the UserStats for a given user
	UserStatsKeyPrefix = "User:"

	// StatsMetadataKey is the key to get the StatsMetadata for the module
	StatsMetadataKey = "Metadata"

	// GlobalStatsKey is the key to get the GlobalStats for the module
	GlobalStatsKey = "Global"

	// BlockStatsKey is the key to get the BlockStats for the module
	BlockStatsKey = "Block"

	// ParamsKey defines the key for the params
	ParamsKey = "Params"

	// CachedStakeAmountKey is the key to get the cached stake amount
	CachedStakeAmountKeyPrefix = "CachedStakeAmount:"
)
