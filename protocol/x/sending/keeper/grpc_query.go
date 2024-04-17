package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
)

var _ types.QueryServer = Keeper{}
