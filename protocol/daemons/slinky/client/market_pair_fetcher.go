package client

import (
	"context"
	"cosmossdk.io/log"
	"fmt"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"google.golang.org/grpc"
	"sync"
)

// MarketPairFetcher is a lightweight process run in a goroutine by the slinky client.
// Its purpose is to periodically query the prices module for MarketParams and create
// an easily indexed mapping between Slinky's CurrencyPair type and the corresponding ID
// in the x/prices module.
type MarketPairFetcher struct {
	logger            log.Logger
	queryConn         *grpc.ClientConn
	pricesQueryClient pricetypes.QueryClient

	// compatMappings stores a mapping between CurrencyPair and the corresponding market(param|price) ID
	compatMappings map[oracletypes.CurrencyPair]uint32
	compatMu       sync.RWMutex
}

func NewMarketPairFetcher(logger log.Logger) *MarketPairFetcher {
	return &MarketPairFetcher{
		logger:         logger,
		compatMappings: make(map[oracletypes.CurrencyPair]uint32),
	}
}

// Start opens the grpc connections for the fetcher.
func (m *MarketPairFetcher) Start(
	ctx context.Context,
	appFlags appflags.Flags,
	grpcClient daemontypes.GrpcClient) error {
	// Create the query client connection
	queryConn, err := grpcClient.NewTcpConnection(ctx, appFlags.GrpcAddress)
	if err != nil {
		m.logger.Error(
			"Failed to establish gRPC connection",
			"gRPC address", appFlags.GrpcAddress,
			"error", err,
		)
		return err
	}
	m.pricesQueryClient = pricetypes.NewQueryClient(queryConn)
	return nil
}

// Stop closes all existing connections.
func (m *MarketPairFetcher) Stop() {
	if m.queryConn != nil {
		_ = m.queryConn.Close()
	}
}

// GetIDForPair returns the ID corresponding to the currency pair in the x/prices module.
// If the currency pair is not found it will return an error.
func (m *MarketPairFetcher) GetIDForPair(cp oracletypes.CurrencyPair) (uint32, error) {
	var id uint32
	m.compatMu.RLock()
	defer m.compatMu.RUnlock()
	id, found := m.compatMappings[cp]
	if !found {
		return id, fmt.Errorf("pair %s not found in compatMappings", cp.String())
	}
	return id, nil
}

// FetchIdMappings is run periodically to refresh the cache of known mappings between
// CurrencyPair and MarketParam ID.
func (m *MarketPairFetcher) FetchIdMappings(ctx context.Context) error {
	// fetch all market params
	resp, err := m.pricesQueryClient.AllMarketParams(ctx, &pricetypes.QueryAllMarketParamsRequest{})
	if err != nil {
		return err
	}
	// Exit early if there are no changes
	// This assumes there will not be an addition and a removal of markets in the same block
	if len(resp.MarketParams) == len(m.compatMappings) {
		return nil
	}
	var compatMappings = make(map[oracletypes.CurrencyPair]uint32, len(resp.MarketParams))
	for _, mp := range resp.MarketParams {
		cp, err := slinky.MarketPairToCurrencyPair(mp.Pair)
		if err != nil {
			return err
		}
		m.logger.Info("Mapped market to pair", "market id", mp.Id, "currency pair", cp.String())
		compatMappings[cp] = mp.Id
	}
	m.compatMu.Lock()
	defer m.compatMu.Unlock()
	m.compatMappings = compatMappings
	return nil
}
