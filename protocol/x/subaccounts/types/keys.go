package types

const (
	// ModuleName defines the module name
	ModuleName = "subaccounts"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "transient_" + ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
