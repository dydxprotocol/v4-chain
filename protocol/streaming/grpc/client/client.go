package client

import (
	"context"
	"sync"

	"cosmossdk.io/log"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// Example client to consume data from a gRPC server.
type GrpcClient struct {
	Logger      log.Logger
	Orderbook   map[uint32]*LocalOrderbook
	initialized bool
}

type LocalOrderbook struct {
	sync.Mutex

	OrderIdToOrder map[v1types.IndexerOrderId]v1types.IndexerOrder
	// todo remove this since we have fills
	OrderRemainingAmount map[v1types.IndexerOrderId]uint64
	Bids                 map[uint64][]v1types.IndexerOrder
	Asks                 map[uint64][]v1types.IndexerOrder
	FillAmounts          map[v1types.IndexerOrderId]uint64

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

// Write method for stream orderbook updates.
func (c *GrpcClient) Update(updates *clobtypes.StreamOrderbookUpdatesResponse) {
	for _, update := range updates.GetUpdates() {
		if orderPlace := update.GetOrderPlace(); orderPlace != nil {
			c.ProcessOrderPlace(orderPlace)
		}
		if orderFill := update.GetOrderFill(); orderFill != nil {
			c.ProcessFill(orderFill)
		}
	}
}

// Write method for order placement updates (indexer offchain events)
func (c *GrpcClient) ProcessOrderPlace(orderPlace *clobtypes.StreamOrderbookUpdatesOrderPlacement) {
	if orderPlace.Snapshot {
		c.Logger.Info("Received orderbook snapshot")
		c.Orderbook = make(map[uint32]*LocalOrderbook)
	}

	for _, update := range orderPlace.Updates {
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

			// Order updates contain the absolute amount of filled quantums.
			orderbook.SetOrderAbsoluteFillAmount(orderId, orderUpdate.TotalFilledQuantums)
		}
	}

	if orderPlace.Snapshot {
		c.initialized = true
	}
}

// Write method for orderbook fills update (clob match)
func (c *GrpcClient) ProcessFill(orderFills *clobtypes.StreamOrderbookUpdatesOrderFills) {
	// Need to wait for the snapshot message first
	if !c.initialized {
		return
	}

	for _, fill := range orderFills.GetFills() {
		orderMap := orderListToMap(fill.Orders)
		clobMatch := fill.ClobMatch

		if matchOrders := clobMatch.GetMatchOrders(); matchOrders != nil {
			c.ProcessMatchOrders(matchOrders, orderMap)
		}

		if matchPerpLiquidation := clobMatch.GetMatchPerpetualLiquidation(); matchPerpLiquidation != nil {
			c.ProcessMatchPerpetualLiquidation(matchPerpLiquidation, orderMap)
		}
	}
}

func (c *GrpcClient) ProcessMatchPerpetualLiquidation(
	perpLiquidation *clobtypes.MatchPerpetualLiquidation,
	orderMap map[clobtypes.OrderId]clobtypes.Order,
) {
	localOrderbook := c.Orderbook[perpLiquidation.ClobPairId]
	for _, fill := range perpLiquidation.GetFills() {
		makerOrder := orderMap[fill.MakerOrderId]
		deltaFillAmount := fill.FillAmount
		localOrderbook.AdjustOrderDeltaFillAmount(makerOrder, deltaFillAmount)
	}
}

func (c *GrpcClient) ProcessMatchOrders(
	matchOrders *clobtypes.MatchOrders,
	orderMap map[clobtypes.OrderId]clobtypes.Order,
) {
	takerOrderId := matchOrders.TakerOrderId
	clobPairId := takerOrderId.GetClobPairId()
	localOrderbook := c.Orderbook[clobPairId]

	takerOrder := orderMap[takerOrderId]

	for _, fill := range matchOrders.Fills {
		c.Logger.Info("recieved an order fills", fill)
		makerOrder := orderMap[fill.MakerOrderId]
		deltaFillAmount := fill.FillAmount
		localOrderbook.AdjustOrderDeltaFillAmount(makerOrder, deltaFillAmount)
		localOrderbook.AdjustOrderDeltaFillAmount(takerOrder, deltaFillAmount)
	}
}

// orderListToMap converts a list of orders to a map of OrderId to Order.
func orderListToMap(orders []clobtypes.Order) map[clobtypes.OrderId]clobtypes.Order {
	orderMap := make(map[clobtypes.OrderId]clobtypes.Order)
	for _, order := range orders {
		orderMap[order.OrderId] = order
	}
	return orderMap
}

func (c *GrpcClient) GetOrderbook(pairId uint32) *LocalOrderbook {
	if _, ok := c.Orderbook[pairId]; !ok {
		c.Orderbook[pairId] = &LocalOrderbook{
			OrderIdToOrder:       make(map[v1types.IndexerOrderId]v1types.IndexerOrder),
			OrderRemainingAmount: make(map[v1types.IndexerOrderId]uint64),
			Bids:                 make(map[uint64][]v1types.IndexerOrder),
			Asks:                 make(map[uint64][]v1types.IndexerOrder),
			FillAmounts:          make(map[v1types.IndexerOrderId]uint64),

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

func (l *LocalOrderbook) SetOrderAbsoluteFillAmount(
	orderId *v1types.IndexerOrderId,
	absoluteFillAmount uint64,
) {
	l.Lock()
	defer l.Unlock()

	l.FillAmounts[*orderId] = absoluteFillAmount
}

func (l *LocalOrderbook) AdjustOrderDeltaFillAmount(
	order clobtypes.Order,
	deltaFillAmount uint64,
) {
	l.Lock()
	defer l.Unlock()

	indexerOrderId := v1.OrderIdToIndexerOrderId(order.OrderId)
	currentFillAmount := l.FillAmounts[indexerOrderId]
	// Increment fill amount by the delta value
	l.FillAmounts[indexerOrderId] = currentFillAmount + deltaFillAmount
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
	delete(l.FillAmounts, orderId)
}
