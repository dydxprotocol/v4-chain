package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var _ types.QueryServer = Keeper{}
