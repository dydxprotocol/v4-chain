package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
)

var _ types.QueryServer = Keeper{}
