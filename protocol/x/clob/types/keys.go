package types

const (
	// ModuleName defines the module name
	ModuleName = "clob"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "transient_" + ModuleName

	// InsuranceFundName defines the root string for the insurance fund account address
	InsuranceFundName = "insurance_fund"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
