package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

var _ types.QueryServer = Keeper{}
