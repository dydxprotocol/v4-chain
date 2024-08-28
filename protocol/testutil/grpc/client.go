package grpc

import (
	deleveragingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/api"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/api"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// QueryClient combines all the query clients used in testing into a single mock interface for testing convenience.
type QueryClient interface {
	blocktimetypes.QueryClient
	satypes.QueryClient
	clobtypes.QueryClient
	perptypes.QueryClient
	pricetypes.QueryClient
	deleveragingtypes.DeleveragingServiceClient
	pricefeedtypes.PriceFeedServiceClient
}
