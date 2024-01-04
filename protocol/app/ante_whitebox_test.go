package app_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobante "github.com/dydxprotocol/v4-chain/protocol/x/clob/ante"
	"github.com/stretchr/testify/require"
)

func newTestHandlerOptions(t *testing.T) app.HandlerOptions {
	tApp := testApp.NewTestAppBuilder(t).Build()
	tApp.InitChain()

	return app.HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   tApp.App.AccountKeeper,
			BankKeeper:      tApp.App.BankKeeper,
			SignModeHandler: tApp.App.TxConfig().SignModeHandler(),
			FeegrantKeeper:  tApp.App.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		ClobKeeper: tApp.App.ClobKeeper,
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
	handlerOptions := newTestHandlerOptions(t)
	decoratorChain := app.NewAnteDecoratorChain(handlerOptions)
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
