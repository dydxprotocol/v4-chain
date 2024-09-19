package types

// Module name and store keys.
const (
	// ModuleName defines the module name.
	ModuleName = "vault"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)

// State.
const (
	// TotalSharesKey is the key to retrieve total shares.
	TotalSharesKey = "TotalShares"

	// OwnerSharesKeyPrefix is the prefix to retrieve all OwnerShares.
	// OwnerShares store: owner string -> shares NumShares.
	OwnerSharesKeyPrefix = "OwnerShares:"

	// OwnerShareUnlocksKeyPrefix is the prefix to retrieve all OwnerShareUnlocks.
	// OwnerShareUnlocks store: owner string -> ownerShareUnlocks OwnerShareUnlocks.
	OwnerShareUnlocksKeyPrefix = "OwnerShareUnlocks:"

	// DefaultQuotingParams is the key to retrieve DefaultQuotingParams.
	// A vault uses DefaultQuotingParams if it does not have its own QuotingParams.
	DefaultQuotingParamsKey = "DefaultQuotingParams"

	// VaultParamsKeyPrefix is the prefix to retrieve all VaultParams.
	// VaultParams store: vaultId VaultId -> VaultParams.
	VaultParamsKeyPrefix = "VaultParams:"

	// VaultAddressKeyPrefix is the prefix to retrieve all vault addresses.
	VaultAddressKeyPrefix = "VaultAddress:"

	// MostRecentClientIdsKeyPrefix is the prefix to retrieve all most recent client IDs.
	// MostRecentClientIdsStore: vaultId VaultId -> clientIds []uint32
	MostRecentClientIdsKeyPrefix = "MostRecentClientIds:"

	// OperatorParamsKey is the key to retrieve OperatorParams.
	OperatorParamsKey = "OperatorParams"
)

// Module accounts
const (
	// MegavaultAccountName defines the root string for megavault module account.
	MegavaultAccountName = "megavault"
)
