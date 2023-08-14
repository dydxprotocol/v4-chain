package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

var _ types.QueryServer = Keeper{}
