package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

var (
	Transfer_Carl_Num0_Dave_Num0_Quote_500 = types.Transfer{
		Sender:    Carl_Num0,
		Recipient: Dave_Num0,
		AssetId:   lib.UsdcAssetId,
		Amount:    500_000_000, // $500
	}
	Transfer_Carl_Num0_Dave_Num0_Quote_600 = types.Transfer{
		Sender:    Carl_Num0,
		Recipient: Dave_Num0,
		AssetId:   lib.UsdcAssetId,
		Amount:    600_000_000, // $600
	}
	Transfer_Carl_Num0_Dave_Num0_Asset_600 = types.Transfer{
		Sender:    Carl_Num0,
		Recipient: Dave_Num0,
		AssetId:   lib.UsdcAssetId,
		Amount:    600_000_000, // $600
	}
	Transfer_Dave_Num0_Carl_Num0_Asset_500 = types.Transfer{
		Sender:    Dave_Num0,
		Recipient: Carl_Num0,
		AssetId:   lib.UsdcAssetId,
		Amount:    500_000_000, // $500
	}
	Transfer_Dave_Num0_Carl_Num0_Asset_500_GTB_20 = types.Transfer{
		Sender:    Dave_Num0,
		Recipient: Carl_Num0,
		AssetId:   lib.UsdcAssetId,
		Amount:    500_000_000, // $500
	}
)

// Test constants for deposit-to-subaccount messages.
var (
	MsgDepositToSubaccount_Alice_To_Alice_Num0_500 = types.MsgDepositToSubaccount{
		Sender:    AliceAccAddress.String(),
		Recipient: Alice_Num0,
		AssetId:   lib.UsdcAssetId,
		Quantums:  500_000_000, // $500
	}
	MsgDepositToSubaccount_Alice_To_Carl_Num0_750 = types.MsgDepositToSubaccount{
		Sender:    AliceAccAddress.String(),
		Recipient: Carl_Num0,
		AssetId:   lib.UsdcAssetId,
		Quantums:  750_000_000, // $750
	}
)

// Test constants for withdraw-from-subaccount messages.
var (
	MsgWithdrawFromSubaccount_Alice_Num0_To_Alice_500 = types.MsgWithdrawFromSubaccount{
		Sender:    Alice_Num0,
		Recipient: AliceAccAddress.String(),
		AssetId:   lib.UsdcAssetId,
		Quantums:  500_000_000, // $500
	}
	MsgWithdrawFromSubaccount_Carl_Num0_To_Alice_750 = types.MsgWithdrawFromSubaccount{
		Sender:    Carl_Num0,
		Recipient: AliceAccAddress.String(),
		AssetId:   lib.UsdcAssetId,
		Quantums:  750_000_000, // $750
	}
)
