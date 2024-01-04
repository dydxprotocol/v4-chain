package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(
		context context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
}
