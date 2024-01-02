package clob

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func GetOpenPositionsFromSubaccounts(
	subaccounts []satypes.Subaccount,
) []clobtypes.SubaccountOpenPositionInfo {
	positionMap := make(map[uint32]*clobtypes.SubaccountOpenPositionInfo)
	for _, subaccount := range subaccounts {
		for _, position := range subaccount.PerpetualPositions {
			info, ok := positionMap[position.PerpetualId]
			if !ok {
				info = &clobtypes.SubaccountOpenPositionInfo{
					PerpetualId:                  position.PerpetualId,
					SubaccountsWithLongPosition:  make([]satypes.SubaccountId, 0),
					SubaccountsWithShortPosition: make([]satypes.SubaccountId, 0),
				}
				positionMap[position.PerpetualId] = info
			}
			if position.GetIsLong() {
				info.SubaccountsWithLongPosition = append(
					info.SubaccountsWithLongPosition,
					*subaccount.Id,
				)
			} else {
				info.SubaccountsWithShortPosition = append(
					info.SubaccountsWithShortPosition,
					*subaccount.Id,
				)
			}
		}
	}

	positionSlice := make([]clobtypes.SubaccountOpenPositionInfo, 0)
	for _, info := range positionMap {
		positionSlice = append(positionSlice, *info)
	}
	return positionSlice
}
