package types

import (
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SubaccountsKeeper interface {
	GetAllSubaccount(ctx sdk.Context) (list []Subaccount)
	GetRandomSubaccount(ctx sdk.Context, rand *rand.Rand) (Subaccount, error)
	GetNetCollateralAndMarginRequirements(
		ctx sdk.Context,
		update Update,
	) (
		bigNetCollateral *big.Int,
		bigInitialMargin *big.Int,
		bigMaintenanceMargin *big.Int,
		err error,
	)
	CanUpdateSubaccounts(
		ctx sdk.Context,
		updates []Update,
	) (
		success bool,
		successPerUpdate []UpdateResult,
		err error,
	)
	UpdateSubaccounts(
		ctx sdk.Context,
		updates []Update,
	) (
		success bool,
		successPerUpdate []UpdateResult,
		err error,
	)
	DepositFundsFromAccountToSubaccount(
		ctx sdk.Context,
		fromAccount sdk.AccAddress,
		toSubaccountId SubaccountId,
		assetId uint32,
		amount *big.Int,
	) (err error)
	WithdrawFundsFromSubaccountToAccount(
		ctx sdk.Context,
		fromSubaccountId SubaccountId,
		toAccount sdk.AccAddress,
		assetId uint32,
		amount *big.Int,
	) (err error)
	SetSubaccount(ctx sdk.Context, subaccount Subaccount)
	GetSubaccount(
		ctx sdk.Context,
		id SubaccountId,
	) (val Subaccount)
}
