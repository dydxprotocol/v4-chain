package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
)

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

type BlockTimeKeeper interface {
	GetPreviousBlockInfo(ctx sdk.Context) blocktimetypes.BlockInfo
}
