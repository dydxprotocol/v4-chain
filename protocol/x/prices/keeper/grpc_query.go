package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

var _ types.QueryServer = Keeper{}
