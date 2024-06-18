package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

var _ types.QueryServer = Keeper{}
