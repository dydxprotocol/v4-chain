package types

import (
	"sync"

	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// DaemonLiquidationInfo maintains the list of subaccount ids to be liquidated
// in the next block. Methods are goroutine safe.
type DaemonLiquidationInfo struct {
	sync.Mutex               // lock
	subaccountsWithPositions map[uint32]*clobtypes.SubaccountOpenPositionInfo
}

// NewDaemonLiquidationInfo creates a new `NewDaemonLiquidationInfo` struct.
func NewDaemonLiquidationInfo() *DaemonLiquidationInfo {
	return &DaemonLiquidationInfo{
		subaccountsWithPositions: make(map[uint32]*clobtypes.SubaccountOpenPositionInfo),
	}
}

// UpdateSubaccountsWithPositions updates the struct with the given a list of subaccount ids with open positions.
func (ls *DaemonLiquidationInfo) UpdateSubaccountsWithPositions(
	subaccountsWithPositions []clobtypes.SubaccountOpenPositionInfo,
) {
	ls.Lock()
	defer ls.Unlock()
	ls.subaccountsWithPositions = make(map[uint32]*clobtypes.SubaccountOpenPositionInfo)
	for _, info := range subaccountsWithPositions {
		clone := &clobtypes.SubaccountOpenPositionInfo{
			PerpetualId:                  info.PerpetualId,
			SubaccountsWithLongPosition:  make([]satypes.SubaccountId, len(info.SubaccountsWithLongPosition)),
			SubaccountsWithShortPosition: make([]satypes.SubaccountId, len(info.SubaccountsWithShortPosition)),
		}
		copy(clone.SubaccountsWithLongPosition, info.SubaccountsWithLongPosition)
		copy(clone.SubaccountsWithShortPosition, info.SubaccountsWithShortPosition)
		ls.subaccountsWithPositions[info.PerpetualId] = clone
	}
}

// GetSubaccountsWithOpenPositions returns the list of subaccount ids with open positions for a perpetual.
func (ls *DaemonLiquidationInfo) GetSubaccountsWithOpenPositions(
	perpetualId uint32,
) []satypes.SubaccountId {
	ls.Lock()
	defer ls.Unlock()

	result := make([]satypes.SubaccountId, 0)
	if info, ok := ls.subaccountsWithPositions[perpetualId]; ok {
		result = append(result, info.SubaccountsWithLongPosition...)
		result = append(result, info.SubaccountsWithShortPosition...)
	}
	return result
}

// GetSubaccountsWithOpenPositionsOnSide returns the list of subaccount ids with open positions
// on a specific side for a perpetual.
func (ls *DaemonLiquidationInfo) GetSubaccountsWithOpenPositionsOnSide(
	perpetualId uint32,
	isLong bool,
) []satypes.SubaccountId {
	ls.Lock()
	defer ls.Unlock()

	result := make([]satypes.SubaccountId, 0)
	if info, ok := ls.subaccountsWithPositions[perpetualId]; ok {
		if isLong {
			result = append(result, info.SubaccountsWithLongPosition...)
		} else {
			result = append(result, info.SubaccountsWithShortPosition...)
		}
	}
	return result
}
