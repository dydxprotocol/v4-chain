package types

import (
	"context"
	"math/big"
)

// EthClient is an interface that encapsulates querying an Ethereum JSON-RPC endpoint.
type EthClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
}
