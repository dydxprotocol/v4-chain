package grpc

import pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"

type QueryServer interface {
	pricetypes.QueryServer
}
