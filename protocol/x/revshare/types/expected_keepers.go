package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

type PerpetualsKeeper interface {
	// Function to get perpetual details from the module store
	GetPerpetual(ctx sdk.Context, perpetualId uint32) (perpetual perptypes.Perpetual, err error)
}
