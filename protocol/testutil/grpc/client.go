package grpc

import (
	liquidationtypes "github.com/dydxprotocol/v4/daemons/liquidation/api"
	pricefeedtypes "github.com/dydxprotocol/v4/daemons/pricefeed/api"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

type QueryClient interface {
	satypes.QueryClient
	clobtypes.QueryClient
	liquidationtypes.LiquidationServiceClient
	pricefeedtypes.PriceFeedServiceClient
}
