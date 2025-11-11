package types

import (
	"context"
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type SubaccountsKeeper interface {
	GetAllSubaccount(ctx sdk.Context) (list []satypes.Subaccount)
	GetRandomSubaccount(ctx sdk.Context, rand *rand.Rand) (satypes.Subaccount, error)
	GetNetCollateralAndMarginRequirements(
		ctx sdk.Context,
		update satypes.Update,
	) (
		risk margin.Risk,
		err error,
	)
	CanUpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
		updateType satypes.UpdateType,
	) (
		success bool,
		successPerUpdate []satypes.UpdateResult,
		err error,
	)
	UpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
		updateType satypes.UpdateType,
	) (
		success bool,
		successPerUpdate []satypes.UpdateResult,
		err error,
	)
	DepositFundsFromAccountToSubaccount(
		ctx sdk.Context,
		fromAccount sdk.AccAddress,
		toSubaccountId satypes.SubaccountId,
		assetId uint32,
		amount *big.Int,
	) (err error)
	WithdrawFundsFromSubaccountToAccount(
		ctx sdk.Context,
		fromSubaccountId satypes.SubaccountId,
		toAccount sdk.AccAddress,
		assetId uint32,
		amount *big.Int,
	) (err error)
	TransferFundsFromSubaccountToSubaccount(
		ctx sdk.Context,
		senderSubaccountId satypes.SubaccountId,
		recipientSubaccountId satypes.SubaccountId,
		assetId uint32,
		quantums *big.Int,
	) (err error)
	SetSubaccount(ctx sdk.Context, subaccount satypes.Subaccount)
	GetSubaccount(
		ctx sdk.Context,
		id satypes.SubaccountId,
	) (val satypes.Subaccount)
}

// AccountKeeper defines the expected account keeper used for simulations.
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoins(
		ctx context.Context,
		fromAddr sdk.AccAddress,
		toAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
}
