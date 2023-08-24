package types

import (
	fmt "fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CLOB module event types.
const (
	EventTypeMatch = "match"

	AttributeKeyTakerSubaccount                         = "taker_subaccount"
	AttributeKeyTakerSubaccountNumber                   = "taker_subaccount_number"
	AttributeKeyMakerSubaccount                         = "maker_subaccount"
	AttributeKeyMakerSubaccountNumber                   = "maker_subaccount_number"
	AttributeKeyTakerOrderFeeQuoteQuantums              = "taker_order_fee_quote_quantums"
	AttributeKeyMakerOrderFeeQuoteQuantums              = "maker_order_fee_quote_quantums"
	AttributeKeyMakerQuoteBalanceDeltaQuoteQuantums     = "maker_quote_balance_delta_quote_quantums"
	AttributeKeyTakerQuoteBalanceDeltaQuoteQuantums     = "taker_quote_balance_delta_quote_quantums"
	AttributeKeyMakerPerpetualQuantumsDeltaBaseQuantums = "maker_perpetual_quantums_delta_base_quantums"
	AttributeKeyTakerPerpetualQuantumsDeltaBaseQuantums = "taker_perpetual_quantums_delta_base_quantums"
	AttributeKeyInsuranceFundDeltaQuoteQuantums         = "insurance_fund_delta_quote_quantums"
	AttributeKeyIsLiquidation                           = "is_liquidation"
	AttributeKeyPerpetualId                             = "perpetual_id"
)

// NewCreateMatchEvent constructs a new match sdk.Event.
func NewCreateMatchEvent(
	taker satypes.SubaccountId,
	maker satypes.SubaccountId,
	takerOrderFee *big.Int,
	makerOrderFee *big.Int,
	takerQuoteBalanceDelta *big.Int,
	makerQuoteBalanceDelta *big.Int,
	takerPerpetualQuantumsDelta *big.Int,
	makerPerpetualQuantumsDelta *big.Int,
	insuranceFundDelta *big.Int,
	isLiquidation bool,
	perpetualId uint32,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeMatch,
		sdk.NewAttribute(AttributeKeyTakerSubaccount, taker.Owner),
		sdk.NewAttribute(AttributeKeyTakerSubaccountNumber, fmt.Sprint(taker.Number)),
		sdk.NewAttribute(AttributeKeyMakerSubaccount, maker.Owner),
		sdk.NewAttribute(AttributeKeyMakerSubaccountNumber, fmt.Sprint(maker.Number)),
		sdk.NewAttribute(AttributeKeyTakerOrderFeeQuoteQuantums, fmt.Sprint(takerOrderFee)),
		sdk.NewAttribute(AttributeKeyMakerOrderFeeQuoteQuantums, fmt.Sprint(makerOrderFee)),
		sdk.NewAttribute(AttributeKeyTakerQuoteBalanceDeltaQuoteQuantums, takerQuoteBalanceDelta.String()),
		sdk.NewAttribute(AttributeKeyMakerQuoteBalanceDeltaQuoteQuantums, makerQuoteBalanceDelta.String()),
		sdk.NewAttribute(AttributeKeyTakerPerpetualQuantumsDeltaBaseQuantums, takerPerpetualQuantumsDelta.String()),
		sdk.NewAttribute(AttributeKeyMakerPerpetualQuantumsDeltaBaseQuantums, makerPerpetualQuantumsDelta.String()),
		sdk.NewAttribute(AttributeKeyInsuranceFundDeltaQuoteQuantums, insuranceFundDelta.String()),
		sdk.NewAttribute(AttributeKeyIsLiquidation, fmt.Sprint(isLiquidation)),
		sdk.NewAttribute(AttributeKeyPerpetualId, fmt.Sprint(perpetualId)),
	)
}
