package types

import (
	"sync"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// LiquidatableSubaccountIds maintains the list of subaccount ids to be liquidated
// in the next block. Methods are goroutine safe.
type LiquidatableSubaccountIds struct {
	sync.Mutex                           // lock
	subaccountIds []satypes.SubaccountId // liquidatable subaccount ids
}

// NewLiquidatableSubaccountIds creates a new `LiquidatableSubaccountIds` struct.
func NewLiquidatableSubaccountIds() *LiquidatableSubaccountIds {
	return &LiquidatableSubaccountIds{
		subaccountIds: make([]satypes.SubaccountId, 0),
	}
}

// UpdateSubaccountIds updates the struct with the given a list of potentially liquidatable subaccount ids.
func (ls *LiquidatableSubaccountIds) UpdateSubaccountIds(updates []satypes.SubaccountId) {
	ls.Lock()
	defer ls.Unlock()
	ls.subaccountIds = make([]satypes.SubaccountId, len(updates))
	copy(ls.subaccountIds, updates)
}

// GetSubaccountIds returns the list of potentially liquidatable subaccount ids
// reported by the liquidation daemon.
func (ls *LiquidatableSubaccountIds) GetSubaccountIds() []satypes.SubaccountId {
	ls.Lock()
	defer ls.Unlock()
	results := make([]satypes.SubaccountId, len(ls.subaccountIds))
	copy(results, ls.subaccountIds)
	return results
}
