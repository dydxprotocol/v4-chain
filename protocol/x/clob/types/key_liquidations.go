package types

import (
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	// SubaccountLiquidationInfoKeyPrefix is the prefix to retrieve the liquidation information
	// for a subaccount within the last block.
	SubaccountLiquidationInfoKeyPrefix = "SubaccountLiquidations/value/"
)

// SubaccountLiquidationInfoKey returns the store key to retrieve the liquidation information
// for a subaccount within the last block.
func SubaccountLiquidationInfoKey(
	id satypes.SubaccountId,
) []byte {
	idBytes, err := id.Marshal()
	if err != nil {
		panic(err)
	}

	return idBytes
}
