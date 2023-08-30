package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
)

type EthClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethcoretypes.Log, error)
}
