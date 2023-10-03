package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "stats"

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "transient_" + ModuleName

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// EpochStatsKeyPrefix is the prefix to retrieve the EpochStats for a given epoch
	EpochStatsKeyPrefix = "epoch_stats/"

	// UserStatsKeyPrefix is the prefix to retrieve the UserStats for a given user
	UserStatsKeyPrefix = "user_stats/"

	// StatsMetadataKey is the key to get the StatsMetadata for the module
	StatsMetadataKey = "stats_metadata"

	// GlobalStatsKey is the key to get the GlobalStats for the module
	GlobalStatsKey = "global_stats"

	// BlockStatsKey is the key to get the BlockStats for the module
	BlockStatsKey = "block_stats"

	// ParamsKey defines the key for the params
	ParamsKey = "params"
)
