package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	"github.com/stretchr/testify/require"
)

func TestNewUpsertVaultEvent_Success(t *testing.T) {
	upsertVaultEvent := events.NewUpsertVaultEvent(
		constants.Alice_Num0.Owner,
		0,
		vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
	)
	expectedUpsertVaultEventProto := &events.UpsertVaultEventV1{
		Address:    constants.Alice_Num0.Owner,
		ClobPairId: 0,
		Status:     v1types.VaultStatus_VAULT_STATUS_QUOTING,
	}
	require.Equal(t, expectedUpsertVaultEventProto, upsertVaultEvent)
}
