package grpc

import (
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	liquidationtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// QueryClient combines all the query clients used in testing into a single mock interface for testing convenience.
type QueryClient interface {
	blocktimetypes.QueryClient
	satypes.QueryClient
	clobtypes.QueryClient
	perptypes.QueryClient
	pricetypes.QueryClient
	bridgetypes.BridgeServiceClient
	liquidationtypes.LiquidationServiceClient
	pricefeedtypes.PriceFeedServiceClient
}
