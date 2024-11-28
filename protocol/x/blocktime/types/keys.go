package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "blocktime"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// DowntimeParamsKey defines the key for the DowntimeParams
	DowntimeParamsKey = "DowntimeParams"

	// AllDowntimeInfoKey defines the key for AllDowntimeInfo
	AllDowntimeInfoKey = "AllDowntimeInfo"

	// PreviousBlockInfoKey defines the key for PreviousBlockInfo
	PreviousBlockInfoKey = "PreviousBlockInfo"

	// SynchronyParamsKey defines the key for the SynchronyParams
	SynchronyParamsKey = "SP:"
)
