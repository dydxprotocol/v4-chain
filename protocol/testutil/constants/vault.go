package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

type VaultState struct {
	VaultId     types.VaultId
	OwnerShares []*types.OwnerShare
}

var (
	Vault_Clob0 = types.VaultId{
		Type:   types.VaultType_VAULT_TYPE_CLOB,
		Number: 0,
	}
	Vault_Clob1 = types.VaultId{
		Type:   types.VaultType_VAULT_TYPE_CLOB,
		Number: 1,
	}

	Vault_Clob0_SingleOwner_Alice0_1000 = VaultState{
		VaultId: Vault_Clob0,
		OwnerShares: []*types.OwnerShare{
			{
				Owner:  Alice_Num0.Owner,
				Shares: &types.NumShares{NumShares: dtypes.NewInt(1000)},
			},
		},
	}
	Vault_Clob0_MultiOwner_Alice0_1000_Bob0_2500 = VaultState{
		VaultId: Vault_Clob0,
		OwnerShares: []*types.OwnerShare{
			{
				Owner:  Alice_Num0.Owner,
				Shares: &types.NumShares{NumShares: dtypes.NewInt(1000)},
			},
			{
				Owner:  Bob_Num0.Owner,
				Shares: &types.NumShares{NumShares: dtypes.NewInt(2500)},
			},
		},
	}

	MsgDepositToVault_Clob0_Alice0_100 = &types.MsgDepositToVault{
		VaultId:       &Vault_Clob0,
		SubaccountId:  &Alice_Num0,
		QuoteQuantums: dtypes.NewInt(100),
	}
)
