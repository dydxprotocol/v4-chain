package types

import (
	"sync"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// DaemonLiquidationInfo maintains the list of subaccount ids to be liquidated
// in the next block. Methods are goroutine safe.
type DaemonLiquidationInfo struct {
	sync.Mutex                                       // lock
	blockHeight               uint32                 // block height of the last update
	liquidatableSubaccountIds []satypes.SubaccountId // liquidatable subaccount ids
	negativeTncSubaccountIds  []satypes.SubaccountId // negative total net collateral subaccount ids
	subaccountsWithPositions  map[uint32]*clobtypes.SubaccountOpenPositionInfo
}

// NewDaemonLiquidationInfo creates a new `NewDaemonLiquidationInfo` struct.
func NewDaemonLiquidationInfo() *DaemonLiquidationInfo {
	return &DaemonLiquidationInfo{
		liquidatableSubaccountIds: make([]satypes.SubaccountId, 0),
		negativeTncSubaccountIds:  make([]satypes.SubaccountId, 0),
		subaccountsWithPositions:  make(map[uint32]*clobtypes.SubaccountOpenPositionInfo),
	}
}

// UpdateBlockHeight updates the struct with the given block height.
func (ls *DaemonLiquidationInfo) UpdateBlockHeight(blockHeight uint32) {
	ls.Lock()
	defer ls.Unlock()
	ls.blockHeight = blockHeight
}

// GetBlockHeight returns the block height of the last update.
func (ls *DaemonLiquidationInfo) GetBlockHeight() uint32 {
	ls.Lock()
	defer ls.Unlock()
	return ls.blockHeight
}

// UpdateLiquidatableSubaccountIds updates the struct with the given a list of potentially
// liquidatable subaccount ids.
func (ls *DaemonLiquidationInfo) UpdateLiquidatableSubaccountIds(updates []satypes.SubaccountId) {
	ls.Lock()
	defer ls.Unlock()
	ls.liquidatableSubaccountIds = make([]satypes.SubaccountId, len(updates))
	copy(ls.liquidatableSubaccountIds, updates)
}

// GetLiquidatableSubaccountIds returns the list of potentially liquidatable subaccount ids
// reported by the liquidation daemon.
func (ls *DaemonLiquidationInfo) GetLiquidatableSubaccountIds() []satypes.SubaccountId {
	ls.Lock()
	defer ls.Unlock()
	results := make([]satypes.SubaccountId, len(ls.liquidatableSubaccountIds))
	copy(results, ls.liquidatableSubaccountIds)
	return results
}

// UpdateNegativeTncSubaccountIds updates the struct with the given a list of subaccount ids
// with negative total net collateral.
func (ls *DaemonLiquidationInfo) UpdateNegativeTncSubaccountIds(updates []satypes.SubaccountId) {
	ls.Lock()
	defer ls.Unlock()
	ls.negativeTncSubaccountIds = make([]satypes.SubaccountId, len(updates))
	copy(ls.negativeTncSubaccountIds, updates)
}

// GetNegativeTncSubaccountIds returns the list of subaccount ids with negative total net collateral
// reported by the liquidation daemon.
func (ls *DaemonLiquidationInfo) GetNegativeTncSubaccountIds() []satypes.SubaccountId {
	ls.Lock()
	defer ls.Unlock()
	results := make([]satypes.SubaccountId, len(ls.negativeTncSubaccountIds))
	copy(results, ls.negativeTncSubaccountIds)
	return results
}

// UpdateSubaccountsWithPositions updates the struct with the given a list of subaccount ids with open positions.
func (ls *DaemonLiquidationInfo) UpdateSubaccountsWithPositions(
	subaccountsWithPositions map[uint32]*clobtypes.SubaccountOpenPositionInfo,
) {
	ls.Lock()
	defer ls.Unlock()
	ls.subaccountsWithPositions = make(map[uint32]*clobtypes.SubaccountOpenPositionInfo)
	for perpetualId, info := range subaccountsWithPositions {
		clone := &clobtypes.SubaccountOpenPositionInfo{
			PerpetualId:                  perpetualId,
			SubaccountsWithLongPosition:  make([]satypes.SubaccountId, len(info.SubaccountsWithLongPosition)),
			SubaccountsWithShortPosition: make([]satypes.SubaccountId, len(info.SubaccountsWithShortPosition)),
		}
		copy(clone.SubaccountsWithLongPosition, info.SubaccountsWithLongPosition)
		copy(clone.SubaccountsWithShortPosition, info.SubaccountsWithShortPosition)
		ls.subaccountsWithPositions[perpetualId] = clone
	}
}

// GetSubaccountsWithPositions returns the list of subaccount ids with open positions.
func (ls *DaemonLiquidationInfo) GetSubaccountsWithPositions() map[uint32]*clobtypes.SubaccountOpenPositionInfo {
	ls.Lock()
	defer ls.Unlock()

	result := make(map[uint32]*clobtypes.SubaccountOpenPositionInfo)
	for perpetualId, info := range ls.subaccountsWithPositions {
		clone := &clobtypes.SubaccountOpenPositionInfo{
			PerpetualId:                  perpetualId,
			SubaccountsWithLongPosition:  make([]satypes.SubaccountId, len(info.SubaccountsWithLongPosition)),
			SubaccountsWithShortPosition: make([]satypes.SubaccountId, len(info.SubaccountsWithShortPosition)),
		}
		copy(clone.SubaccountsWithLongPosition, info.SubaccountsWithLongPosition)
		copy(clone.SubaccountsWithShortPosition, info.SubaccountsWithShortPosition)
		result[perpetualId] = clone
	}
	return result
}
