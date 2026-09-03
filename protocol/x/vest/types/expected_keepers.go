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

// AccountKeeper defines the expected interface for checking module account registration.
type AccountKeeper interface {
	// GetModuleAddress returns nil iff the module account is not registered.
	GetModuleAddress(name string) sdk.AccAddress
}

type BlockTimeKeeper interface {
	GetPreviousBlockInfo(ctx sdk.Context) blocktimetypes.BlockInfo
}
