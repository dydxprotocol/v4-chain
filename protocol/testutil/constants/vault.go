package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

var (
	Vault_Clob0 = types.VaultId{
		Type:   types.VaultType_VAULT_TYPE_CLOB,
		Number: 0,
	}
	Vault_Clob1 = types.VaultId{
		Type:   types.VaultType_VAULT_TYPE_CLOB,
		Number: 1,
	}
	Vault_Clob7 = types.VaultId{
		Type:   types.VaultType_VAULT_TYPE_CLOB,
		Number: 7,
	}

	MsgDepositToMegavault_Alice0_100 = &types.MsgDepositToMegavault{
		SubaccountId:  &Alice_Num0,
		QuoteQuantums: dtypes.NewInt(100),
	}

	QuotingParams = types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     4_321,
		SpreadBufferPpm:                  1_789,
		SkewFactorPpm:                    767_323,
		OrderSizePctPpm:                  234_567,
		OrderExpirationSeconds:           111,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(9_876_543),
	}
	VaultParams = types.VaultParams{
		Status:        types.VaultStatus_VAULT_STATUS_QUOTING,
		QuotingParams: &QuotingParams,
	}
	InvalidQuotingParams = types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     4_321,
		SpreadBufferPpm:                  1_789,
		SkewFactorPpm:                    767_323,
		OrderSizePctPpm:                  234_567,
		OrderExpirationSeconds:           0,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(9_876_543),
	}
)
