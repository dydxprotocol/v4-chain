package types

import (
	"math/big"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
)

type SubaccountsKeeper interface {
	GetAllSubaccount(ctx sdk.Context) (list []Subaccount)
	GetRandomSubaccount(ctx sdk.Context, rand *rand.Rand) (Subaccount, error)
	GetNetCollateralAndMarginRequirements(
		ctx sdk.Context,
		update Update,
	) (
		risk margin.Risk,
		err error,
	)
	CanUpdateSubaccounts(
		ctx sdk.Context,
		updates []Update,
		updateType UpdateType,
	) (
		success bool,
		successPerUpdate []UpdateResult,
		err error,
	)
	UpdateSubaccounts(
		ctx sdk.Context,
		updates []Update,
		updateType UpdateType,
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
	TransferFundsFromSubaccountToSubaccount(
		ctx sdk.Context,
		senderSubaccountId SubaccountId,
		recipientSubaccountId SubaccountId,
		assetId uint32,
		quantums *big.Int,
	) (err error)
	SetSubaccount(ctx sdk.Context, subaccount Subaccount)
	GetSubaccount(
		ctx sdk.Context,
		id SubaccountId,
	) (val Subaccount)
	GetStreamSubaccountUpdate(
		ctx sdk.Context,
		id SubaccountId,
		snapshot bool,
	) (val StreamSubaccountUpdate)
	LegacyGetNegativeTncSubaccountSeenAtBlock(ctx sdk.Context) (uint32, bool)
	GetNegativeTncSubaccountSeenAtBlock(
		ctx sdk.Context,
		perpetualId uint32,
	) (uint32, bool, error)
	SetNegativeTncSubaccountSeenAtBlock(
		ctx sdk.Context,
		perpetualId uint32,
		blockHeight uint32,
	) error
}
