package client

import (
	"context"

	"github.com/cometbft/cometbft/rpc/client/mock"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
)

type MockClient struct {
	mock.Client
	Err error
}

func (c MockClient) BroadcastTxSync(
	ctx context.Context,
	tx tmtypes.Tx,
) (*ctypes.ResultBroadcastTx, error) {
	return nil, c.Err
}
