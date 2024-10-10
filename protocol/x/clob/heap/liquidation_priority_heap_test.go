package heap

import (
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/assert"
)

func TestLiquidationPriorityHeap(t *testing.T) {
	t.Run("NewLiquidationPriorityHeap", func(t *testing.T) {
		h := NewLiquidationPriorityHeap()
		assert.NotNil(t, h)
		assert.Equal(t, 0, h.Len())
	})

	t.Run("AddSubaccount", func(t *testing.T) {
		h := NewLiquidationPriorityHeap()
		h.AddSubaccount(types.SubaccountId{Owner: "owner1", Number: 1}, big.NewFloat(1.0))
		assert.Equal(t, 1, h.Len())
	})

	t.Run("PopLowestPriority", func(t *testing.T) {
		h := NewLiquidationPriorityHeap()
		h.AddSubaccount(types.SubaccountId{Owner: "owner1", Number: 1}, big.NewFloat(2.0))
		h.AddSubaccount(types.SubaccountId{Owner: "owner2", Number: 2}, big.NewFloat(1.0))

		lowest := h.PopLowestPriority()
		assert.NotNil(t, lowest)
		assert.Equal(t, "owner2", lowest.SubaccountId.Owner)
		assert.Equal(t, uint32(2), lowest.SubaccountId.Number)
		assert.Equal(t, big.NewFloat(1.0), lowest.Priority)
		assert.Equal(t, 1, h.Len())
	})

	t.Run("UpdatePriority", func(t *testing.T) {
		h := NewLiquidationPriorityHeap()
		h.AddSubaccount(types.SubaccountId{Owner: "owner1", Number: 1}, big.NewFloat(2.0))
		h.AddSubaccount(types.SubaccountId{Owner: "owner2", Number: 2}, big.NewFloat(1.0))

		item := (*h)[1]
		success := h.UpdatePriority(item, big.NewFloat(0.5))
		assert.True(t, success)

		lowest := h.PopLowestPriority()
		assert.NotNil(t, lowest)
		assert.Equal(t, "owner1", lowest.SubaccountId.Owner)
		assert.Equal(t, uint32(1), lowest.SubaccountId.Number)
		assert.Equal(t, big.NewFloat(0.5), lowest.Priority)
	})

	t.Run("UpdatePriority_InvalidItem", func(t *testing.T) {
		h := NewLiquidationPriorityHeap()
		h.AddSubaccount(types.SubaccountId{Owner: "owner1", Number: 1}, big.NewFloat(1.0))

		invalidItem := &LiquidationPriority{
			SubaccountId: types.SubaccountId{Owner: "invalid", Number: 999},
			Priority:     big.NewFloat(1.0),
			Index:        -1,
		}
		success := h.UpdatePriority(invalidItem, big.NewFloat(0.5))
		assert.False(t, success)
	})

	t.Run("Ordering", func(t *testing.T) {
		h := NewLiquidationPriorityHeap()
		h.AddSubaccount(types.SubaccountId{Owner: "owner1", Number: 1}, big.NewFloat(3.0))
		h.AddSubaccount(types.SubaccountId{Owner: "owner2", Number: 2}, big.NewFloat(1.0))
		h.AddSubaccount(types.SubaccountId{Owner: "owner3", Number: 3}, big.NewFloat(2.0))

		expected := []string{"owner2", "owner3", "owner1"}
		for i := 0; i < 3; i++ {
			item := h.PopLowestPriority()
			assert.Equal(t, expected[i], item.SubaccountId.Owner)
		}
		assert.Equal(t, 0, h.Len())
	})
}
