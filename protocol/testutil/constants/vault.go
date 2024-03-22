package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

var (
	Vault_Clob_0 = vaulttypes.VaultId{
		Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
		Number: 0,
	}
	Vault_Clob_1 = vaulttypes.VaultId{
		Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
		Number: 1,
	}

	MsgDepositToVault_Clob0_Alice0_100 = &types.MsgDepositToVault{
		VaultId:       &Vault_Clob_0,
		SubaccountId:  &Alice_Num0,
		QuoteQuantums: dtypes.NewInt(100),
	}
)
