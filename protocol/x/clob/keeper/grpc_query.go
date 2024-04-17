package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
)

var _ types.QueryServer = Keeper{}
