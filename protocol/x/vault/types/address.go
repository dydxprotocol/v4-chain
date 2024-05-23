package types

var (
	vaultAddressSet map[string]struct{}
)

const (
	// Number of vault supported
	SupportedNumVaults = 200
	SupportedVaultType = VaultType_VAULT_TYPE_CLOB
)

func init() {
	vaultAddressSet = make(map[string]struct{})
	for i := 0; i < SupportedNumVaults; i++ {
		vaultId := VaultId{
			Type:   SupportedVaultType,
			Number: uint32(i),
		}
		address := vaultId.ToModuleAccountAddress()
		vaultAddressSet[address] = struct{}{}
	}
}

// Returns true if the given address is a vault address.
func IsVaultAddress(address string) bool {
	_, ok := vaultAddressSet[address]
	return ok
}
