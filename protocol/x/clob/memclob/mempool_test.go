package memclob_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestOrderPool(t *testing.T) {
	pool := memclob.NewOrderPool()

	// Get an order from the pool
	order := pool.Get()
	require.NotNil(t, order)

	// Set some values on the order
	order.Side = types.Order_SIDE_BUY
	order.Quantums = 1000
	order.Subticks = 5000
	order.ReduceOnly = true
	order.ClientMetadata = 42

	// Return the order to the pool
	pool.Put(order)

	// Get another order from the pool
	order2 := pool.Get()
	require.NotNil(t, order2)

	// Verify that the fields were reset to zero values
	require.Equal(t, types.Order_SIDE_UNSPECIFIED, order2.Side)
	require.Equal(t, uint64(0), order2.Quantums)
	require.Equal(t, uint64(0), order2.Subticks)
	require.Equal(t, false, order2.ReduceOnly)
	require.Equal(t, uint32(0), order2.ClientMetadata)
}

func TestClobOrderPool(t *testing.T) {
	pool := memclob.NewClobOrderPool()

	// Get a ClobOrder from the pool
	clobOrder := pool.Get()
	require.NotNil(t, clobOrder)

	// Set some values on the ClobOrder
	clobOrder.Order.Side = types.Order_SIDE_BUY
	clobOrder.Order.Quantums = 1000
	clobOrder.Signature = []byte("test signature")

	// Return the ClobOrder to the pool
	pool.Put(clobOrder)

	// Get another ClobOrder from the pool
	clobOrder2 := pool.Get()
	require.NotNil(t, clobOrder2)

	// Verify that the fields were reset to zero values
	require.Equal(t, types.Order{}, clobOrder2.Order)
	require.Nil(t, clobOrder2.Signature)
}

func TestLevelOrderPool(t *testing.T) {
	pool := memclob.NewLevelOrderPool()

	// Get a level order from the pool
	levelOrder := pool.Get()
	require.NotNil(t, levelOrder)

	// Set some values on the level order
	levelOrder.Value.Order.Side = types.Order_SIDE_BUY
	levelOrder.Value.Order.Quantums = 1000
	levelOrder.Value.Signature = []byte("test signature")

	// Set up a linked list structure
	nextLevelOrder := pool.Get()
	nextLevelOrder.Value.Order.Side = types.Order_SIDE_SELL
	levelOrder.Next = nextLevelOrder
	nextLevelOrder.Prev = levelOrder

	// Return the level order to the pool
	pool.Put(levelOrder)

	// Get another level order from the pool
	levelOrder2 := pool.Get()
	require.NotNil(t, levelOrder2)

	// Verify that the fields were reset to zero values
	require.Equal(t, types.ClobOrder{}, levelOrder2.Value)
	require.Nil(t, levelOrder2.Next)
	require.Nil(t, levelOrder2.Prev)
}

func TestMakerFillPool(t *testing.T) {
	pool := memclob.NewMakerFillPool()

	// Get a maker fill from the pool
	makerFill := pool.Get()
	require.NotNil(t, makerFill)

	// Set some values on the maker fill
	makerFill.FillAmount = 1000
	makerFill.MakerOrderId.ClientId = 42

	// Return the maker fill to the pool
	pool.Put(makerFill)

	// Get another maker fill from the pool
	makerFill2 := pool.Get()
	require.NotNil(t, makerFill2)

	// Verify that the fields were reset to zero values
	require.Equal(t, uint64(0), makerFill2.FillAmount)
	require.Equal(t, types.OrderId{}, makerFill2.MakerOrderId)
}

func TestMakerFillWithOrderPool(t *testing.T) {
	pool := memclob.NewMakerFillWithOrderPool()

	// Get a MakerFillWithOrder from the pool
	makerFillWithOrder := pool.Get()
	require.NotNil(t, makerFillWithOrder)

	// Set some values on the MakerFillWithOrder
	makerFillWithOrder.MakerFill.FillAmount = 1000
	makerFillWithOrder.Order.Side = types.Order_SIDE_BUY

	// Return the MakerFillWithOrder to the pool
	pool.Put(makerFillWithOrder)

	// Get another MakerFillWithOrder from the pool
	makerFillWithOrder2 := pool.Get()
	require.NotNil(t, makerFillWithOrder2)

	// Verify that the fields were reset to zero values
	require.Equal(t, types.MakerFill{}, makerFillWithOrder2.MakerFill)
	require.Equal(t, types.Order{}, makerFillWithOrder2.Order)
}

func TestGlobalMemPools(t *testing.T) {
	// Verify that the global pools are initialized
	require.NotNil(t, memclob.GlobalMemPools.OrderPool)
	require.NotNil(t, memclob.GlobalMemPools.ClobOrderPool)
	require.NotNil(t, memclob.GlobalMemPools.LevelOrderPool)
	require.NotNil(t, memclob.GlobalMemPools.MakerFillPool)
	require.NotNil(t, memclob.GlobalMemPools.MakerFillWithOrderPool)

	// Verify that we can get objects from the global pools
	order := memclob.GlobalMemPools.OrderPool.Get()
	require.NotNil(t, order)

	clobOrder := memclob.GlobalMemPools.ClobOrderPool.Get()
	require.NotNil(t, clobOrder)

	levelOrder := memclob.GlobalMemPools.LevelOrderPool.Get()
	require.NotNil(t, levelOrder)

	makerFill := memclob.GlobalMemPools.MakerFillPool.Get()
	require.NotNil(t, makerFill)

	makerFillWithOrder := memclob.GlobalMemPools.MakerFillWithOrderPool.Get()
	require.NotNil(t, makerFillWithOrder)
}

// BenchmarkOrderCreationWithPool measures the performance of creating and reusing Order objects with a memory pool
func BenchmarkOrderCreationWithPool(b *testing.B) {
	pool := memclob.NewOrderPool()
	
	b.ResetTimer()
	
	// Create a slice to store references so they don't get garbage collected during the benchmark
	orders := make([]*types.Order, 0, 1000)
	
	for i := 0; i < b.N; i++ {
		// Get a new order from the pool
		order := pool.Get()
		order.Side = types.Order_SIDE_BUY
		order.Quantums = 1000
		order.Subticks = 5000
		order.ReduceOnly = true
		order.OrderId.ClientId = uint32(i)
		
		// Every 1000 iterations, return all orders to the pool
		if i%1000 == 999 {
			for _, o := range orders {
				pool.Put(o)
			}
			orders = orders[:0] // Clear the slice
		} else {
			orders = append(orders, order)
		}
	}
}

// BenchmarkOrderCreationWithoutPool measures the performance of creating Order objects without a memory pool
func BenchmarkOrderCreationWithoutPool(b *testing.B) {
	b.ResetTimer()
	
	// Create a slice to store references so they don't get garbage collected during the benchmark
	orders := make([]*types.Order, 0, 1000)
	
	for i := 0; i < b.N; i++ {
		// Create a new order directly
		order := new(types.Order)
		order.Side = types.Order_SIDE_BUY
		order.Quantums = 1000
		order.Subticks = 5000
		order.ReduceOnly = true
		order.OrderId.ClientId = uint32(i)
		
		// Every 1000 iterations, clear the references
		if i%1000 == 999 {
			orders = orders[:0] // Clear the slice
		} else {
			orders = append(orders, order)
		}
	}
}

// BenchmarkMakerFillCreationWithPool measures the performance of creating and reusing MakerFill objects with a memory pool
func BenchmarkMakerFillCreationWithPool(b *testing.B) {
	pool := memclob.NewMakerFillPool()
	
	b.ResetTimer()
	
	// Create a slice to store references so they don't get garbage collected during the benchmark
	fills := make([]*types.MakerFill, 0, 1000)
	
	for i := 0; i < b.N; i++ {
		// Get a new maker fill from the pool
		makerFill := pool.Get()
		makerFill.FillAmount = 1000
		makerFill.MakerOrderId.ClientId = uint32(i)
		
		// Every 1000 iterations, return all maker fills to the pool
		if i%1000 == 999 {
			for _, fill := range fills {
				pool.Put(fill)
			}
			fills = fills[:0] // Clear the slice
		} else {
			fills = append(fills, makerFill)
		}
	}
}

// BenchmarkMakerFillCreationWithoutPool measures the performance of creating MakerFill objects without a memory pool
func BenchmarkMakerFillCreationWithoutPool(b *testing.B) {
	b.ResetTimer()
	
	// Create a slice to store references so they don't get garbage collected during the benchmark
	fills := make([]*types.MakerFill, 0, 1000)
	
	for i := 0; i < b.N; i++ {
		// Create a new maker fill directly
		makerFill := new(types.MakerFill)
		makerFill.FillAmount = 1000
		makerFill.MakerOrderId.ClientId = uint32(i)
		
		// Every 1000 iterations, clear the references
		if i%1000 == 999 {
			fills = fills[:0] // Clear the slice
		} else {
			fills = append(fills, makerFill)
		}
	}
} 