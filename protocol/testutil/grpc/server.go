package grpc

import pricetypes "github.com/dydxprotocol/v4/x/prices/types"

type QueryServer interface {
	pricetypes.QueryServer
}
