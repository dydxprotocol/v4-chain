package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "subaccounts"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// SubaccountKeyPrefix is the prefix to retrieve all Subaccounts
	SubaccountKeyPrefix = "SA:"
	// NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix is the prefix for the store key that
	// stores the last block a negative TNC subaccount was seen in state for a specific collateral pool.
	NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix = "NegSA:"
	// Suffix for the store key to the last block a negative TNC subaccount was seen in state for the
	// cross collateral pool.
	CrossCollateralSuffix = "cross"

	// Safety Heap
	SafetyHeapStorePrefix             = "SH"
	SafetyHeapSubaccountIdsPrefix     = "Heap/"
	SafetyHeapSubaccountToIndexPrefix = "Idx/"
	SafetyHeapLengthPrefix            = "Len/"

	// Leverage
	LeverageKeyPrefix = "Lev:"
)
