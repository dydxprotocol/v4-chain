package types

import (
	fmt "fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (id *VaultId) ToStateKey() []byte {
	b, err := id.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

// ToModuleAccountAddress returns the module account address for the vault
// (generated from string "vault-<type>-<number>")
func (id *VaultId) ToModuleAccountAddress() string {
	return authtypes.NewModuleAddress(
		fmt.Sprintf("vault-%s-%d", id.Type, id.Number),
	).String()
}
