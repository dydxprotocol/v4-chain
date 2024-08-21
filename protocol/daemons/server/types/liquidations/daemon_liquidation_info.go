package types

import (
	"fmt"
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

func (ls *DaemonLiquidationInfo) Update(
	blockHeight uint32,
	liquidatableSubaccountIds []satypes.SubaccountId,
	negativeTncSubaccountIds []satypes.SubaccountId,
	subaccountsWithPositions []clobtypes.SubaccountOpenPositionInfo,
) {
	ls.Lock()
	defer ls.Unlock()

	if blockHeight > ls.blockHeight {
		ls.liquidatableSubaccountIds = make([]satypes.SubaccountId, 0)
		ls.negativeTncSubaccountIds = make([]satypes.SubaccountId, 0)
		ls.subaccountsWithPositions = make(map[uint32]*clobtypes.SubaccountOpenPositionInfo)
	} else if blockHeight < ls.blockHeight {
		panic(
			fmt.Sprintf(
				"UpdateLiquidatableSubaccountIds: block height %d cannot be less than the current block height %d",
				blockHeight,
				ls.blockHeight,
			),
		)
	}
	ls.UpdateBlockHeight(blockHeight)

	ls.UpdateLiquidatableSubaccountIds(liquidatableSubaccountIds, blockHeight)
	ls.UpdateNegativeTncSubaccountIds(negativeTncSubaccountIds, blockHeight)
	ls.UpdateSubaccountsWithPositions(subaccountsWithPositions, blockHeight)
}

// UpdateBlockHeight updates the struct with the given block height.
func (ls *DaemonLiquidationInfo) UpdateBlockHeight(blockHeight uint32) {
	ls.blockHeight = blockHeight
}

// UpdateLiquidatableSubaccountIds updates the struct with the given a list of potentially
// liquidatable subaccount ids.
func (ls *DaemonLiquidationInfo) UpdateLiquidatableSubaccountIds(
	updates []satypes.SubaccountId,
	blockHeight uint32,
) {
	ls.liquidatableSubaccountIds = append(ls.liquidatableSubaccountIds, updates...)
}

// UpdateNegativeTncSubaccountIds updates the struct with the given a list of subaccount ids
// with negative total net collateral.
func (ls *DaemonLiquidationInfo) UpdateNegativeTncSubaccountIds(
	updates []satypes.SubaccountId,
	blockHeight uint32,
) {
	ls.negativeTncSubaccountIds = append(ls.negativeTncSubaccountIds, updates...)
}

// UpdateSubaccountsWithPositions updates the struct with the given a list of subaccount ids with open positions.
func (ls *DaemonLiquidationInfo) UpdateSubaccountsWithPositions(
	subaccountsWithPositions []clobtypes.SubaccountOpenPositionInfo,
	blockHeight uint32,
) {
	// Append to the current map if the block height not changed.
	for _, info := range subaccountsWithPositions {
		if _, ok := ls.subaccountsWithPositions[info.PerpetualId]; !ok {
			ls.subaccountsWithPositions[info.PerpetualId] = &clobtypes.SubaccountOpenPositionInfo{
				PerpetualId:                  info.PerpetualId,
				SubaccountsWithLongPosition:  make([]satypes.SubaccountId, 0),
				SubaccountsWithShortPosition: make([]satypes.SubaccountId, 0),
			}
		}
		ls.subaccountsWithPositions[info.PerpetualId].SubaccountsWithLongPosition = append(
			ls.subaccountsWithPositions[info.PerpetualId].SubaccountsWithLongPosition,
			info.SubaccountsWithLongPosition...,
		)
		ls.subaccountsWithPositions[info.PerpetualId].SubaccountsWithShortPosition = append(
			ls.subaccountsWithPositions[info.PerpetualId].SubaccountsWithShortPosition,
			info.SubaccountsWithShortPosition...,
		)
	}
}

// GetBlockHeight returns the block height of the last update.
func (ls *DaemonLiquidationInfo) GetBlockHeight() uint32 {
	ls.Lock()
	defer ls.Unlock()
	return ls.blockHeight
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

// GetNegativeTncSubaccountIds returns the list of subaccount ids with negative total net collateral
// reported by the liquidation daemon.
func (ls *DaemonLiquidationInfo) GetNegativeTncSubaccountIds() []satypes.SubaccountId {
	ls.Lock()
	defer ls.Unlock()
	results := make([]satypes.SubaccountId, len(ls.negativeTncSubaccountIds))
	copy(results, ls.negativeTncSubaccountIds)
	return results
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
