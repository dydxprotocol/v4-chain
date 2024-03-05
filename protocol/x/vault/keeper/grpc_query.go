package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

var _ types.QueryServer = Keeper{}
