package memclob

import (
	"sync"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// OrderPool is a memory pool for Order objects to reduce GC pressure
type OrderPool struct {
	pool sync.Pool
}

// NewOrderPool creates a new memory pool for Order objects
func NewOrderPool() *OrderPool {
	return &OrderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &types.Order{}
			},
		},
	}
}

// Get retrieves an Order object from the pool
func (p *OrderPool) Get() *types.Order {
	return p.pool.Get().(*types.Order)
}

// Put returns an Order object to the pool after resetting its fields
func (p *OrderPool) Put(order *types.Order) {
	if order == nil {
		return
	}
	
	// Reset all fields to zero values to prevent memory leaks
	order.OrderId = types.OrderId{}
	order.Side = types.Order_SIDE_UNSPECIFIED
	order.Quantums = 0
	order.Subticks = 0
	order.GoodTilOneof = nil
	order.TimeInForce = types.Order_TIME_IN_FORCE_UNSPECIFIED
	order.ReduceOnly = false
	order.ClientMetadata = 0
	order.ConditionType = types.Order_CONDITION_TYPE_UNSPECIFIED
	order.ConditionalOrderTriggerSubticks = 0
	order.TwapParameters = nil

	p.pool.Put(order)
}

// ClobOrderPool is a memory pool for ClobOrder objects
type ClobOrderPool struct {
	pool sync.Pool
}

// NewClobOrderPool creates a new memory pool for ClobOrder objects
func NewClobOrderPool() *ClobOrderPool {
	return &ClobOrderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &types.ClobOrder{}
			},
		},
	}
}

// Get retrieves a ClobOrder object from the pool
func (p *ClobOrderPool) Get() *types.ClobOrder {
	return p.pool.Get().(*types.ClobOrder)
}

// Put returns a ClobOrder object to the pool after resetting its fields
func (p *ClobOrderPool) Put(clobOrder *types.ClobOrder) {
	if clobOrder == nil {
		return
	}
	
	// Reset all fields to zero values
	clobOrder.Order = types.Order{}
	clobOrder.Signature = nil

	p.pool.Put(clobOrder)
}

// LevelOrderPool is a memory pool for LevelOrder objects
// Since LevelOrder is a type alias for list.Node[ClobOrder], we need to handle this carefully
type LevelOrderPool struct {
	pool sync.Pool
}

// NewLevelOrderPool creates a new memory pool for LevelOrder objects
func NewLevelOrderPool() *LevelOrderPool {
	return &LevelOrderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(types.LevelOrder)
			},
		},
	}
}

// Get retrieves a LevelOrder object from the pool
func (p *LevelOrderPool) Get() *types.LevelOrder {
	return p.pool.Get().(*types.LevelOrder)
}

// Put returns a LevelOrder object to the pool after resetting its fields
func (p *LevelOrderPool) Put(levelOrder *types.LevelOrder) {
	if levelOrder == nil {
		return
	}
	
	// Reset fields - since LevelOrder is a list.Node[ClobOrder], we reset its fields
	// including Value (ClobOrder), Next and Prev pointers
	levelOrder.Value = types.ClobOrder{}
	levelOrder.Next = nil
	levelOrder.Prev = nil

	p.pool.Put(levelOrder)
}

// MakerFillPool is a memory pool for MakerFill objects
type MakerFillPool struct {
	pool sync.Pool
}

// NewMakerFillPool creates a new memory pool for MakerFill objects
func NewMakerFillPool() *MakerFillPool {
	return &MakerFillPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &types.MakerFill{}
			},
		},
	}
}

// Get retrieves a MakerFill object from the pool
func (p *MakerFillPool) Get() *types.MakerFill {
	return p.pool.Get().(*types.MakerFill)
}

// Put returns a MakerFill object to the pool after resetting its fields
func (p *MakerFillPool) Put(makerFill *types.MakerFill) {
	if makerFill == nil {
		return
	}
	
	// Reset all fields to zero values
	makerFill.FillAmount = 0
	makerFill.MakerOrderId = types.OrderId{}

	p.pool.Put(makerFill)
}

// MakerFillWithOrderPool is a memory pool for MakerFillWithOrder objects
type MakerFillWithOrderPool struct {
	pool sync.Pool
}

// NewMakerFillWithOrderPool creates a new memory pool for MakerFillWithOrder objects
func NewMakerFillWithOrderPool() *MakerFillWithOrderPool {
	return &MakerFillWithOrderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &types.MakerFillWithOrder{}
			},
		},
	}
}

// Get retrieves a MakerFillWithOrder object from the pool
func (p *MakerFillWithOrderPool) Get() *types.MakerFillWithOrder {
	return p.pool.Get().(*types.MakerFillWithOrder)
}

// Put returns a MakerFillWithOrder object to the pool after resetting its fields
func (p *MakerFillWithOrderPool) Put(makerFillWithOrder *types.MakerFillWithOrder) {
	if makerFillWithOrder == nil {
		return
	}
	
	// Reset all fields to zero values
	makerFillWithOrder.MakerFill = types.MakerFill{}
	makerFillWithOrder.Order = types.Order{}

	p.pool.Put(makerFillWithOrder)
}

// GlobalMemPools contains the global instance of memory pools
var GlobalMemPools = struct {
	OrderPool            *OrderPool
	ClobOrderPool        *ClobOrderPool
	LevelOrderPool       *LevelOrderPool
	MakerFillPool        *MakerFillPool
	MakerFillWithOrderPool *MakerFillWithOrderPool
}{
	OrderPool:            NewOrderPool(),
	ClobOrderPool:        NewClobOrderPool(),
	LevelOrderPool:       NewLevelOrderPool(),
	MakerFillPool:        NewMakerFillPool(),
	MakerFillWithOrderPool: NewMakerFillWithOrderPool(),
} 