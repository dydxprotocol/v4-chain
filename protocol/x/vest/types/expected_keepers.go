package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

type BlockTimeKeeper interface {
	GetPreviousBlockInfo(ctx sdk.Context) blocktimetypes.BlockInfo
}
