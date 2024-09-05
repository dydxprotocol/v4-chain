package ethqueryclienttypes

import (
	"context"
	"math/big"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/contract"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SDAIEthClient is an interface that encapsulates querying an Ethereum JSON-RPC endpoint.
type EthQueryClient interface {
	ChainID(ctx context.Context, client *ethclient.Client) (*big.Int, error)
	QueryDaiConversionRate(client *ethclient.Client) (string, error)
}

// EthQueryClientImpl is a concrete implementation of the EthQueryClient interface.
type EthQueryClientImpl struct{}

// ChainID wraps the existing ChainID function.
func (e *EthQueryClientImpl) ChainID(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	return client.ChainID(ctx)
}

// QueryDaiConversionRate wraps the existing QueryDaiConversionRate function.
func (e *EthQueryClientImpl) QueryDaiConversionRate(client *ethclient.Client) (string, error) {
	return store.QueryDaiConversionRate(client)
}
