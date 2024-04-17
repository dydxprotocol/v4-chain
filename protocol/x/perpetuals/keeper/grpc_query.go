package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
)

var _ types.QueryServer = Keeper{}
