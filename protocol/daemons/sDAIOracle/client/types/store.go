package types

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient is an interface that encapsulates querying an Ethereum JSON-RPC endpoint.
type StoreMock interface {
	QueryDaiConversionRate(client *ethclient.Client) (string, string, error)
}
