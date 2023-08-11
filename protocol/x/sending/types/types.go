package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SendingKeeper interface {
	ProcessTransfer(ctx sdk.Context, transfer *Transfer) error
	ProcessDepositToSubaccount(
		ctx sdk.Context,
		msgDepositToSubaccount *MsgDepositToSubaccount,
	) error
	ProcessWithdrawFromSubaccount(
		ctx sdk.Context,
		msgWithdrawFromSubaccount *MsgWithdrawFromSubaccount,
	) error
}
