package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// NewUpsertVaultEvent creates a UpsertVaultEventV1
// representing an create / update of a vault.
func NewUpsertVaultEvent(
	vaultAddress string,
	clobPairId uint32,
	status types.VaultStatus,
) *UpsertVaultEventV1 {
	return &UpsertVaultEventV1{
		Address:    vaultAddress,
		ClobPairId: clobPairId,
		Status:     v1.VaultStatusToIndexerVaultStatus(status),
	}
}
