package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest/types"
)

var _ types.QueryServer = Keeper{}
