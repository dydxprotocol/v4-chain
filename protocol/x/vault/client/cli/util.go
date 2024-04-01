package cli

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetVaultTypeFromString returns a vault type from a string.
func GetVaultTypeFromString(rawType string) (vaultType types.VaultType, err error) {
	switch rawType {
	case "clob":
		return types.VaultType_VAULT_TYPE_CLOB, nil
	default:
		return vaultType, fmt.Errorf("invalid vault type: %s", rawType)
	}
}
