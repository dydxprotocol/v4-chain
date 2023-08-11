package mempool

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
)

// TODO(DEC-1194): this is copied from SDK. Remove this forked impl in favor of SDK's NoOpMempool.
var _ mempool.Mempool = (*noOpMempool)(nil)

// NoOpMempool defines a no-op mempool. Transactions are completely discarded and
// ignored when BaseApp interacts with the mempool.
//
// Note: When this mempool is used, it assumed that an application will rely
// on Tendermint's transaction ordering defined in `RequestPrepareProposal`, which
// is FIFO-ordered by default.
type noOpMempool struct{}

func NewNoOpMempool() mempool.Mempool {
	return &noOpMempool{}
}

func (noOpMempool) Insert(context.Context, sdk.Tx) error              { return nil }
func (noOpMempool) Select(context.Context, [][]byte) mempool.Iterator { return nil }
func (noOpMempool) CountTx() int                                      { return 0 }
func (noOpMempool) Remove(sdk.Tx) error                               { return nil }
