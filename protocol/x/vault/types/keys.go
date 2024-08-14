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
	// TotalSharesKeyPrefix is the prefix to retrieve all TotalShares.
	TotalSharesKeyPrefix = "TotalShares:"

	// OwnerSharesKeyPrefix is the prefix to retrieve all OwnerShares.
	// OwnerShares store: vaultId VaultId -> owner string -> shares NumShares.
	OwnerSharesKeyPrefix = "OwnerShares:"

	// DefaultQuotingParams is the key to retrieve DefaultQuotingParams.
	// A vault uses DefaultQuotingParams if it does not have its own QuotingParams.
	DefaultQuotingParamsKey = "DefaultQuotingParams"

	// QuotingParamsKeyPrefix is the prefix to retrieve all QuotingParams.
	// QuotingParams store: vaultId VaultId -> QuotingParams.
	QuotingParamsKeyPrefix = "QuotingParams:"

	// VaultAddressKeyPrefix is the prefix to retrieve all vault addresses.
	VaultAddressKeyPrefix = "VaultAddress:"

	// MostRecentClientIdsKeyPrefix is the prefix to retrieve all most recent client IDs.
	// MostRecentClientIdsStore: vaultId VaultId -> clientIds []uint32
	MostRecentClientIdsKeyPrefix = "MostRecentClientIds:"
)
