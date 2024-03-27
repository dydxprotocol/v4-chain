package client

import (
	"context"
	"sync"

	"cosmossdk.io/log"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// Example client to consume data from a gRPC server.
type GrpcClient struct {
	Logger    log.Logger
	Orderbook map[uint32]*LocalOrderbook
}

type LocalOrderbook struct {
	sync.Mutex

	OrderIdToOrder       map[v1types.IndexerOrderId]v1types.IndexerOrder
	OrderRemainingAmount map[v1types.IndexerOrderId]uint64
	Bids                 map[uint64][]v1types.IndexerOrder
	Asks                 map[uint64][]v1types.IndexerOrder

	Logger log.Logger
}

func NewGrpcClient(appflags appflags.Flags, logger log.Logger) *GrpcClient {
	logger = logger.With("module", "grpc-example-client")

	client := &GrpcClient{
		Logger:    logger,
		Orderbook: make(map[uint32]*LocalOrderbook),
	}

	// Subscribe to grpc orderbook updates.
	go func() {
		grpcClient := daemontypes.GrpcClientImpl{}

		// Make a connection to the Cosmos gRPC query services.
		queryConn, err := grpcClient.NewTcpConnection(context.Background(), appflags.GrpcAddress)
		if err != nil {
			logger.Error("Failed to establish gRPC connection to Cosmos gRPC query services", "error", err)
			return
		}
		defer func() {
			if err := grpcClient.CloseConnection(queryConn); err != nil {
				logger.Error("Failed to close gRPC connection", "error", err)
			}
		}()

		clobQueryClient := clobtypes.NewQueryClient(queryConn)
		updateClient, err := clobQueryClient.StreamOrderbookUpdates(
			context.Background(),
			&clobtypes.StreamOrderbookUpdatesRequest{
				ClobPairId: []uint32{0, 1},
			},
		)
		if err != nil {
			logger.Error("Failed to stream orderbook updates", "error", err)
			return
		}

		for {
			update, err := updateClient.Recv()
			if err != nil {
				logger.Error("Failed to receive orderbook update", "error", err)
				return
			}

			logger.Info("Received orderbook update", "update", update)
			client.Update(update)
		}
	}()
	return client
}

// Read method
func (c *GrpcClient) GetOrderbookSnapshot(pairId uint32) *LocalOrderbook {
	return c.GetOrderbook(pairId)
}

// Write method
func (c *GrpcClient) Update(updates *clobtypes.StreamOrderbookUpdatesResponse) {
	if updates.Snapshot {
		c.Logger.Info("Received orderbook snapshot")
		c.Orderbook = make(map[uint32]*LocalOrderbook)
	}

	for _, update := range updates.Updates {
		if orderPlace := update.GetOrderPlace(); orderPlace != nil {
			order := orderPlace.GetOrder()
			orderbook := c.GetOrderbook(order.OrderId.ClobPairId)
			orderbook.AddOrder(*order)
		}

		if orderRemove := update.GetOrderRemove(); orderRemove != nil {
			orderId := orderRemove.RemovedOrderId
			orderbook := c.GetOrderbook(orderId.ClobPairId)
			orderbook.RemoveOrder(*orderId)
		}

		if orderUpdate := update.GetOrderUpdate(); orderUpdate != nil {
			orderId := orderUpdate.OrderId
			orderbook := c.GetOrderbook(orderId.ClobPairId)
			orderbook.SetOrderRemainingAmount(*orderId, orderUpdate.TotalFilledQuantums)
		}
	}
}

func (c *GrpcClient) GetOrderbook(pairId uint32) *LocalOrderbook {
	if _, ok := c.Orderbook[pairId]; !ok {
		c.Orderbook[pairId] = &LocalOrderbook{
			OrderIdToOrder:       make(map[v1types.IndexerOrderId]v1types.IndexerOrder),
			OrderRemainingAmount: make(map[v1types.IndexerOrderId]uint64),
			Bids:                 make(map[uint64][]v1types.IndexerOrder),
			Asks:                 make(map[uint64][]v1types.IndexerOrder),

			Logger: c.Logger,
		}
	}
	return c.Orderbook[pairId]
}

func (l *LocalOrderbook) SetOrderRemainingAmount(orderId v1types.IndexerOrderId, totalFilledQuantums uint64) {
	l.Lock()
	defer l.Unlock()

	order := l.OrderIdToOrder[orderId]
	if totalFilledQuantums > order.Quantums {
		l.Logger.Error("totalFilledQuantums > order.Quantums")
	}
	l.OrderRemainingAmount[orderId] = order.Quantums - totalFilledQuantums
}

func (l *LocalOrderbook) AddOrder(order v1types.IndexerOrder) {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.OrderIdToOrder[order.OrderId]; ok {
		l.Logger.Error("order already exists in orderbook")
	}

	subticks := order.GetSubticks()
	if order.Side == v1types.IndexerOrder_SIDE_BUY {
		if _, ok := l.Bids[subticks]; !ok {
			l.Bids[subticks] = make([]v1types.IndexerOrder, 0)
		}
		l.Bids[subticks] = append(
			l.Bids[subticks],
			order,
		)
	} else {
		if _, ok := l.Asks[subticks]; !ok {
			l.Asks[subticks] = make([]v1types.IndexerOrder, 0)
		}
		l.Asks[subticks] = append(
			l.Asks[subticks],
			order,
		)
	}

	l.OrderIdToOrder[order.OrderId] = order
	l.OrderRemainingAmount[order.OrderId] = 0
}

func (l *LocalOrderbook) RemoveOrder(orderId v1types.IndexerOrderId) {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.OrderIdToOrder[orderId]; !ok {
		l.Logger.Error("order not found in orderbook")
	}

	order := l.OrderIdToOrder[orderId]
	subticks := order.GetSubticks()

	if order.Side == v1types.IndexerOrder_SIDE_BUY {
		for i, o := range l.Bids[subticks] {
			if o.OrderId == order.OrderId {
				l.Bids[subticks] = append(
					l.Bids[subticks][:i],
					l.Bids[subticks][i+1:]...,
				)
				break
			}
		}
		if len(l.Bids[subticks]) == 0 {
			delete(l.Bids, subticks)
		}
	} else {
		for i, o := range l.Asks[subticks] {
			if o.OrderId == order.OrderId {
				l.Asks[subticks] = append(
					l.Asks[subticks][:i],
					l.Asks[subticks][i+1:]...,
				)
				break
			}
		}
		if len(l.Asks[subticks]) == 0 {
			delete(l.Asks, subticks)
		}
	}

	delete(l.OrderRemainingAmount, orderId)
	delete(l.OrderIdToOrder, orderId)
}
