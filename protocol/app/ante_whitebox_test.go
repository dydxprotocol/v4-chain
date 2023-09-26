package app

import (
	"reflect"
	"testing"

	delaymsgmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	clobante "github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	clobmodulekeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	clobmodulememclob "github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func newTestHandlerOptions() HandlerOptions {
	encodingConfig := GetEncodingConfig()
	appCodec := encodingConfig.Codec
	txConfig := encodingConfig.TxConfig

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		nil,
		authtypes.ProtoBaseAccount,
		nil,
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		nil,
		accountKeeper,
		BlockedAddresses(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	feeGrantKeeper := feegrantkeeper.NewKeeper(appCodec, nil, accountKeeper)

	memClob := clobmodulememclob.NewMemClobPriceTimePriority(false)
	clobKeeper := clobmodulekeeper.NewKeeper(
		appCodec,
		nil,
		nil,
		nil,
		[]string{
			authtypes.NewModuleAddress(delaymsgmoduletypes.ModuleName).String(),
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		},
		memClob,
		nil,
		nil,
		nil,
		bankKeeper,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		flags.GetDefaultClobFlags(),
		rate_limit.NewNoOpRateLimiter[*types.MsgPlaceOrder](),
		rate_limit.NewNoOpRateLimiter[*types.MsgCancelOrder](),
	)
	return HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   accountKeeper,
			BankKeeper:      bankKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			FeegrantKeeper:  feeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		ClobKeeper: clobKeeper,
	}
}

func wrapDecoratorStr(decoratorStr string) string {
	return "(" + decoratorStr + ")"
}

func humanReadableDecoratorTypes(decoratorChain []sdk.AnteDecorator) []string {
	var dTypes []string

	for _, decorator := range decoratorChain {
		decoratorType := reflect.TypeOf(decorator).String()
		switch decorator := decorator.(type) {
		case libante.AppInjectedMsgAnteWrapper:
			switch nestedDecorator := decorator.GetAnteHandler().(type) {
			// The SingleMsgClobTxnAnteWrapper has an additional layer of nesting
			case clobante.SingleMsgClobTxAnteWrapper:
				dTypes = append(dTypes, decoratorType+
					wrapDecoratorStr(
						reflect.TypeOf(nestedDecorator).String()+
							wrapDecoratorStr(
								reflect.TypeOf(nestedDecorator.GetAnteHandler()).String(),
							),
					),
				)
			// The ShortTermSingleMsgClobTxnAnteWrapper has an additional layer of nesting
			case clobante.ShortTermSingleMsgClobTxAnteWrapper:
				dTypes = append(dTypes, decoratorType+
					wrapDecoratorStr(
						reflect.TypeOf(nestedDecorator).String()+
							wrapDecoratorStr(
								reflect.TypeOf(nestedDecorator.GetAnteHandler()).String(),
							),
					),
				)
			default:
				dTypes = append(dTypes, decoratorType+wrapDecoratorStr(reflect.TypeOf(nestedDecorator).String()))
			}
		default:
			dTypes = append(dTypes, decoratorType)
		}
	}
	return dTypes
}

func TestAnteHandlerChainOrder_Valid(t *testing.T) {
	handlerOptions := newTestHandlerOptions()
	decoratorChain := newAnteDecoratorChain(handlerOptions)
	decoratorTypes := humanReadableDecoratorTypes(decoratorChain)

	expectedDecoratorTypes := []string{
		"ante.AppInjectedMsgAnteWrapper(ante.SingleMsgClobTxAnteWrapper(ante.SetUpContextDecorator))",
		"ante.FreeInfiniteGasDecorator",
		"ante.RejectExtensionOptionsDecorator",
		"ante.ValidateMsgTypeDecorator",
		"ante.AppInjectedMsgAnteWrapper(ante.ValidateBasicDecorator)",
		"ante.TxTimeoutHeightDecorator",
		"ante.ValidateMemoDecorator",
		"ante.ConsumeTxSizeGasDecorator",
		"ante.AppInjectedMsgAnteWrapper(ante.SingleMsgClobTxAnteWrapper(ante.DeductFeeDecorator))",
		"ante.AppInjectedMsgAnteWrapper(ante.SetPubKeyDecorator)",
		"ante.ValidateSigCountDecorator",
		"ante.AppInjectedMsgAnteWrapper(ante.SigGasConsumeDecorator)",
		"ante.AppInjectedMsgAnteWrapper(ante.SigVerificationDecorator)",
		"ante.AppInjectedMsgAnteWrapper(ante.ShortTermSingleMsgClobTxAnteWrapper(ante.IncrementSequenceDecorator))",
		"ante.ClobRateLimitDecorator",
		"ante.ClobDecorator",
	}

	require.Equal(t, expectedDecoratorTypes, decoratorTypes, "Decorator order does not match expected")
}
