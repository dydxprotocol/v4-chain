package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
)

var _ types.QueryServer = Keeper{}
