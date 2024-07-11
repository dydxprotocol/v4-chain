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

	// ParamsKey is the key to retrieve Params.
	ParamsKey = "Params"

	// VaultParamsKeyPrefix is the prefix to retrieve all VaultParams.
	VaultParamsKeyPrefix = "VaultParams:"

	// VaultAddressKeyPrefix is the prefix to retrieve all vault addresses.
	VaultAddressKeyPrefix = "VaultAddress:"

	// MostRecentClientIdsKeyPrefix is the prefix to retrieve all most recent client IDs.
	// MostRecentClientIdsStore: vaultId VaultId -> clientIds []uint32
	MostRecentClientIdsKeyPrefix = "MostRecentClientIds:"
)
