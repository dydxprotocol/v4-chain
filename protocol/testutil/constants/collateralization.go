package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var (
	// Collateralization check functions.
	CollatCheck_EmptyUpdateResults_Success = func(
		subaccountMatchedOrders map[satypes.SubaccountId][]clobtypes.PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	) {
		return true, make(map[satypes.SubaccountId]satypes.UpdateResult)
	}
	CollatCheck_EmptyUpdateResults_Failure = func(
		subaccountMatchedOrders map[satypes.SubaccountId][]clobtypes.PendingOpenOrder,
	) (
		success bool,
		successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
	) {
		saMap := make(map[satypes.SubaccountId]satypes.UpdateResult)
		for a := range subaccountMatchedOrders {
			saMap[a] = satypes.NewlyUndercollateralized
		}
		return false, saMap
	}
)
