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

// GetVaultStatusFromString returns a vault status from a string.
func GetVaultStatusFromString(rawStatus string) (vaultStatus types.VaultStatus, err error) {
	switch rawStatus {
	case "deactivated":
		return types.VaultStatus_VAULT_STATUS_DEACTIVATED, nil
	case "stand_by":
		return types.VaultStatus_VAULT_STATUS_STAND_BY, nil
	case "quoting":
		return types.VaultStatus_VAULT_STATUS_QUOTING, nil
	case "close_only":
		return types.VaultStatus_VAULT_STATUS_CLOSE_ONLY, nil
	default:
		return vaultStatus, fmt.Errorf(`invalid vault status: %s.
										options are:
										- deactivated
										- stand_by
										- quoting
										- close_only`, rawStatus)
	}
}
