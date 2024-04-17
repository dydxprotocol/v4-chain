package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/types"
)

var _ types.QueryServer = Keeper{}
