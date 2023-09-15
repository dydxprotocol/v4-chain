package types

import (
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type SubaccountsKeeper interface {
	GetAllSubaccount(ctx sdk.Context) (list []satypes.Subaccount)
	GetRandomSubaccount(ctx sdk.Context, rand *rand.Rand) (satypes.Subaccount, error)
	GetNetCollateralAndMarginRequirements(
		ctx sdk.Context,
		update satypes.Update,
	) (
		bigNetCollateral *big.Int,
		bigInitialMargin *big.Int,
		bigMaintenanceMargin *big.Int,
		err error,
	)
	CanUpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
	) (
		success bool,
		successPerUpdate []satypes.UpdateResult,
		err error,
	)
	UpdateSubaccounts(
		ctx sdk.Context,
		updates []satypes.Update,
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
	SetSubaccount(ctx sdk.Context, subaccount satypes.Subaccount)
	GetSubaccount(
		ctx sdk.Context,
		id satypes.SubaccountId,
	) (val satypes.Subaccount)
}

// AccountKeeper defines the expected account keeper used for simulations.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	SetAccount(ctx sdk.Context, acc types.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
